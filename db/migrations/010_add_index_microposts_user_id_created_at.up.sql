CREATE INDEX IF NOT EXISTS idx_microposts_user_id_created_at
ON microposts (user_id, created_at);