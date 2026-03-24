package middleware

import "net/http"

// MethodOverride はPOSTフォームからPATCH/DELETEを扱うためのミドルウェア
func MethodOverride(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err == nil {
				switch r.Form.Get("_method") {
				case http.MethodPatch, http.MethodDelete:
					r.Method = r.Form.Get("_method")
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}
