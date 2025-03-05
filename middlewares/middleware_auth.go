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
			log.Println("Couldn't find token error: ", err)
			return
		}

		claims, err := security.ValidateJWTToken(tokenString.Value, cfg.JWTSecret)
		if err != nil {
			log.Printf("token validation error: %v\n", err)
			if err == jwt.ErrTokenExpired {
				log.Println("Token expired error: ", err)
				return
			}

			RespondWithError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		user, err := cfg.DB.GetUserByID(r.Context(), claims.UserID.String())
		if err != nil {
			log.Println("Couldn't get user error: ", err)
			return
		}

		if isAPIKeyExpired(user) {
			log.Println("Api key expired error: ", err)
			return
		}

		handler(cfg, w, r, user)
	}
}

func isAPIKeyExpired(user database.User) bool {
	return user.ApiKeyExpiresAt.Before(time.Now().Local())
}
