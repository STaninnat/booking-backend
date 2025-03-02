// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users.sql

package database

import (
	"context"
	"database/sql"
	"time"
)

const checkUserExistsByEmail = `-- name: CheckUserExistsByEmail :one

SELECT EXISTS (SELECT email FROM users WHERE email = $1)
`

func (q *Queries) CheckUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	row := q.db.QueryRowContext(ctx, checkUserExistsByEmail, email)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkUserExistsByFullname = `-- name: CheckUserExistsByFullname :one

SELECT EXISTS (SELECT full_name FROM users WHERE full_name = $1)
`

func (q *Queries) CheckUserExistsByFullname(ctx context.Context, fullName string) (bool, error) {
	row := q.db.QueryRowContext(ctx, checkUserExistsByFullname, fullName)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkUserExistsByUsername = `-- name: CheckUserExistsByUsername :one

SELECT EXISTS (SELECT username FROM users WHERE username = $1)
`

func (q *Queries) CheckUserExistsByUsername(ctx context.Context, username string) (bool, error) {
	row := q.db.QueryRowContext(ctx, checkUserExistsByUsername, username)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createUser = `-- name: CreateUser :exec
INSERT INTO users (id, created_at, updated_at, full_name, email, username, password, api_key, api_key_expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
`

type CreateUserParams struct {
	ID              string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	FullName        string
	Email           string
	Username        string
	Password        string
	ApiKey          string
	ApiKeyExpiresAt time.Time
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.ExecContext(ctx, createUser,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.FullName,
		arg.Email,
		arg.Username,
		arg.Password,
		arg.ApiKey,
		arg.ApiKeyExpiresAt,
	)
	return err
}

const getUserByID = `-- name: GetUserByID :one

SELECT id, created_at, updated_at, full_name, email, phone, username, password, api_key, api_key_expires_at FROM users 
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetUserByID(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.FullName,
		&i.Email,
		&i.Phone,
		&i.Username,
		&i.Password,
		&i.ApiKey,
		&i.ApiKeyExpiresAt,
	)
	return i, err
}

const getUserByKey = `-- name: GetUserByKey :one

SELECT id, created_at, updated_at, full_name, email, phone, username, password, api_key, api_key_expires_at FROM users 
WHERE api_key = $1
LIMIT 1
`

func (q *Queries) GetUserByKey(ctx context.Context, apiKey string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByKey, apiKey)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.FullName,
		&i.Email,
		&i.Phone,
		&i.Username,
		&i.Password,
		&i.ApiKey,
		&i.ApiKeyExpiresAt,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one

SELECT id, created_at, updated_at, full_name, email, phone, username, password, api_key, api_key_expires_at FROM users 
WHERE username = $1
LIMIT 1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.FullName,
		&i.Email,
		&i.Phone,
		&i.Username,
		&i.Password,
		&i.ApiKey,
		&i.ApiKeyExpiresAt,
	)
	return i, err
}

const updateUserInfo = `-- name: UpdateUserInfo :exec

UPDATE users
SET updated_at = $1, full_name = $2, email = $3, phone = $4
WHERE id = $5
`

type UpdateUserInfoParams struct {
	UpdatedAt time.Time
	FullName  string
	Email     string
	Phone     sql.NullString
	ID        string
}

func (q *Queries) UpdateUserInfo(ctx context.Context, arg UpdateUserInfoParams) error {
	_, err := q.db.ExecContext(ctx, updateUserInfo,
		arg.UpdatedAt,
		arg.FullName,
		arg.Email,
		arg.Phone,
		arg.ID,
	)
	return err
}

const updateUserKey = `-- name: UpdateUserKey :exec

UPDATE users
SET updated_at = $1, api_key = $2, api_key_expires_at = $3
WHERE id = $4
`

type UpdateUserKeyParams struct {
	UpdatedAt       time.Time
	ApiKey          string
	ApiKeyExpiresAt time.Time
	ID              string
}

func (q *Queries) UpdateUserKey(ctx context.Context, arg UpdateUserKeyParams) error {
	_, err := q.db.ExecContext(ctx, updateUserKey,
		arg.UpdatedAt,
		arg.ApiKey,
		arg.ApiKeyExpiresAt,
		arg.ID,
	)
	return err
}
