package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sofuejin0121/toy_app_go/internal/middleware"
	"github.com/sofuejin0121/toy_app_go/internal/model"
	"github.com/sofuejin0121/toy_app_go/internal/store"
)

// logIn は渡されたユーザーでログインする。
// セッション固定攻撃対策として、古いセッションCookieを破棄してから
// 新しいセッションCookieにユーザーIDを保存する。
func logIn(w http.ResponseWriter, userID int64, remember bool, s *store.Store) {
	// セッション固定攻撃対策: 古いセッションCookieを破棄
	middleware.ClearSessionCookie(w)
	// 新しいセッションCookieにユーザーIDを保存
	middleware.SetSessionValue(w, "user_id", strconv.FormatInt(userID, 10))
	user, err := s.GetUser(userID)
	if err != nil {
		return
	}
	middleware.SetSessionValue(w, "session_token", user.SessionToken(s))
	if remember {
		rememberUser(w, user, s)
	}
}

// isCurrentUser は渡されたユーザーIDが現在のユーザーと一致するばtrueを返す
func isCurrentUser(r *http.Request, userID int64) bool {
	user := currentUser(r)
	return user != nil && user.ID == userID
}

// forgetUser は永続的セッションのためにユーザーのログイン情報を破棄する
func forgetUser(w http.ResponseWriter, user *model.User, s *store.Store) {
	_ = user.Forget(s)
	// remember_user_idとremember_tokenを削除
	middleware.DeleteCookie(w, "remember_user_id")
	middleware.DeleteCookie(w, "remember_token")
}

// logOut は現在のユーザーをログアウトする。
func logOut(w http.ResponseWriter, user *model.User, s *store.Store) {
	forgetUser(w, user, s)
	middleware.ClearSessionCookie(w)
}

// rememberUser は永続的セッションのためにユーザーをデータベースに記憶する
func rememberUser(w http.ResponseWriter, user *model.User, s *store.Store) {
	if err := user.Remember(s); err != nil {
		log.Printf("rememberUser: %v", err)
		return
	}
	// ユーザーIDを署名付き永続cookieに保存
	http.SetCookie(w, &http.Cookie{
		Name:     "remember_user_id",
		Value:    middleware.SignUserID(user.ID),
		Expires:  time.Now().Add(middleware.PermanentCookieExpiry),
		Path:     "/",
		HttpOnly: true,
		Secure:   middleware.IsCrossOrigin(),
		SameSite: middleware.CookieSameSite(),
	})
	// 記憶トークンを永続cookieに保存
	http.SetCookie(w, &http.Cookie{
		Name:     "remember_token",
		Value:    user.RememberToken,
		Expires:  time.Now().Add(middleware.PermanentCookieExpiry),
		Path:     "/",
		HttpOnly: true,
		Secure:   middleware.IsCrossOrigin(),
		SameSite: middleware.CookieSameSite(),
	})
}

// RequireLogin はログインしていないユーザーをリダイレクトするミドルウェアです。
// /api/ プレフィックスのリクエストには JSON エラーを返します。
func RequireLogin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isLoggedIn(r) {
			if strings.HasPrefix(r.URL.Path, "/api/") {
				writeError(w, http.StatusUnauthorized, "unauthorized")
				return
			}
			storeLocation(w, r)
			setFlash(w, "danger", "Please log in.")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

// storeLocation はアクセスしようとしたURLをセッションに保存する
// GETリクエストの場合のみ保存する
func storeLocation(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		middleware.SetSessionValue(w, "forwarding_url", r.URL.String())
	}
}

// getForwardingURL は保存しておいた転送先URLを返し、読みだしたら削除する
func getForwardingURL(w http.ResponseWriter, r *http.Request) string {
	url := middleware.GetSessionValue(r, "forwarding_url")
	if url == "" {
		return ""
	}
	middleware.DeleteCookie(w, "forwarding_url")
	return url

}
