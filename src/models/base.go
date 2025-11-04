package models

import (
	"database/sql"

	"github.com/google/uuid"
)

type BaseModel struct {
	ID        uuid.UUID    `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CreatedAt sql.NullTime `json:"created_at" gorm:"type:TIMESTAMP with time zone;not null"`
	UpdatedAt sql.NullTime `json:"updated_at" gorm:"type:TIMESTAMP with time zone;not null"`
	DeletedAt sql.NullTime `json:"-" gorm:"type:TIMESTAMP with time zone;null"`
}
