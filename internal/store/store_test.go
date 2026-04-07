package store

import (
    "strings"
    "testing"

    "github.com/sofuejin0121/toy_app_go/internal/model"
)

func newTestStore(t *testing.T) *Store {
    t.Helper()

    s, err := New(":memory:")
    if err != nil {
        t.Fatalf("new store: %v", err)
    }
    t.Cleanup(func() {
        s.Close()
    })

    _, err = s.db.Exec(`
        CREATE TABLE users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            email TEXT NOT NULL,
            password_digest TEXT NOT NULL DEFAULT '',
            created_at TEXT NOT NULL,
            updated_at TEXT NOT NULL
        );
    `)
    if err != nil {
        t.Fatalf("create users table: %v", err)
    }

    _, err = s.db.Exec(`CREATE UNIQUE INDEX index_users_on_email ON users (email);`)
    if err != nil {
        t.Fatalf("create email index: %v", err)
    }

    return s
}

func TestEmailSavedAsLowercase(t *testing.T) {
	s := newTestStore(t)

	mixedCaseEmail := "Foo@ExAMPLe.CoM"
	user := &model.User{Name: "Example User", Email: mixedCaseEmail}
	if err := s.CreateUser(user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	reloaded, err := s.GetUser(user.ID)
	if err != nil {
		t.Fatalf("get user: %v", err)
	}
	if reloaded.Email != strings.ToLower(mixedCaseEmail) {
        t.Errorf("email not lowercased: got %q, want %q", reloaded.Email, strings.ToLower(mixedCaseEmail))
    }
}