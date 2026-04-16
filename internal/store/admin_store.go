package store

// DailySignup は日別の新規登録数を表す構造体
type DailySignup struct {
	Date  string
	Count int
}

// AdminStats は管理ダッシュボード用の集計データをまとめた構造体
type AdminStats struct {
	TotalUsers   int
	TotalPosts   int
	TodaySignups int
	DailySignups []DailySignup // 日別なので配列
}

func (s *Store) GetAdminStats() (AdminStats, error) {
	var stats AdminStats

	// ユーザー総数
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&stats.TotalUsers); err != nil {
		return stats, err
	}
	// 投稿総数
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM microposts`).Scan(&stats.TotalPosts); err != nil {
		return stats, err
	}
	// 本日の新規登録数
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM users WHERE DATE(created_at) = DATE('now')`).Scan(&stats.TodaySignups); err != nil {
		return stats, err
	}
	// ④ 過去7日間の日別新規登録数 (GROUP BY の新知識!)
	rows, err := s.db.Query(`
        SELECT DATE(created_at) AS date, COUNT(*) AS count
        FROM users
        WHERE created_at >= DATE('now', '-6 days')
        GROUP BY DATE(created_at)
        ORDER BY date ASC`)
	if err != nil {
		return stats, err
	}
	defer rows.Close()
	for rows.Next() {
		var d DailySignup
		if err := rows.Scan(&d.Date, &d.Count); err != nil {
			return stats, err
		}
		stats.DailySignups = append(stats.DailySignups, d)
	}
	return stats, rows.Err()
}
