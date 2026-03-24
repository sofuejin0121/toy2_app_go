package store

import (
	"database/sql"
	"fmt"

	"github.com/sofuejin0121/toy_app_go/internal/model"
)

type userScanner interface {
	Scan(dest ...any) error
}

func scanUser(scanner userScanner) (model.User, error) {
	var user model.User
	var createdAt string
	var updatedAt string
	if err := scanner.Scan(&user.ID, &user.Name, &user.Email, &createdAt, &updatedAt); err != nil {
		return model.User{}, err
	}
	user.CreatedAt = parseTime(createdAt)
	user.UpdatedAt = parseTime(updatedAt)
	return user, nil
}

// AllUsers はすべてのユーザーを返します
func (s *Store) AllUsers() ([]model.User, error) {
	rows, err := s.db.Query(
		"SELECT id, name, email, created_at, updated_at FROM users ORDER BY id",
	)
	if err != nil {
		return nil, fmt.Errorf("query all users: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}
	return users, nil
}

// GetUser は指定したIDのユーザーを返します
func (s *Store) GetUser(id int64) (*model.User, error) {
	row := s.db.QueryRow(
		"SELECT id, name, email, created_at, updated_at FROM users WHERE id = ?",
		id,
	)
	user, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user %d not found", id)
		}
		return nil, fmt.Errorf("get user %d: %w", id, err)
	}
	return &user, nil
}

// CreateUser は新しいユーザーを作成
func (s *Store) CreateUser(user *model.User) error {
	now := nowString()
	result, err := s.db.Exec(
		"INSERT INTO users (name, email, created_at, updated_at) VALUES (?,?,?,?)",
		user.Name,
		user.Email,
		now,
		now,
	)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("fetch created user id: %w", err)
	}
	user.ID = id
	user.CreatedAt = parseTime(now)
	user.UpdatedAt = parseTime(now)
	return nil
}

// UpdateUser は既存ユーザーを更新します。
func (s *Store) UpdateUser(user *model.User) error {
	now := nowString()
	result, err := s.db.Exec(
		"UPDATE users SET name = ?, email = ?, updated_at = ? WHERE id = ?",
		user.Name,
		user.Email,
		now,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("update user %d: %w", user.ID, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected for user %d: %w", user.ID, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user %d not found", user.ID)
	}
	user.UpdatedAt = parseTime(now)
	return nil
}

// DeleteUser はユーザーを削除します。
func (s *Store) DeleteUser(id int64) error {
	result, err := s.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete user %d: %w", id, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected for user %d: %w", id, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user %d not found", id)
	}
	return nil
}

// GetMicropostsByUserID は指定ユーザーのマイクロポスト一覧を返します。
func (s *Store) GetMicropostsByUserID(userID int64) ([]model.Micropost, error) {
	rows, err := s.db.Query(
		"SELECT id, content, user_id, created_at, updated_at FROM microposts WHERE user_id = ? ORDER BY id",
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("query microposts for user %d: %w", userID, err)
	}
	defer rows.Close()

	var microposts []model.Micropost
	for rows.Next() {
		micropost, err := scanMicropost(rows)
		if err != nil {
			return nil, fmt.Errorf("scan micropost for user %d: %w", userID, err)
		}
		microposts = append(microposts, micropost)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate microposts for user %d: %w", userID, err)
	}
	return microposts, nil
}