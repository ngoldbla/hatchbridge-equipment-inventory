-- +goose Up
-- Add kiosk_action field to loans table
ALTER TABLE loans ADD COLUMN kiosk_action bool NOT NULL DEFAULT false;

-- Add self_registered field to borrowers table
ALTER TABLE borrowers ADD COLUMN self_registered bool NOT NULL DEFAULT false;

-- Create kiosk_sessions table for tracking active kiosk states
CREATE TABLE IF NOT EXISTS kiosk_sessions (
    id              uuid         NOT NULL PRIMARY KEY,
    created_at      datetime     NOT NULL,
    updated_at      datetime     NOT NULL,
    is_active       bool         NOT NULL DEFAULT true,
    unlocked_until  datetime,
    user_kiosk_session uuid
        CONSTRAINT kiosk_sessions_users_kiosk_session
            REFERENCES users(id)
            ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS kiosk_sessions_user_kiosk_session_key ON kiosk_sessions(user_kiosk_session);
CREATE INDEX IF NOT EXISTS kiosk_sessions_is_active_idx ON kiosk_sessions(is_active);

-- +goose Down
DROP INDEX IF EXISTS kiosk_sessions_is_active_idx;
DROP INDEX IF EXISTS kiosk_sessions_user_kiosk_session_key;
DROP TABLE IF EXISTS kiosk_sessions;
-- SQLite doesn't support DROP COLUMN, would need table recreation for full rollback
