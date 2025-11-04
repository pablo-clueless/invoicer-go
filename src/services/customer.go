package services

import "gorm.io/gorm"

type CustomerService struct {
	database *gorm.DB
}

func NewCustomerService(database *gorm.DB) *CustomerService {
	return &CustomerService{
		database: database,
	}
}

func (s *CustomerService) CreateCustomer() {}

func (s *CustomerService) UpdateCustomer() {}

func (s *CustomerService) DeleteCustomer() {}

func (s *CustomerService) GetCustomers() {}

func (s *CustomerService) GetCustomer() {}
