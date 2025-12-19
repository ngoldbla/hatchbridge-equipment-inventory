-- +goose Up
-- Fix column names to match Ent ORM expectations
-- PostgreSQL supports RENAME COLUMN directly
ALTER TABLE loans RENAME COLUMN checked_out_by_checkouts TO user_checkouts;
ALTER TABLE loans RENAME COLUMN returned_by_returns TO user_returns;

-- +goose Down
-- Revert to old column names
ALTER TABLE loans RENAME COLUMN user_checkouts TO checked_out_by_checkouts;
ALTER TABLE loans RENAME COLUMN user_returns TO returned_by_returns;
