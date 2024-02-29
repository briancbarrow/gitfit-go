-- name: GetUser :one
SELECT * FROM users
WHERE stytch_id = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY name;

-- name: CreateUser :one
INSERT INTO users (stytch_id, database_id, created)
	VALUES(?, ?, datetime('now'))
RETURNING *;

-- name: UserExists :one
SELECT EXISTS(SELECT true FROM users WHERE stytch_id = ?);



-- -- name: UpdateAuthor :exec
-- UPDATE authors
-- set name = ?,
-- bio = ?
-- WHERE id = ?
-- RETURNING *;

-- -- name: DeleteAuthor :exec
-- DELETE FROM authors
-- WHERE id = ?;