package store

import (
	"time"
)

// User model of user
type User struct {
	ID        int32     `db:"id"`
	Login     string    `db:"login"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	Super     bool      `db:"is_super"`
}
