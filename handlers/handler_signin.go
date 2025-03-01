package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/STaninnat/booking-backend/internal/config"
	"github.com/STaninnat/booking-backend/internal/database"
	"github.com/STaninnat/booking-backend/middlewares"
	"github.com/STaninnat/booking-backend/security"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HandlerSignin(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			UserName string `json:"username"`
			Password string `json:"password"`
		}

		defer r.Body.Close()
		params := parameters{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&params); err != nil {
			middlewares.RespondWithError(w, http.StatusBadRequest, "couldn't decode parameters")
			return
		}

		user, err := cfg.DB.GetUserByUsername(r.Context(), params.UserName)
		if err != nil {
			if err == sql.ErrNoRows {
				middlewares.RespondWithError(w, http.StatusBadRequest, "username not found")
			} else {
				middlewares.RespondWithError(w, http.StatusInternalServerError, "error retrieving user")
			}
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
		if err != nil {
			middlewares.RespondWithError(w, http.StatusBadRequest, "incorrect password")
			return
		}

		if user.ApiKeyExpiresAt.Before(time.Now().Local()) {
			middlewares.RespondWithError(w, http.StatusUnauthorized, "apikey expired")
			return
		}

		jwtExpiresAt := time.Now().Local().Add(1 * time.Hour)

		userID, err := uuid.Parse(user.ID)
		if err != nil {
			log.Printf("error parsing user id: %v", err)
			middlewares.RespondWithError(w, http.StatusInternalServerError, "invalid user ID")
			return
		}

		tokenString, err := security.GenerateJWTToken(userID, cfg.JWTSecret, jwtExpiresAt)
		if err != nil {
			middlewares.RespondWithError(w, http.StatusInternalServerError, "couldn't generate access token")
			return
		}

		tx, err := cfg.DBConn.BeginTx(r.Context(), nil)
		if err != nil {
			middlewares.RespondWithError(w, http.StatusInternalServerError, "failed to start transaction")
			return
		}

		defer func() {
			if p := recover(); p != nil {
				if err := tx.Rollback(); err != nil {
					log.Printf("failed to rollback transaction: %v", err)
				}
				panic(p)
			} else if err != nil {
				if err := tx.Rollback(); err != nil {
					log.Printf("failed to rollback transaction: %v", err)
				}
			} else {
				err = tx.Commit()
				if err != nil {
					log.Printf("failed to commit transaction: %v", err)
					middlewares.RespondWithError(w, http.StatusInternalServerError, "failed to commit transaction")
					return
				}
			}
		}()

		queriesTx := database.New(tx)

		_, hashedApiKey, err := security.GenerateAndHashAPIKey()
		if err != nil {
			middlewares.RespondWithError(w, http.StatusInternalServerError, "couldn't generate apikey")
			return
		}

		keyExpiresAt := time.Now().Local().Add(30 * 24 * time.Hour)

		err = queriesTx.UpdateUserKey(r.Context(), database.UpdateUserKeyParams{
			UpdatedAt:       time.Now().Local(),
			ApiKey:          hashedApiKey,
			ApiKeyExpiresAt: keyExpiresAt,
			ID:              user.ID,
		})
		if err != nil {
			middlewares.RespondWithError(w, http.StatusInternalServerError, "failed to update new api key")
			return
		}

		refreshToken, err := security.GenerateJWTToken(userID, cfg.RefreshSecret, keyExpiresAt)
		if err != nil {
			middlewares.RespondWithError(w, http.StatusInternalServerError, "couldn't generate refresh token")
			return
		}

		err = queriesTx.UpdateUserToken(r.Context(), database.UpdateUserTokenParams{
			UpdatedAt:             time.Now().Local(),
			AccessTokenExpiresAt:  jwtExpiresAt,
			RefreshToken:          refreshToken,
			RefreshTokenExpiresAt: keyExpiresAt,
			UserID:                user.ID,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				err = queriesTx.CreateUserRfKey(r.Context(), database.CreateUserRfKeyParams{
					ID:                    uuid.New().String(),
					CreatedAt:             time.Now().Local(),
					UpdatedAt:             time.Now().Local(),
					AccessTokenExpiresAt:  jwtExpiresAt,
					RefreshToken:          refreshToken,
					RefreshTokenExpiresAt: keyExpiresAt,
					UserID:                user.ID,
				})
				if err != nil {
					middlewares.RespondWithError(w, http.StatusInternalServerError, "failed to create new refresh token")
					return
				}
			} else {
				middlewares.RespondWithError(w, http.StatusInternalServerError, "failed to update new refresh token")
				return
			}
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    tokenString,
			Expires:  jwtExpiresAt,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			SameSite: http.SameSiteStrictMode,
			// SameSite: http.SameSiteLaxMode,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			Expires:  keyExpiresAt,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			SameSite: http.SameSiteStrictMode,
			// SameSite: http.SameSiteLaxMode,
		})

		userResp := map[string]string{
			"message": "Signed in successfully",
		}

		middlewares.RespondWithJSON(w, http.StatusOK, userResp)
	}
}
