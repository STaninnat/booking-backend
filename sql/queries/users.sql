-- name: CreateUser :exec
INSERT INTO users (id, created_at, updated_at, name, password, api_key, api_key_expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7);
--

-- name: GetUser :one
SELECT * FROM users WHERE api_key = $1;
--

-- name: GetUserByName :one
SELECT * FROM users WHERE name = $1;
--

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;
--

-- name: UpdateUser :exec
UPDATE users
SET updated_at = $1, api_key = $2, api_key_expires_at = $3
WHERE id = $4;
--