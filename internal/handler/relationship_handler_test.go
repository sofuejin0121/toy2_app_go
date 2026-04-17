package handler

import (
	"html"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/sofuejin0121/toy_app_go/internal/mailer"
	"github.com/sofuejin0121/toy_app_go/internal/middleware"
	"github.com/sofuejin0121/toy_app_go/internal/model"
	"github.com/sofuejin0121/toy_app_go/internal/storage"
	"github.com/sofuejin0121/toy_app_go/internal/store"
)

func buildTestMux(s *store.Store) *http.ServeMux {
	mux := http.NewServeMux()

	mockMailer := &mailer.LogMailer{}
	userHandler := NewUserHandler(s, mockMailer)
	sessionHandler := NewSessionHandler(s)
	staticHandler := NewStaticHandler(s)
	localStorage, _ := storage.NewLocalStorage(os.TempDir(), "http://localhost")
	micropostHandler := NewMicropostHandler(s, localStorage)
	relationshipHandler := NewRelationshipHandler(s, mockMailer)

	mux.HandleFunc("GET /{$}", staticHandler.Home)
	mux.HandleFunc("GET /signup", userHandler.New)
	mux.HandleFunc("POST /users", userHandler.Create)
	mux.HandleFunc("GET /users", RequireLogin(userHandler.Index))
	mux.HandleFunc("GET /users/{id}", userHandler.Show)
	mux.HandleFunc("GET /users/{id}/edit",
		RequireLogin(userHandler.RequireCorrectUser(userHandler.Edit)))
	mux.HandleFunc("PATCH /users/{id}",
		RequireLogin(userHandler.RequireCorrectUser(userHandler.Update)))
	mux.HandleFunc("DELETE /users/{id}",
		RequireLogin(userHandler.RequireAdmin(userHandler.Destroy)))
	mux.HandleFunc("GET /users/{id}/following",
		RequireLogin(userHandler.Following))
	mux.HandleFunc("GET /users/{id}/followers",
		RequireLogin(userHandler.Followers))
	mux.HandleFunc("GET /login", sessionHandler.New)
	mux.HandleFunc("POST /login", sessionHandler.Create)
	mux.HandleFunc("POST /logout", sessionHandler.Destroy)
	mux.HandleFunc("POST /microposts",
		RequireLogin(micropostHandler.Create))
	mux.HandleFunc("DELETE /microposts/{id}",
		RequireLogin(micropostHandler.RequireCorrectUser(micropostHandler.Destroy)))
	mux.HandleFunc("POST /relationships",
		RequireLogin(relationshipHandler.Create))
	mux.HandleFunc("DELETE /relationships/{id}",
		RequireLogin(relationshipHandler.Destroy))

	return mux
}

