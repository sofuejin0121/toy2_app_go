package store

import (
	"database/sql"
	"fmt"

	"github.com/sofuejin0121/toy_app_go/internal/model"
)

// FeedItem はフィードの1項目。
// リプライの場合は ParentMicropost と ParentUser に元投稿の情報が入る。
type FeedItem struct {
	Micropost       model.Micropost
	User            model.User
	LikeCount       int
	IsLiked         bool
	IsBookmarked    bool
	ParentMicropost *model.Micropost // nil = リプライでない
	ParentUser      *model.User      // nil = リプライでない
}

// scanFeedItem は SELECT 結果の行を FeedItem にスキャンする共通ヘルパー。
// SELECT カラム順: m.id, m.content, m.user_id, image_path, m.in_reply_to_id,
//   m.created_at, m.updated_at, u.id, u.name, u.email,
//   like_count, is_liked,
//   parent.id, parent.content, parent.user_id, pu.name, pu.email
func scanFeedItem(rows interface {
	Scan(dest ...any) error
}, item *FeedItem) error {
	var (
		createdAt, updatedAt               string
		inReplyToID                         sql.NullInt64
		parentID, parentUserID              sql.NullInt64
		parentContent                       sql.NullString
		parentUserName, parentUserEmail     sql.NullString
	)
	if err := rows.Scan(
		&item.Micropost.ID, &item.Micropost.Content,
		&item.Micropost.UserID, &item.Micropost.ImagePath,
		&inReplyToID,
		&createdAt, &updatedAt,
		&item.User.ID, &item.User.Name, &item.User.Email,
		&item.LikeCount, &item.IsLiked, &item.IsBookmarked,
		&parentID, &parentContent, &parentUserID,
		&parentUserName, &parentUserEmail,
	); err != nil {
		return err
	}
	item.Micropost.CreatedAt = parseTime(createdAt)
	item.Micropost.UpdatedAt = parseTime(updatedAt)
	if inReplyToID.Valid {
		id := inReplyToID.Int64
		item.Micropost.InReplyToID = &id
	}
	if parentID.Valid {
		item.ParentMicropost = &model.Micropost{
			ID:      parentID.Int64,
			Content: parentContent.String,
			UserID:  parentUserID.Int64,
		}
		item.ParentUser = &model.User{
			ID:    parentUserID.Int64,
			Name:  parentUserName.String,
			Email: parentUserEmail.String,
		}
	}
	return nil
}

// Follow はfollowerIDのユーザーがfollowedIDのユーザーをフォローする。
func (s *Store) Follow(followerID, followedID int64) error {
	if followerID == followedID {
		return nil
	}
	_, err := s.db.Exec(
		`INSERT INTO relationships (follower_id, followed_id) VALUES (?, ?)`,
		followerID, followedID,
	)
	if err != nil {
		return fmt.Errorf("follow %d -> %d: %w", followerID, followedID, err)
	}
	s.CreateNotification(followedID, followerID, "follow", nil)
	return nil
}

// Unfollow はfollowerIDのユーザーがfollowedIDのユーザーをフォロー解除する。
func (s *Store) Unfollow(followerID, followedID int64) error {
	_, err := s.db.Exec(
		`DELETE FROM relationships WHERE follower_id = ? AND followed_id = ?`,
		followerID, followedID,
	)
	if err != nil {
		return fmt.Errorf("unfollow %d -> %d: %w", followerID, followedID, err)
	}
	return nil
}

