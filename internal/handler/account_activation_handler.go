package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sofuejin0121/toy_app_go/internal/store"
)

// AccountActivationHandler はアカウント有効化を処理するハンドラー
type AccountActivationHandler struct {
	store *store.Store
}

func NewAccountActivationHandler(s *store.Store) *AccountActivationHandler {
	return &AccountActivationHandler{store: s}
}

func (h *AccountActivationHandler) Edit(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("id")
	email := r.URL.Query().Get("email")

	user, err := h.store.GetUserByEmail(email)
	if err != nil || user == nil {
		setFlash(w, "danger", "Invalid activation link")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if !user.Activated && user.Authenticated("activation", token) {
		// リファクタリング後: Userモデルのメソッドを使用
		if err := user.Activate(h.store); err != nil {
			log.Printf("AccountActivation: activate user %d: %v", user.ID, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		logIn(w, user.ID, false, h.store)
		setFlash(w, "success", "Account activated!")
		http.Redirect(w, r, fmt.Sprintf("/users/%d", user.ID), http.StatusSeeOther)
	} else {
		setFlash(w, "danger", "Invalid activation link")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}	
}
