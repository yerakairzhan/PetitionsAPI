-- name: CreateUser :exec
INSERT INTO users (username, password_hash) VALUES ($1, $2);

-- name: GetUser :one
SELECT id, username, created_at FROM Users WHERE id = $1;

-- name: ListUsers :many
SELECT id, username, created_at FROM Users;

-- name: GetUserByUsername :one
SELECT id, username, password_hash, created_at
FROM users
WHERE username = $1;



