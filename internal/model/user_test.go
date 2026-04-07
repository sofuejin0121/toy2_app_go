package model

import "testing"

func TestUserValidate(t *testing.T) {
	u:= User {
		Name: "Example User",
		Email: "user@example.com",
	}
	if err := u.Validate(); err != nil {
		t.Errorf("Validate() returned error: %v", err)
	}
}