package model

import "time"

// Userはユーザーを表す構造体
type User struct {
	ID        int64
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// validate はユーザーのバリデーションを実行します
func (u *User) Validate() []string {
	var errors []string
	if u.Name == "" {
		errors = append(errors, "Name can't be blank")
	}
	if u.Email == "" {
		errors = append(errors, "Email can't be blank")
	}

	return errors
}
