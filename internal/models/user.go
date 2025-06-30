package models

import "time"

type User struct {
	ID        int        `db:"id"`
	FullName  string     `db:"full_name"`
	Password  string     `db:"password"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
	Active    bool       `db:"active"`
}
