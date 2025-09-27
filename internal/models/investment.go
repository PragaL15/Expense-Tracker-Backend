package models

import (
	"time"

	"github.com/google/uuid"
)

type Investment struct {
	InvestmentID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"investment_id"`
	UserID         uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Type           string    `json:"type" validate:"required"`
	AmountInvested float64   `json:"amount_invested" validate:"required,gte=0"`
	CurrentValue   float64   `json:"current_value"`
	DateInvested   time.Time `json:"date_invested" validate:"required"`
	ReminderDate   *time.Time `json:"reminder_date,omitempty"`
	Notes          string    `json:"notes,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

