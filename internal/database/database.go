package database

import (
	"context"

	"tracker/internal/types/user"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrNotFound = Error("not found")
)

// Database
type Database interface {
	Users() UsersDatabase
}

// UserDatabase abstracts the user interaction with the database.
type UsersDatabase interface {
	// Insert the user into the database
	Create(context.Context, *user.User) error
	// Details the user based on the email address.
	Details(ctx context.Context, email string) (*user.User, error)
}
