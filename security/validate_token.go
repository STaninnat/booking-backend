package security

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"slices"

	"github.com/golang-jwt/jwt/v5"
)

func IsValidUserName(name string) bool {
	var usernameRegex = `^[a-zA-Z0-9]+([-._]?[a-zA-Z0-9]+)*$`

	re := regexp.MustCompile(usernameRegex)

	return len(name) >= 3 && len(name) <= 30 && re.MatchString(name)
}

func ValidateJWTToken(tokenString string, secret string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not parse token: %w", err)
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.Issuer != "my-api-service" {
		return nil, fmt.Errorf("invalid issuer: expected 'my-api-service', got '%s'", claims.Issuer)
	}

	if !contains(claims.Audience, "my-frontend-app") {
		return nil, fmt.Errorf("invalid audience: expected 'my-frontend-app', got '%s'", claims.Audience)
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("token expired")
	}

	if claims.NotBefore.After(time.Now()) {
		return nil, fmt.Errorf("token not valid yet")
	}

	return claims, nil
}

func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}
