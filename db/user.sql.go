// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package db

import (
	"context"
	"database/sql"
)

const createUser = `-- name: CreateUser :exec
INSERT INTO users (username, password_hash) VALUES ($1, $2)
`

type CreateUserParams struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.ExecContext(ctx, createUser, arg.Username, arg.PasswordHash)
	return err
}

const getUser = `-- name: GetUser :one
SELECT id, username, created_at FROM Users WHERE id = $1
`

type GetUserRow struct {
	ID        int32        `json:"id"`
	Username  string       `json:"username"`
	CreatedAt sql.NullTime `json:"created_at"`
}

func (q *Queries) GetUser(ctx context.Context, id int32) (GetUserRow, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i GetUserRow
	err := row.Scan(&i.ID, &i.Username, &i.CreatedAt)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, username, password_hash, created_at
FROM users
WHERE username = $1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.PasswordHash,
		&i.CreatedAt,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, username, created_at FROM Users
`

type ListUsersRow struct {
	ID        int32        `json:"id"`
	Username  string       `json:"username"`
	CreatedAt sql.NullTime `json:"created_at"`
}

func (q *Queries) ListUsers(ctx context.Context) ([]ListUsersRow, error) {
	rows, err := q.db.QueryContext(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListUsersRow
	for rows.Next() {
		var i ListUsersRow
		if err := rows.Scan(&i.ID, &i.Username, &i.CreatedAt); err != nil {
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
