package store

import (
	"database/sql"
	"errors"

	"github.com/sofuejin0121/toy_app_go/internal/model"
)

// Like はユーザーがマイクロポストにいいねする。
// 既にいいね済みの場合は UNIQUE 制約エラーとなるが、冪等に扱う。
func (s *Store) Like(userID, micropostID int64) error {
	var ownerID int64
	s.db.QueryRow(`
	SELECT user_id FROM microposts WHERE id = ?`, micropostID).Scan(&ownerID)

	_, err := s.db.Exec(`
		INSERT OR IGNORE INTO likes (user_id, micropost_id, created_at, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
		userID, micropostID)
	// いいねした人と投稿のオーナーが別人なら通知を作成する
	if ownerID != userID {
		s.CreateNotification(ownerID, userID, "like", &micropostID)
	}
	return err
}

// Unlike はユーザーのいいねを取り消す。
func (s *Store) Unlike(userID, micropostID int64) error {
	_, err := s.db.Exec(
		"DELETE FROM likes WHERE user_id = ? AND micropost_id = ?",
		userID, micropostID)
	return err
}

// IsLiked は userID が micropostID にいいね済みかどうかを返す。
func (s *Store) IsLiked(userID, micropostID int64) (bool, error) {
	var id int64
	err := s.db.QueryRow(
		"SELECT id FROM likes WHERE user_id = ? AND micropost_id = ?",
		userID, micropostID).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

// CountLikes は micropostID に対するいいね数を返す。
func (s *Store) CountLikes(micropostID int64) (int, error) {
	var count int
	err := s.db.QueryRow(
		"SELECT COUNT(*) FROM likes WHERE micropost_id = ?", micropostID).Scan(&count)
	return count, err
}

// GetLike は userID と micropostID で like レコードを取得する。
// 存在しない場合は sql.ErrNoRows を返す。
func (s *Store) GetLike(userID, micropostID int64) (*model.Like, error) {
	like := &model.Like{}
	err := s.db.QueryRow(`
		SELECT id, user_id, micropost_id, created_at, updated_at
		FROM likes WHERE user_id = ? AND micropost_id = ?`,
		userID, micropostID).Scan(
		&like.ID, &like.UserID, &like.MicropostID, &like.CreatedAt, &like.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return like, nil
}

// GetLikeByID は likes.id で like レコードを取得する。
func (s *Store) GetLikeByID(id int64) (*model.Like, error) {
	like := &model.Like{}
	err := s.db.QueryRow(`
		SELECT id, user_id, micropost_id, created_at, updated_at
		FROM likes WHERE id = ?`, id).Scan(
		&like.ID, &like.UserID, &like.MicropostID, &like.CreatedAt, &like.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return like, nil
}

// CountLikedMicroposts は userID がいいねしたマイクロポストの総数を返す。
func (s *Store) CountLikedMicroposts(userID int64) (int, error) {
	var count int
	err := s.db.QueryRow(
		"SELECT COUNT(*) FROM likes WHERE user_id = ?", userID).Scan(&count)
	return count, err
}

// PaginateLikedMicropostsAsFeedItems は userID がいいねしたマイクロポストを
// FeedItem として返す。viewerID は IsLiked フラグ計算に使う（0 なら常に false）。
func (s *Store) PaginateLikedMicropostsAsFeedItems(userID, viewerID int64, page, perPage int) ([]FeedItem, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * perPage
	rows, err := s.db.Query(`
		SELECT`+feedSelectCols+`
		FROM microposts m
		JOIN likes l ON l.micropost_id = m.id
		JOIN users u ON m.user_id = u.id`+feedJoinClauses+`
		WHERE l.user_id = ?
		ORDER BY l.created_at DESC
		LIMIT ? OFFSET ?`, viewerID, viewerID, userID, perPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []FeedItem
	for rows.Next() {
		var item FeedItem
		if err := scanFeedItem(rows, &item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
