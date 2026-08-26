-- name: GetUserByEmail :one
SELECT id, email, password_hash, created_at
FROM users
WHERE lower(email) = lower(sqlc.arg(email));

-- name: GetUserByID :one
SELECT id, email, password_hash, created_at
FROM users
WHERE id = $1;

-- name: CountUsers :one
SELECT count(*) FROM users;

-- name: CreateUser :one
INSERT INTO users (email, password_hash)
VALUES (sqlc.arg(email), sqlc.arg(password_hash))
RETURNING id, email, password_hash, created_at;

-- name: UpdateUserPassword :exec
UPDATE users SET password_hash = sqlc.arg(password_hash) WHERE id = $1;
