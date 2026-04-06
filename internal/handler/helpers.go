package handler

import (
	"net/http"

	"github.com/sofuejin0121/toy_app_go/internal/model"
)

// getFlash はリクエストコンテキストからフラッシュメッセージを取得する
// 第７章以降で実装
func getFlash(r *http.Request) map[string]string {
	return nil
}


// isLoggedIn はユーザーがログイン中かどうかを返す
func isLoggedIn(r *http.Request) bool {
	return false
}

// currentUser はリクエストコンテキストからログイン中のユーザーを返す
func currentUser(r *http.Request) *model.User {
	return nil
}