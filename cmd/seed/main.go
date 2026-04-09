package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sofuejin0121/toy_app_go/internal/model"
	"github.com/sofuejin0121/toy_app_go/internal/store"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "db/toy.db"
	}

	s, err := store.New(dbURL)
	if err != nil {
		log.Fatalf("open store: %v", err)
	}
	defer s.Close()

	// 管理者ユーザーを作成
	admin := &model.User{
		Name:  "Example User",
		Email: "example@example.com",
		Admin: true,
	}
	if err := admin.SetPassword("foobar"); err != nil {
		log.Fatal(err)
	}
	if err := s.CreateUser(admin); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created admin user: %s (%s)\n",
		admin.Name, admin.Email)

	// サンプルユーザーを99人作成
	for i := 1; i <= 99; i++ {
		name := fmt.Sprintf("User %d", i)
		email := fmt.Sprintf("user-%d@example.com", i)
		user := &model.User{
			Name:  name,
			Email: strings.ToLower(email),
		}
		if err := user.SetPassword("password"); err != nil {
			log.Fatal(err)
		}
		if err := s.CreateUser(user); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Created 99 sample users")
	fmt.Println("Seed completed!")
}
