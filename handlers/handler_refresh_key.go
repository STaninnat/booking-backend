package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/STaninnat/booking-backend/internal/config"
	"github.com/STaninnat/booking-backend/internal/database"
	"github.com/STaninnat/booking-backend/middlewares"
	"github.com/STaninnat/booking-backend/security"
	"github.com/google/uuid"
)

func HandlerRefreshKey(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("refresh_token")
		if err != nil {
			log.Println("Couldn't find token error:", err)
			return
		}
		refreshToken := cookie.Value

		user, err := cfg.DB.GetUserByRfKey(r.Context(), refreshToken)
		if err != nil {
			log.Println("Couldn't get user error:", err)
			return
		}

		_, newHashedApiKey, err := security.GenerateAndHashAPIKey()
		if err != nil {
			log.Println("Couldn't generate new key error:", err)
			return
		}

		userID, err := uuid.Parse(user.UserID)
		if err != nil {
			log.Printf("Error parsing user ID: %v", err)
			return
		}

		newApiKeyExpiresAt := time.Now().Local().AddDate(0, 3, 0)
		newAccessTokenExpiresAt := time.Now().Local().Add(1 * time.Hour)

		newAccessToken, err := security.GenerateJWTToken(userID, cfg.JWTSecret, newAccessTokenExpiresAt)
		if err != nil {
			log.Println("Couldn't generate new token error:", err)
			return
		}

		err = cfg.DB.UpdateUserKey(r.Context(), database.UpdateUserKeyParams{
			UpdatedAt:       time.Now().Local(),
			ApiKey:          newHashedApiKey,
			ApiKeyExpiresAt: newApiKeyExpiresAt,
			ID:              user.UserID,
		})
		if err != nil {
			log.Println("Failed to update apikey error:", err)
			return
		}

		newRefreshTokenExpiresAt := time.Now().Local().Add(30 * 24 * time.Hour)
		err = cfg.DB.UpdateUserTK(r.Context(), database.UpdateUserTKParams{
			UpdatedAt:             time.Now().Local(),
			AccessTokenExpiresAt:  newAccessTokenExpiresAt,
			RefreshToken:          refreshToken,
			RefreshTokenExpiresAt: newRefreshTokenExpiresAt,
			UserID:                user.UserID,
		})
		if err != nil {
			log.Println("Failed to update refresh token error:", err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    newAccessToken,
			Expires:  newAccessTokenExpiresAt,
			HttpOnly: true,
			Path:     "/",
			Secure:   true,
			// SameSite: http.SameSiteStrictMode,
			SameSite: http.SameSiteLaxMode,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			Expires:  newRefreshTokenExpiresAt,
			HttpOnly: true,
			Path:     "/",
			Secure:   true,
			// SameSite: http.SameSiteStrictMode,
			SameSite: http.SameSiteLaxMode,
		})

		userResp := map[string]string{
			"message": "Token refreshed successfully",
		}

		middlewares.RespondWithJSON(w, http.StatusOK, userResp)
	}
}
