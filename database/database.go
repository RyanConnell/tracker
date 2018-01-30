package database

import (
	"database/sql"
	"fmt"
)

func Open(name string) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("rhino:@/%s", name))
	if err != nil {
		return nil, err
	}
	return db, nil
}
