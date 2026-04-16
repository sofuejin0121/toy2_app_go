package store

import (
	"fmt"
	"strings"
	"time"

	"github.com/sofuejin0121/toy_app_go/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type userScanner interface {
	Scan(dest ...any) error
}

func scanUser(scanner userScanner) (model.User, error) {
	var user model.User
	var createdAt string
	var updatedAt string
	if err := scanner.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordDigest, &user.RememberDigest, &user.Admin, &createdAt, &updatedAt); err != nil {
		return model.User{}, err
	}
	user.CreatedAt = parseTime(createdAt)
	user.UpdatedAt = parseTime(updatedAt)
	return user, nil
}

func scanPaginatedUser(scanner userScanner) (model.User, error) {
	var user model.User
	var createdAt string
	var updatedAt string
	// TODO: Adminを後で追加
	if err := scanner.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordDigest, &user.RememberDigest, &user.Admin, &createdAt, &updatedAt); err != nil {
		return model.User{}, err
	}
	user.CreatedAt = parseTime(createdAt)
	user.UpdatedAt = parseTime(updatedAt)
	return user, nil
}
func (s *Store) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	var createdAt string
	var updatedAt string
	var activatedAtStr *string
	var resetSentAtStr *string
	err := s.db.QueryRow(`
        SELECT id, name, email, COALESCE(bio, ''), password_digest, COALESCE(remember_digest, ''),
               admin, COALESCE(activation_digest, ''), activated, activated_at,
               COALESCE(reset_digest, ''), reset_sent_at,
               created_at, updated_at
        FROM users WHERE email = ?`, strings.ToLower(email)).Scan(
		&user.ID, &user.Name, &user.Email, &user.Bio, &user.PasswordDigest,
		&user.RememberDigest, &user.Admin, &user.ActivationDigest,
		&user.Activated, &activatedAtStr, &user.ResetDigest, &resetSentAtStr, &createdAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}
	if activatedAtStr != nil {
		t := parseTime(*activatedAtStr)
		user.ActivatedAt = &t
	}
	if resetSentAtStr != nil {
		t := parseTime(*resetSentAtStr)
		user.ResetSentAt = &t
	}
	user.CreatedAt = parseTime(createdAt)
	user.UpdatedAt = parseTime(updatedAt)
	return &user, nil
}

