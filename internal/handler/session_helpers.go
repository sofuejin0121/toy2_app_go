package handler

import (
	"net/http"
	"strconv"

	"github.com/sofuejin0121/toy_app_go/internal/middleware"
)

// logIn は渡されたユーザーでログインする。
// セッション固定攻撃対策として、古いセッションCookieを破棄してから
// 新しいセッションCookieにユーザーIDを保存する。
func logIn(w http.ResponseWriter, r *http.Request, userID int64) {
	middleware.ClearSessionCookie(w)
	middleware.SetSessionValue(w, "user_id", strconv.FormatInt(userID, 10))
}

// logOut は現在のユーザーをログアウトする。
func logOut(w http.ResponseWriter) {
	middleware.ClearSessionCookie(w)
}
