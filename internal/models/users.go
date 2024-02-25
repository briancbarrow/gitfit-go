package models

import (
	"database/sql"
	"time"

	"github.com/mattn/go-sqlite3"
)

type User struct {
	ID             int
	FirstName      string
	LastName       string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(databaseID, stytchId string) error {
	stmt := `INSERT INTO users (stytch_id, database_id, created)
	VALUES(?, ?, datetime('now'))`
	_, err := m.DB.Exec(stmt, stytchId, databaseID)

	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.Code == sqlite3.ErrConstraint {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"

	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}
