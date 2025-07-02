package models

import "time"

type Transfer struct {
	ID            int       `db:"id"`
	FromAccountID int       `db:"from_account_id"`
	ToAccountID   int       `db:"to_account_id"`
	Amount        int       `db:"amount"`
	Currency      string    `db:"currency"`
	CreatedAt     time.Time `db:"created_at"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}
