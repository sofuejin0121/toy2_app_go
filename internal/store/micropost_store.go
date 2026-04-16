package store

import (
	"database/sql"
	"fmt"

	"github.com/sofuejin0121/toy_app_go/internal/model"
)

// CreateMicropost はマイクロポストをデータベースに保存します。
// m.InReplyToID が nil のとき in_reply_to_id に NULL を挿入します。
func (s *Store) CreateMicropost(m *model.Micropost) error {
	now := nowString()
	result, err := s.db.Exec(`
		INSERT INTO microposts (content, user_id, image_path, in_reply_to_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		m.Content, m.UserID, m.ImagePath, m.InReplyToID, now, now)
	if err != nil {
		return fmt.Errorf("insert micropost: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	m.ID = id
	m.CreatedAt = parseTime(now)
	m.UpdatedAt = parseTime(now)
	return nil
}

// GetMicropost はIDでマイクロポストを取得します。
func (s *Store) GetMicropost(id int64) (*model.Micropost, error) {
	var m model.Micropost
	var createdAt, updatedAt string
	var inReplyToID sql.NullInt64
	err := s.db.QueryRow(`
		SELECT id, content, user_id, COALESCE(image_path, ''), in_reply_to_id, created_at, updated_at
		FROM microposts WHERE id = ?`, id).Scan(
		&m.ID, &m.Content, &m.UserID, &m.ImagePath, &inReplyToID, &createdAt, &updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get micropost %d: %w", id, err)
	}
	m.CreatedAt = parseTime(createdAt)
	m.UpdatedAt = parseTime(updatedAt)
	if inReplyToID.Valid {
		id := inReplyToID.Int64
		m.InReplyToID = &id
	}
	return &m, nil
}

// UpdateMicropost は既存マイクロポストを更新します。
func (s *Store) UpdateMicropost(m *model.Micropost) error {
	now := nowString()
	result, err := s.db.Exec(`
		UPDATE microposts SET content = ?, user_id = ?, image_path = ?, updated_at = ?
		WHERE id = ?`,
		m.Content, m.UserID, m.ImagePath, now, m.ID,
	)
	if err != nil {
		return fmt.Errorf("update micropost %d: %w", m.ID, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected for micropost %d: %w", m.ID, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("micropost %d not found", m.ID)
	}
	m.UpdatedAt = parseTime(now)
	return nil
}

// DeleteMicropost はマイクロポストを削除します。
func (s *Store) DeleteMicropost(id int64) error {
	_, err := s.db.Exec("DELETE FROM microposts WHERE id = ?", id)
	return err
}

// AllMicroposts はすべてのマイクロポストを返します。
func (s *Store) AllMicroposts() ([]model.Micropost, error) {
	rows, err := s.db.Query(`
		SELECT id, content, user_id, COALESCE(image_path, ''), created_at, updated_at
		FROM microposts ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("query all microposts: %w", err)
	}
	defer rows.Close()

	var microposts []model.Micropost
	for rows.Next() {
		var m model.Micropost
		var createdAt, updatedAt string
		if err := rows.Scan(&m.ID, &m.Content, &m.UserID, &m.ImagePath, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("scan micropost: %w", err)
		}
		m.CreatedAt = parseTime(createdAt)
		m.UpdatedAt = parseTime(updatedAt)
		microposts = append(microposts, m)
	}
	return microposts, rows.Err()
}

// GetUserByMicropostID はマイクロポストに紐付いたユーザーを取得します。
func (s *Store) GetUserByMicropostID(micropostID int64) (*model.User, error) {
	var user model.User
	var createdAt, updatedAt string
	err := s.db.QueryRow(`
		SELECT u.id, u.name, u.email, u.created_at, u.updated_at
		FROM users u
		INNER JOIN microposts m ON m.user_id = u.id
		WHERE m.id = ?`, micropostID).Scan(
		&user.ID, &user.Name, &user.Email, &createdAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}
	user.CreatedAt = parseTime(createdAt)
	user.UpdatedAt = parseTime(updatedAt)
	return &user, nil
}

