-- name: CreateSession :one
INSERT INTO sessions (token, expires_at, user_id, ip_address, user_agent)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, token, expires_at, user_id;

-- name: UpdateSessionByUserId :one
UPDATE sessions SET token = $1, expires_at = $2, ip_address = $3, user_agent = $4
WHERE user_id = $5
RETURNING id, token, expires_at, user_id;

-- name: FindSessionByUserId :one
SELECT id, token, expires_at, user_id
FROM sessions
WHERE user_id = $1;

-- name: FindSessionByToken :one
SELECT id, token, expires_at, user_id
FROM sessions
WHERE token = $1;

-- name: DeleteSessionByUserId :exec
DELETE FROM sessions WHERE user_id = $1;

-- name: DeleteSessionByToken :exec
DELETE FROM sessions WHERE token = $1;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions WHERE expires_at < NOW();
