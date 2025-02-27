// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users_key.sql

package database

import (
	"context"
	"time"
)

const createUserRfKey = `-- name: CreateUserRfKey :exec
INSERT INTO users_key (id, created_at, updated_at ,access_token_expires_at, refresh_token, refresh_token_expires_at, user_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
`

type CreateUserRfKeyParams struct {
	ID                    string
	CreatedAt             time.Time
	UpdatedAt             time.Time
	AccessTokenExpiresAt  time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
	UserID                string
}

func (q *Queries) CreateUserRfKey(ctx context.Context, arg CreateUserRfKeyParams) error {
	_, err := q.db.ExecContext(ctx, createUserRfKey,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.AccessTokenExpiresAt,
		arg.RefreshToken,
		arg.RefreshTokenExpiresAt,
		arg.UserID,
	)
	return err
}

const getRfKeyByUserID = `-- name: GetRfKeyByUserID :one

SELECT id, created_at, updated_at, access_token_expires_at, refresh_token, refresh_token_expires_at, user_id FROM users_key WHERE user_id = $1
LIMIT 1
`

func (q *Queries) GetRfKeyByUserID(ctx context.Context, userID string) (UsersKey, error) {
	row := q.db.QueryRowContext(ctx, getRfKeyByUserID, userID)
	var i UsersKey
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.AccessTokenExpiresAt,
		&i.RefreshToken,
		&i.RefreshTokenExpiresAt,
		&i.UserID,
	)
	return i, err
}

const getUserByRfKey = `-- name: GetUserByRfKey :one

SELECT id, created_at, updated_at, access_token_expires_at, refresh_token, refresh_token_expires_at, user_id FROM users_key WHERE refresh_token = $1 
LIMIT 1
`

func (q *Queries) GetUserByRfKey(ctx context.Context, refreshToken string) (UsersKey, error) {
	row := q.db.QueryRowContext(ctx, getUserByRfKey, refreshToken)
	var i UsersKey
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.AccessTokenExpiresAt,
		&i.RefreshToken,
		&i.RefreshTokenExpiresAt,
		&i.UserID,
	)
	return i, err
}

const updateUserRfKey = `-- name: UpdateUserRfKey :exec

UPDATE users_key
SET updated_at = $1, access_token_expires_at = $2, refresh_token = $3, refresh_token_expires_at = $4
WHERE user_id = $5
`

type UpdateUserRfKeyParams struct {
	UpdatedAt             time.Time
	AccessTokenExpiresAt  time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
	UserID                string
}

func (q *Queries) UpdateUserRfKey(ctx context.Context, arg UpdateUserRfKeyParams) error {
	_, err := q.db.ExecContext(ctx, updateUserRfKey,
		arg.UpdatedAt,
		arg.AccessTokenExpiresAt,
		arg.RefreshToken,
		arg.RefreshTokenExpiresAt,
		arg.UserID,
	)
	return err
}
