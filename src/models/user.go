package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type User struct {
	BaseModel
	BankInformation *BankInformation `json:"bankInformation" gorm:"embedded"`
	CompanyLogo     string           `json:"companyLogo" gorm:"type:varchar(255);not null"`
	CompanyName     string           `json:"companyName" gorm:"type:varchar(255);not null"`
	Email           string           `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Name            string           `json:"name" gorm:"type:varchar(255);not null"`
	Phone           string           `json:"phone" gorm:"type:varchar(255);uniqueIndex;not null"`
	Provider        string           `json:"provider" gorm:"type:varchar(255);not null"`
	RcNumber        string           `json:"rcNumber" gorm:"type:varchar(255);uniqueIndex;not null"`
	TaxId           string           `json:"taxId" gorm:"type:varchar(255);uniqueIndex;not null"`
	Website         string           `json:"website" gorm:"type:varchar(255);uniqueIndex;not null"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	u.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	return nil
}

type BankInformation struct {
	AccountName   string `json:"accountName" gorm:"type:varchar(255);not null"`
	AccountNumber string `json:"accountNumber" gorm:"type:varchar(255);not null"`
	BankName      string `json:"bankName" gorm:"type:varchar(255);not null"`
	BankSwiftCode string `json:"bankSwiftCode" gorm:"type:varchar(255);not null"`
	Iban          string `json:"iban" gorm:"type:varchar(255);not null"`
	RoutingNumber string `json:"routingNumber" gorm:"type:varchar(255);not null"`
}
