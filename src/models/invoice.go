package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InvoiceStatus string

const (
	Pending InvoiceStatus = "pending"
	Paid    InvoiceStatus = "paid"
	Overdue InvoiceStatus = "overdue"
	Draft   InvoiceStatus = "draft"
)

type DiscountType string

const (
	Fixed      DiscountType = "fixed"
	Percentage DiscountType = "percentage"
)

type Invoice struct {
	BaseModel
	Currency     string        `json:"currency" gorm:"type:varchar(3)"`
	CustomerID   uuid.UUID     `json:"customerId" gorm:"index"`
	Customer     Customer      `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	DateDue      time.Time     `json:"dateDue" gorm:"index"`
	DateIssued   time.Time     `json:"dateIssued"`
	Discount     float64       `json:"discount"`
	DiscountType DiscountType  `json:"discountType" gorm:"type:varchar(10)"`
	Items        []InvoiceItem `json:"items,omitempty" gorm:"foreignKey:InvoiceID"`
	Note         string        `json:"note" gorm:"type:text"`
	ReferenceNo  string        `json:"referenceNo" gorm:"type:varchar(100);uniqueIndex"`
	Status       InvoiceStatus `json:"status" gorm:"type:varchar(10);index"`
	SubTotal     float64       `json:"subTotal"`
	Tax          float64       `json:"tax"`
	TaxType      DiscountType  `json:"taxType" gorm:"type:varchar(10)"`
	Title        string        `json:"title" gorm:"type:varchar(255)"`
	Total        float64       `json:"total"`
}

type InvoiceItem struct {
	BaseModel
	InvoiceID   uuid.UUID `json:"invoiceId" gorm:"index"`
	Description string    `json:"description" gorm:"type:text"`
	LineTotal   float64   `json:"lineTotal"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
}

func (u *Invoice) BeforeCreate(tx *gorm.DB) error {
	u.DateIssued = time.Now()
	u.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	u.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	u.SubTotal = 0
	for _, item := range u.Items {
		u.SubTotal += item.LineTotal
	}
	DiscountAmount := 0
	TaxAmount := 0

	switch u.TaxType {
	case Fixed:
		TaxAmount = int(u.Tax)
	case Percentage:
		TaxAmount = int(u.Tax * u.SubTotal / 100)
	}

	switch u.DiscountType {
	case Fixed:
		DiscountAmount = int(u.Discount)
	case Percentage:
		DiscountAmount = int(u.Discount * u.SubTotal / 100)
	}

	u.Total = u.SubTotal + float64(TaxAmount) - float64(DiscountAmount)
	return nil
}

func (u *Invoice) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}

	u.SubTotal = 0
	for _, item := range u.Items {
		u.SubTotal += item.LineTotal
	}

	DiscountAmount := 0
	TaxAmount := 0

	switch u.TaxType {
	case Fixed:
		TaxAmount = int(u.Tax)
	case Percentage:
		TaxAmount = int(u.Tax * u.SubTotal / 100)
	}

	switch u.DiscountType {
	case Fixed:
		DiscountAmount = int(u.Discount)
	case Percentage:
		DiscountAmount = int(u.Discount * u.SubTotal / 100)
	}

	u.Total = u.SubTotal + float64(TaxAmount) - float64(DiscountAmount)

	return nil
}
