-- +goose Up
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users ( 
  "stytch_id" TEXT NOT NULL UNIQUE,
  "first_name" TEXT NOT NULL,
  "last_name" TEXT,
  "email" TEXT NOT NULL UNIQUE,
  "created" TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP 
);
-- +goose StatementEnd
