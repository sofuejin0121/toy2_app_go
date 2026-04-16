-- db/migrations/008_add_reset_to_users.down.sql
ALTER TABLE users DROP COLUMN IF EXISTS reset_digest;
ALTER TABLE users DROP COLUMN IF EXISTS reset_sent_at;
