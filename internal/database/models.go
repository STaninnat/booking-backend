// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package database

import (
	"time"
)

type User struct {
	ID              string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	FullName        string
	LastName        string
	Username        string
	Password        string
	ApiKey          string
	ApiKeyExpiresAt time.Time
}

type UsersKey struct {
	ID                    string
	CreatedAt             time.Time
	UpdatedAt             time.Time
	AccessTokenExpiresAt  time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
	UserID                string
}
