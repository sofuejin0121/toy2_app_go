package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/sofuejin0121/toy_app_go/internal/middleware"
	"github.com/sofuejin0121/toy_app_go/internal/store"
	"github.com/sofuejin0121/toy_app_go/web/components"
)

// SessionHandler はセッションリソースのHTTPハンドラーです。
type SessionHandler struct {
	store *store.Store
}

// NewSessionHandler は新しいSessionHandlerを返します。
func NewSessionHandler(store *store.Store) *SessionHandler {
	return &SessionHandler{store: store}
}

// ログインフォームの表示
func (h *SessionHandler) New(w http.ResponseWriter, r *http.Request) {
	data := components.SessionPageData{
		Title:     "Log in",
		CSRFToken: middleware.CSRFTokenFromContext(r),
	}
	_ = components.SessionNew(data).Render(r.Context(), w)
}

// Create はユーザーの認証を行い、ログインを処理します。
// フレンドリーフォワーディング: ログイン前にアクセスしようとしたURLがあれば
// そこにリダイレクトします。
func (h *SessionHandler) Create(w http.ResponseWriter, r *http.Request) {
	email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
	password := r.FormValue("password")
	remember := r.FormValue("remember_me") == "1"

	user, err := h.store.FindUserByEmail(email)
	if err != nil || !user.Authenticate(password) {
		data := components.SessionPageData{
			Title:       "Log in",
			Flash:       map[string]string{"danger": "Invalid email/password combination"},
			LoggedIn:    false,
			CurrentUser: nil,
			CSRFToken:   middleware.CSRFTokenFromContext(r),
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		components.SessionNew(data).Render(r.Context(), w)
		return
	}

	// フレンドリーフォワーディング: ログイン前の転送先URLを取得
	forwardingURL := getForwardingURL(w, r)

	// ［Remember me］チェックボックスの値で分岐
	if remember {
		rememberUser(w, user, h.store)
	} else {
		forgetUser(w, user, h.store)
	}
	logIn(w, r, user.ID, h.store)

	// 転送先URLがあればそこにリダイレクト、なければプロフィールページへ
	if forwardingURL != "" {
		http.Redirect(w, r, forwardingURL, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/users/%d", user.ID),
			http.StatusSeeOther)
	}
}

func (h *SessionHandler) Destroy(w http.ResponseWriter, r *http.Request) {
	if loggedIn(r) {
		logOut(w, currentUser(r), h.store)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *SessionHandler) renderNewWithFlash(w http.ResponseWriter, r *http.Request, level string, message string) {
	setFlash(w, level, message)
	h.New(w, r)
}
