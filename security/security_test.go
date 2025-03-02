package security

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestIsValidUserNameFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid username", "user123", true},
		{"too short", "ab", false},
		{"too long", "abcdefghijklmnopqrstuvwxyz12345", false},
		{"invalid chars", "user@name", false},
		{"valid with dash", "user-name", true},
		{"valid with underscore", "user_name", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidUserNameFormat(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsValidateEmailFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid email", "test@example.com", true},
		{"missing @", "testexample.com", false},
		{"missing domain", "test@.com", false},
		{"missing .", "test@examplecom", false},
		{"valid with numbers", "user123@example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidateEmailFormat(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateJWTToken(t *testing.T) {
	secret := "mysecret"
	expirationTime := time.Now().Add(1 * time.Hour)
	userID := uuid.New()
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "my-api-service",
			Audience:  []string{"my-frontend-app"},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)

	result, err := ValidateJWTToken(tokenString, secret)
	assert.NoError(t, err)
	assert.Equal(t, claims.UserID, result.UserID)
	assert.Equal(t, claims.Issuer, result.Issuer)
	assert.Equal(t, claims.Audience, result.Audience)
	assert.Equal(t, claims.ExpiresAt, result.ExpiresAt)
	assert.Equal(t, claims.NotBefore, result.NotBefore)

	_, err = ValidateJWTToken("invalid.token.here", secret)
	assert.Error(t, err)

	expiredClaims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "my-api-service",
			Audience:  []string{"my-frontend-app"},
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		},
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredTokenString, err := expiredToken.SignedString([]byte(secret))
	assert.NoError(t, err)

	_, err = ValidateJWTToken(expiredTokenString, secret)
	assert.Error(t, err)
}

func TestGenerateAndHashAPIKey(t *testing.T) {
	apiKey, hashedKey, err := GenerateAndHashAPIKey()
	assert.NoError(t, err)
	assert.NotEmpty(t, apiKey)
	assert.NotEmpty(t, hashedKey)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(hashedKey), []byte(apiKey)))
}

func TestGenerateRandomSHA256HASH(t *testing.T) {
	hash, err := GenerateRandomSHA256HASH()
	assert.NoError(t, err)
	assert.Len(t, hash, 64)
}

func TestGenerateJWTToken(t *testing.T) {
	secret := "test-secret"
	userID := uuid.New()
	expiresAt := time.Now().Add(1 * time.Hour)
	tokenString, err := GenerateJWTToken(userID, secret, expiresAt)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)
}
