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

func (s *InvoiceService) CreateInvoice(payload dto.CreateInvoiceDto) (*models.Invoice, error) {
	customerService := NewCustomerService(s.database)
	if _, err := customerService.FindCustomerById(payload.CustomerID); err != nil {
		return nil, err
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
		Status:       models.Draft,
	}

	if !payload.IsDraft {
		invoice.Status = models.Pending
	}

	invoice.Items = make([]models.InvoiceItem, 0, len(payload.Items))
	for _, item := range payload.Items {
		invoice.Items = append(invoice.Items, models.InvoiceItem{
			Description: item.Description,
			Quantity:    item.Quantity,
			Price:       item.Price,
			LineTotal:   item.LineTotal,
		})
	}

	if err := s.database.Create(invoice).Error; err != nil {
		return nil, err
	}
	return invoice, nil
}

func (s *InvoiceService) UpdateInvoice(id string, payload dto.UpdateInvoiceDto) (*models.Invoice, error) {
	invoice, err := s.FindInvoiceById(id)
	if err != nil {
		return nil, err
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

	if payload.Items != nil {
		if err = s.database.Where("invoice_id = ?", invoice.ID).Delete(&models.InvoiceItem{}).Error; err != nil {
			return nil, err
		}

		invoice.Items = make([]models.InvoiceItem, len(payload.Items))
		for i, item := range payload.Items {
			invoice.Items[i] = models.InvoiceItem{
				Description: item.Description,
				Quantity:    item.Quantity,
				Price:       item.Price,
				LineTotal:   item.LineTotal,
			}
		}
	}

	err = s.database.Transaction(func(tx *gorm.DB) error {
		if err = tx.Save(invoice).Error; err != nil {
			return err
		}

		if payload.Items != nil {
			for i := range invoice.Items {
				invoice.Items[i].InvoiceID = invoice.ID
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

	query := s.database.Model(&models.Invoice{}).Preload("Items")

	if params.ReferenceNo != nil && strings.TrimSpace(*params.ReferenceNo) != "" {
		referenceNoPattern := "%" + strings.ToLower(strings.TrimSpace(*params.ReferenceNo)) + "%"
		query = query.Where("LOWER(reference_no) LIKE ?", referenceNoPattern)
	}

	if params.CustomerId != nil && strings.TrimSpace(*params.CustomerId) != "" {
		query = query.Where("customer_id = ?", strings.TrimSpace(*params.CustomerId))
	}

	if params.Status != nil && strings.TrimSpace(*params.Status) != "" {
		query = query.Where("status = ?", strings.TrimSpace(*params.Status))
	}

	if params.Title != nil && strings.TrimSpace(*params.Title) != "" {
		titlePattern := "%" + strings.ToLower(strings.TrimSpace(*params.Title)) + "%"
		query = query.Where("LOWER(title) LIKE ?", titlePattern)
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
	return s.FindInvoiceById(id)
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
