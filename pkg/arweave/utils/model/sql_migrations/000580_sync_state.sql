-- +migrate Down
ALTER TABLE sync_state DROP COLUMN IF EXISTS updated_at;
-- +migrate Up
ALTER TABLE sync_state ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;