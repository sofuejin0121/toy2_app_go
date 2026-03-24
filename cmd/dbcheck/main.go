package main

import (
	"fmt"
	"log"

	"github.com/sofuejin0121/toy_app_go/internal/store"
)

func main() {
	// DBに接続
	s, err := store.New("db/toy.db")
	if err != nil {
		log.Fatal(err)
	}

	// 最初のユーザーを取得
	firstUser, err := s.GetUser(1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("first_user: %+v\n\n", firstUser)

	// ユーザーのマイクロポストを取得
	microposts, err := s.GetMicropostsByUserID(firstUser.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("first_user.microposts: %+v\n\n", microposts)

	// 最初のマイクロポストからユーザーを取得
	if len(microposts) > 0 {
		micropost := microposts[0]
		fmt.Printf("micropost: %+v\n\n", micropost)

		user, err := s.GetUserByMicropostID(micropost.ID)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("micropost.user: %+v\n", user)
	} else {
		fmt.Println("No microposts found for this user.")
	}
}
