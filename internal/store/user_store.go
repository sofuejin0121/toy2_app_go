package store

import (
	"database/sql"
	"fmt"
	"github.com/sofuejin0121/toy_app_go/internal/model"
	"strings"
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
		"SELECT id, name, email, password_digest, created_at, updated_at FROM users ORDER BY id",
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
		"SELECT id, name, email, password_digest, created_at, updated_at FROM users WHERE id = ?",
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

// FindUserByEmail はメールアドレスに一致するユーザーを返します。
func (s *Store) FindUserByEmail(email string) (*model.User, error) {
	row := s.db.QueryRow(
		"SELECT id, name, email, password_digest, created_at, updated_at FROM users WHERE email = ?",
		strings.ToLower(email),
	)
	user, err := scanUser(row)
	if err != nil {
		return nil, fmt.Errorf("find user by email %q: %w", email, err)
	}
	return &user, nil
}
// CreateUser は新しいユーザーをデータベースに作成する
// UniQue制約違反の場合はユーザーフレンドリーなエラーを返す
func (s *Store) CreateUser(u *model.User) error {
	// ① 現在時刻を文字列として取得する
	// データベースのcreated_atとupdated_atカラムに保存するための文字列形式の時刻
	// Go の time.Time 型をそのまま SQLite に渡すと型変換の問題が起きやすいため、
	// あらかじめ文字列に変換しておく。
	now := nowString()

	// ② メールアドレスを小文字に統一
	// "User@Example.com" と "user@example.com" は同じアドレスだが、
	// 大文字・小文字が違うと別のユーザーとして登録されてしまう
	// 保存前に必ず小文字化することで、重複登録を防ぐ
	u.Email = strings.ToLower(u.Email)

	// ③ SQL のINSERT 文を実行してDBにレコードを作成
	// ?はプレースホルダー(SQL インジェクション対策)
	// 直接文字列を埋め込むと悪意のある入力でDBを操作される危険があるため
	// 必ずプレースホルダーを使って値を渡す
	result, err := s.db.Exec(
		"INSERT INTO users (name, email, password_digest, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		u.Name, u.Email, u.PasswordDigest, now, now,
	)

	// ④ INSERT が失敗した場合のエラーハンドリング
	if err != nil {
		// メールアドレスの重複チェック
		// DBのUNIQUE制約によって同じメールがすでに存在する場合
		// SQLiteは"UNIQUE constraint failed"というエラーを返す
		// "Email has already been taken" に変換する。
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return fmt.Errorf("Email has already been taken")
		}
		return fmt.Errorf("create user: %w", err)
	}

	// ⑤採番されたIDを取得する
	// INSERT が成功するとDBが自動的にIDを採番する
	// LastInsertId() メソッドで採番されたIDを取得する
	// これにより、呼び出し元は作成後すぐに u.ID を参照できる。
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("fetch created user id: %w", err)
	}

	// ⑥ 構造体のフィールドをDBに保存した値を更新する
	// INSERT 後Go側のオブジェクトuはIDや時刻が空のまま
	// DB に保存された値（ID・作成日時・更新日時）をオブジェクトに書き戻すことで、
	// 呼び出し元が改めて DB を再取得しなくても最新の状態を使えるようにする。
	u.ID = id
	u.CreatedAt = parseTime(now)
	u.UpdatedAt = parseTime(now)

	// ⑦ 成功を通知する
	// エラーがなければ、作成したユーザーの情報を呼び出し元に返す。
	// このメソッドは、ユーザー作成後に ID を参照したい場合に使う。
	return nil
}

// UpdateUser は既存ユーザーを更新します。
func (s *Store) UpdateUser(user *model.User) error {
	now := nowString()
	user.Email = strings.ToLower(user.Email) // 保存前に小文字化
	result, err := s.db.Exec(
		"UPDATE users SET name = ?, email = ?, password_digest = ?, updated_at = ? WHERE id = ?",
		user.Name,
		user.Email,
		user.PasswordDigest,
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

