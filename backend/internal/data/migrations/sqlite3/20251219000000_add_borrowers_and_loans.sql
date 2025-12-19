-- +goose Up
-- Create borrowers table for tracking individuals who can borrow equipment
CREATE TABLE IF NOT EXISTS borrowers (
    id                  uuid         NOT NULL PRIMARY KEY,
    created_at          datetime     NOT NULL,
    updated_at          datetime     NOT NULL,
    name                text         NOT NULL,
    email               text         NOT NULL,
    phone               text,
    organization        text,
    student_id          text,
    notes               text,
    is_active           bool         NOT NULL DEFAULT true,
    group_borrowers     uuid         NOT NULL
        CONSTRAINT borrowers_groups_borrowers
            REFERENCES groups(id)
            ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS borrowers_name_idx ON borrowers(name);
CREATE INDEX IF NOT EXISTS borrowers_email_idx ON borrowers(email);
CREATE INDEX IF NOT EXISTS borrowers_is_active_idx ON borrowers(is_active);
CREATE INDEX IF NOT EXISTS borrowers_group_borrowers_idx ON borrowers(group_borrowers);

-- Create loans table for tracking equipment checkouts
CREATE TABLE IF NOT EXISTS loans (
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

CREATE INDEX IF NOT EXISTS loans_checked_out_at_idx ON loans(checked_out_at);
CREATE INDEX IF NOT EXISTS loans_due_at_idx ON loans(due_at);
CREATE INDEX IF NOT EXISTS loans_returned_at_idx ON loans(returned_at);
CREATE INDEX IF NOT EXISTS loans_item_loans_idx ON loans(item_loans);
CREATE INDEX IF NOT EXISTS loans_borrower_loans_idx ON loans(borrower_loans);
CREATE INDEX IF NOT EXISTS loans_group_loans_idx ON loans(group_loans);

-- +goose Down
DROP TABLE IF EXISTS loans;
DROP TABLE IF EXISTS borrowers;
