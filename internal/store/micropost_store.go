package store

import (
	"database/sql"
	"fmt"

	"github.com/sofuejin0121/toy_app_go/internal/model"
)

type micropostScanner interface {
	Scan(dest ...any) error
}

func scanMicropost(scanner micropostScanner) (model.Micropost, error) {
	var micropost model.Micropost
	var createdAt string
	var updatedAt string
	if err := scanner.Scan(&micropost.ID, &micropost.Content, &micropost.UserID, &createdAt, &updatedAt); err != nil {
		return model.Micropost{}, err
	}
	micropost.CreatedAt = parseTime(createdAt)
	micropost.UpdatedAt = parseTime(updatedAt)
	return micropost, nil
}

// AllMicroposts はすべてのマイクロポストを返します。
func (s *Store) AllMicroposts() ([]model.Micropost, error) {
	rows, err := s.db.Query(
		"SELECT id, content, user_id, created_at, updated_at FROM microposts ORDER BY id",
	)
	if err != nil {
		return nil, fmt.Errorf("query all microposts: %w", err)
	}
	defer rows.Close()

	var microposts []model.Micropost
	for rows.Next() {
		micropost, err := scanMicropost(rows)
		if err != nil {
			return nil, fmt.Errorf("scan micropost: %w", err)
		}
		microposts = append(microposts, micropost)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate microposts: %w", err)
	}
	return microposts, nil
}

// GetMicropost は指定したIDのマイクロポストを返します。
func (s *Store) GetMicropost(id int64) (*model.Micropost, error) {
	row := s.db.QueryRow(
		"SELECT id, content, user_id, created_at, updated_at FROM microposts WHERE id = ?",
		id,
	)
	micropost, err := scanMicropost(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("micropost %d not found", id)
		}
		return nil, fmt.Errorf("get micropost %d: %w", id, err)
	}
	return &micropost, nil
}

// CreateMicropost は新しいマイクロポストを作成します。
func (s *Store) CreateMicropost(micropost *model.Micropost) error {
	now := nowString()
	result, err := s.db.Exec(
		"INSERT INTO microposts (content, user_id, created_at, updated_at) VALUES (?, ?, ?, ?)",
		micropost.Content,
		micropost.UserID,
		now,
		now,
	)
	if err != nil {
		return fmt.Errorf("create micropost: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("fetch created micropost id: %w", err)
	}
	micropost.ID = id
	micropost.CreatedAt = parseTime(now)
	micropost.UpdatedAt = parseTime(now)
	return nil
}

// UpdateMicropost は既存マイクロポストを更新します。
func (s *Store) UpdateMicropost(micropost *model.Micropost) error {
	now := nowString()
	result, err := s.db.Exec(
		"UPDATE microposts SET content = ?, user_id = ?, updated_at = ? WHERE id = ?",
		micropost.Content,
		micropost.UserID,
		now,
		micropost.ID,
	)
	if err != nil {
		return fmt.Errorf("update micropost %d: %w", micropost.ID, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected for micropost %d: %w", micropost.ID, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("micropost %d not found", micropost.ID)
	}
	micropost.UpdatedAt = parseTime(now)
	return nil
}

// DeleteMicropost はマイクロポストを削除します。
func (s *Store) DeleteMicropost(id int64) error {
	result, err := s.db.Exec("DELETE FROM microposts WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete micropost %d: %w", id, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected for micropost %d: %w", id, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("micropost %d not found", id)
	}
	return nil
}

// GetUserByMicropostID はマイクロポストの投稿者を返します。
func (s *Store) GetUserByMicropostID(micropostID int64) (*model.User, error) {
	var userID int64
	if err := s.db.QueryRow(
		"SELECT user_id FROM microposts WHERE id = ?",
		micropostID,
	).Scan(&userID); err != nil {
		return nil, fmt.Errorf("get user_id for micropost %d: %w", micropostID, err)
	}
	return s.GetUser(userID)
}
