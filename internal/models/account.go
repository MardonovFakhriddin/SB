package models

import "time"

type Account struct {
	ID        int        `db:"id"`
	UserID    int        `db:"user_id"`
	Balance   float64    `db:"balance"`
	Currency  string     `db:"currency"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
	Active    bool       `db:"active"`
}
