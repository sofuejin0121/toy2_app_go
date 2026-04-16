package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/sofuejin0121/toy_app_go/internal/model"
)

// GET /api/me
func (h *APIHandler) Me(w http.ResponseWriter, r *http.Request) {
	cu := currentUser(r)
	if cu == nil {
		writeError(w, http.StatusUnauthorized, "ログインしていません")
		return
	}
	writeJSON(w, http.StatusOK, userToJSON(*cu))
}

// POST /api/login
func (h *APIHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Remember bool   `json:"remember"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "リクエストが不正です")
		return
	}
	email := strings.ToLower(strings.TrimSpace(body.Email))
	user, err := h.store.Authenticate(email, body.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "メールアドレスまたはパスワードが正しくありません")
		return
	}
	if !user.Activated {
		writeError(w, http.StatusForbidden, "アカウントが有効化されていません。メールをご確認ください。")
		return
	}
	logIn(w, user.ID, body.Remember, h.store)
	writeJSON(w, http.StatusOK, userToJSON(*user))
}

// DELETE /api/logout
func (h *APIHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if loggedIn(r) {
		logOut(w, currentUser(r), h.store)
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "ログアウトしました"})
}

// GET /api/account_activations/{token}/edit?email=
func (h *APIHandler) ActivateAccount(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	email := r.URL.Query().Get("email")

	user, err := h.store.GetUserByEmail(email)
	if err != nil || user == nil {
		writeError(w, http.StatusBadRequest, "無効なアカウント有効化リンクです")
		return
	}
	if !user.Activated && user.Authenticated("activation", token) {
		if err := user.Activate(h.store); err != nil {
			writeError(w, http.StatusInternalServerError, "内部エラーが発生しました")
			return
		}
		logIn(w, user.ID, false, h.store)
		writeJSON(w, http.StatusOK, map[string]any{
			"message": "アカウントが有効化されました！",
			"user":    userToJSON(*user),
		})
	} else {
		writeError(w, http.StatusBadRequest, "無効なアカウント有効化リンクです")
	}
}

// POST /api/password_resets
func (h *APIHandler) CreatePasswordReset(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "リクエストが不正です")
		return
	}
	email := strings.ToLower(strings.TrimSpace(body.Email))
	user, err := h.store.GetUserByEmail(email)
	if err != nil || user == nil {
		writeError(w, http.StatusNotFound, "メールアドレスが見つかりません")
		return
	}
	if err := user.CreateResetDigest(); err != nil {
		writeError(w, http.StatusInternalServerError, "内部エラーが発生しました")
		return
	}
	if err := h.store.UpdateResetDigest(user.ID, user.ResetDigest, time.Now()); err != nil {
		writeError(w, http.StatusInternalServerError, "内部エラーが発生しました")
		return
	}
	if err := user.SendPasswordResetEmail(h.mailer); err != nil {
		log.Printf("SendPasswordResetEmail: %v", err)
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "パスワード再設定の手順をメールで送信しました"})
}

// GET /api/password_resets/{token}/edit?email=
func (h *APIHandler) GetPasswordReset(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	email := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("email")))
	user, err := h.store.GetUserByEmail(email)
	if err != nil || user == nil || !user.Activated || !user.Authenticated("reset", token) {
		writeError(w, http.StatusBadRequest, "無効なパスワード再設定リンクです")
		return
	}
	if user.PasswordResetExpired() {
		writeError(w, http.StatusBadRequest, "パスワード再設定リンクの有効期限が切れています")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"email": email, "token": token})
}

// PATCH /api/password_resets/{token}
func (h *APIHandler) UpdatePasswordReset(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	var body struct {
		Email                string `json:"email"`
		Password             string `json:"password"`
		PasswordConfirmation string `json:"password_confirmation"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "リクエストが不正です")
		return
	}
	email := strings.ToLower(strings.TrimSpace(body.Email))
	user, err := h.store.GetUserByEmail(email)
	if err != nil || user == nil || !user.Activated || !user.Authenticated("reset", token) {
		writeError(w, http.StatusBadRequest, "無効なパスワード再設定リンクです")
		return
	}
	if user.PasswordResetExpired() {
		writeError(w, http.StatusBadRequest, "パスワード再設定リンクの有効期限が切れています")
		return
	}
	var errs []string
	if body.Password == "" {
		errs = append(errs, "パスワードを入力してください")
	}
	if body.Password != body.PasswordConfirmation {
		errs = append(errs, "パスワード確認が一致しません")
	}
	if body.Password != "" {
		if err := model.ValidatePassword(body.Password); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"errors": errs})
		return
	}
	if err := h.store.UpdatePassword(user.ID, body.Password); err != nil {
		writeError(w, http.StatusInternalServerError, "内部エラーが発生しました")
		return
	}
	if err := h.store.ClearResetDigest(user.ID); err != nil {
		log.Printf("ClearResetDigest: %v", err)
	}
	_ = user.Forget(h.store)
	logIn(w, user.ID, false, h.store)
	writeJSON(w, http.StatusOK, map[string]any{
		"message": "パスワードが再設定されました。",
		"user":    userToJSON(*user),
	})
}

