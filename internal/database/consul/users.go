package consul

import (
	"context"
	"fmt"
	"path"

	"tracker/internal/types/user"
)

type UsersDatabase struct {
	db     *Database
	prefix string
}

func (db *Database) Users() *UsersDatabase {
	return &UsersDatabase{
		db:     db,
		prefix: "users",
	}
}

// Create the data in Consul, where the key is the user email
func (db *UsersDatabase) Create(ctx context.Context, u *user.User) error {
	if err := db.put(ctx, u.Email, u); err != nil {
		return fmt.Errorf("unable to insert user: %w", err)
	}

	return nil
}

// Details of the user
func (db *UsersDatabase) Details(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	if err := db.get(ctx, email, &u); err != nil {
		return nil, err
	}

	return &u, nil
}

func (db *UsersDatabase) get(ctx context.Context, key string, value any) error {
	return db.db.get(ctx, path.Join(db.prefix, key), value)
}

func (db *UsersDatabase) put(ctx context.Context, key string, value any) error {
	return db.db.put(ctx, path.Join(db.prefix, key), value)
}
