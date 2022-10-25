package sql

import (
	"testing"
	"tracker/internal/database"
)

func TestImplements(t *testing.T) {
	var i any = &UsersDatabase{}

	if _, ok := i.(database.UsersDatabase); !ok {
		t.Errorf("UserDatabase doesn't implement database.UserDatabase")
	}
}
