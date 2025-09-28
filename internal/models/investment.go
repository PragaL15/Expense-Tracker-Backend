package models

import (
	"time"

	"github.com/google/uuid"
)

type Investment struct {
	InvestmentID   uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"investment_id"`
	UserID         uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	Type           string     `gorm:"type:text;not null" json:"type" validate:"required"`
	AmountInvested float64    `gorm:"type:numeric(12,2);not null" json:"amount_invested" validate:"required,gte=0"`
	CurrentValue   *float64   `gorm:"type:numeric(12,2)" json:"current_value,omitempty"`
	DateInvested   time.Time  `gorm:"type:date;not null" json:"date_invested" validate:"required"`
	ReminderDate   *time.Time `gorm:"type:date" json:"reminder_date,omitempty"`
	Notes          *string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt      time.Time  `gorm:"type:timestamptz;default:now()" json:"created_at" gorm:"autoCreateTime"`
}
