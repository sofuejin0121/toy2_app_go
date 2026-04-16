package store

import (
	"database/sql"
	"errors"

	"github.com/sofuejin0121/toy_app_go/internal/model"
)

// GetOrCreateUserPreference はユーザーの通知設定を返す。
// レコードがなければデフォルト値(両方true)で新規作成する。
func (s *Store) GetOrCreateUserPreference(userID int64) (*model.UserPreference, error) {
	pref := &model.UserPreference{}
	err := s.db.QueryRow(`
		SELECT id, user_id, email_on_follow, email_on_like
		FROM user_preferences
		WHERE user_id = ?`, userID).
		Scan(&pref.ID, &pref.UserID, &pref.EmailOnFollow, &pref.EmailOnLike)

	if errors.Is(err, sql.ErrNoRows) {
		// レコードがなければデフォルト値(true/true)で作成
		res, err := s.db.Exec(`
			INSERT INTO user_preferences (user_id, email_on_follow, email_on_like)
			VALUES (?, TRUE, TRUE)`, userID)
		if err != nil {
			return nil, err
		}
		id, _ := res.LastInsertId()
		return &model.UserPreference{
			ID:            id,
			UserID:        userID,
			EmailOnFollow: true,
			EmailOnLike:   true,
		}, nil
	}
	if err != nil {
		return nil, err
	}
	return pref, nil
}

// UpdateUserPreference はユーザーの通知設定を更新する。
func (s *Store) UpdateUserPreference(userID int64, emailOnFollow, emailOnLike bool) error {
	_, err := s.db.Exec(`
		UPDATE user_preferences
		SET email_on_follow = ?, email_on_like = ?, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = ?`, emailOnFollow, emailOnLike, userID)
	return err
}
