package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/sofuejin0121/toy_app_go/internal/model"
	"github.com/sofuejin0121/toy_app_go/web/components"
)

// getFlash はリクエストコンテキストからフラッシュメッセージを取得する
// 第７章以降で実装
func getFlash(r *http.Request) map[string]string {
	cookie, err := r.Cookie("flash")
	if err != nil {
		return nil
	}
	raw, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return nil
	}
	var flash map[string]string
	if err := json.Unmarshal([]byte(raw), &flash); err != nil {
		return nil
	}
	return flash
}

// clearFlash はフラッシュメッセージ用Cookieを削除する
func clearFlash(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "flash",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}

// setFlash はフラッシュメッセージをCookieに保存する
func setFlash(w http.ResponseWriter, kind, message string) {
	flash := map[string]string{kind: message}
	b, _ := json.Marshal(flash)
	http.SetCookie(w, &http.Cookie{
		Name:  "flash",
		Value: url.QueryEscape(string(b)),
		Path:  "/",
	})
}

// isLoggedIn はユーザーがログイン中かどうかを返す
func isLoggedIn(r *http.Request) bool {
	return false
}

// currentUser はリクエストコンテキストからログイン中のユーザーを返す
func currentUser(r *http.Request) *model.User {
	return nil
}

// setDebugInfo はデバッグ情報をページデータにセットします。
func (h *UserHandler) setDebugInfo(data *components.UserPageData, r *http.Request) {
	if os.Getenv("APP_ENV") != "production" {
		data.Debug = true
		data.DebugInfo = fmt.Sprintf("%+v", data)
	}
}
