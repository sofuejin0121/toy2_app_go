package model

import (
	"errors"
	"time"
)

// Relationship はフォロー/フォロワーの関係を表す構造体
type Relationship struct {
	ID         int64
	FollowerID int64
	FollowedID int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Validate はRelationshipのバリデーションを実行する
func (r *Relationship) Validate() error {
	if r.FollowerID == 0 {
		return errors.New("follower_id can't be blank")
	}
	if r.FollowedID == 0 {
		return errors.New("followed_id can't be blank")
	}
	return nil
}
