-- +goose Up
-- Fix column names to match Ent ORM expectations
-- SQLite doesn't support RENAME COLUMN directly in older versions, so we need to recreate the table

-- Step 1: Create new loans table with correct column names
CREATE TABLE loans_new (
    id                      uuid         NOT NULL PRIMARY KEY,
    created_at              datetime     NOT NULL,
    updated_at              datetime     NOT NULL,
    checked_out_at          datetime     NOT NULL,
    due_at                  datetime     NOT NULL,
    returned_at             datetime,
    notes                   text,
    return_notes            text,
    quantity                integer      NOT NULL DEFAULT 1,
    item_loans              uuid         NOT NULL
        CONSTRAINT loans_items_loans
            REFERENCES items(id)
            ON DELETE CASCADE,
    borrower_loans          uuid         NOT NULL
        CONSTRAINT loans_borrowers_loans
            REFERENCES borrowers(id)
            ON DELETE CASCADE,
    user_checkouts          uuid
        CONSTRAINT loans_users_checkouts
            REFERENCES users(id)
            ON DELETE SET NULL,
    user_returns            uuid
        CONSTRAINT loans_users_returns
            REFERENCES users(id)
            ON DELETE SET NULL,
    group_loans             uuid         NOT NULL
        CONSTRAINT loans_groups_loans
            REFERENCES groups(id)
            ON DELETE CASCADE
);

-- Step 2: Copy data from old table (if exists and has data)
INSERT OR IGNORE INTO loans_new (id, created_at, updated_at, checked_out_at, due_at, returned_at, notes, return_notes, quantity, item_loans, borrower_loans, user_checkouts, user_returns, group_loans)
SELECT id, created_at, updated_at, checked_out_at, due_at, returned_at, notes, return_notes, quantity, item_loans, borrower_loans, checked_out_by_checkouts, returned_by_returns, group_loans
FROM loans;

-- Step 3: Drop old table
DROP TABLE loans;

-- Step 4: Rename new table to loans
ALTER TABLE loans_new RENAME TO loans;

-- Step 5: Recreate indexes
CREATE INDEX IF NOT EXISTS loans_checked_out_at_idx ON loans(checked_out_at);
CREATE INDEX IF NOT EXISTS loans_due_at_idx ON loans(due_at);
CREATE INDEX IF NOT EXISTS loans_returned_at_idx ON loans(returned_at);
CREATE INDEX IF NOT EXISTS loans_item_loans_idx ON loans(item_loans);
CREATE INDEX IF NOT EXISTS loans_borrower_loans_idx ON loans(borrower_loans);
CREATE INDEX IF NOT EXISTS loans_group_loans_idx ON loans(group_loans);

-- +goose Down
-- Revert to old column names
CREATE TABLE loans_old (
    id                      uuid         NOT NULL PRIMARY KEY,
    created_at              datetime     NOT NULL,
    updated_at              datetime     NOT NULL,
    checked_out_at          datetime     NOT NULL,
    due_at                  datetime     NOT NULL,
    returned_at             datetime,
    notes                   text,
    return_notes            text,
    quantity                integer      NOT NULL DEFAULT 1,
    item_loans              uuid         NOT NULL
        CONSTRAINT loans_items_loans
            REFERENCES items(id)
            ON DELETE CASCADE,
    borrower_loans          uuid         NOT NULL
        CONSTRAINT loans_borrowers_loans
            REFERENCES borrowers(id)
            ON DELETE CASCADE,
    checked_out_by_checkouts uuid
        CONSTRAINT loans_users_checkouts
            REFERENCES users(id)
            ON DELETE SET NULL,
    returned_by_returns     uuid
        CONSTRAINT loans_users_returns
            REFERENCES users(id)
            ON DELETE SET NULL,
    group_loans             uuid         NOT NULL
        CONSTRAINT loans_groups_loans
            REFERENCES groups(id)
            ON DELETE CASCADE
);

INSERT OR IGNORE INTO loans_old (id, created_at, updated_at, checked_out_at, due_at, returned_at, notes, return_notes, quantity, item_loans, borrower_loans, checked_out_by_checkouts, returned_by_returns, group_loans)
SELECT id, created_at, updated_at, checked_out_at, due_at, returned_at, notes, return_notes, quantity, item_loans, borrower_loans, user_checkouts, user_returns, group_loans
FROM loans;

DROP TABLE loans;
ALTER TABLE loans_old RENAME TO loans;

CREATE INDEX IF NOT EXISTS loans_checked_out_at_idx ON loans(checked_out_at);
CREATE INDEX IF NOT EXISTS loans_due_at_idx ON loans(due_at);
CREATE INDEX IF NOT EXISTS loans_returned_at_idx ON loans(returned_at);
CREATE INDEX IF NOT EXISTS loans_item_loans_idx ON loans(item_loans);
CREATE INDEX IF NOT EXISTS loans_borrower_loans_idx ON loans(borrower_loans);
CREATE INDEX IF NOT EXISTS loans_group_loans_idx ON loans(group_loans);
