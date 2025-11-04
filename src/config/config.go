package config

import (
	"net/http"
	"os"
	"strconv"
	"time"
)

type ApiRoute struct {
	Endpoint string
	Method   string
}

type Config struct {
	AccessTokenExpiresIn time.Duration
	AppEmail             string
	ApiUrl               string
	ClientUrl            string
	CloudinaryName       string
	CloudinaryKey        string
	CloudinarySecret     string
	CookieDomain         string
	CurrentUser          string
	CurrentUserId        string
	CurrentUserRole      string
	GinMode              string
	GoogleAuthId         string
	GoogleAuthSecret     string
	IsDevMode            bool
	JWTSecret            []byte
	MaxImageSize         int
	NonAuthRoutes        []ApiRoute
	Port                 string
	PostgresDbUrl        string
	SmtpHost             string
	SmtpPassword         string
	SmtpPort             int
	SmtpUser             string
	Version              string
}

var AppConfig *Config

func InitializeConfig() {
	AppConfig = &Config{
		AccessTokenExpiresIn: time.Hour * 24 * 7,
		AppEmail:             os.Getenv("APP_EMAIL"),
		ApiUrl:               os.Getenv("API_URL"),
		ClientUrl:            os.Getenv("CLIENT_URL"),
		CloudinaryName:       os.Getenv("CLOUDINARY_NAME"),
		CloudinaryKey:        os.Getenv("CLOUDINARY_KEY"),
		CloudinarySecret:     os.Getenv("CLOUDINARY_SECRET"),
		CookieDomain:         os.Getenv("COOKIE_DOMAIN"),
		CurrentUser:          "CURRENT_USER",
		CurrentUserId:        "CURRENT_USER_ID",
		CurrentUserRole:      "CURRENT_USER_ROLE",
		GinMode:              os.Getenv("GIN_MODE"),
		GoogleAuthId:         os.Getenv("GOOGLE_AUTH_ID"),
		GoogleAuthSecret:     os.Getenv("GOOGLE_AUTH_SECRET"),
		IsDevMode:            os.Getenv("IS_DEV_MODE") == "true",
		JWTSecret:            []byte(os.Getenv("JWT_SECRET")),
		MaxImageSize:         1024 * 1024 * 5,
		Port:                 os.Getenv("PORT"),
		PostgresDbUrl:        os.Getenv("POSTGRES_DB_URL"),
		SmtpHost:             os.Getenv("SMTP_HOST"),
		SmtpPassword:         os.Getenv("SMTP_PASSWORD"),
		SmtpPort:             func() int { port, _ := strconv.Atoi(os.Getenv("SMTP_PORT")); return port }(),
		SmtpUser:             os.Getenv("SMTP_USER"),
		Version:              os.Getenv("VERSION"),
		NonAuthRoutes: []ApiRoute{
			{Endpoint: "/api/v1", Method: http.MethodGet},
			{Endpoint: "/api/v1/health", Method: http.MethodGet},
			{Endpoint: "/api/v1/auth/signin", Method: http.MethodPost},
			{Endpoint: "/api/v1/auth/google/callback", Method: http.MethodGet},
		},
	}
}
