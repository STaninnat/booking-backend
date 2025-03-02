package security

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"

	"slices"

	"github.com/golang-jwt/jwt/v5"
)

func IsValidUserNameFormat(name string) bool {
	usernameRegex := `^[a-zA-Z0-9]+([-._]?[a-zA-Z0-9]+)*$`

	re := regexp.MustCompile(usernameRegex)

	return len(name) >= 3 && len(name) <= 30 && re.MatchString(name)
}

func IsValidateEmailFormat(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	re := regexp.MustCompile(emailRegex)

	return re.MatchString(email)
}

func ValidateJWTToken(tokenString string, secret string) (*Claims, error) {
	claims := &Claims{}

	apiServiceName := os.Getenv("API_SERVICE_NAME")
	frontendAppName := os.Getenv("FRONTEND_APP_NAME")

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not parse token: %w", err)
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.Issuer != apiServiceName {
		return nil, fmt.Errorf("invalid issuer: expected 'my-api-service', got '%s'", claims.Issuer)
	}

	if !contains(claims.Audience, frontendAppName) {
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
