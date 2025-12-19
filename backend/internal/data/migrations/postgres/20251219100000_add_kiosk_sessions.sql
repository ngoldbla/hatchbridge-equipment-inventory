-- +goose Up
-- Add kiosk_action field to loans table
ALTER TABLE loans ADD COLUMN kiosk_action BOOLEAN NOT NULL DEFAULT false;

-- Add self_registered field to borrowers table
ALTER TABLE borrowers ADD COLUMN self_registered BOOLEAN NOT NULL DEFAULT false;

-- Create kiosk_sessions table for tracking active kiosk states
CREATE TABLE IF NOT EXISTS kiosk_sessions (
    id              UUID         NOT NULL PRIMARY KEY,
    created_at      TIMESTAMPTZ  NOT NULL,
    updated_at      TIMESTAMPTZ  NOT NULL,
    is_active       BOOLEAN      NOT NULL DEFAULT true,
    unlocked_until  TIMESTAMPTZ,
    user_kiosk_session UUID
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
ALTER TABLE loans DROP COLUMN IF EXISTS kiosk_action;
ALTER TABLE borrowers DROP COLUMN IF EXISTS self_registered;
