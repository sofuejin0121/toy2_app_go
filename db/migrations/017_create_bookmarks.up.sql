CREATE TABLE IF NOT EXISTS bookmarks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    micropost_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (micropost_id) REFERENCES microposts(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS index_bookmarks_on_user_id_and_micropost_id ON bookmarks (user_id, micropost_id);