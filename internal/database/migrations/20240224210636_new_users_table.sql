-- +goose Up
-- +goose StatementBegin
CREATE TABLE users_new (
  "stytch_id" TEXT NOT NULL PRIMARY KEY,
  "database_id" TEXT NOT NULL,
  "created" TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP 
);

INSERT INTO users_new (stytch_id, database_id, created)
SELECT stytch_id, database_id, created FROM users;

DROP TABLE users;

ALTER TABLE users_new RENAME TO users;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users
-- +goose StatementEnd