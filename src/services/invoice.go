package services

import "gorm.io/gorm"

type InvoiceService struct {
	database *gorm.DB
}

func NewInvoiceService(database *gorm.DB) *InvoiceService {
	return &InvoiceService{
		database: database,
	}
}

func (s *InvoiceService) CreateInvoice() {}

func (s *InvoiceService) UpdateInvoice() {}

func (s *InvoiceService) DeleteInvoice() {}

func (s *InvoiceService) GetInvoices() {}

func (s *InvoiceService) GetInvoice() {}
