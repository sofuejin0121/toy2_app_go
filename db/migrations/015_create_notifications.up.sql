CREATE TABLE notifications (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id     INTEGER NOT NULL,
    actor_id    INTEGER NOT NULL,
    action_type TEXT NOT NULL,
    target_id   INTEGER,
    read        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id)    REFERENCES users(id)      ON DELETE CASCADE,
    FOREIGN KEY (actor_id)   REFERENCES users(id)      ON DELETE CASCADE,
    FOREIGN KEY (target_id)  REFERENCES microposts(id) ON DELETE CASCADE
);
CREATE INDEX idx_notifications_user_id ON notifications(user_id);