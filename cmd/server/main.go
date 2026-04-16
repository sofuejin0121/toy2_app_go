package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sofuejin0121/toy_app_go/internal/handler"
	"github.com/sofuejin0121/toy_app_go/internal/mailer"
	"github.com/sofuejin0121/toy_app_go/internal/middleware"
	"github.com/sofuejin0121/toy_app_go/internal/store"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	var m mailer.Mailer
	if os.Getenv("GO_ENV") == "production" {
		switch os.Getenv("MAILER") {
		case "smtp":
			smtpPort := 587
			m = &mailer.SMTPMailer{
				Host:     os.Getenv("SMTP_HOST"),
				Port:     smtpPort,
				Username: os.Getenv("SMTP_USERNAME"),
				Password: os.Getenv("SMTP_PASSWORD"),
				From:     os.Getenv("MAILER_FROM"),
				AppHost:  os.Getenv("APP_HOST"),
			}
		default:
			m = &mailer.ResendMailer{
				APIKey:  os.Getenv("RESEND_API_KEY"),
				From:    os.Getenv("MAILER_FROM"),
				AppHost: os.Getenv("APP_HOST"),
			}
		}
	} else {
		m = &mailer.LogMailer{
			From: mailer.DefaultFrom,
			Host: "localhost:" + port,
		}
	}

	dbPath := os.Getenv("DATABASE_URL")
	if dbPath == "" {
		dbPath = "db/toy.db"
	}
	s, err := store.New(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	imageDir := filepath.Join("web", "static", "images", "microposts")
	apiHandler := handler.NewAPIHandler(s, m, imageDir)

	mux := http.NewServeMux()

	// ---- JSON API ルート (/api/ プレフィックス) ----

	// 認証
	mux.HandleFunc("GET /api/me", apiHandler.Me)
	mux.HandleFunc("POST /api/login", apiHandler.Login)
	mux.HandleFunc("DELETE /api/logout", apiHandler.Logout)

	// ユーザー
	mux.HandleFunc("POST /api/users", apiHandler.CreateUser)
	mux.HandleFunc("GET /api/users", handler.RequireLogin(apiHandler.ListUsers))
	mux.HandleFunc("GET /api/users/{id}", apiHandler.GetUser)
	mux.HandleFunc("PATCH /api/users/{id}", handler.RequireLogin(apiHandler.UpdateUser))
	mux.HandleFunc("DELETE /api/users/{id}", handler.RequireLogin(apiHandler.DeleteUser))
	mux.HandleFunc("GET /api/users/{id}/following", handler.RequireLogin(apiHandler.GetFollowing))
	mux.HandleFunc("GET /api/users/{id}/followers", handler.RequireLogin(apiHandler.GetFollowers))
	mux.HandleFunc("GET /api/users/{id}/likes", handler.RequireLogin(apiHandler.GetUserLikes))
	mux.HandleFunc("GET /api/users/{id}/bookmarks", handler.RequireLogin(apiHandler.GetUserBookmarks))

	// フィード
	mux.HandleFunc("GET /api/feed", handler.RequireLogin(apiHandler.Feed))

	// マイクロポスト
	mux.HandleFunc("GET /api/microposts/{id}", apiHandler.GetMicropost)
	mux.HandleFunc("POST /api/microposts", handler.RequireLogin(apiHandler.CreateMicropost))
	mux.HandleFunc("DELETE /api/microposts/{id}", handler.RequireLogin(apiHandler.DeleteMicropost))

	// フォロー
	mux.HandleFunc("POST /api/relationships", handler.RequireLogin(apiHandler.CreateRelationship))
	mux.HandleFunc("DELETE /api/relationships/{id}", handler.RequireLogin(apiHandler.DeleteRelationship))

	// いいね
	mux.HandleFunc("POST /api/likes", handler.RequireLogin(apiHandler.CreateLike))
	mux.HandleFunc("DELETE /api/likes/{id}", handler.RequireLogin(apiHandler.DeleteLike))

	// ブックマーク
	mux.HandleFunc("POST /api/bookmarks", handler.RequireLogin(apiHandler.CreateBookmark))
	mux.HandleFunc("DELETE /api/bookmarks/{id}", handler.RequireLogin(apiHandler.DeleteBookmark))

	// 通知
	mux.HandleFunc("GET /api/notifications", handler.RequireLogin(apiHandler.ListNotifications))
	mux.HandleFunc("GET /api/notifications/unread_count", apiHandler.UnreadNotificationCount)
	mux.HandleFunc("DELETE /api/notifications/{id}", handler.RequireLogin(apiHandler.DeleteNotification))

	// 管理者
	mux.HandleFunc("GET /api/admin", handler.RequireLogin(apiHandler.AdminStats))

	// 設定
	mux.HandleFunc("GET /api/settings", handler.RequireLogin(apiHandler.GetSettings))
	mux.HandleFunc("PATCH /api/settings", handler.RequireLogin(apiHandler.UpdateSettings))

	// アカウント有効化
	mux.HandleFunc("GET /api/account_activations/{token}/edit", apiHandler.ActivateAccount)

	// パスワードリセット
	mux.HandleFunc("POST /api/password_resets", apiHandler.CreatePasswordReset)
	mux.HandleFunc("GET /api/password_resets/{token}/edit", apiHandler.GetPasswordReset)
	mux.HandleFunc("PATCH /api/password_resets/{token}", apiHandler.UpdatePasswordReset)

	// ---- 静的ファイル ----

	// micropost画像
	imgFS := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", imgFS))

	// React SPAビルド成果物
	distDir := "frontend/dist"
	distFS := http.FileServer(http.Dir(distDir))

	// /assets/ 等の静的ファイルはそのまま配信
	mux.Handle("/assets/", distFS)

	// それ以外の全ルートはindex.htmlを返す（React Routerが処理）
	mux.HandleFunc("/", handler.ServeReact(distDir))

	// CORS → Flash → Auth → MethodOverride → CSRF の順でミドルウェア適用
	// FRONTEND_URL が設定されている場合のみ CORS ヘッダーが付与される
	h := middleware.CORS(middleware.Flash(middleware.Auth(s)(middleware.MethodOverride(middleware.CSRF(mux)))))

	log.Printf("Starting server on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, h))
}
