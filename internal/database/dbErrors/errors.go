package dbErrors

import (
	"errors"
)

var (
	ErrNoRecord           = errors.New("database: no matching record found")
	ErrInvalidCredentials = errors.New("database: invalid credentials")
	ErrDuplicateEmail     = errors.New("database: duplicate email")
)
