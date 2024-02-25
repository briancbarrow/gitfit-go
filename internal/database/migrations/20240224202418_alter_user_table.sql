-- +goose Up
-- +goose StatementBegin
-- Create a new table with only the stytch_id and database_id columns
CREATE TABLE users_new (
  "stytch_id" TEXT NOT NULL UNIQUE,
  "database_id" INTEGER,
  "created" TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP 
);

-- Copy the stytch_id data from the old table to the new table
INSERT INTO users_new (stytch_id)
SELECT stytch_id FROM users;

-- Drop the old table
DROP TABLE users;

-- Rename the new table to the old table's name
ALTER TABLE users_new RENAME TO users;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
