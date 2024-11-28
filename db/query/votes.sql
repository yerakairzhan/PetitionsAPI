-- name: CreateVote :exec
INSERT INTO votes (user_id, petition_id, created_at)
VALUES ($1, $2, $3);

-- name: ListVotes :many
SELECT
    votes.id AS vote_id,
    votes.created_at AS vote_date,
    users.username AS user_username,
    petitions.title AS petition_title
FROM votes
JOIN users ON votes.user_id = users.id
JOIN petitions ON votes.petition_id = petitions.id
ORDER BY votes.created_at DESC;

-- name: GetVoteByUserAndPetition :one
SELECT
    votes.id AS vote_id,
    votes.created_at AS vote_date,
    users.username AS user_username,
    petitions.title AS petition_title
FROM votes
JOIN users ON votes.user_id = users.id
JOIN petitions ON votes.petition_id = petitions.id
WHERE votes.user_id = $1 AND votes.petition_id = $2;


-- name: HasUserVoted :one
SELECT EXISTS (
    SELECT 1
    FROM votes
    WHERE user_id = $1 AND petition_id = $2
) AS exists;


-- name: RecordVote :exec
INSERT INTO votes (user_id, petition_id, created_at)
VALUES ($1, $2, NOW());


-- name: IncrementVoteCount :exec
UPDATE petitions
SET number_votes = number_votes + 1
WHERE id = $1;

