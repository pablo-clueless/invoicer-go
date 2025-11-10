package services

import (
	"errors"
	"invoicer-go/m/src/dto"
	"invoicer-go/m/src/models"
	"strings"

	"github.com/google/uuid"
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

var (
	ErrInvoiceTitleExists = errors.New("an invoice with this title already exists")
)

func (s *InvoiceService) calculateInvoiceTotals(invoice *models.Invoice) {
	invoice.SubTotal = 0
	for i := range invoice.Items {
		invoice.Items[i].LineTotal = float64(invoice.Items[i].Quantity) * invoice.Items[i].Price
		invoice.SubTotal += invoice.Items[i].LineTotal
	}

	var discountAmount, taxAmount float64

	switch invoice.DiscountType {
	case models.Fixed:
		discountAmount = invoice.Discount
	case models.Percentage:
		discountAmount = invoice.Discount * invoice.SubTotal / 100
	}

	switch invoice.TaxType {
	case models.Fixed:
		taxAmount = invoice.Tax
	case models.Percentage:
		taxAmount = invoice.Tax * invoice.SubTotal / 100
	}

	invoice.Total = invoice.SubTotal + taxAmount - discountAmount
}

func (s *InvoiceService) CreateInvoice(payload dto.CreateInvoiceDto) (*models.Invoice, error) {
	customerService := NewCustomerService(s.database)
	if _, err := customerService.FindCustomerById(payload.CustomerID); err != nil {
		return nil, err
	}

	existingInvoice, _ := s.FindInvoiceByTitle(payload.Title)
	if existingInvoice != nil {
		return nil, ErrInvoiceTitleExists
	}

	var status models.InvoiceStatus
	if payload.IsDraft {
		status = models.Draft
	} else {
		status = models.Pending
	}

	invoice := &models.Invoice{
		Currency:     payload.Currency,
		CustomerID:   uuid.MustParse(payload.CustomerID),
		DateDue:      payload.DateDue,
		Discount:     payload.Discount,
		DiscountType: models.DiscountType(payload.DiscountType),
		Note:         payload.Note,
		Tax:          payload.Tax,
		TaxType:      models.DiscountType(payload.TaxType),
		Title:        payload.Title,
		Status:       status,
	}

	invoice.Items = make([]models.InvoiceItem, 0, len(payload.Items))
	for _, item := range payload.Items {
		invoice.Items = append(invoice.Items, models.InvoiceItem{
			Description: item.Description,
			Quantity:    item.Quantity,
			Price:       item.Price,
		})
	}

	s.calculateInvoiceTotals(invoice)

	if err := s.database.Create(invoice).Error; err != nil {
		return nil, err
	}

	if err := s.database.Preload("Customer").Preload("Items").First(invoice, invoice.ID).Error; err != nil {
		return nil, err
	}

	return invoice, nil
}

func (s *InvoiceService) UpdateInvoice(id string, payload dto.UpdateInvoiceDto) (*models.Invoice, error) {
	invoice, err := s.FindInvoiceById(id)
	if err != nil {
		return nil, err
	}

	if payload.Title != nil && !strings.EqualFold(*payload.Title, invoice.Title) {
		existingInvoice, _ := s.FindInvoiceByTitle(*payload.Title)
		if existingInvoice != nil && existingInvoice.ID != invoice.ID {
			return nil, ErrInvoiceTitleExists
		}
	}

	if payload.Status != nil {
		invoice.Status = models.InvoiceStatus(*payload.Status)
	}
	if payload.Currency != nil {
		invoice.Currency = *payload.Currency
	}
	if payload.DateDue != nil {
		invoice.DateDue = *payload.DateDue
	}
	if payload.Discount != nil {
		invoice.Discount = *payload.Discount
	}
	if payload.DiscountType != nil {
		invoice.DiscountType = models.DiscountType(*payload.DiscountType)
	}
	if payload.Note != nil {
		invoice.Note = *payload.Note
	}
	if payload.Tax != nil {
		invoice.Tax = *payload.Tax
	}
	if payload.TaxType != nil {
		invoice.TaxType = models.DiscountType(*payload.TaxType)
	}
	if payload.Title != nil {
		invoice.Title = *payload.Title
	}

	err = s.database.Transaction(func(tx *gorm.DB) error {
		if payload.Items != nil {
			if err = tx.Where("invoice_id = ?", invoice.ID).Delete(&models.InvoiceItem{}).Error; err != nil {
				return err
			}

			invoice.Items = make([]models.InvoiceItem, len(payload.Items))
			for i, item := range payload.Items {
				invoice.Items[i] = models.InvoiceItem{
					InvoiceID:   invoice.ID,
					Description: item.Description,
					Quantity:    item.Quantity,
					Price:       item.Price,
				}
			}
		}

		s.calculateInvoiceTotals(invoice)

		if err = tx.Save(invoice).Error; err != nil {
			return err
		}

		if payload.Items != nil {
			for i := range invoice.Items {
				if err = tx.Create(&invoice.Items[i]).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if err := s.database.Preload("Customer").Preload("Items").First(invoice, invoice.ID).Error; err != nil {
		return nil, err
	}

	return invoice, nil
}

func (s *InvoiceService) DeleteInvoice(id string) error {
	invoice, err := s.FindInvoiceById(id)
	if err != nil {
		return err
	}

	return s.database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("invoice_id = ?", invoice.ID).Delete(&models.InvoiceItem{}).Error; err != nil {
			return err
		}
		return tx.Delete(invoice).Error
	})
}

func (s *InvoiceService) GetInvoices(params dto.InvoicePagination) (*dto.PaginatedResponse[models.Invoice], error) {
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	var invoices []models.Invoice
	var totalItems int64

	query := s.database.Model(&models.Invoice{})

	if params.Query != nil && strings.TrimSpace(*params.Query) != "" {
		search := "%" + strings.ToLower(strings.TrimSpace(*params.Query)) + "%"
		query = query.Joins("JOIN customers ON customers.id = invoices.customer_id").
			Where("LOWER(invoices.reference_no) LIKE ? OR LOWER(invoices.title) LIKE ? OR LOWER(invoices.status) LIKE ? OR LOWER(customers.name) LIKE ?",
				search, search, search, search)
	}

	if err := query.Count(&totalItems).Error; err != nil {
		return &dto.PaginatedResponse[models.Invoice]{
			Data:       []models.Invoice{},
			Limit:      params.Limit,
			Page:       params.Page,
			TotalItems: 0,
			TotalPages: 0,
		}, err
	}

	offset := (params.Page - 1) * params.Limit

	if err := query.Offset(offset).
		Preload("Customer").
		Preload("Items").
		Limit(params.Limit).
		Order("created_at DESC").
		Find(&invoices).Error; err != nil {
		return nil, err
	}

	totalPages := 0
	if totalItems > 0 {
		totalPages = int((totalItems + int64(params.Limit) - 1) / int64(params.Limit))
	}

	return &dto.PaginatedResponse[models.Invoice]{
		Data:       invoices,
		TotalItems: int(totalItems),
		TotalPages: totalPages,
		Page:       params.Page,
		Limit:      params.Limit,
	}, nil
}

func (s *InvoiceService) GetInvoice(id string) (*models.Invoice, error) {
	invoice := &models.Invoice{}
	if err := s.database.Preload("Customer").Preload("Items").Where("id = ?", id).First(invoice).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvoiceNotFound
		}
		return nil, err
	}
	return invoice, nil
}

func (s *InvoiceService) FindInvoiceById(id string) (*models.Invoice, error) {
	invoice := &models.Invoice{}
	if err := s.database.Preload("Items").Where("id = ?", id).First(invoice).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvoiceNotFound
		}
		return nil, err
	}
	return invoice, nil
}

func (s *InvoiceService) FindInvoiceByTitle(title string) (*models.Invoice, error) {
	invoice := &models.Invoice{}
	if err := s.database.Preload("Items").Where("title = ?", title).First(invoice).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvoiceNotFound
		}
		return nil, err
	}
	return invoice, nil
}
