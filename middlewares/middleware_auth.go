package middlewares

import (
	"log"
	"net/http"
	"time"

	"github.com/STaninnat/booking-backend/internal/config"
	"github.com/STaninnat/booking-backend/internal/database"
	"github.com/STaninnat/booking-backend/security"
	"github.com/golang-jwt/jwt/v5"
)

type authhandler func(*config.ApiConfig, http.ResponseWriter, *http.Request, database.User)

func MiddlewareAuth(cfg *config.ApiConfig, handler authhandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := r.Cookie("access_token")
		if err != nil {
			RespondWithError(w, http.StatusUnauthorized, "couldn't find token")
			return
		}

		claims, err := security.ValidateJWTToken(tokenString.Value, cfg.JWTSecret)
		if err != nil {
			log.Printf("token validation error: %v\n", err)
			if err == jwt.ErrTokenExpired {
				RespondWithError(w, http.StatusUnauthorized, "token expired")
				return
			}

			RespondWithError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		user, err := cfg.DB.GetUserByID(r.Context(), claims.UserID.String())
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "couldn't get user")
			return
		}

		if isAPIKeyExpired(user) {
			RespondWithError(w, http.StatusUnauthorized, "api key expired")
			return
		}

		handler(cfg, w, r, user)
	}
}

func isAPIKeyExpired(user database.User) bool {
	return user.ApiKeyExpiresAt.Before(time.Now().Local())
}
