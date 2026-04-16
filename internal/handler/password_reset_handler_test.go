package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/sofuejin0121/toy_app_go/internal/model"
)

// createActivatedTestUser は有効化済みのテストユーザーを作成する。
func createActivatedTestUser(t *testing.T, env *testEnv, name, email, password string) *model.User {
	t.Helper()

	now := time.Now()
	user := &model.User{
		Name:        name,
		Email:       email,
		Activated:   true,
		ActivatedAt: &now,
	}
	if err := user.SetPassword(password); err != nil {
		t.Fatalf("SetPassword: %v", err)
	}
	if err := env.store.CreateUser(user); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	return user
}

// setResetTokenForUser は指定した送信時刻でリセットトークンを付与する。
func setResetTokenForUser(t *testing.T, env *testEnv, user *model.User, sentAt time.Time) string {
	t.Helper()

	if err := user.CreateResetDigest(); err != nil {
		t.Fatalf("CreateResetDigest: %v", err)
	}
	if err := env.store.UpdateResetDigest(user.ID, user.ResetDigest, sentAt); err != nil {
		t.Fatalf("UpdateResetDigest: %v", err)
	}
	return user.ResetToken
}

func TestPasswordResetExpired(t *testing.T) {
	env, ts, client := setupTestServer(t)
	user := createActivatedTestUser(t, env, "Michael Example",
		"michael@example.com", "password")
	token := setResetTokenForUser(t, env, user, time.Now().Add(-3*time.Hour))

	t.Run("should redirect to the password-reset page", func(t *testing.T) {
		csrf := getCSRFToken(t, ts, client)
		form := url.Values{
			"_method":               {"PATCH"},
			"email":                 {user.Email},
			"password":              {"newpassword"},
			"password_confirmation": {"newpassword"},
			"csrf_token":            {csrf},
		}
		resp, err := client.PostForm(
			ts.URL+fmt.Sprintf("/password_resets/%s", token), form)
		if err != nil {
			t.Fatalf("PATCH /password_resets: %v", err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusSeeOther {
			t.Errorf("status: got %d, want %d",
				resp.StatusCode, http.StatusSeeOther)
		}
		loc := resp.Header.Get("Location")
		if loc != "/password_resets/new" {
			t.Errorf("redirect location: got %q, want %q",
				loc, "/password_resets/new")
		}
	})
}