// GetMicropostsByUserID はユーザーに紐付いたマイクロポストを作成日時の降順で取得します。
func (s *Store) GetMicropostsByUserID(userID int64) ([]model.Micropost, error) {
	rows, err := s.db.Query(`
		SELECT id, content, user_id, COALESCE(image_path, ''), created_at, updated_at
		FROM microposts
		WHERE user_id = ?
		ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var microposts []model.Micropost
	for rows.Next() {
		var m model.Micropost
		var createdAt, updatedAt string
		if err := rows.Scan(&m.ID, &m.Content, &m.UserID, &m.ImagePath, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		m.CreatedAt = parseTime(createdAt)
		m.UpdatedAt = parseTime(updatedAt)
		microposts = append(microposts, m)
	}
	return microposts, rows.Err()
}

// PaginateMicropostsByUserID はユーザーに紐付いたマイクロポストをページネーション付きで取得します。
func (s *Store) PaginateMicropostsByUserID(userID int64, page, perPage int) ([]model.Micropost, error) {
	offset := (page - 1) * perPage
	rows, err := s.db.Query(`
		SELECT id, content, user_id, COALESCE(image_path, ''), created_at, updated_at
		FROM microposts
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`, userID, perPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var microposts []model.Micropost
	for rows.Next() {
		var m model.Micropost
		var createdAt, updatedAt string
		if err := rows.Scan(&m.ID, &m.Content, &m.UserID, &m.ImagePath, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		m.CreatedAt = parseTime(createdAt)
		m.UpdatedAt = parseTime(updatedAt)
		microposts = append(microposts, m)
	}
	return microposts, rows.Err()
}

// PaginateMicropostsWithStats はプロフィールページ用。
// profileUserID の投稿を取得しつつ、viewerID によるいいね状態とリプライ情報を含む FeedItem を返す。
// viewerID = 0 の場合（未ログイン時）は IsLiked が常に false になる。
func (s *Store) PaginateMicropostsWithStats(profileUserID, viewerID int64, page, perPage int) ([]FeedItem, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * perPage
	rows, err := s.db.Query(`
		SELECT`+feedSelectCols+`
		FROM microposts m
		JOIN users u ON m.user_id = u.id`+feedJoinClauses+`
		WHERE m.user_id = ?
		ORDER BY m.created_at DESC
		LIMIT ? OFFSET ?`, viewerID, viewerID, profileUserID, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("paginate microposts with stats for user %d: %w", profileUserID, err)
	}
	defer rows.Close()

	var items []FeedItem
	for rows.Next() {
		var item FeedItem
		if err := scanFeedItem(rows, &item); err != nil {
			return nil, fmt.Errorf("scan micropost with stats: %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// CountMicropostsByUserID はユーザーに紐付いたマイクロポスト数を返します。
func (s *Store) CountMicropostsByUserID(userID int64) (int, error) {
	var count int
	err := s.db.QueryRow(
		"SELECT COUNT(*) FROM microposts WHERE user_id = ?", userID).Scan(&count)
	return count, err
}

// GetMicropostByUserIDAndID はユーザーIDとマイクロポストIDからマイクロポストを検索します。
func (s *Store) GetMicropostByUserIDAndID(userID, micropostID int64) (*model.Micropost, error) {
	var m model.Micropost
	var createdAt, updatedAt string
	err := s.db.QueryRow(`
		SELECT id, content, user_id, COALESCE(image_path, ''), created_at, updated_at
		FROM microposts
		WHERE user_id = ? AND id = ?`, userID, micropostID).Scan(
		&m.ID, &m.Content, &m.UserID, &m.ImagePath, &createdAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}
	m.CreatedAt = parseTime(createdAt)
	m.UpdatedAt = parseTime(updatedAt)
	return &m, nil
}

// CountAllMicroposts はマイクロポストの総数を返します。
func (s *Store) CountAllMicroposts() (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM microposts").Scan(&count)
	return count, err
}

// GetMicropostAsFeedItem は1件の投稿をリプライ情報・いいね情報付き FeedItem で返す。
func (s *Store) GetMicropostAsFeedItem(id, viewerID int64) (*FeedItem, error) {
	var item FeedItem
	row := s.db.QueryRow(`
		SELECT`+feedSelectCols+`
		FROM microposts m
		JOIN users u ON m.user_id = u.id`+feedJoinClauses+`
		WHERE m.id = ?`, viewerID, viewerID, id)
	if err := scanFeedItem(row, &item); err != nil {
		return nil, fmt.Errorf("get micropost feed item %d: %w", id, err)
	}
	return &item, nil
}

// GetReplies は micropostID へのリプライ一覧を作成日時昇順で返す。
func (s *Store) GetReplies(micropostID, viewerID int64) ([]FeedItem, error) {
	rows, err := s.db.Query(`
		SELECT`+feedSelectCols+`
		FROM microposts m
		JOIN users u ON m.user_id = u.id`+feedJoinClauses+`
		WHERE m.in_reply_to_id = ?
		ORDER BY m.created_at ASC`, viewerID, viewerID, micropostID)
	if err != nil {
		return nil, fmt.Errorf("get replies for %d: %w", micropostID, err)
	}
	defer rows.Close()

	var items []FeedItem
	for rows.Next() {
		var item FeedItem
		if err := scanFeedItem(rows, &item); err != nil {
			return nil, fmt.Errorf("scan reply: %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// CountReplies は micropostID へのリプライ数を返す。
func (s *Store) CountReplies(micropostID int64) (int, error) {
	var count int
	err := s.db.QueryRow(
		"SELECT COUNT(*) FROM microposts WHERE in_reply_to_id = ?", micropostID).Scan(&count)
	return count, err
}

