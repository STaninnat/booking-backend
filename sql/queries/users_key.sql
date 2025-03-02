-- name: CreateUserRfKey :exec
INSERT INTO users_token (id, created_at, updated_at, access_token_expires_at, refresh_token, refresh_token_expires_at, user_id)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: UpdateUserTK :exec
UPDATE users_token
SET updated_at = $1, access_token_expires_at = $2, refresh_token = $3, refresh_token_expires_at = $4
WHERE user_id = $5;

-- name: GetRfKeyByUserID :one
SELECT * FROM users_token WHERE user_id = $1
LIMIT 1;

-- name: GetUserByRfKey :one
SELECT * FROM users_token WHERE refresh_token = $1 
LIMIT 1;