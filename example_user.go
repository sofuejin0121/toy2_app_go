//go:build ignore
package main

import (
	"fmt"
)
// User はユーザーを表す構造体
type User struct {
	Name string
	Email string
}

// NewUser は新しいUserを作成するファクトリ関数
func NewUser(name, email string) *User {
	return &User{
		Name: name,
		Email: email,
	}
}

// FormattedEmail はフォーマット済みのメールアドレスを返す
func (u *User) FormattedEmail() string {
	return fmt.Sprintf("%s <%s>", u.Name, u.Email)
}

func main() {
	example := &User{}
	fmt.Println(example.Name) // ゼロ値(空文字列)

	example.Name = "Example User"
	example.Email = "user@example.com"
	fmt.Println(example.FormattedEmail())
}