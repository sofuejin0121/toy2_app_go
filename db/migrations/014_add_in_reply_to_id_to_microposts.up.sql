ALTER TABLE microposts
    ADD COLUMN in_reply_to_id INTEGER DEFAULT NULL
    REFERENCES microposts(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS index_microposts_on_in_reply_to_id
    ON microposts (in_reply_to_id);