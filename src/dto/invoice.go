package dto

import "time"

type CreateInvoiceDto struct {
	Currency     string                 `json:"currency"`
	CustomerID   string                 `json:"customerId"`
	DateDue      time.Time              `json:"dateDue"`
	DateIssued   time.Time              `json:"dateDssued"`
	Discount     float64                `json:"discount"`
	DiscountType string                 `json:"discountType"`
	Items        []CreateInvoiceItemDto `json:"items,omitempty"`
	Note         string                 `json:"note"`
	Tax          float64                `json:"tax"`
	TaxType      string                 `json:"taxType"`
	Title        string                 `json:"title"`
}

type CreateInvoiceItemDto struct {
	Description string  `json:"description"`
	LineTotal   float64 `json:"lineTotal"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

type UpdateInvoiceDto struct {
	Currency     string                 `json:"currency"`
	CustomerID   string                 `json:"customerId"`
	DateDue      time.Time              `json:"dateDue"`
	DateIssued   time.Time              `json:"dateDssued"`
	Discount     float64                `json:"discount"`
	DiscountType string                 `json:"discountType"`
	Items        []CreateInvoiceItemDto `json:"items"`
	Note         string                 `json:"note"`
	Tax          float64                `json:"tax"`
	TaxType      string                 `json:"taxType"`
	Title        string                 `json:"title"`
	Status       string                 `json:"status"`
}
