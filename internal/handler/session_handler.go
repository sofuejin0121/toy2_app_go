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

// ログイン処理
func (h *SessionHandler) Create(w http.ResponseWriter, r *http.Request) {
	email := strings.ToLower(r.FormValue("email"))
	password := r.FormValue("password")

	user, err := h.store.FindUserByEmail(email)
	if err == nil && user.Authenticate(password) {
		logIn(w, r, user.ID)
		setFlash(w, "success", "Logged in!")
		http.Redirect(w, r, fmt.Sprintf("/users/%d", user.ID), http.StatusSeeOther)
		return
	}
	data := components.SessionPageData{
		Title:     "Log in",
		Flash:     map[string]string{"danger": "Invalid email/password combination"},
		CSRFToken: middleware.CSRFTokenFromContext(r),
	}
	w.WriteHeader(http.StatusUnprocessableEntity)
	_ = components.SessionNew(data).Render(r.Context(), w)
}

func (h *SessionHandler) Destroy(w http.ResponseWriter, r *http.Request) {
	logOut(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
