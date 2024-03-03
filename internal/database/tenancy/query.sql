-- name: GetWorkoutSetCounts :many
SELECT date, COUNT(*) as count 
FROM workout_sets
GROUP BY date;

-- name: ListWorkoutSets :many
SELECT workout_sets.*, exercises.name as exercise_name, exercises.target_area FROM workout_sets
JOIN exercises ON workout_sets.exercise = exercises.id
WHERE workout_sets.date = ?;


-- name: CreateWorkoutSet :one
INSERT INTO workout_sets (date, exercise, reps, note)
  VALUES(?, ?, ?, ?)
RETURNING *;

-- name: DeleteWorkoutSet :exec
DELETE from workout_sets WHERE id = ?;

-- name: ListExercises :many
SELECT * FROM exercises;

INSERT INTO users (stytch_id, database_id, created)
	VALUES(?, ?, datetime('now'))
RETURNING *;