package models

import "time"

type Entry struct {
	ID        int       `json:"id"`
	AccountID int       `json:"account_id"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
