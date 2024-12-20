// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: votes.sql

package db

import (
	"context"
	"database/sql"
)

const createVote = `-- name: CreateVote :exec
INSERT INTO votes (user_id, petition_id, created_at)
VALUES ($1, $2, $3)
`

type CreateVoteParams struct {
	UserID     int32        `json:"user_id"`
	PetitionID int32        `json:"petition_id"`
	CreatedAt  sql.NullTime `json:"created_at"`
}

func (q *Queries) CreateVote(ctx context.Context, arg CreateVoteParams) error {
	_, err := q.db.ExecContext(ctx, createVote, arg.UserID, arg.PetitionID, arg.CreatedAt)
	return err
}

const getVoteByUserAndPetition = `-- name: GetVoteByUserAndPetition :one
SELECT
    votes.id AS vote_id,
    votes.created_at AS vote_date,
    users.username AS user_username,
    petitions.title AS petition_title
FROM votes
JOIN users ON votes.user_id = users.id
JOIN petitions ON votes.petition_id = petitions.id
WHERE votes.user_id = $1 AND votes.petition_id = $2
`

type GetVoteByUserAndPetitionParams struct {
	UserID     int32 `json:"user_id"`
	PetitionID int32 `json:"petition_id"`
}

type GetVoteByUserAndPetitionRow struct {
	VoteID        int32        `json:"vote_id"`
	VoteDate      sql.NullTime `json:"vote_date"`
	UserUsername  string       `json:"user_username"`
	PetitionTitle string       `json:"petition_title"`
}

func (q *Queries) GetVoteByUserAndPetition(ctx context.Context, arg GetVoteByUserAndPetitionParams) (GetVoteByUserAndPetitionRow, error) {
	row := q.db.QueryRowContext(ctx, getVoteByUserAndPetition, arg.UserID, arg.PetitionID)
	var i GetVoteByUserAndPetitionRow
	err := row.Scan(
		&i.VoteID,
		&i.VoteDate,
		&i.UserUsername,
		&i.PetitionTitle,
	)
	return i, err
}

const hasUserVoted = `-- name: HasUserVoted :one
SELECT EXISTS (
    SELECT 1
    FROM votes
    WHERE user_id = $1 AND petition_id = $2
) AS exists
`

type HasUserVotedParams struct {
	UserID     int32 `json:"user_id"`
	PetitionID int32 `json:"petition_id"`
}

func (q *Queries) HasUserVoted(ctx context.Context, arg HasUserVotedParams) (bool, error) {
	row := q.db.QueryRowContext(ctx, hasUserVoted, arg.UserID, arg.PetitionID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const incrementVoteCount = `-- name: IncrementVoteCount :exec
UPDATE petitions
SET number_votes = number_votes + 1
WHERE id = $1
`

func (q *Queries) IncrementVoteCount(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, incrementVoteCount, id)
	return err
}

const listVotes = `-- name: ListVotes :many
SELECT
    votes.id AS vote_id,
    votes.created_at AS vote_date,
    users.username AS user_username,
    petitions.title AS petition_title
FROM votes
JOIN users ON votes.user_id = users.id
JOIN petitions ON votes.petition_id = petitions.id
ORDER BY votes.created_at DESC
`

type ListVotesRow struct {
	VoteID        int32        `json:"vote_id"`
	VoteDate      sql.NullTime `json:"vote_date"`
	UserUsername  string       `json:"user_username"`
	PetitionTitle string       `json:"petition_title"`
}

func (q *Queries) ListVotes(ctx context.Context) ([]ListVotesRow, error) {
	rows, err := q.db.QueryContext(ctx, listVotes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListVotesRow
	for rows.Next() {
		var i ListVotesRow
		if err := rows.Scan(
			&i.VoteID,
			&i.VoteDate,
			&i.UserUsername,
			&i.PetitionTitle,
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

const recordVote = `-- name: RecordVote :exec
INSERT INTO votes (user_id, petition_id, created_at)
VALUES ($1, $2, NOW())
`

type RecordVoteParams struct {
	UserID     int32 `json:"user_id"`
	PetitionID int32 `json:"petition_id"`
}

func (q *Queries) RecordVote(ctx context.Context, arg RecordVoteParams) error {
	_, err := q.db.ExecContext(ctx, recordVote, arg.UserID, arg.PetitionID)
	return err
}
