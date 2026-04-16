package model

import (
	"strings"
	"time"
)

// Micropost はマイクロポストを表す構造体です。
// InReplyToID は nil のとき通常投稿、non-nil のときリプライ（リプライ元のID）。
type Micropost struct {
	ID            int64
	Content       string
	UserID        int64
	ImagePath     string
	InReplyToID   *int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Validate はマイクロポストのバリデーションを実行します。
func (m *Micropost) Validate() []string {
	var errors []string
	if m.UserID == 0 {
		errors = append(errors, "ユーザーが存在しません")
	}
	if len([]rune(m.Content)) > 140 {
		errors = append(errors, "本文は140文字以内で入力してください")
	}
	if strings.TrimSpace(m.Content) == "" {
		errors = append(errors, "本文を入力してください")
	}
	return errors
}
