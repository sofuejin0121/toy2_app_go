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
		m = &mailer.SMTPMailer{
			Host:     os.Getenv("SMTP_HOST"),
			Port:     587,
			Username: os.Getenv("SMTP_USERNAME"),
			Password: os.Getenv("SMTP_PASSWORD"),
			From:     os.Getenv("MAILER_FROM"),
			AppHost:  os.Getenv("APP_HOST"),
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

	userHandler := handler.NewUserHandler(s, m)
	sessionHandler := handler.NewSessionHandler(s)
	staticHandler := handler.NewStaticHandler(s)
	accountActivationHandler := handler.NewAccountActivationHandler(s)
	passwordResetHandler := handler.NewPasswordResetHandler(s, m)
	imageDir := filepath.Join("web", "static", "images", "microposts")
	micropostHandler := handler.NewMicropostHandler(s, imageDir)
	relationshipHandler := handler.NewRelationshipHandler(s, m)
	likeHandler := handler.NewLikeHandler(s, m)
	userPreferenceHandler := handler.NewUserPreferenceHandler(s, m)
	notificationHandler := handler.NewNotificationHandler(s)
	adminHandler := handler.NewAdminHandler(s)
	bookmarkHandler := handler.NewBookmarkHandler(s)

	mux := http.NewServeMux()

	// 静的ページ
	mux.HandleFunc("GET /{$}", staticHandler.Home)
	mux.HandleFunc("GET /help", staticHandler.Help)
	mux.HandleFunc("GET /about", staticHandler.About)
	mux.HandleFunc("GET /contact", staticHandler.Contact)

	// ユーザー
	mux.HandleFunc("GET /signup", userHandler.New)
	mux.HandleFunc("POST /users", userHandler.Create)
	mux.HandleFunc("GET /users", handler.RequireLogin(userHandler.Index))
	mux.HandleFunc("GET /users/{id}", userHandler.Show)
	mux.HandleFunc("GET /users/{id}/edit",
		handler.RequireLogin(userHandler.RequireCorrectUser(userHandler.Edit)))
	mux.HandleFunc("PATCH /users/{id}",
		handler.RequireLogin(userHandler.RequireCorrectUser(userHandler.Update)))
	mux.HandleFunc("DELETE /users/{id}",
		handler.RequireLogin(userHandler.RequireAdmin(userHandler.Destroy)))
	mux.HandleFunc("GET /users/{id}/following",
		handler.RequireLogin(userHandler.Following))
	mux.HandleFunc("GET /users/{id}/followers",
		handler.RequireLogin(userHandler.Followers))
	mux.HandleFunc("GET /users/{id}/likes",
		handler.RequireLogin(userHandler.LikedPosts))
	mux.HandleFunc("GET /users/{id}/bookmarks",
		handler.RequireLogin(userHandler.BookmarkedPosts))

	// セッション
	mux.HandleFunc("GET /login", sessionHandler.New)
	mux.HandleFunc("POST /login", sessionHandler.Create)
	mux.HandleFunc("DELETE /logout", sessionHandler.Destroy)

	// アカウント有効化
	mux.HandleFunc("GET /account_activation/{id}/edit", accountActivationHandler.Edit)

	// パスワード再設定
	mux.HandleFunc("GET /password_resets/new", passwordResetHandler.New)
	mux.HandleFunc("POST /password_resets", passwordResetHandler.Create)
	mux.HandleFunc("GET /password_resets/{id}/edit", passwordResetHandler.Edit)
	mux.HandleFunc("PATCH /password_resets/{id}", passwordResetHandler.Update)

	// マイクロポスト
	mux.HandleFunc("POST /microposts",
		handler.RequireLogin(micropostHandler.Create))
	mux.HandleFunc("GET /microposts/{id}", micropostHandler.Show)
	mux.HandleFunc("DELETE /microposts/{id}",
		handler.RequireLogin(
			micropostHandler.RequireCorrectUser(micropostHandler.Destroy)))
	mux.HandleFunc("GET /microposts", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusFound)
	})

	// リレーションシップ
	mux.HandleFunc("POST /relationships",
		handler.RequireLogin(relationshipHandler.Create))
	mux.HandleFunc("DELETE /relationships/{id}",
		handler.RequireLogin(relationshipHandler.Destroy))

	// いいね: POST /likes でいいね作成、DELETE /likes/{micropost_id} でいいね解除
	mux.HandleFunc("POST /likes",
		handler.RequireLogin(likeHandler.Create))
	mux.HandleFunc("DELETE /likes/{id}",
		handler.RequireLogin(likeHandler.Destroy))

	// ブックマーク
	mux.HandleFunc("POST /bookmarks",
		handler.RequireLogin(bookmarkHandler.Create))
	mux.HandleFunc("DELETE /bookmarks/{id}",
		handler.RequireLogin(bookmarkHandler.Destroy))
	mux.HandleFunc("GET /notifications", handler.RequireLogin(notificationHandler.Index))
	mux.HandleFunc("GET /admin", handler.RequireLogin(userHandler.RequireAdmin(adminHandler.Index)))

	// 通知設定
	mux.HandleFunc("GET /settings", handler.RequireLogin(userPreferenceHandler.Edit))
	mux.HandleFunc("PATCH /settings", handler.RequireLogin(userPreferenceHandler.Update))
	mux.HandleFunc("DELETE /notifications/{id}", handler.RequireLogin(notificationHandler.Destroy))

	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Flash → Auth → MethodOverride の順で適用
	// Flash を先に挟むことで、リクエストごとに flash Cookie を読み取った後に削除し
	// フラッシュメッセージが次のリクエストに持ち越されないようにする
	h := middleware.Flash(middleware.Auth(s)(middleware.MethodOverride(mux)))

	log.Printf("Starting server on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, h))
}
