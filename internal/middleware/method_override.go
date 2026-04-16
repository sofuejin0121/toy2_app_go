package middleware

import (
	"net/http"
	"strings"
)

// MethodOverride はPOSTフォームからPATCH/DELETEを扱うためのミドルウェア
func MethodOverride(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			ct := r.Header.Get("Content-Type")
			if strings.HasPrefix(ct, "multipart/form-data") {
				// multipart の場合: ParseForm() だと r.Form が空マップで初期化されるだけで
				// ボディが読まれない。その後 r.FormValue() が r.Form != nil を見て
				// ParseMultipartForm をスキップしてしまうので、ここで先に呼んでおく
				r.ParseMultipartForm(32 << 20)
			} else {
				r.ParseForm()
			}
			switch r.FormValue("_method") {
			case http.MethodPatch, http.MethodDelete:
				r.Method = r.FormValue("_method")
			}
		}
		next.ServeHTTP(w, r)
	})
}
