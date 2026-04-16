package model

import "time"

// Like はユーザーがマイクロポストに「いいね」した関係を表す。
type Like struct {
	ID          int64
	UserID      int64
	MicropostID int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Validate は Like の整合性を検証する。
func (l *Like) Validate() map[string]string {
	errors := make(map[string]string)
	if l.UserID == 0 {
		errors["user_id"] = "User ID is required"
	}
	if l.MicropostID == 0 {
		errors["micropost_id"] = "Micropost ID is required"
	}
	return errors
}
