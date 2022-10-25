package consul

import (
	"context"
	"errors"
	"testing"

	"tracker/internal/database"
	"tracker/internal/types/user"

	"github.com/go-test/deep"
)

func TestImplements(t *testing.T) {
	var i any = &UsersDatabase{}

	if _, ok := i.(database.UsersDatabase); !ok {
		t.Errorf("UserDatabase does not implement database.UserDatabase")
	}
}

func TestE2E(t *testing.T) {
	kv := &testKV{
		m: make(map[string][]byte),
	}
	db, err := NewDatabase("tracker", KVClient(kv))
	if err != nil {
		t.Fatalf("unable to create user database client: %v", err)
	}

	usersDB := db.Users()

	_, err = usersDB.Details(context.Background(), "user@example.com")
	if !errors.Is(err, errNotFound) {
		t.Fatalf("(Empty) Get() err = %v, want %v", err, errNotFound)
	}

	want := &user.User{
		Email: "user@example.com",
		Name:  "Example User",
	}

	if err := usersDB.Create(context.Background(), want); err != nil {
		t.Fatalf("Insert() err = %v, want %v", err, nil)
	}

	got, err := usersDB.Details(context.Background(), "user@example.com")
	if err != nil {
		t.Fatalf("Get() err = %v, want %v", got, nil)
	}

	if diff := deep.Equal(got, want); diff != nil {
		t.Fatalf("Get() = %v, got %v, diff = %v", got, want, diff)
	}
}
