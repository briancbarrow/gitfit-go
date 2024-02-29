-- name: ListWorkoutSets :many
SELECT workout_sets.*, exercises.name as exercise_name, exercises.target_area FROM workout_sets
JOIN exercises ON workout_sets.exercise = exercises.id;


-- name: CreateWorkoutSet :one
INSERT INTO workout_sets (date, exercise, reps, note)
  VALUES(?, ?, ?, ?)
RETURNING *;

-- name: ListExercises :many
SELECT * FROM exercises;

INSERT INTO users (stytch_id, database_id, created)
	VALUES(?, ?, datetime('now'))
RETURNING *;