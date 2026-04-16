package store

import (
	"strings"
	"testing"
	"time"

	"github.com/sofuejin0121/toy_app_go/internal/model"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	return NewTestStore(t)
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

func createTestUser(t *testing.T, s *Store, name, email string, admin bool) *model.User {
	t.Helper()
	password := "password"
	digest, _ := model.Digest(password)
	activationDigest, _ := model.Digest("dummy-token")
	activatedAt := time.Now()
	user := &model.User{
		Name:             name,
		Email:            email,
		PasswordDigest:   digest,
		Admin:            admin,
		Activated:        true,
		ActivatedAt:      &activatedAt,
		ActivationDigest: activationDigest,
	}
	if err := s.CreateUser(user); err != nil {
		t.Fatalf("create test user: %v", err)
	}
	return user
}
