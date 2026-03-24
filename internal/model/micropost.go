package model

import "time"

// Micropost はマイクロポストを表す構造体です。
type Micropost struct {
	ID        int64
	Content   string
	UserID    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Validate はマイクロポストのバリデーションを実行します。
// 第2章の最初の状態では、まだ何も検証しません。
func (m *Micropost) Validate() []string {
	var errors []string
	if len([]rune(m.Content)) > 140 {
		errors = append(errors, "Content is too long (maximum is 140 characters)")
	}
	if m.Content == "" {
		errors = append(errors, "Content can't be blank")
	}
	return errors
}