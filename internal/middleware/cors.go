package middleware

import (
	"net/http"
	"os"
)

// CORS は指定したオリジンからのクロスオリジンリクエストを許可するミドルウェアです。
// FRONTEND_URL 環境変数が設定されている場合のみ有効になります。
func CORS(next http.Handler) http.Handler {
	frontendURL := os.Getenv("FRONTEND_URL")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if frontendURL != "" {
			w.Header().Set("Access-Control-Allow-Origin", frontendURL)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With")
		}
		// preflight リクエストはここで終了
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
