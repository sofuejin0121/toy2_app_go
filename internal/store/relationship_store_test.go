package store

import "testing"

func TestFeedByJoinNoDuplicates(t *testing.T) {
	s := newTestStore(t)

	if _, err := s.db.Exec(`
		INSERT INTO users (name, email, password_digest, created_at, updated_at)
		VALUES
		    ('Michael', 'michael@example.com', 'x', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
		    ('Lana', 'lana@example.com', 'x', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
		INSERT INTO microposts (content, user_id, image_path, created_at, updated_at)
		VALUES
		    ('self post', 1, '', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
		    ('followed post', 2, '', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
	`); err != nil {
		t.Fatalf("seed feed data: %v", err)
	}

	if err := s.Follow(1, 2); err != nil {
		t.Fatalf("Follow: %v", err)
	}

	items, err := s.FeedByJoin(1, 1, 50)
	if err != nil {
		t.Fatalf("FeedByJoin: %v", err)
	}

	seen := map[int64]bool{}
	for _, item := range items {
		if seen[item.Micropost.ID] {
			t.Fatalf("duplicate micropost in feed: ID=%d", item.Micropost.ID)
		}
		seen[item.Micropost.ID] = true
	}

	if len(items) != 2 {
		t.Errorf("expected 2 feed items, got %d", len(items))
	}
}
