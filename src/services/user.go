package services

import "gorm.io/gorm"

type UserService struct {
	database *gorm.DB
}

func NewUserService(database *gorm.DB) *UserService {
	return &UserService{
		database: database,
	}
}

func (s *UserService) UpdateUser() {}

func (s *UserService) DeleteUser() {}

func (s *UserService) GetUser() {}
