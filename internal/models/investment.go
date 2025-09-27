package models

import (
	"time"

	"github.com/google/uuid"
)

type Investment struct {
	InvestmentID   uuid.UUID `db:"investment_id" json:"investment_id"`
	UserID         uuid.UUID `db:"user_id" json:"user_id"`
	Type           string    `db:"type" json:"type"`
	AmountInvested float64   `db:"amount_invested" json:"amount_invested"`
	CurrentValue   float64   `db:"current_value" json:"current_value"`
	DateInvested   time.Time `db:"date_invested" json:"date_invested"`
	ReminderDate   *time.Time `db:"reminder_date" json:"reminder_date,omitempty"`
	Notes          *string    `db:"notes" json:"notes,omitempty"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}
