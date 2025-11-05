package services

import (
	"errors"
	"fmt"
	"invoicer-go/m/src/config"
	"invoicer-go/m/src/lib"
	"invoicer-go/m/src/models"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"gorm.io/gorm"
)

var (
	ErrCustomerNotFound = errors.New("customer not found")
	ErrEmailSendFailed  = errors.New("failed to send email")
	ErrInvalidProvider  = errors.New("invalid provider")
	ErrInvoiceNotFound  = errors.New("invoice not found")
	ErrRecordExists     = errors.New("this record exists already")
	ErrUserNotFound     = errors.New("user not found")
)

type AuthService struct {
	database *gorm.DB
}

func NewAuthService(database *gorm.DB) *AuthService {
	return &AuthService{
		database: database,
	}
}

type SigninResponse struct {
	RedirectUrl string      `json:"redirectUrl"`
	Token       string      `json:"token"`
	User        models.User `json:"user"`
}

func InitializeProvider() error {
	appConfig := config.AppConfig

	if appConfig.GoogleClientId == "" || appConfig.GoogleClientSecret == "" {
		return errors.New("google OAuth credentials are required")
	}

	store := sessions.NewCookieStore([]byte(appConfig.GoogleClientId))

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   int(30 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   !appConfig.IsDevMode,
		SameSite: http.SameSiteLaxMode,
	}

	gothic.Store = store

	base := strings.TrimSuffix(appConfig.ApiUrl, "/")
	version := strings.Trim(appConfig.Version, "/")
	if version != "" && !strings.HasSuffix(base, "/"+version) {
		base = base + "/" + version
	}

	goth.UseProviders(
		google.New(
			appConfig.GoogleClientId,
			appConfig.GoogleClientSecret,
			fmt.Sprintf("%s/auth/google/callback", base),
			"email", "profile",
		),
	)

	return nil
}

func (s *AuthService) GetFrontendRedirectURL(token, userID string) string {
	frontendURL := strings.TrimSuffix(config.AppConfig.ClientUrl, "/")
	redirectURL := fmt.Sprintf("%s/auth/success", frontendURL)

	params := url.Values{}
	params.Add("token", token)
	params.Add("user_id", userID)

	return redirectURL + "?" + params.Encode()
}

func (s *AuthService) GetFrontendErrorRedirectURL(message string) string {
	frontendURL := strings.TrimSuffix(config.AppConfig.ClientUrl, "/")

	redirectURL := fmt.Sprintf("%s/auth/error", frontendURL)

	params := url.Values{}
	params.Add("error", message)

	return redirectURL + "?" + params.Encode()
}

func (s *AuthService) HandleOAuthCallback(res http.ResponseWriter, req *http.Request) (string, error) {
	gothUser, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		return "", fmt.Errorf("failed to complete OAuth: %w", err)
	}

	signinResponse, err := s.SigninWithOauth(&gothUser)
	if err != nil {
		return "", fmt.Errorf("failed to sign in with OAuth: %w", err)
	}

	redirectURL := s.GetFrontendRedirectURL(signinResponse.Token, signinResponse.User.ID.String())
	return redirectURL, nil
}

func (s *AuthService) SigninWithOauth(payload *goth.User) (*SigninResponse, error) {
	if payload.Email == "" {
		return nil, errors.New("email is required from OAuth provider")
	}

	var user *models.User
	var err error

	user, err = s.FindUserByEmail(payload.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			user, err = s.createUserFromOAuth(payload)
			if err != nil {
				return nil, fmt.Errorf("failed to create user: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to find user: %w", err)
		}
	} else {

		user = s.updateUserFromOAuth(user, payload)
		if err = s.database.Save(user).Error; err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
	}

	token, err := lib.GenerateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &SigninResponse{
		RedirectUrl: s.GetFrontendRedirectURL(token, user.ID.String()),
		User:        *user,
		Token:       token,
	}, nil
}

func (s *AuthService) createUserFromOAuth(payload *goth.User) (*models.User, error) {
	user := &models.User{
		Email:    strings.ToLower(strings.TrimSpace(payload.Email)),
		Name:     s.formatUserName(payload.FirstName, payload.LastName),
		Provider: payload.Provider,
	}

	if err := s.database.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user in database: %w", err)
	}

	return user, nil
}

func (s *AuthService) updateUserFromOAuth(user *models.User, payload *goth.User) *models.User {

	newName := s.formatUserName(payload.FirstName, payload.LastName)
	if user.Name == "" || user.Name != newName {
		user.Name = newName
	}

	if user.Provider == "" {
		user.Provider = payload.Provider
	}

	return user
}

func (s *AuthService) formatUserName(firstName, lastName string) string {
	name := strings.TrimSpace(fmt.Sprintf("%s %s", firstName, lastName))
	if name == "" {
		return "User"
	}
	return name
}

func (s *AuthService) FindUserByEmail(email string) (*models.User, error) {
	var user models.User

	err := s.database.Where("LOWER(email) = LOWER(?)", strings.TrimSpace(email)).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("database error while finding user by email: %w", err)
	}

	return &user, nil
}

func (s *AuthService) FindUserById(id string) (*models.User, error) {
	var user models.User

	err := s.database.Where("id = ?", strings.TrimSpace(id)).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("database error while finding user by ID: %w", err)
	}

	return &user, nil
}

func (s *AuthService) GetUserProfile(userID string) (*models.User, error) {
	user, err := s.FindUserById(userID)
	if err != nil {
		return nil, err
	}

	user.Provider = ""

	return user, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*models.User, error) {
	claims, err := lib.ValidateToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	user, err := s.FindUserById(claims.ID)
	if err != nil {
		return nil, fmt.Errorf("user not found for token: %w", err)
	}

	return user, nil
}

func (s *AuthService) SignOut(userID string) error {
	_, err := s.FindUserById(userID)
	if err != nil {
		return err
	}
	store := sessions.NewCookieStore([]byte(config.AppConfig.GoogleClientId))
	session, err := store.Get(nil, gothic.SessionName)
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1
	return store.Save(nil, nil, session)
}
