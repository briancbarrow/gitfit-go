// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package tenant_database

import (
	"context"
	"database/sql"
)

const createWorkoutSet = `-- name: CreateWorkoutSet :one
INSERT INTO workout_sets (date, exercise, reps, note)
  VALUES(?, ?, ?, ?)
RETURNING id, date, exercise, reps, note
`

type CreateWorkoutSetParams struct {
	Date     string
	Exercise int64
	Reps     int64
	Note     sql.NullString
}

func (q *Queries) CreateWorkoutSet(ctx context.Context, arg CreateWorkoutSetParams) (WorkoutSet, error) {
	row := q.db.QueryRowContext(ctx, createWorkoutSet,
		arg.Date,
		arg.Exercise,
		arg.Reps,
		arg.Note,
	)
	var i WorkoutSet
	err := row.Scan(
		&i.ID,
		&i.Date,
		&i.Exercise,
		&i.Reps,
		&i.Note,
	)
	return i, err
}

const deleteWorkoutSet = `-- name: DeleteWorkoutSet :exec
DELETE from workout_sets WHERE id = ?
`

func (q *Queries) DeleteWorkoutSet(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteWorkoutSet, id)
	return err
}

const getWorkoutSetCounts = `-- name: GetWorkoutSetCounts :many
SELECT date, COUNT(*) as count 
FROM workout_sets
GROUP BY date
`

type GetWorkoutSetCountsRow struct {
	Date  string
	Count int64
}

func (q *Queries) GetWorkoutSetCounts(ctx context.Context) ([]GetWorkoutSetCountsRow, error) {
	rows, err := q.db.QueryContext(ctx, getWorkoutSetCounts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetWorkoutSetCountsRow
	for rows.Next() {
		var i GetWorkoutSetCountsRow
		if err := rows.Scan(&i.Date, &i.Count); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listExercises = `-- name: ListExercises :many
SELECT id, name, target_area FROM exercises
`

func (q *Queries) ListExercises(ctx context.Context) ([]Exercise, error) {
	rows, err := q.db.QueryContext(ctx, listExercises)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Exercise
	for rows.Next() {
		var i Exercise
		if err := rows.Scan(&i.ID, &i.Name, &i.TargetArea); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listWorkoutSets = `-- name: ListWorkoutSets :many
SELECT workout_sets.id, workout_sets.date, workout_sets.exercise, workout_sets.reps, workout_sets.note, exercises.name as exercise_name, exercises.target_area FROM workout_sets
JOIN exercises ON workout_sets.exercise = exercises.id
WHERE workout_sets.date = ?
`

type ListWorkoutSetsRow struct {
	ID           int64
	Date         string
	Exercise     int64
	Reps         int64
	Note         sql.NullString
	ExerciseName string
	TargetArea   string
}

func (q *Queries) ListWorkoutSets(ctx context.Context, date string) ([]ListWorkoutSetsRow, error) {
	rows, err := q.db.QueryContext(ctx, listWorkoutSets, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListWorkoutSetsRow
	for rows.Next() {
		var i ListWorkoutSetsRow
		if err := rows.Scan(
			&i.ID,
			&i.Date,
			&i.Exercise,
			&i.Reps,
			&i.Note,
			&i.ExerciseName,
			&i.TargetArea,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
