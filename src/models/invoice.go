package models

import (
	"database/sql"
	"time"

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
	CustomerID   uint          `json:"customer_id" gorm:"index"`
	Customer     Customer      `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	DateDue      time.Time     `json:"date_due" gorm:"index"`
	DateIssued   time.Time     `json:"date_issued"`
	Discount     float64       `json:"discount"`
	DiscountType DiscountType  `json:"discount_type" gorm:"type:varchar(10)"`
	Items        []InvoiceItem `json:"items,omitempty" gorm:"foreignKey:InvoiceID"`
	Note         string        `json:"note" gorm:"type:text"`
	ReferenceNo  string        `json:"reference_no" gorm:"type:varchar(100);uniqueIndex"`
	Status       InvoiceStatus `json:"status" gorm:"type:varchar(10);index"`
	SubTotal     float64       `json:"sub_total"`
	Tax          float64       `json:"tax"`
	TaxType      DiscountType  `json:"tax_type" gorm:"type:varchar(10)"`
	Title        string        `json:"title" gorm:"type:varchar(255)"`
	Total        float64       `json:"total"`
}

type InvoiceItem struct {
	BaseModel
	InvoiceID   uint    `json:"invoice_id" gorm:"index"`
	Description string  `json:"description" gorm:"type:text"`
	LineTotal   float64 `json:"line_total"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

func (u *Invoice) BeforeCreate(tx *gorm.DB) error {
	u.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	u.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	return nil
}

func (u *Invoice) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	return nil
}
