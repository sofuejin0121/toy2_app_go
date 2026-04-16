CREATE TABLE relationships (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    follower_id INTEGER NOT NULL,
    followed_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (followed_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX idx_relationships_follower_id ON relationships(follower_id);
CREATE INDEX idx_relationships_followed_id ON relationships(followed_id);
CREATE UNIQUE INDEX idx_relationships_follower_followed
    ON relationships(follower_id, followed_id);
