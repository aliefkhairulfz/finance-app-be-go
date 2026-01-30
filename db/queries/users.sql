-- name: CreateUser :one
INSERT INTO users (name, email)
VALUES ($1, $2)
RETURNING id, name, email, email_verified, image;

-- name: FindUserByEmail :one
SELECT id, name, email, email_verified, image
FROM users
WHERE email = $1;

-- name: FindUserById :one
SELECT id, name, email, email_verified, image
FROM users
WHERE id = $1;

-- name: FindUsers :many
SELECT id, name, email, email_verified, image
FROM users
LIMIT $1 OFFSET $2;
