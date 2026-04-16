package store

import (
	"database/sql"

	"github.com/sofuejin0121/toy_app_go/internal/model"
)

// NotificationItem は通知1件 + 関連データをまとめた構造体
// DBからJOINで取得したデータをGoで扱いやすい形にまとめる
type NotificationItem struct {
	Notification model.Notification // 通知本体
	Actor        model.User         // アクションしたユーザー(JOINで取得)
	Target       *model.Micropost   // いいねされた投稿(フォロー通知はnil)
}

// 通知作成(Like時・Follow時に呼ぶ)
func (s *Store) CreateNotification(userID, actorID int64, actionType string, targetID *int64) error {
	// targetIDは*int64(ポインタ)なので、nilの場合がある
	// SQLの?にnilを渡すとNULL として挿入される
	_, err := s.db.Exec(`
	  INSERT INTO notifications (user_id, actor_id, action_type, target_id)
	  VALUES (?, ?, ?, ?)`, userID, actorID, actionType, targetID)
	return err
}

// 通知一覧取得(actorとtargetとJOINで一度に取得)
func (s *Store) GetNotifications(userID int64) ([]NotificationItem, error) {
	rows, err := s.db.Query(`
	  SELECT n.id, n.action_type, n.read, n.created_at,
	  actor.id, actor.name,
	  m.id, m.content
	  FROM notifications n
	  JOIN users actor ON actor.id = n.actor_id
	  LEFT JOIN microposts m ON m.id = n.target_id
	  WHERE n.user_id = ?
	  ORDER BY n.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []NotificationItem
	for rows.Next() {
		var item NotificationItem
		var createdAtStr string
		// m.id と m.content はNULL の可能性があるのでsql.NullInt64 / sql.NullString を使う
		var targetID sql.NullInt64
		var targetContent sql.NullString

		err := rows.Scan(
			&item.Notification.ID, &item.Notification.ActionType, &item.Notification.Read, &createdAtStr,
			&item.Actor.ID, &item.Actor.Name, &targetID, &targetContent,
		)
		if err != nil {
			return nil, err
		}
		item.Notification.CreatedAt = parseTime(createdAtStr)

		// NULL でなければ Target を作る
		if targetID.Valid {
			item.Target = &model.Micropost{
				ID: targetID.Int64,
				Content: targetContent.String,
			}
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// 既読にする
func (s *Store) MarkAllRead(userID int64) error {
	// WHERE user_id = ? で自分の通知だけ既読する
	_, err := s.db.Exec(`
	UPDATE notifications SET read = TRUE WHERE user_id = ?`, userID)
	return err
}

// 通知削除
func (s *Store) DeleteNotification(id, userID int64) error {
	_, err := s.db.Exec(`
	DELETE FROM notifications WHERE id = ? AND user_id = ?`,id, userID)
	return err
}