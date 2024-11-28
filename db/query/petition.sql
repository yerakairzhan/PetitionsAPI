-- name: CreatePetition :one
INSERT INTO petitions (title, description, created_at, number_votes, user_id)
VALUES ($1, $2, NOW(), 0, $3)
RETURNING id, title, description, created_at, number_votes, user_id;

-- name: GetPetition :one
SELECT id, title, description, created_at, number_votes
FROM petitions
WHERE id = $1;

-- name: ListPetitions :many
SELECT id, title, description, created_at, number_votes
FROM petitions
ORDER BY created_at DESC;

-- name: UpdatePetitionVotes :exec
UPDATE petitions
SET number_votes = $1
WHERE id = $2;

-- name: DeletePetition :exec
DELETE FROM petitions
WHERE id = $1;


-- name: ListPetitionsByCreatedAtAsc :many
SELECT id, title, description, created_at, number_votes, user_id
FROM petitions
ORDER BY created_at ASC;


-- name: ListPetitionsByCreatedAtDesc :many
SELECT id, title, description, created_at, number_votes, user_id
FROM petitions
ORDER BY created_at DESC;

-- name: ListPetitionsByVotesAsc :many
SELECT id, title, description, created_at, number_votes, user_id
FROM petitions
ORDER BY number_votes ASC;

-- name: ListPetitionsByVotesDesc :many
SELECT id, title, description, created_at, number_votes, user_id
FROM petitions
ORDER BY number_votes DESC;

-- name: GetPetitionByID :one
SELECT id, title, description, created_at, number_votes, user_id
FROM petitions
WHERE id = $1;








                                 


