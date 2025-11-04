package services

import (
	"errors"
	"invoicer-go/m/src/dto"
	"invoicer-go/m/src/models"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrEmailSendFailed = errors.New("failed to send email")
)

type AuthService struct {
	database *gorm.DB
}

func NewAuthService(database *gorm.DB) *AuthService {
	return &AuthService{
		database: database,
	}
}

func (s *AuthService) Signin(payload dto.CreateUserDto) error {
	return nil
}

func (s *AuthService) FindUserById(id string) (*models.User, error) {
	var user models.User
	if err := s.database.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
