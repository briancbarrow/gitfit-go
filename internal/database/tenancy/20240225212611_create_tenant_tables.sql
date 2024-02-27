-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS exercises ( 
  "id" INTEGER PRIMARY KEY,
  "name" TEXT NOT NULL,
  "target_area" TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS workout_sets (
  "id" INTEGER PRIMARY KEY,
  "date" TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "exercise" INTEGER NOT NULL,
  "reps" INTEGER NOT NULL DEFAULT 0,
  "note" TEXT,
  FOREIGN KEY(exercise) REFERENCES exercises(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE exercises;
DROP TABLE workout_sets;
-- +goose StatementEnd
