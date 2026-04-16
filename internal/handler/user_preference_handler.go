package handler

import (
	"net/http"
	"os"

	"github.com/sofuejin0121/toy_app_go/internal/mailer"
	"github.com/sofuejin0121/toy_app_go/internal/middleware"
	"github.com/sofuejin0121/toy_app_go/internal/model"
	"github.com/sofuejin0121/toy_app_go/internal/store"
	"github.com/sofuejin0121/toy_app_go/web/components"
)

type UserPreferenceHandler struct {
	store  userPreferenceStore
	mailer mailer.Mailer
}

type userPreferenceStore interface {
	GetOrCreateUserPreference(userID int64) (*model.UserPreference, error)
	UpdateUserPreference(userID int64, emailOnFollow, emailOnLike bool) error
}

func NewUserPreferenceHandler(s userPreferenceStore, m mailer.Mailer) *UserPreferenceHandler {
	return &UserPreferenceHandler{store: s, mailer: m}
}

// Edit は通知設定ページを表示する
func (h *UserPreferenceHandler) Edit(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	pref, err := h.store.GetOrCreateUserPreference(user.ID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := components.UserPreferencePageData{
		Title:       "Notification Settings",
		Flash:       getFlash(r),
		LoggedIn:    isLoggedIn(r),
		CurrentUser: currentUser(r),
		CSRFToken:   middleware.CSRFTokenFromContext(r),
		Debug:       os.Getenv("APP_ENV") != "production",
		Pref:        pref,
	}
	_ = components.UserPreferenceEdit(data).Render(r.Context(), w)
}

// Update は通知設定を保存する
func (h *UserPreferenceHandler) Update(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// チェックボックスはチェックされていると値が送られ、未チェックは送られない
	emailOnFollow := r.FormValue("email_on_follow") == "1"
	emailOnLike := r.FormValue("email_on_like") == "1"

	if err := h.store.UpdateUserPreference(user.ID, emailOnFollow, emailOnLike); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	setFlash(w, "success", "Settings saved!")
	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

// GetUserPreferenceForStore は外部からストアを取得するためのヘルパー
// relationship_handler / like_handler からメール送信判定に使う
func GetUserPreferenceForStore(s *store.Store, userID int64) (*model.UserPreference, error) {
	return s.GetOrCreateUserPreference(userID)
}
