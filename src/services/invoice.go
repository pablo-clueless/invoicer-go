package services

import (
	"invoicer-go/m/src/dto"
	"invoicer-go/m/src/models"

	"gorm.io/gorm"
)

type InvoiceService struct {
	database *gorm.DB
}

func NewInvoiceService(database *gorm.DB) *InvoiceService {
	return &InvoiceService{
		database: database,
	}
}

func (s *InvoiceService) CreateInvoice(payload dto.CreateInvoiceDto) (*models.Invoice, error) {
	return nil, nil
}

func (s *InvoiceService) UpdateInvoice(id string, payload dto.CreateInvoiceDto) (*models.Invoice, error) {
	return nil, nil
}

func (s *InvoiceService) DeleteInvoice(id string) error {
	return nil
}

func (s *InvoiceService) GetInvoices(params dto.Pagination) *dto.PaginatedResponse[models.Invoice] {
	return nil
}

func (s *InvoiceService) GetInvoice(id string) (*models.Invoice, error) {
	return nil, nil
}
