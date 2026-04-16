package store

import "testing"

// NewTestStore creates an in-memory Store with the full application schema.
// Intended for use in tests across packages.
func NewTestStore(t testing.TB) *Store {
	t.Helper()
	s, err := New(":memory:")
	if err != nil {
		t.Fatalf("new test store: %v", err)
	}
	t.Cleanup(func() { s.Close() })

	queries := []string{
		`CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			bio TEXT NOT NULL DEFAULT '',
			password_digest TEXT NOT NULL DEFAULT '',
			remember_digest TEXT,
			admin BOOLEAN NOT NULL DEFAULT FALSE,
			activation_digest TEXT,
			activated BOOLEAN NOT NULL DEFAULT FALSE,
			activated_at TEXT,
			reset_digest TEXT,
			reset_sent_at TEXT,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS index_users_on_email ON users (email)`,
		`CREATE TABLE microposts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT NOT NULL,
			user_id INTEGER NOT NULL,
			image_path TEXT DEFAULT '',
			in_reply_to_id INTEGER DEFAULT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE relationships (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			follower_id INTEGER NOT NULL,
			followed_id INTEGER NOT NULL,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE likes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id    INTEGER NOT NULL,
			micropost_id INTEGER NOT NULL,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS index_likes_on_user_id_and_micropost_id ON likes (user_id, micropost_id)`,
	}
	for _, q := range queries {
		if _, err := s.db.Exec(q); err != nil {
			t.Fatalf("init test schema: %v", err)
		}
	}
	return s
}
