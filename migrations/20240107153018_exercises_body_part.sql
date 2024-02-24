-- +goose Up
-- +goose StatementBegin
ALTER table exercises ADD COLUMN
target_area TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE exercises_new (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL
  -- include all other columns from the original table, except 'target_area'
);

-- Copy data from old table to new table
INSERT INTO exercises_new (id, name)
SELECT id, name FROM exercises;

-- Delete old table
DROP TABLE exercises;

-- Rename new table to old table name
ALTER TABLE exercises_new RENAME TO exercises;
-- +goose StatementEnd
