package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"tracker/database"
)

type User struct {
	Username string
	Email    string
}

func (u *User) Scan(rows *sql.Rows) error {
	return rows.Scan(&u.Username, &u.Email)
}

func CurrentUser(r *http.Request) (*User, error) {
	fmt.Println("[Auth] Getting current user")
	session, err := GetSession(r, "tracker")
	if err != nil {
		fmt.Printf("[Auth] CurrentUser: Error: %v\n", err)
	}
	id, ok := session.Values["user-id"]
	fmt.Printf("[Auth] CurrentUser: ID (%v), Ok (%v)\n", id, ok)
	if !ok {
		return nil, nil
	}
	user, err := LoadUser(id.(string))
	return user, err
}

func LoadUser(email string) (*User, error) {
	db, err := database.Open("accounts")
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT username,email FROM users WHERE email=?", email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := &User{}
		err = user.Scan(rows)
		if err != nil {
			return nil, err
		}
		return user, err
	}

	return nil, nil
}

func CreateUser(userinfo *GoogleUserInfo) (*User, error) {
	db, err := database.Open("accounts")
	if err != nil {
		return nil, err
	}

	username := strings.Split(userinfo.Email, "@")[0]
	_, err = db.Exec(`INSERT INTO users(username, email) VALUES(?, ?)`, username, userinfo.Email)
	if err != nil {
		return nil, err
	}
	return LoadUser(userinfo.Email)
}
