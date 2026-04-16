package middleware 

import (
	"context"
	"net/http"
)

// CSRF はCSRFトークンの生成・検証を行うミドルウェア
func CSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csrfToken := getOrCreateCSRFToken(w, r)
		ctx := context.WithValue(r.Context(), csrfTokenKey, csrfToken)
		r = r.WithContext(ctx)

		if r.Method == "POST" || r.Method == "PATCH" || r.Method == "DELETE" {
			formToken := r.FormValue("csrf_token")
			if formToken != csrfToken {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}