func setupRelationshipTestServer(t *testing.T) (*testEnv, *httptest.Server, *http.Client) {
	t.Helper()

	f, err := os.CreateTemp("", "handler-test-*.db")
	if err != nil {
		t.Fatalf("create temp db: %v", err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	f.Close()

	s, err := store.New(f.Name())
	if err != nil {
		t.Fatalf("new store: %v", err)
	}
	t.Cleanup(func() { s.Close() })

	_, err = s.DB().Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			bio TEXT NOT NULL DEFAULT '',
			password_digest TEXT NOT NULL DEFAULT '',
			remember_digest VARCHAR(255) NOT NULL DEFAULT '',
			admin BOOLEAN NOT NULL DEFAULT FALSE,
			activation_digest VARCHAR(255),
			activated BOOLEAN NOT NULL DEFAULT FALSE,
			activated_at TIMESTAMP,
			reset_digest VARCHAR(255),
			reset_sent_at TIMESTAMP,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE TABLE microposts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT NOT NULL,
			user_id INTEGER NOT NULL,
			image_path TEXT DEFAULT '',
			in_reply_to_id INTEGER DEFAULT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE TABLE relationships (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			follower_id INTEGER NOT NULL,
			followed_id INTEGER NOT NULL,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE likes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id    INTEGER NOT NULL,
			micropost_id INTEGER NOT NULL,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		CREATE UNIQUE INDEX index_likes_on_user_id_and_micropost_id ON likes (user_id, micropost_id);
	`)
	if err != nil {
		t.Fatalf("create tables: %v", err)
	}
	_, err = s.DB().Exec(`CREATE UNIQUE INDEX index_users_on_email ON users (email);`)
	if err != nil {
		t.Fatalf("create email index: %v", err)
	}

	env := &testEnv{store: s}

	mux := buildTestMux(s)
	accountActivationHandler := NewAccountActivationHandler(s)
	passwordResetHandler := NewPasswordResetHandler(s, &mailer.LogMailer{})
	mux.HandleFunc("GET /account_activations/{id}/edit", accountActivationHandler.Edit)
	mux.HandleFunc("GET /password_resets/new", passwordResetHandler.New)
	mux.HandleFunc("POST /password_resets", passwordResetHandler.Create)
	mux.HandleFunc("GET /password_resets/{id}/edit", passwordResetHandler.Edit)
	mux.HandleFunc("PATCH /password_resets/{id}", passwordResetHandler.Update)

	var h http.Handler = mux
	h = middleware.MethodOverride(h)
	h = middleware.Flash(h)
	h = middleware.CSRF(h)
	h = middleware.Auth(s)(h)

	ts := httptest.NewServer(h)
	t.Cleanup(ts.Close)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("new cookie jar: %v", err)
	}
	client := &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return env, ts, client
}

func createTestUser(t *testing.T, env *testEnv, name, email, password string) *model.User {
	t.Helper()
	u := &model.User{Name: name, Email: email, Activated: true}
	if err := u.SetPassword(password); err != nil {
		t.Fatalf("SetPassword %s: %v", name, err)
	}
	if err := env.store.CreateUser(u); err != nil {
		t.Fatalf("CreateUser %s: %v", name, err)
	}
	return u
}

func loginAs(t *testing.T, ts *httptest.Server, client *http.Client, user *model.User, remember string) {
	t.Helper()
	csrfToken := getCSRFToken(t, ts, client)
	form := url.Values{
		"email":       {user.Email},
		"password":    {user.Password},
		"csrf_token":  {csrfToken},
		"remember_me": {remember},
	}
	resp, err := client.PostForm(ts.URL+"/login", form)
	if err != nil {
		t.Fatalf("POST /login: %v", err)
	}
	resp.Body.Close()
}

func readBody(t *testing.T, resp *http.Response) string {
	t.Helper()
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return string(body)
}

func TestFeedOnHomePage(t *testing.T) {
	env, ts, client := setupRelationshipTestServer(t)

	michael := createTestUser(t, env, "Michael", "michael@example.com", "password")
	lana := createTestUser(t, env, "Lana", "lana@example.com", "password")

	if err := env.store.CreateMicropost(&model.Micropost{Content: "<b>self post</b>", UserID: michael.ID}); err != nil {
		t.Fatalf("CreateMicropost self: %v", err)
	}
	if err := env.store.CreateMicropost(&model.Micropost{Content: "<i>followed post</i>", UserID: lana.ID}); err != nil {
		t.Fatalf("CreateMicropost followed: %v", err)
	}
	if err := env.store.Follow(michael.ID, lana.ID); err != nil {
		t.Fatalf("Follow: %v", err)
	}

	loginAs(t, ts, client, michael, "0")

	resp, err := client.Get(ts.URL + "/")
	if err != nil {
		t.Fatalf("GET /: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d, want %d", resp.StatusCode, http.StatusOK)
	}

	body := readBody(t, resp)
	feed, err := env.store.Feed(michael.ID, 1, 30)
	if err != nil {
		t.Fatalf("Feed: %v", err)
	}
	for _, item := range feed {
		if !strings.Contains(body, html.EscapeString(item.Micropost.Content)) {
			t.Errorf("expected feed to contain micropost content: %s", item.Micropost.Content)
		}
	}
}
