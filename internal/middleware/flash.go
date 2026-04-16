package middleware

import (
	"context"
	"net/http"
)

const flashKey = "flash"

// Flash はリクエストからフラッシュCookieを読み取り、
// レスポンス後に削除するミドルウェア
func Flash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie("flash"); err == nil {
			ctx := context.WithValue(r.Context(), flashKey, cookie.Value)
			r = r.WithContext(ctx)
			DeleteCookie(w, "flash")
		}
		next.ServeHTTP(w, r)
	})
}