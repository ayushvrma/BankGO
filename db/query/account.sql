-- name: CreateAccount :one
INSERT INTO accounts (
    owner,
    balance,
    currency
) VALUES (
  "ayush", 2, "USD"
)RETURNING *;