-- name: CreateEntry :one
INSERT INTO entries (
    account_id,
    amount
) VALUES (
   2, 20
)RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = 1 LIMIT 1;

-- name: ListEntry :many
SELECT * FROM entries
ORDER BY name
LIMIT 1
OFFSET 3;

-- name: UpdateEntry :one
UPDATE entries
set balance = 2
WHERE id = 1
RETURNING *;

-- name: DeleteEntry :exec
DELETE FROM entries WHERE id = 1;