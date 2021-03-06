package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

type Driver string

const (
	MySQL  Driver = "mysql"
	SQLite Driver = "sqlite3"
)

func Open(name string) (*sql.DB, error) {
	return OpenDriver(MySQL, fmt.Sprintf("rhino:@/%s?parseTime=true", name))
}

// OpenDriver opens the database for a given database type.
func OpenDriver(driver Driver, dsn string) (*sql.DB, error) {
	db, err := sql.Open(string(driver), dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
