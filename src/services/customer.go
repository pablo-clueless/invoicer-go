package services

import (
	"invoicer-go/m/src/dto"
	"invoicer-go/m/src/models"

	"gorm.io/gorm"
)

type CustomerService struct {
	database *gorm.DB
}

func NewCustomerService(database *gorm.DB) *CustomerService {
	return &CustomerService{
		database: database,
	}
}

func (s *CustomerService) CreateCustomer(payload dto.CreateCustomerDto) (*models.Customer, error) {
	return nil, nil
}

func (s *CustomerService) UpdateCustomer(id string, payload dto.CreateCustomerDto) (*models.Customer, error) {
	return nil, nil
}

func (s *CustomerService) DeleteCustomer(id string) error {
	return nil
}

func (s *CustomerService) GetCustomers(params *dto.Pagination) *dto.PaginatedResponse[models.Customer] {
	return nil
}

func (s *CustomerService) GetCustomer(id string) (*models.Customer, error) {
	return nil, nil
}
