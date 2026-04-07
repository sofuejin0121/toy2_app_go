package model

import (
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Userはユーザーを表す構造体
var validEmailRegex = regexp.MustCompile(`(?i)^[\w+\-.]+@[a-z\d\-.]+\.[a-z]+$`)
type User struct {
	ID                   int64
	Name                 string
	Email                string
	PasswordDigest       string
	Password             string // 仮想フィールド(DBに保存しない)
	PasswordConfirmation string // 仮想フィールド(DBに保存しない)
	CreatedAt            time.Time
	UpdatedAt            time.Time
}


// SetPassword はパスワードをbcryptでハッシュ化し、PasswordDigestに格納
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordDigest = string(hash)
	u.Password = password
	return nil
}

// Authenticate はパスワードを照合し、一致すればtrueを返す
func (u *User) Authenticate(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordDigest), []byte(password))
	return err == nil
}


// validate はユーザーのバリデーションを実行します
func (u *User) Validate() []string {
	var errors []string
	// 名前のバリデーション
	// strings.TrimSpaceは文字列の両端の空白を削除します。
	if strings.TrimSpace(u.Name) == "" {
		errors = append(errors, "Name can't be blank")
	}
	// 名前の長さに制限をかける
	if len(u.Name) > 50 {
		errors = append(errors, "Name is too long (maximum is 50 characters)")
	}
	// メールアドレスのバリデーション
	if strings.TrimSpace(u.Email) == "" {
		errors = append(errors, "Email can't be blank")
	}
	if len(u.Email) > 255 {
		errors = append(errors, "Email is too long (maximum is 255 characters)")
	}
	if strings.TrimSpace(u.Email) != "" && !validEmailRegex.MatchString(u.Email) {
		errors = append(errors, "Email is invalid")
	}
	// パスワードのバリデーション: 新規作成時(PasswordDigest)
	// またはパスワード変更時(Password非空)のみ実行
	if u.PasswordDigest == "" || u.Password != "" {
		if strings.TrimSpace(u.Password) == "" {
			errors = append(errors, "Password can't be blank")
		}
		if len(u.Password) < 8 {
			errors = append(errors, "Password is too short (minimum is 8 characters)")
		}
		if u.Password != u.PasswordConfirmation {
			errors = append(errors, "Password confirmation doesn't match")
		}
	}
	return errors
}
