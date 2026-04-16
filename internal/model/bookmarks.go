package model

import "time"

type Bookmark struct {
	ID          int64
	UserID      int64
	MicropostID int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
