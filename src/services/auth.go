package services

import (
	"errors"
	"fmt"
	"invoicer-go/m/src/config"
	"invoicer-go/m/src/models"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrEmailSendFailed = errors.New("failed to send email")
	ErrInvalidProvider = errors.New("invalid provider")
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
	User  models.User `json:"user"`
	Token string      `json:"token"`
}

func InitializeProvider() {
	appConfig := config.AppConfig

	store := sessions.NewCookieStore([]byte(appConfig.GoogleAuthId))
	store.Options.SameSite = http.SameSiteNoneMode
	store.MaxAge(int(time.Hour) * 30)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = !appConfig.IsDevMode

	gothic.Store = store
	goth.UseProviders(
		google.New(
			appConfig.GoogleAuthId,
			appConfig.GoogleAuthSecret,
			fmt.Sprintf("%s/auth/google/callback", appConfig.ApiUrl),
			"email", "profile"),
	)
}

func (s *AuthService) SigninWithOauth(payload *goth.User) (*SigninResponse, error) {
	if payload.Provider != "google" {
		return nil, ErrInvalidProvider
	}

	return nil, nil
}

func (s *AuthService) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := s.database.First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) FindUserById(id string) (*models.User, error) {
	var user models.User
	if err := s.database.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}
