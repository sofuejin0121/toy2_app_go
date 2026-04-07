package main

import (
    "fmt"
    "log"

    "github.com/sofuejin0121/toy_app_go/internal/store"
)

func main() {
    s, err := store.New("db/toy.db")
    if err != nil {
        log.Fatal(err)
    }
    defer s.Close()

    users, err := s.AllUsers()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("User count: %d\n", len(users))

    if len(users) > 0 {
        user, err := s.GetUser(1)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("User: %+v\n", *user)
    }
}