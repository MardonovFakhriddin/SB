package models

import "time"

type Credit struct {
	ID             int        `db:"id"`
	UserID         int        `db:"user_id"`
	Amount         int64      `db:"amount"`
	Currency       string     `db:"currency"`
	DurationMonths int        `db:"duration_months"`
	InterestRate   float64    `db:"interest_rate"`
	CreatedAt      time.Time  `db:"created_at"`
	ApprovedAt     *time.Time `db:"approved_at"`
	Active         bool       `db:"active"`
}
