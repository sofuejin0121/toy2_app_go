package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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

	now := time.Now()

	// 管理者ユーザーを作成
	admin := &model.User{
		Name:        "Example User",
		Email:       "example@example.com",
		Admin:       true,
		Activated:   true,
		ActivatedAt: &now,
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
			Name:        name,
			Email:       strings.ToLower(email),
			Activated:   true,
			ActivatedAt: &now,
		}
		if err := user.SetPassword("password"); err != nil {
			log.Fatal(err)
		}
		if err := s.CreateUser(user); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Created 99 sample users")

	// ユーザーの一部を対象にマイクロポストを生成する
	db := s.DB()
	rows, err := db.Query("SELECT id FROM users ORDER BY created_at LIMIT 6")
	if err != nil {
		log.Printf("query users for microposts: %v", err)
	} else {
		defer rows.Close()

		var userIDs []int64
		for rows.Next() {
			var id int64
			rows.Scan(&id)
			userIDs = append(userIDs, id)
		}

		sentences := []string{
			"Lorem ipsum dolor sit amet",
			"Consectetur adipiscing elit",
			"Sed do eiusmod tempor incididunt",
			"Ut labore et dolore magna aliqua",
			"Ut enim ad minim veniam",
		}

		for i := 0; i < 50; i++ {
			content := sentences[i%len(sentences)]
			for _, userID := range userIDs {
				db.Exec(`INSERT INTO microposts (content, user_id, created_at, updated_at)
					VALUES (?, ?, datetime('now'), datetime('now'))`,
					content, userID)
			}
		}
		fmt.Printf("Created microposts for %d users\n", len(userIDs))
	}

	// フォロー関係を作成
	// ユーザー1（admin）がユーザー3〜51をフォロー
	for i := int64(3); i <= 51; i++ {
		if err := s.Follow(admin.ID, i); err != nil {
			log.Printf("Follow %d -> %d: %v", admin.ID, i, err)
		}
	}
	// ユーザー4〜41がユーザー1（admin）をフォロー
	for i := int64(4); i <= 41; i++ {
		if err := s.Follow(i, admin.ID); err != nil {
			log.Printf("Follow %d -> %d: %v", i, admin.ID, err)
		}
	}
	fmt.Println("Created follow relationships")

	fmt.Println("Seed completed!")
}
