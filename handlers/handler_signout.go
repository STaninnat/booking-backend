package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/STaninnat/booking-backend/internal/config"
	"github.com/STaninnat/booking-backend/internal/database"
	"github.com/STaninnat/booking-backend/middlewares"
	"github.com/google/uuid"
)

func HandlerSignout(cfg *config.ApiConfig, w http.ResponseWriter, r *http.Request, user database.User) {
	newKeyExpiredAt := time.Now().Local().AddDate(-1, 0, 0)
	newTokenExpired := "expired-" + uuid.New().String()[:28]

	if err := cfg.DB.UpdateUserKey(r.Context(), database.UpdateUserKeyParams{
		UpdatedAt:       time.Now().Local(),
		ApiKey:          newTokenExpired,
		ApiKeyExpiresAt: newKeyExpiredAt,
		ID:              user.ID,
	}); err != nil {
		log.Println("Couldn't update user and signout error: ", err)
		return
	}

	if err := cfg.DB.UpdateUserTK(r.Context(), database.UpdateUserTKParams{
		UpdatedAt:             time.Now().Local(),
		AccessTokenExpiresAt:  newKeyExpiredAt,
		RefreshToken:          newTokenExpired,
		RefreshTokenExpiresAt: newKeyExpiredAt,
		UserID:                user.ID,
	}); err != nil {
		log.Println("Couldn't update user key and signout error: ", err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  newKeyExpiredAt,
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		// SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  newKeyExpiredAt,
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		// SameSite: http.SameSiteLaxMode,
	})

	resp := map[string]string{
		"message": "Signed out successfully",
	}

	middlewares.RespondWithJSON(w, http.StatusOK, resp)
}
