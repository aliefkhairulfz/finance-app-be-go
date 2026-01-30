-- name: CreateAccount :one
INSERT INTO accounts (provider_id, password, user_id)
VALUES ($1, $2, $3)
RETURNING id, provider_id, password, user_id;

-- name: FindAccountByUserID :one
SELECT id, provider_id, password, user_id
FROM accounts
WHERE user_id = $1;

-- name: FindAccountByID :one
SELECT id, provider_id, password, user_id
FROM accounts
WHERE id = $1;

-- name: FindAccountByUserIdAndProviderID :one
SELECT id, provider_id, password, user_id
FROM accounts
WHERE user_id = $1 AND provider_id = $2;
