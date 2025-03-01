-- name: CreateUser :exec
INSERT INTO users (id, created_at, updated_at, full_name, email, username, password, api_key, api_key_expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);
--

-- name: CheckUserExistsByUsername :one
SELECT EXISTS (SELECT username FROM users WHERE username = $1);
--

-- name: CheckUserExistsByFullname :one
SELECT EXISTS (SELECT full_name FROM users WHERE full_name = $1);
--

-- name: CheckUserExistsByEmail :one
SELECT EXISTS (SELECT email FROM users WHERE email = $1);
--

-- name: GetUserByUsername :one
SELECT * FROM users 
WHERE api_key = $1
LIMIT 1;
--

-- name: GetUserByKey :one
SELECT * FROM users 
WHERE api_key = $1
LIMIT 1;
--

-- name: GetUserByID :one
SELECT * FROM users 
WHERE id = $1
LIMIT 1;
--

-- name: UpdateUserInfo :exec
UPDATE users
SET updated_at = $1, full_name = $2, email = $3, phone = $4
WHERE id = $5;
--

-- name: UpdateUserKey :exec
UPDATE users
SET updated_at = $1, api_key = $2, api_key_expires_at = $3
WHERE id = $4;
--
