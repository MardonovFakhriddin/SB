package models

import "time"

type Deposit struct {
	ID             int       `db:"id"`
	UserID         int       `db:"user_id"`
	Amount         float64   `db:"amount"`
	Currency       string    `db:"currency"`
	InterestRate   float64   `db:"interest_rate"`
	DurationMonths int       `db:"duration_months"`
	CreatedAt      time.Time `db:"created_at"`
	ExpiresAt      time.Time `db:"expires_at"`
	Active         bool      `db:"active"`
}
