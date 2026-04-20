package store

import (
	"database/sql"
	"strings"
	"time"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

// Store はデータベース接続を管理します。
type Store struct {
	db *sql.DB
}

// New はデータベースに接続し、Storeを返します。
// dbPath が "libsql://" で始まる場合は Turso(libsql) ドライバを使用します。
// それ以外はローカルの SQLite ファイルに接続します。
func New(dbPath string) (*Store, error) {
	driver := "sqlite"
	if strings.HasPrefix(dbPath, "libsql://") {
		driver = "libsql"
	}

	db, err := sql.Open(driver, dbPath)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		db.Close()
		return nil, err
	}
	// ローカル SQLite のみ。短時間の同時アクセスで SQLITE_BUSY → 500 になり得るため待機を入れる。
	// libsql / Turso では PRAGMA が無効なためスキップする。
	if driver == "sqlite" {
		if _, err := db.Exec("PRAGMA busy_timeout = 8000"); err != nil {
			db.Close()
			return nil, err
		}
		_, _ = db.Exec("PRAGMA journal_mode = WAL")
	}
	return &Store{db: db}, nil
}

// Close はデータベース接続を閉じます。
func (s *Store) Close() error {
	return s.db.Close()
}

func nowString() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func parseTime(value string) time.Time {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
	}
	for _, format := range formats {
		if t, err := time.Parse(format, value); err == nil {
			return t
		}
	}
	return time.Time{}
}

// DB はデータベース接続を返します。
func (s *Store) DB() *sql.DB {
	return s.db
}


