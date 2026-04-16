package model

import "time"

type Notification struct {
	ID int64
	UserID int64
	ActorID int64
	ActionType string // "like", "follow", "reply"
	TargetID *int64 // いいね投稿ID(フォローはnil)
	Read bool
	CreatedAt time.Time
}