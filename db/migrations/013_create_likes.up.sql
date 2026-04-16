CREATE TABLE IF NOT EXISTS likes (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id    INTEGER NOT NULL,
    micropost_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id)      REFERENCES users(id)      ON DELETE CASCADE,
    FOREIGN KEY (micropost_id) REFERENCES microposts(id) ON DELETE CASCADE
);

-- 同一ユーザーが同一マイクロポストに複数回いいねするのを防ぐ複合ユニーク制約
CREATE UNIQUE INDEX IF NOT EXISTS index_likes_on_user_id_and_micropost_id
    ON likes (user_id, micropost_id);

-- マイクロポスト別のいいね数を高速に集計するためのインデックス
CREATE INDEX IF NOT EXISTS index_likes_on_micropost_id
    ON likes (micropost_id);
