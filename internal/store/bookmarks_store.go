package store

import (
	"database/sql"
	"errors"
)

// Bookmark はユーザーがマイクロポストをブックマークする。
func (s *Store) Bookmark(userID, micropostID int64) error {
	_, err := s.db.Exec(`
	INSERT OR IGNORE INTO bookmarks (user_id, micropost_id, created_at, updated_at)
	VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`, userID, micropostID)
	return err
}

// UnBookmark はユーザーのブックマークを取り消す。
func (s *Store) UnBookmark(userID, micropostID int64) error {
	_, err := s.db.Exec(`
	DELETE FROM bookmarks WHERE user_id = ? AND micropost_id = ?`, userID, micropostID)
	return err
}

// IsBookmarked は userID が micropostID をブックマーク済みかどうかを返す。
func (s *Store) IsBookmarked(userID, micropostID int64) (bool, error) {
	var id int64
	err := s.db.QueryRow(`
	SELECT id FROM bookmarks WHERE user_id = ? AND micropost_id = ?`, userID, micropostID).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

// CountBookmarkedMicroposts は userID がブックマークした投稿の総数を返す。
func (s *Store) CountBookmarkedMicroposts(userID int64) (int, error) {
	var count int
	err := s.db.QueryRow(
		"SELECT COUNT(*) FROM bookmarks WHERE user_id = ?", userID).Scan(&count)
	return count, err
}

// GetBookmarkedPosts は userID がブックマークした投稿を FeedItem として返す。
// viewerID は IsLiked フラグ計算に使う（0 なら常に false）。
func (s *Store) GetBookmarkedPosts(userID, viewerID int64, page, perPage int) ([]FeedItem, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * perPage
	rows, err := s.db.Query(`
		SELECT`+feedSelectCols+`
		FROM microposts m
		JOIN bookmarks b ON b.micropost_id = m.id
		JOIN users u ON m.user_id = u.id`+feedJoinClauses+`
		WHERE b.user_id = ?
		ORDER BY b.created_at DESC
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
