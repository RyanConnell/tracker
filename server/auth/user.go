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
	Admin    bool
}

func (u *User) Scan(rows *sql.Rows) error {
	return rows.Scan(&u.Username, &u.Email, &u.Admin)
}

func CurrentUser(r *http.Request) (User, error) {
	session, err := GetSession(r, "tracker")
	if err != nil {
		fmt.Printf("[Auth] Error getting session: %v\n", err)
	}
	id, ok := session.Values["user-id"]
	if !ok {
		return User{}, nil
	}
	user, err := LoadUser(id.(string))
	fmt.Printf("[Auth] Current User: %v\n", user)
	return *user, err
}

func LoadUser(email string) (*User, error) {
	db, err := database.Open("accounts")
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT username,email,admin FROM users WHERE email=?", email)
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
