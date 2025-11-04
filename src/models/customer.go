package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	BaseModel
	Email string `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Name  string `json:"name" gorm:"type:varchar(255);not null"`
	Phone string `json:"phone" gorm:"type:varchar(255);uniqueIndex;not null"`
}

func (u *Customer) BeforeCreate(tx *gorm.DB) error {
	u.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	u.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	return nil
}

func (u *Customer) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	return nil
}
