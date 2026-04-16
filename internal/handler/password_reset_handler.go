package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sofuejin0121/toy_app_go/internal/mailer"
	"github.com/sofuejin0121/toy_app_go/internal/middleware"
	"github.com/sofuejin0121/toy_app_go/internal/model"
	"github.com/sofuejin0121/toy_app_go/internal/store"
	"github.com/sofuejin0121/toy_app_go/web/components"
)

// PasswordResetHandler はパスワード再設定を処理するハンドラーです。
type PasswordResetHandler struct {
	store  *store.Store
	mailer mailer.Mailer
}

// NewPasswordResetHandler は新しいPasswordResetHandlerを作成します。
func NewPasswordResetHandler(s *store.Store, m mailer.Mailer) *PasswordResetHandler {
	return &PasswordResetHandler{store: s, mailer: m}
}

// setDebugInfo はデバッグ情報をページデータにセットします。
func (h *PasswordResetHandler) setDebugInfo(data *components.PasswordResetPageData, r *http.Request) {
	if os.Getenv("APP_ENV") != "production" {
		data.Debug = true
		data.DebugInfo = fmt.Sprintf("%+v", data)
	}
}

// getUser はメールアドレスからユーザーを検索するヘルパー。
// クエリパラメータ（GET）またはフォーム値（POST/PATCH）からメールアドレスを取得する。
func (h *PasswordResetHandler) getUser(r *http.Request) *model.User {
	email := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("email")))
	if email == "" {
		email = strings.ToLower(strings.TrimSpace(r.FormValue("email")))
	}
	user, err := h.store.GetUserByEmail(email)
	if err != nil || user == nil {
		return nil
	}
	return user
}

// validUser は正しいユーザーかどうか確認するヘルパー。
// ユーザーが存在し、有効化されており、トークンが認証済みであることを確認する。
func (h *PasswordResetHandler) validUser(user *model.User, token string) bool {
	return user != nil && user.Activated &&
		user.Authenticated("reset", token)
}

// New はパスワード再設定リクエストフォームを表示します。
func (h *PasswordResetHandler) New(w http.ResponseWriter, r *http.Request) {
	data := components.PasswordResetPageData{
		Title:       "Forgot password",
		Flash:       getFlash(r),
		CSRFToken:   middleware.CSRFTokenFromContext(r),
		LoggedIn:    isLoggedIn(r),
		CurrentUser: currentUser(r),
	}
	h.setDebugInfo(&data, r)
	_ = components.PasswordResetNew(data).Render(r.Context(), w)
}

func (h *PasswordResetHandler) Create(w http.ResponseWriter, r *http.Request) {
	email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
	user, err := h.store.GetUserByEmail(email)
	if err != nil || user == nil {
		data := components.PasswordResetPageData{
			Title:       "Forgot password",
			Flash:       map[string]string{"danger": "Email address not found"},
			CSRFToken:   middleware.CSRFTokenFromContext(r),
			LoggedIn:    isLoggedIn(r),
			CurrentUser: currentUser(r),
		}
		h.setDebugInfo(&data, r)
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = components.PasswordResetNew(data).Render(r.Context(), w)
		return
	}
	if err := user.CreateResetDigest(); err != nil {
		log.Printf("CreateResetDigest: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := h.store.UpdateResetDigest(user.ID, user.ResetDigest, time.Now()); err != nil {
		log.Printf("UpdateResetDigest: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := user.SendPasswordResetEmail(h.mailer); err != nil {
		log.Printf("SendPasswordResetEmail: %v", err)
	}
	setFlash(w, "info", "Email sent with password reset instructions")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Edit はパスワード再設定フォームを表示する。
// URLパスのトークンとクエリパラメータのメールアドレスでユーザーを認証する。
func (h *PasswordResetHandler) Edit(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("id")
	user := h.getUser(r)
	if !h.validUser(user, token) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// 期限切れチェック
	if user.PasswordResetExpired() {
		setFlash(w, "danger", "Password reset has expired.")
		http.Redirect(w, r, "/password_resets/new", http.StatusSeeOther)
		return
	}
	// パスワード再設定フォームを表示
	data := components.PasswordResetPageData{
		Title:       "Reset password",
		Flash:       getFlash(r),
		CSRFToken:   middleware.CSRFTokenFromContext(r),
		LoggedIn:    isLoggedIn(r),
		CurrentUser: currentUser(r),
		User:        user,
		Token:       token,
	}
	h.setDebugInfo(&data, r)
	_ = components.PasswordResetEdit(data).Render(r.Context(), w)
}

// Update はパスワード再設定を実行する。
// 4つのケースに対応:
//  1. 期限切れチェック
//  2. 無効なパスワード（バリデーションエラー）
//  3. 空のパスワード
//  4. 正しいパスワードで更新成功
func (h *PasswordResetHandler) Update(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("id")
	user := h.getUser(r)

	if !h.validUser(user, token) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// （1）への対応: 期限切れチェック
	if user.PasswordResetExpired() {
		setFlash(w, "danger", "Password reset has expired.")
		http.Redirect(w, r, "/password_resets/new", http.StatusSeeOther)
		return
	}

	password := r.FormValue("password")
	passwordConfirmation := r.FormValue("password_confirmation")
	var errs []string

	// （3）への対応: 空パスワードのチェック
	if password == "" {
		errs = append(errs, "Password can't be empty")
	}
	if password != passwordConfirmation {
		errs = append(errs, "Password confirmation doesn't match Password")
	}
	if password != "" {
		if err := model.ValidatePassword(password); err != nil {
			errs = append(errs, err.Error())
		}
	}

	// （2）と（3）への対応: バリデーションエラー
	if len(errs) > 0 {
		data := components.PasswordResetPageData{
			Title:       "Reset password",
			Flash:       getFlash(r),
			CSRFToken:   middleware.CSRFTokenFromContext(r),
			LoggedIn:    isLoggedIn(r),
			CurrentUser: currentUser(r),
			User:        user,
			Token:       token,
			Errors:      errs,
		}
		h.setDebugInfo(&data, r)
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = components.PasswordResetEdit(data).Render(r.Context(), w)
		return
	}

	if err := h.store.UpdatePassword(user.ID, password); err != nil {
		log.Printf("UpdatePassword(%d): %v", user.ID, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// （4）への対応: 更新成功
	// リセットダイジェストをクリア（セキュリティ対策: 同じリンクの再利用を防ぐ）
	if err := h.store.ClearResetDigest(user.ID); err != nil {
		log.Printf("ClearResetDigest: %v", err)
	}
	// 既存セッションをすべて無効化（セッションハイジャック対策）
	_ = user.Forget(h.store)
	logIn(w, user.ID, false, h.store)
	setFlash(w, "success", "Password has been reset.")
	http.Redirect(w, r, fmt.Sprintf("/users/%d", user.ID), http.StatusSeeOther)
}
