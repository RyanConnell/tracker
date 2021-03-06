// Package sql provides the database access for user data.
package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"tracker/internal/database"
	"tracker/internal/types/user"
)

type UsersDatabase struct {
	db *Database

	getUserStmt    *sql.Stmt
	insertUserStmt *sql.Stmt
}

func (db *Database) Users() *UsersDatabase {
	getUserStmt, err := db.db.Prepare(getUserQuery)
	if err != nil {
		panic(fmt.Sprintf("unable to prepare query to get user: %v", err))
	}
	insertUserStmt, err := db.db.Prepare(insertUserQuery)
	if err != nil {
		panic(fmt.Sprintf("unable to prepare query to insert user: %v", err))
	}

	return &UsersDatabase{
		db: db,

		getUserStmt:    getUserStmt,
		insertUserStmt: insertUserStmt,
	}
}

func (db *UsersDatabase) Create(ctx context.Context, u *user.User) error {
	if _, err := db.insertUserStmt.ExecContext(ctx, u.Name, u.Email); err != nil {
		return fmt.Errorf("unable to insert user: %w", err)
	}

	return nil
}

func (db *UsersDatabase) Details(ctx context.Context, email string) (*user.User, error) {
	u := &user.User{}
	if err := db.getUserStmt.QueryRowContext(ctx, email).Scan(
		&u.Name,
		&u.Email,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, database.ErrNotFound
		}
		return nil, fmt.Errorf("unable to get user: %w", err)
	}

	return u, nil
}

// Username is actually user's name.
// TODO: Update the schema to reflect that.
const getUserQuery = `
SELECT
	username,
	email,
FROM users
WHERE
	email=?
LIMIT 1;
`

const insertUserQuery = `
INSERT INTO users (
	username,
	email,
) VALUES (
	?,
	?
);
`
