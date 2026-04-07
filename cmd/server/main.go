package main

import (
	"log"
	"net/http"
	"os"


	"github.com/sofuejin0121/toy_app_go/internal/handler"
	"github.com/sofuejin0121/toy_app_go/internal/middleware"
	"github.com/sofuejin0121/toy_app_go/internal/store"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	s, err := store.New("db/toy.db")
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	userHandler := handler.NewUserHandler(s)
	micropostHandler := handler.NewMicropostHandler(s)

	// StaticPagesハンドラーを作成
	staticHandler := handler.NewStaticHandler()

	mux := http.NewServeMux()

	// Micropostsリソース
	mux.HandleFunc("GET /microposts", micropostHandler.Index)
	mux.HandleFunc("GET /microposts/new", micropostHandler.New)
	mux.HandleFunc("GET /microposts/{id}", micropostHandler.Show)
	mux.HandleFunc("GET /microposts/{id}/edit", micropostHandler.Edit)
	mux.HandleFunc("POST /microposts", micropostHandler.Create)
	mux.HandleFunc("PATCH /microposts/{id}", micropostHandler.Update)
	mux.HandleFunc("DELETE /microposts/{id}", micropostHandler.Destroy)

	// Usersリソース
	mux.HandleFunc("GET /users", userHandler.Index) // ユーザー一覧ページ
	mux.HandleFunc("GET /users/new", userHandler.New) // ユーザー新規作成するページ
	mux.HandleFunc("GET /users/{id}", userHandler.Show) // 特定のユーザーを表示するページ
	mux.HandleFunc("GET /users/{id}/edit", userHandler.Edit) // 特定のユーザーを編集するページ
	mux.HandleFunc("POST /users", userHandler.Create) // ユーザーを作成する
	mux.HandleFunc("PATCH /users/{id}", userHandler.Update) // 特定のユーザーを更新する
	mux.HandleFunc("DELETE /users/{id}", userHandler.Destroy) // 特定のユーザーを削除する

	// StaticPages用ルーティング
	mux.HandleFunc("GET /{$}", staticHandler.Home)
	mux.HandleFunc("GET /help", staticHandler.Help)
	mux.HandleFunc("GET /about", staticHandler.About)
	mux.HandleFunc("GET /contact", staticHandler.Contact)

	// Users用ルーティング
	mux.HandleFunc("GET /signup", userHandler.New)

	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))


	h := middleware.MethodOverride(mux)

	log.Printf("Starting server on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, h))
}
