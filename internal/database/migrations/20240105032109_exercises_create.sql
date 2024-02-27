-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS exercises ( 
  "id" INTEGER PRIMARY KEY,
  "name" TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE exercises;
-- +goose StatementEnd
