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
	mux.HandleFunc("GET /users", userHandler.Index)
	mux.HandleFunc("GET /users/new", userHandler.New)
	mux.HandleFunc("GET /users/{id}", userHandler.Show)
	mux.HandleFunc("GET /users/{id}/edit", userHandler.Edit)
	mux.HandleFunc("POST /users", userHandler.Create)
	mux.HandleFunc("PATCH /users/{id}", userHandler.Update)
	mux.HandleFunc("DELETE /users/{id}", userHandler.Destroy)

	// ルートURLをユーザー一覧に変更
	mux.HandleFunc("GET /", userHandler.Index)

	handler := middleware.MethodOverride(mux)

	log.Printf("Starting server on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
