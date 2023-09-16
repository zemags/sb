-- name: CreateTransfer :one
INSERT INTO transfers (
    from_account_id,
    to_account_id,
    amount
) VALUES (
    $1,
    $2,
    $3
) RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers WHERE id = $1;

-- name: GetTransfers :many
SELECT * FROM transfers WHERE from_account_id = $1 OR to_account_id = $1 ORDER BY created_at DESC;

-- name: GetTransfersByAccount :many
SELECT * FROM transfers WHERE from_account_id = $1 OR to_account_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: UpdateTransfer :one
UPDATE transfers SET
    from_account_id = $1,
    to_account_id = $2,
    amount = $3
WHERE id = $4
RETURNING *;

-- name: DeleteTransfer :exec
DELETE FROM transfers WHERE id = $1;

