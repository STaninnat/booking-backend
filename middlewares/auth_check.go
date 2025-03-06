package middlewares

import (
	"log"
	"net/http"
	"time"

	"github.com/STaninnat/booking-backend/internal/config"
	"github.com/STaninnat/booking-backend/security"
	"github.com/golang-jwt/jwt/v5"
)

func HandlerCheckAuth(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("access_token")
		if err != nil {
			log.Println("Couldn't find token error:", err)
			RespondWithJSON(w, http.StatusUnauthorized, map[string]bool{"isAuthenticated": false})
			return
		}

		tokenString := tokenCookie.Value

		claims, err := security.ValidateJWTToken(tokenString, cfg.JWTSecret)
		if err != nil {
			if err == jwt.ErrTokenExpired {
				log.Println("Token expired error:", err)
				RespondWithJSON(w, http.StatusUnauthorized, map[string]bool{"isAuthenticated": false})
				return
			}

			log.Printf("Token validation error: %v\n", err)
			RespondWithJSON(w, http.StatusUnauthorized, map[string]bool{"isAuthenticated": false})
			return
		}

		if claims.ExpiresAt.Before(time.Now()) {
			http.Error(w, "Token expired", http.StatusUnauthorized)
			RespondWithJSON(w, http.StatusUnauthorized, map[string]bool{"isAuthenticated": false})
			return
		}

		RespondWithJSON(w, http.StatusOK, map[string]bool{"isAuthenticated": true})
	}
}
