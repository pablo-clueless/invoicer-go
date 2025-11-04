package services

import (
	"errors"
	"invoicer-go/m/src/dto"
	"invoicer-go/m/src/models"
	"strings"

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
	existingCustomer, err := s.FindCustomerByEmail(payload.Email)
	if err != nil && !errors.Is(err, ErrCustomerNotFound) {
		return nil, err
	}
	if existingCustomer != nil {
		return nil, ErrRecordExists
	}

	newCustomer := &models.Customer{
		Name:  payload.Name,
		Email: payload.Email,
		Phone: payload.Phone,
	}

	if err := s.database.Create(newCustomer).Error; err != nil {
		return nil, err
	}

	return newCustomer, nil
}

func (s *CustomerService) UpdateCustomer(id string, payload dto.UpdateCustomerDto) (*models.Customer, error) {
	customer, err := s.FindCustomerById(id)
	if err != nil {
		return nil, err
	}

	if payload.Name != nil {
		customer.Name = *payload.Name
	}
	if payload.Phone != nil {
		customer.Phone = *payload.Phone
	}

	if err := s.database.Save(customer).Error; err != nil {
		return nil, err
	}

	return customer, nil
}

func (s *CustomerService) DeleteCustomer(id string) error {
	customer, err := s.FindCustomerById(id)
	if err != nil {
		return err
	}

	return s.database.Transaction(func(tx *gorm.DB) error {
		var invoiceCount int64
		if err := tx.Model(&models.Invoice{}).Where("customer_id = ?", id).Count(&invoiceCount).Error; err != nil {
			return err
		}

		if invoiceCount > 0 {
			return errors.New("cannot delete customer with existing invoices")
		}

		return tx.Delete(customer).Error
	})
}

func (s *CustomerService) GetCustomers(params dto.CustomerPagination) (*dto.PaginatedResponse[models.Customer], error) {
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	var customers []models.Customer
	var totalItems int64

	query := s.database.Model(&models.Customer{})

	if params.Name != nil && strings.TrimSpace(*params.Name) != "" {
		namePattern := "%" + strings.ToLower(strings.TrimSpace(*params.Name)) + "%"
		query = query.Where("LOWER(name) LIKE ?", namePattern)
	}

	if params.Email != nil && strings.TrimSpace(*params.Email) != "" {
		emailPattern := "%" + strings.ToLower(strings.TrimSpace(*params.Email)) + "%"
		query = query.Where("LOWER(email) LIKE ?", emailPattern)
	}

	if err := query.Count(&totalItems).Error; err != nil {
		return &dto.PaginatedResponse[models.Customer]{
			Data:       []models.Customer{},
			Limit:      params.Limit,
			Page:       params.Page,
			TotalItems: 0,
			TotalPages: 0,
		}, err
	}

	offset := (params.Page - 1) * params.Limit

	if err := query.Offset(offset).
		Limit(params.Limit).
		Order("created_at DESC").
		Find(&customers).Error; err != nil {
		return nil, err
	}

	totalPages := 0
	if totalItems > 0 && params.Limit > 0 {
		totalPages = int((totalItems + int64(params.Limit) - 1) / int64(params.Limit))
	}

	return &dto.PaginatedResponse[models.Customer]{
		Data:       customers,
		TotalItems: int(totalItems),
		TotalPages: totalPages,
		Page:       params.Page,
		Limit:      params.Limit,
	}, nil
}

func (s *CustomerService) GetCustomer(id string) (*models.Customer, error) {
	return s.FindCustomerById(id)
}

func (s *CustomerService) FindCustomerByEmail(email string) (*models.Customer, error) {
	customer := &models.Customer{}
	err := s.database.Where("email = ?", email).First(customer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCustomerNotFound
		}
		return nil, err
	}

	return customer, nil
}

func (s *CustomerService) FindCustomerById(id string) (*models.Customer, error) {
	customer := &models.Customer{}
	err := s.database.Where("id = ?", id).First(customer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCustomerNotFound
		}
		return nil, err
	}

	return customer, nil
}

func (s *CustomerService) FindCustomersByIds(ids []string) ([]models.Customer, error) {
	var customers []models.Customer
	if err := s.database.Where("id IN ?", ids).Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}
