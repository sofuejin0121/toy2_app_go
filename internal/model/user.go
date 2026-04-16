package model

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
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
	Bio                  string
	PasswordDigest       string
	Password             string // 仮想フィールド(DBに保存しない)
	PasswordConfirmation string // 仮想フィールド(DBに保存しない)
	RememberDigest       string
	RememberToken        string // 仮想フィールド(DBに保存しない)
	Admin                bool
	ActivationToken      string // 仮想フィールド(DBに保存しない)
	ActivationDigest     string
	Activated            bool
	ActivatedAt          *time.Time
	ResetToken           string // 仮想フィールド(DBに保存しない)
	ResetDigest          string
	ResetSentAt          *time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// UserStoreInterfaceはUserモデルが永続化時に必要とする最小限のStore機能
// 具体的なStore型を参照するとimport cycleになるため、この小さなinterfaceで受ける
type UserStoreInterface interface {
	UpdateRememberDigest(userID int64, digest string) error
	UpdateActivation(UserID int64, activated bool, activatedAt time.Time) error
}

// MailerInterface はメール送信のインターフェイス
// model →　mailerの循環インポートを避けるため、modelパッケージ内に定義する
// mailer.LogMailerとmailer.SMTPMailer はこのインターフェイスを暗黙的に満たす
type MailerInterface interface {
	SendAccountActivation(user *User) error
	SendPasswordReset(user *User) error
}

// SendPasswordResetEmail はパスワード再設定用のメールを送信する
func (u *User) SendPasswordResetEmail(m MailerInterface) error {
	return m.SendPasswordReset(u)
}

// Activate はアカウントを有効にする
// DB更新成功後にローカルフィールドを更新し、メモリとDBの不整合を防ぐ
func (u *User) Activate(store UserStoreInterface) error {
	now := time.Now()
	if err := store.UpdateActivation(u.ID, true, now); err != nil {
		return err
	}
	u.Activated = true
	u.ActivatedAt = &now
	return nil
}

// SendActivationEmail は有効化用のメールを送信する
func (u *User) SendActivationEmail(m MailerInterface) error {
	return m.SendAccountActivation(u)
}

// PasswordResetExpired はパスワード再設定リンクの有効期限（2時間）が切れていればtrueを返す
func (u *User) PasswordResetExpired() bool {
	if u.ResetSentAt == nil {
		return true
	}
	return time.Since(*u.ResetSentAt) > 2*time.Hour
}

// CreateResetDigest はパスワード再設定用のトークンとダイジェストを作成および代入する
func (u *User) CreateResetDigest() error {
	token, err := NewToken()
	if err != nil {
		return err
	}
	u.ResetToken = token
	digest, err := Digest(token)
	if err != nil {
		return err
	}
	u.ResetDigest = digest
	return nil
}

// CreateActivationDigest は有効化トークンとダイジェストを作成および代入する
func (u *User) CreateActivationDigest() error {
	token, err := NewToken()
	if err != nil {
		return err
	}
	u.ActivationToken = token
	digest, err := Digest(token)
	if err != nil {
		return err
	}
	u.ActivationDigest = digest
	return nil
}

// Remember は永続的セッションのためにユーザーをデータベースに記憶する
func (u *User) Remember(store UserStoreInterface) error {
	token, err := NewToken()
	if err != nil {
		return err
	}
	u.RememberToken = token
	digest, err := Digest(token)
	if err != nil {
		return err
	}
	u.RememberDigest = digest
	if err := store.UpdateRememberDigest(u.ID, digest); err != nil {
		return err
	}
	return nil // ダイジェストはu.RememberDigestに保持される
}

// SessionToken はセッションハイジャック防止のためにセッショントークンを返す
// この記憶ダイジェストを再利用しているのは単に利便性のため
func (u *User) SessionToken(store UserStoreInterface) string {
	if u.RememberDigest != "" {
		return u.RememberDigest
	}
	u.Remember(store)
	return u.RememberDigest
}

// Digest は渡された文字列のハッシュ値を返す
// bcrypt.GenerateFromPasswordでパスワードやremember tokenなどの文字列をハッシュ化
// 入力: "mysecretpassword 出力: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"
func Digest(s string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// NewToken はランダムなトークンを返す
func NewToken() (string, error) {
	// 16バイト分の空のバイト列を作る
	b := make([]byte, 16)
	// rand.Read(b) 箱に暗号論的乱数を詰め込む
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// バイト列を文字列に変換
	return base64.URLEncoding.EncodeToString(b), nil
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

// Authenticated は渡されたトークンがダイジェストと一致すればtrueを返す
// Authenticated メソッドにresetケースを追加（internal/model/user.go）
func (u *User) Authenticated(attribute, token string) bool {
	var digest string
	switch attribute {
	case "remember":
		digest = u.RememberDigest
	case "activation":
		digest = u.ActivationDigest
	case "reset":
		digest = u.ResetDigest
	default:
		return false
	}
	if digest == "" {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(digest), []byte(token))
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
	// bioの長さに制限をかける（160文字以内）
	if len([]rune(u.Bio)) > 160 {
		errors = append(errors, "Bio is too long (maximum is 160 characters)")
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

// ValidatePassword はパスワードの単体バリデーションを行う
func ValidatePassword(password string) error {
	if strings.TrimSpace(password) == "" {
		return errors.New("Password can't be blank")
	}
	if len(password) < 6 {
		return errors.New("Password is too short (minimum is 6 characters)")
	}
	return nil
}

// GravatarURL はユーザー名から UI Avatars のイニシャルアイコン URL を返す
// 登録不要・完全無料・名前からイニシャルを自動生成する
func (u *User) GravatarURL(size int) string {
	name := url.QueryEscape(u.Name)
	return fmt.Sprintf(
		"https://ui-avatars.com/api/?name=%s&size=%d&background=random&color=fff&rounded=true&bold=true",
		name, size)
}

// Forget はユーザーのログイン情報を破棄する
func (u *User) Forget(store UserStoreInterface) error {
	u.RememberDigest = ""
	return store.UpdateRememberDigest(u.ID, "")
}
