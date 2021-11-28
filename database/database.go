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
	db, err := OpenDriver(MySQL, fmt.Sprintf("tracker:@tcp(mysql:3306)/%s", name))
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(fmt.Sprintf("USE %s", name))
	if err != nil {
		return nil, err
	}
	return db, err
}

// OpenDriver opens the database for a given database type.
func OpenDriver(driver Driver, dsn string) (*sql.DB, error) {
	db, err := sql.Open(string(driver), dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
