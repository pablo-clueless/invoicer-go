package services

import (
	"invoicer-go/m/src/dto"
	"invoicer-go/m/src/models"

	"gorm.io/gorm"
)

type UserService struct {
	database *gorm.DB
}

func NewUserService(database *gorm.DB) *UserService {
	return &UserService{
		database: database,
	}
}

func (s *UserService) UpdateUser(dto.UpdateUserDto) error {
	return nil
}

func (s *UserService) DeleteUser(id string) error {
	return nil
}

func (s *UserService) GetUser(id string) (*models.User, error) {
	return nil, nil
}
