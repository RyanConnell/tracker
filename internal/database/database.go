package database

import (
	"context"

	"tracker/internal/types/show"
	"tracker/internal/types/user"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrNotFound = Error("not found")
)

// Database is the shared abstraction over all databases
type Database interface {
	Users() UsersDatabase
	Shows() ShowsDatabase
}

// UserDatabase abstracts the user interaction with the database.
type UsersDatabase interface {
	// Insert the user into the database
	Create(context.Context, *user.User) error
	// Details the user based on the email address.
	Details(ctx context.Context, email string) (*user.User, error)
}

type ShowsDatabase interface {
	// List returns the list of shows. The episode list is empty in each show.
	List(context.Context) ([]*show.Show, error)
	// Details gives details about a show by a show ID, including the episodes.
	Details(context.Context, int) (*show.Show, error)
}
