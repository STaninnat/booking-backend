-- name: CreateUser :exec
INSERT INTO users (id, created_at, updated_at, full_name, last_name, username, password, api_key, api_key_expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);
--

-- name: GetUser :one
SELECT * FROM users WHERE api_key = $1;
--

-- name: GetUserByName :one
SELECT * FROM users WHERE username = $1;
--

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;
--

-- name: UpdateUser :exec
UPDATE users
SET updated_at = $1, api_key = $2, api_key_expires_at = $3
WHERE id = $4;
--