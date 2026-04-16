package middleware

import (
	"context"
	"net/http"
	"strings"
)

// CSRF はCSRFトークンの生成・検証を行うミドルウェア
func CSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// /api/ プレフィックスのルートはCSRF検証をスキップ
		// 同一オリジンのReact SPAからのリクエストはSameSiteCookieで保護される
		if strings.HasPrefix(r.URL.Path, "/api/") {
			next.ServeHTTP(w, r)
			return
		}

		csrfToken := getOrCreateCSRFToken(w, r)
		ctx := context.WithValue(r.Context(), csrfTokenKey, csrfToken)
		r = r.WithContext(ctx)

		if r.Method == "POST" || r.Method == "PATCH" || r.Method == "DELETE" {
			formToken := r.FormValue("csrf_token")
		if formToken != csrfToken {
			http.Error(w, "アクセスが拒否されました", http.StatusForbidden)
			return
		}
		}
		next.ServeHTTP(w, r)
	})
}