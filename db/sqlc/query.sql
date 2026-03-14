-- name: DeleteAuthor :exec
DELETE FROM authors WHERE id = $1;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (id, name, email, created_at)
VALUES ($1, $2, $3, NOW())
    RETURNING *;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
    LIMIT $1 OFFSET $2;

-- name: UpdateUserEmail :exec
UPDATE users SET email = $2 WHERE id = $1;