package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"tracker/database"
)

var (
	ErrUserNotFound = errors.New("auth: user not found")
)

func init() {
	db, err := database.Open("accounts")
	if err != nil {
		log.Fatalf("unable to open database for authentication: %v", err)
	}

	loadUserStmt, err = db.Prepare("SELECT username,email FROM users WHERE email=?")
	if err != nil {
		log.Fatalf("unable to prepare loading users: %v", err)
	}

	createUserStmt, err = db.Prepare("INSERT INTO users(username, email) VALUES(?, ?)")
	if err != nil {
		log.Fatalf("unable to prepare creating users: %v", err)
	}
}

// TODO: Convert from functions to a struct and move this global variable
//       into the new struct
var loadUserStmt *sql.Stmt
var createUserStmt *sql.Stmt

type User struct {
	Username string
	Email    string
}

func (u *User) Scan(rows *sql.Row) error {
	return rows.Scan(&u.Username, &u.Email)
}

func CurrentUser(r *http.Request) (User, error) {
	fmt.Println("Getting current user")
	session, err := GetSession(r, "tracker")
	if err != nil {
		fmt.Printf("CurrentUser: Error: %v\n", err)
	}
	id, ok := session.Values["user-id"]
	fmt.Printf("ID: %v. Ok: %v", id, ok)
	if !ok {
		return User{}, nil
	}
	user, err := LoadUser(id.(string))
	return *user, err
}

func LoadUser(email string) (*User, error) {
	user := new(User)
	switch err := loadUserStmt.QueryRow(email).Scan(user); err {
	case sql.ErrNoRows:
		return nil, ErrUserNotFound
	case nil:
		return user, nil
	default:
		return nil, fmt.Errorf("unknown error when getting user: %w", err)
	}
}

func CreateUser(userinfo *GoogleUserInfo) (*User, error) {
	username := strings.Split(userinfo.Email, "@")[0]
	if _, err := createUserStmt.Exec(username, userinfo.Email); err != nil {
		return nil, err
	}
	return LoadUser(userinfo.Email)
}
