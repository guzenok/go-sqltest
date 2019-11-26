package store

import (
	"context"
)

// Store data access layer
type Store interface {
	Close() error
	Ping() error

	UsersStore
}

// UsersStore repository of user objects
type UsersStore interface {
	CreateUser(ctx context.Context, in *User) (out *User, err error)
	DeleteUser(ctx context.Context, id int32) (err error)
	GetUserByID(ctx context.Context, id int32) (out *User, err error)
	GetUserByLogin(ctx context.Context, login string) (out *User, err error)
	SetPassword(ctx context.Context, userID int32, newPassword string) (err error)
}