// IsFollowing はfollowerIDのユーザーがfollowedIDのユーザーをフォローしていればtrueを返す。
func (s *Store) IsFollowing(followerID, followedID int64) (bool, error) {
	var exists bool
	err := s.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM relationships
		 WHERE follower_id = ? AND followed_id = ?)`,
		followerID, followedID,
	).Scan(&exists)
	return exists, err
}

// PaginateFollowing はuserIDがフォローしているユーザーをページネーション付きで返す。
func (s *Store) PaginateFollowing(userID int64, page, perPage int) ([]model.User, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * perPage
	rows, err := s.db.Query(`
		SELECT u.id, u.name, u.email, u.created_at, u.updated_at
		FROM users u
		INNER JOIN relationships r ON r.followed_id = u.id
		WHERE r.follower_id = ?
		ORDER BY r.created_at DESC
		LIMIT ? OFFSET ?`, userID, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("paginate following for user %d: %w", userID, err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		var createdAt, updatedAt string
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("scan following user: %w", err)
		}
		u.CreatedAt = parseTime(createdAt)
		u.UpdatedAt = parseTime(updatedAt)
		users = append(users, u)
	}
	return users, rows.Err()
}

// PaginateFollowers はuserIDのフォロワーをページネーション付きで返す。
func (s *Store) PaginateFollowers(userID int64, page, perPage int) ([]model.User, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * perPage
	rows, err := s.db.Query(`
		SELECT u.id, u.name, u.email, u.created_at, u.updated_at
		FROM users u
		INNER JOIN relationships r ON r.follower_id = u.id
		WHERE r.followed_id = ?
		ORDER BY r.created_at DESC
		LIMIT ? OFFSET ?`, userID, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("paginate followers for user %d: %w", userID, err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		var createdAt, updatedAt string
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("scan follower user: %w", err)
		}
		u.CreatedAt = parseTime(createdAt)
		u.UpdatedAt = parseTime(updatedAt)
		users = append(users, u)
	}
	return users, rows.Err()
}

// CountFollowing はuserIDがフォローしているユーザー数を返す。
func (s *Store) CountFollowing(userID int64) (int, error) {
	var count int
	err := s.db.QueryRow(
		`SELECT COUNT(*) FROM relationships WHERE follower_id = ?`, userID,
	).Scan(&count)
	return count, err
}

// CountFollowers はuserIDのフォロワー数を返す。
func (s *Store) CountFollowers(userID int64) (int, error) {
	var count int
	err := s.db.QueryRow(
		`SELECT COUNT(*) FROM relationships WHERE followed_id = ?`, userID,
	).Scan(&count)
	return count, err
}

// GetRelationship はIDでリレーションシップを取得する。
func (s *Store) GetRelationship(id int64) (*model.Relationship, error) {
	var r model.Relationship
	var createdAt, updatedAt string
	err := s.db.QueryRow(
		`SELECT id, follower_id, followed_id, created_at, updated_at
		 FROM relationships WHERE id = ?`, id,
	).Scan(&r.ID, &r.FollowerID, &r.FollowedID, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	r.CreatedAt = parseTime(createdAt)
	r.UpdatedAt = parseTime(updatedAt)
	return &r, nil
}

// GetRelationshipByUsers はfollowerIDとfollowedIDからリレーションシップを取得する。
func (s *Store) GetRelationshipByUsers(followerID, followedID int64) (*model.Relationship, error) {
	var r model.Relationship
	var createdAt, updatedAt string
	err := s.db.QueryRow(
		`SELECT id, follower_id, followed_id, created_at, updated_at
		 FROM relationships
		 WHERE follower_id = ? AND followed_id = ?`,
		followerID, followedID,
	).Scan(&r.ID, &r.FollowerID, &r.FollowedID, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	r.CreatedAt = parseTime(createdAt)
	r.UpdatedAt = parseTime(updatedAt)
	return &r, nil
}

// CountRelationships はリレーションシップの総数を返す。
func (s *Store) CountRelationships() (int, error) {
	var count int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM relationships`).Scan(&count)
	return count, err
}

// feedSelectCols は Feed/FeedByJoin/PaginateMicropostsWithStats で共通のSELECTカラム定義。
// 自己参照 LEFT JOIN でリプライ元情報（parent.*）も一度に取得する。
const feedSelectCols = `
		m.id, m.content, m.user_id, COALESCE(m.image_path, ''),
		m.in_reply_to_id,
		m.created_at, m.updated_at,
		u.id, u.name, u.email,
		(SELECT COUNT(*) FROM likes WHERE micropost_id = m.id) AS like_count,
		EXISTS(SELECT 1 FROM likes WHERE micropost_id = m.id AND user_id = ?) AS is_liked,
		EXISTS(SELECT 1 FROM bookmarks WHERE micropost_id = m.id AND user_id = ?) AS is_bookmarked,
		parent.id, parent.content, parent.user_id,
		pu.name, pu.email`

// feedJoinClauses は parent への自己参照 LEFT JOIN クロス節。
const feedJoinClauses = `
		LEFT JOIN microposts parent ON parent.id = m.in_reply_to_id
		LEFT JOIN users pu ON pu.id = parent.user_id`

// Feed はuserIDのステータスフィードを返す（eager loading版）。
// サブセレクトでフォロー中ユーザーと自分自身の投稿を取得し、JOINでN+1クエリ問題を解決する。
// userIDは閲覧者IDも兼ねるため、IsLikedはホームページ（自分のフィード）に適切に反映される。
func (s *Store) Feed(userID int64, page, perPage int) ([]FeedItem, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * perPage
	rows, err := s.db.Query(`
		SELECT`+feedSelectCols+`
		FROM microposts m
		JOIN users u ON m.user_id = u.id`+feedJoinClauses+`
		WHERE m.user_id IN (
		    SELECT followed_id FROM relationships WHERE follower_id = ?
		) OR m.user_id = ?
		ORDER BY m.created_at DESC
		LIMIT ? OFFSET ?`, userID, userID, userID, userID, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("feed for user %d: %w", userID, err)
	}
	defer rows.Close()

	var items []FeedItem
	for rows.Next() {
		var item FeedItem
		if err := scanFeedItem(rows, &item); err != nil {
			return nil, fmt.Errorf("scan feed item: %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// FeedByJoin はLEFT OUTER JOINを使ったフィード取得（SELECT DISTINCTで重複解消済み）。
func (s *Store) FeedByJoin(userID int64, page, perPage int) ([]FeedItem, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * perPage
	rows, err := s.db.Query(`
		SELECT DISTINCT`+feedSelectCols+`
		FROM microposts m
		JOIN users u ON m.user_id = u.id`+feedJoinClauses+`
		LEFT OUTER JOIN relationships r
		    ON r.followed_id = m.user_id AND r.follower_id = ?
		WHERE r.follower_id = ? OR m.user_id = ?
		ORDER BY m.created_at DESC
		LIMIT ? OFFSET ?`, userID, userID, userID, userID, userID, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("feed (join) for user %d: %w", userID, err)
	}
	defer rows.Close()

	var items []FeedItem
	for rows.Next() {
		var item FeedItem
		if err := scanFeedItem(rows, &item); err != nil {
			return nil, fmt.Errorf("scan feed item (join): %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