// GetAllUsers はすべてのユーザーを返します
func (s *Store) GetAllUsers() ([]model.User, error) {
	rows, err := s.db.Query(`
        SELECT id, name, email, password_digest, COALESCE(remember_digest, ''),
               admin, created_at, updated_at
        FROM users ORDER BY created_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		var createdAt, updatedAt string
		if err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.PasswordDigest,
			&user.RememberDigest, &user.Admin, &createdAt, &updatedAt,
		); err != nil {
			return nil, err
		}
		user.CreatedAt = parseTime(createdAt)
		user.UpdatedAt = parseTime(updatedAt)
		users = append(users, user)
	}
	return users, rows.Err()
}

// 既存ハンドラー互換のAllUsersラッパーを追加する
func (s *Store) AllUsers() ([]model.User, error) {
	return s.GetAllUsers()
}

// GetUser は指定したIDのユーザーを返します
func (s *Store) GetUser(id int64) (*model.User, error) {
	var user model.User
	var createdAt string
	var updatedAt string
	var activatedAtStr *string
	var resetSentAtStr *string
	err := s.db.QueryRow(`
        SELECT id, name, email, COALESCE(bio, ''), password_digest, COALESCE(remember_digest, ''),
               admin, COALESCE(activation_digest, ''), activated, activated_at,
               COALESCE(reset_digest, ''), reset_sent_at,
               created_at, updated_at
        FROM users WHERE id = ?`, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.Bio, &user.PasswordDigest,
		&user.RememberDigest, &user.Admin, &user.ActivationDigest,
		&user.Activated, &activatedAtStr, &user.ResetDigest, &resetSentAtStr, &createdAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}
	if activatedAtStr != nil {
		t := parseTime(*activatedAtStr)
		user.ActivatedAt = &t
	}
	if resetSentAtStr != nil {
		t := parseTime(*resetSentAtStr)
		user.ResetSentAt = &t
	}
	user.CreatedAt = parseTime(createdAt)
	user.UpdatedAt = parseTime(updatedAt)
	return &user, nil
}


// FindUserByEmail はメールアドレスに一致するユーザーを返します。
func (s *Store) FindUserByEmail(email string) (*model.User, error) {
	return s.GetUserByEmail(email)
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
	var activatedAt *string
	if u.ActivatedAt != nil {
		s := u.ActivatedAt.Format(time.RFC3339)
		activatedAt = &s
	}
	result, err := s.db.Exec(
		"INSERT INTO users (name, email, password_digest, admin, activation_digest, activated, activated_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		u.Name, u.Email, u.PasswordDigest, u.Admin, u.ActivationDigest, u.Activated, activatedAt, now, now,
	)

	// ④ INSERT が失敗した場合のエラーハンドリング
	if err != nil {
		// メールアドレスの重複チェック
		// DBのUNIQUE制約によって同じメールがすでに存在する場合
		// SQLiteは"UNIQUE constraint failed"というエラーを返す
		// "Email has already been taken" に変換する。
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return fmt.Errorf("このメールアドレスはすでに使用されています")
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
func (s *Store) UpdateUser(u *model.User) error {
	now := time.Now()
	_, err := s.db.Exec(`
        UPDATE users
        SET name = ?, email = ?, bio = ?, password_digest = ?, remember_digest = ?,
            admin = ?, activation_digest = ?, activated = ?, activated_at = ?,
            updated_at = ?
        WHERE id = ?`,
		u.Name, strings.ToLower(u.Email), u.Bio, u.PasswordDigest, u.RememberDigest,
		u.Admin, u.ActivationDigest, u.Activated, u.ActivatedAt, now, u.ID,
	)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	u.UpdatedAt = now
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

// UpdatePassword はユーザーのパスワードを更新します。
func (s *Store) UpdatePassword(userID int64, password string) error {
	digest, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	_, err = s.db.Exec(
		"UPDATE users SET password_digest = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		string(digest), userID,
	)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	return nil
}

// UpdateRememberDigest はremember_digestを更新します。
func (s *Store) UpdateRememberDigest(userID int64, digest string) error {
	now := nowString()
	_, err := s.db.Exec(
		"UPDATE users SET remember_digest = ?, updated_at = ? WHERE id = ?",
		digest,
		now,
		userID,
	)
	return err
}

// CountUsers はユーザーの総数を返します。
func (s *Store) CountUsers() (int, error) {
	row := s.db.QueryRow("SELECT COUNT(*) FROM users")
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("count users: %w", err)
	}
	return count, nil
}

// CountActivatedUsers は有効化済みユーザーの総数を返します。
func (s *Store) CountActivatedUsers() (int, error) {
	row := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE activated = TRUE")
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("count activated users: %w", err)
	}
	return count, nil
}

// PaginateUsers はユーザーをページングして返します。
func (s *Store) PaginateUsers(page, perPage int) ([]model.User, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * perPage
	rows, err := s.db.Query(`
        SELECT id, name, email, admin, activated, created_at, updated_at
        FROM users ORDER BY created_at LIMIT ? OFFSET ?`, perPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		var createdAt, updatedAt string
		if err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.Admin,
			&user.Activated, &createdAt, &updatedAt,
		); err != nil {
			return nil, err
		}
		user.CreatedAt = parseTime(createdAt)
		user.UpdatedAt = parseTime(updatedAt)
		users = append(users, user)
	}
	return users, rows.Err()
}

// Authenticate はメールアドレスとパスワードでユーザーを認証します
// 認証に成功した場合はユーザーを返し、失敗した場合はエラーを返す
func (s *Store) Authenticate(email, password string) (*model.User, error) {
	user, err := s.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}
	if !user.Authenticate(password) {
		return nil, fmt.Errorf("authenticate: invalid password")
	}
	return user, nil
}

// UpdateActivation はユーザーの有効化ステータスを更新する
func (s *Store) UpdateActivation(userID int64, activated bool, activatedAt time.Time) error {
	_, err := s.db.Exec(
		"Update users SET activated = ?, activated_at = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		activated, activatedAt, userID,
	)
	return err
}

func (s *Store) GetActivatedUsers(page int) ([]model.User, error) {
	if page < 1 {
		page = 1
	}
	const perPage = 30
	offset := (page - 1) * perPage
	rows, err := s.db.Query(`
	SELECT id, name, email, admin, activated, created_at, updated_at
	FROM users
	WHERE activated = TRUE
	ORDER BY created_at LIMIT ? OFFSET ?`, perPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		var createdAt, updatedAt string
		if err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.Admin,
			&user.Activated, &createdAt, &updatedAt,
		); err != nil {
			return nil, err
		}
		user.CreatedAt = parseTime(createdAt)
		user.UpdatedAt = parseTime(updatedAt)
		users = append(users, user)
	}
	return users, rows.Err()
}

// SearchActivatedUsers は名前かメールにqueryを含む有効ユーザーをページング付きで返す。
// SQL の LIKE 句で前後にワイルドカード(%)を付けた部分一致検索を行う。
func (s *Store) SearchActivatedUsers(query string, page, perPage int) ([]model.User, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * perPage
	like := "%" + query + "%"
	rows, err := s.db.Query(`
		SELECT id, name, email, admin, activated, created_at, updated_at
		FROM users
		WHERE activated = TRUE AND (name LIKE ? OR email LIKE ?)
		ORDER BY created_at
		LIMIT ? OFFSET ?`, like, like, perPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		var createdAt, updatedAt string
		if err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.Admin,
			&user.Activated, &createdAt, &updatedAt,
		); err != nil {
			return nil, err
		}
		user.CreatedAt = parseTime(createdAt)
		user.UpdatedAt = parseTime(updatedAt)
		users = append(users, user)
	}
	return users, rows.Err()
}

// CountSearchActivatedUsers は検索条件に一致する有効ユーザー数を返す。
func (s *Store) CountSearchActivatedUsers(query string) (int, error) {
	like := "%" + query + "%"
	var count int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM users
		WHERE activated = TRUE AND (name LIKE ? OR email LIKE ?)`, like, like).Scan(&count)
	return count, err
}

// UpdateResetDigest はリセットダイジェストと送信時刻を更新する
func (s *Store) UpdateResetDigest(userID int64, digest string, sentAt time.Time) error {
	_, err := s.db.Exec(
		"UPDATE users SET reset_digest = ?, reset_sent_at = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		digest, sentAt, userID,
	)
	return err
}

// ClearResetDigest は使用済みのリセットダイジェストを無効化する
func (s *Store) ClearResetDigest(userID int64) error {
	_, err := s.db.Exec(
		"UPDATE users SET reset_digest = NULL, reset_sent_at = NULL, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		userID,
	)
	return err
}
