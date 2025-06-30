package models

import "time"

type Transaction struct {
	ID            int       `db:"id"`
	FromAccountID int       `db:"from_account_id"`
	ToAccountID   int       `db:"to_account_id"`
	Amount        float64   `db:"amount"`
	Commission    float64   `db:"commission"`
	CreatedAt     time.Time `db:"created_at"`
}
