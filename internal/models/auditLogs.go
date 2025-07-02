package models

import "time"

type AuditLog struct {
	ID        int       `db:"id"`
	Action    string    `db:"action"`
	Entity    string    `db:"entity"`
	EntityID  int       `db:"entity_id"`
	UserID    int       `db:"user_id"`
	Timestamp time.Time `db:"timestamp"`
}
