// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package tenant_database

import (
	"database/sql"
)

type Exercise struct {
	ID         int64
	Name       string
	TargetArea string
}

type WorkoutSet struct {
	ID       int64
	Date     string
	Exercise int64
	Reps     int64
	Note     sql.NullString
}
