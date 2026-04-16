package handler

import (
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/sofuejin0121/toy_app_go/internal/mailer"
	"github.com/sofuejin0121/toy_app_go/internal/middleware"
	"github.com/sofuejin0121/toy_app_go/internal/store"
)

type testEnv struct {
	store *store.Store
}

// setupTestServer はハンドラーテスト用のテストサーバーを起動する。
// 全ルートとミドルウェアチェーンを含む実際のサーバーに近い構成で起動する。
func setupTestServer(t *testing.T) (*testEnv, *httptest.Server, *http.Client) {
	t.Helper()

	s := store.NewTestStore(t)
	m := &mailer.LogMailer{From: "test@example.com", Host: "localhost"}
	env := &testEnv{store: s}

	mux := http.NewServeMux()

	userHandler := NewUserHandler(s, m)
	sessionHandler := NewSessionHandler(s)
	staticHandler := NewStaticHandler()
	accountActivationHandler := NewAccountActivationHandler(s)
	passwordResetHandler := NewPasswordResetHandler(s, m)

	mux.HandleFunc("GET /{$}", staticHandler.Home)
	mux.HandleFunc("GET /login", sessionHandler.New)
	mux.HandleFunc("POST /login", sessionHandler.Create)
	mux.HandleFunc("DELETE /logout", sessionHandler.Destroy)
	mux.HandleFunc("GET /users/{id}", userHandler.Show)
	mux.HandleFunc("GET /account_activation/{id}/edit", accountActivationHandler.Edit)
	mux.HandleFunc("GET /password_resets/new", passwordResetHandler.New)
	mux.HandleFunc("POST /password_resets", passwordResetHandler.Create)
	mux.HandleFunc("GET /password_resets/{id}/edit", passwordResetHandler.Edit)
	mux.HandleFunc("PATCH /password_resets/{id}", passwordResetHandler.Update)

	h := middleware.CSRF(middleware.Auth(s)(middleware.MethodOverride(mux)))

	ts := httptest.NewServer(h)
	t.Cleanup(ts.Close)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("cookiejar: %v", err)
	}
	client := &http.Client{
		Jar: jar,
		// リダイレクトを追わない（レスポンスのLocationヘッダーをテストで確認するため）
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return env, ts, client
}

// getCSRFToken は /login ページへのGETリクエストでCSRFトークンを取得する。
func getCSRFToken(t *testing.T, ts *httptest.Server, client *http.Client) string {
	t.Helper()

	resp, err := client.Get(ts.URL + "/login")
	if err != nil {
		t.Fatalf("GET /login: %v", err)
	}
	resp.Body.Close()

	u, _ := url.Parse(ts.URL)
	for _, c := range client.Jar.Cookies(u) {
		if c.Name == "csrf_token" {
			return c.Value
		}
	}
	t.Fatal("csrf_token cookie not found")
	return ""
}
