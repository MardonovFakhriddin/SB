package models

import "time"

type User struct {
	ID        int        `db:"id"`
	FullName  string     `db:"full_name"`
	Password  string     `db:"password" json:"-"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
	Active    bool       `db:"active"`
}

type UpdateUser struct {
	ID       int     `json:"id"`
	FullName *string `json:"full_name"`
	Password *string `json:"password"`
}
