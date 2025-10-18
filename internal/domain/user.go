package domain

import "time"

type User struct {
	ID        int
	FullName  string
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
