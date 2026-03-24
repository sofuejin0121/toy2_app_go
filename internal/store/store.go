package store

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

// Store はデータベース接続を管理します。
type Store struct {
	db *sql.DB
}

// New はSQLiteに接続し、Storeを返します。
func New(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		db.Close()
		return nil, err
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
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}
	}
	return t
}
