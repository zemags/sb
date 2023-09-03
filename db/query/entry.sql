-- name: CreateEntry :one
INSERT INTO entries (
    account_id,
    amount
) VALUES (
    $1,
    $2
)
RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries WHERE id = $1;

-- name: GetEntryByAccount :many
SELECT * FROM entries WHERE account_id = $1;

-- name: UpdateEntry :one
UPDATE entries SET
    account_id = $1,
    amount = $2
WHERE id = $3
RETURNING *;

-- name: UpdateEntryByAccount :many
UPDATE entries SET
    account_id = $1,
    amount = $2
WHERE account_id = $3
RETURNING *;

-- name: DeleteEntry :exec
DELETE FROM entries WHERE id = $1;

-- name: DeleteEntryByAccount :exec
DELETE FROM entries WHERE account_id = $1;


