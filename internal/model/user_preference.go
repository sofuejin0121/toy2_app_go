package model

import "time"

// UserPreference はユーザーのメール通知設定を表す構造体
type UserPreference struct {
	ID            int64
	UserID        int64
	EmailOnFollow bool // フォロー時にメールを送るか
	EmailOnLike   bool // いいね時にメールを送るか
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
