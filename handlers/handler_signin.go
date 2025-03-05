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
			log.Println("Decode error: ", err)
			return
		}

		user, err := cfg.DB.GetUserByUsername(r.Context(), params.UserName)
		if err != nil {
			if err == sql.ErrNoRows {
				middlewares.RespondWithError(w, http.StatusBadRequest, "Username not found")
			} else {
				log.Println("Retrieving user error: ", err)
			}
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
		if err != nil {
			middlewares.RespondWithError(w, http.StatusBadRequest, "Incorrect password")
			return
		}

		if user.ApiKeyExpiresAt.Before(time.Now().Local()) {
			log.Println("Apikey expired error")
			return
		}

		jwtExpiresAt := time.Now().Local().Add(1 * time.Hour)

		userID, err := uuid.Parse(user.ID)
		if err != nil {
			log.Printf("Error parsing user id: %v\n", err)
			return
		}

		tokenString, err := security.GenerateJWTToken(userID, cfg.JWTSecret, jwtExpiresAt)
		if err != nil {
			log.Println("Couldn't generate access token error: ", err)
			return
		}

		tx, err := cfg.DBConn.BeginTx(r.Context(), nil)
		if err != nil {
			log.Println("Failed to start transaction error: ", err)
			return
		}

		defer func() {
			if p := recover(); p != nil {
				if err := tx.Rollback(); err != nil {
					log.Printf("Failed to rollback transaction: %v\n", err)
				}
				panic(p)
			} else if err != nil {
				if err := tx.Rollback(); err != nil {
					log.Printf("Failed to rollback transaction: %v\n", err)
				}
			} else {
				err = tx.Commit()
				if err != nil {
					log.Printf("Failed to commit transaction: %v\n", err)
					return
				}
			}
		}()

		queriesTx := database.New(tx)

		_, hashedApiKey, err := security.GenerateAndHashAPIKey()
		if err != nil {
			log.Println("Couldn't generate apikey error: ", err)
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
			log.Println("Failed to update new apikey error: ", err)
			return
		}

		refreshToken, err := security.GenerateJWTToken(userID, cfg.RefreshSecret, keyExpiresAt)
		if err != nil {
			log.Println("Couldn't generate refresh token error: ", err)
			return
		}

		err = queriesTx.UpdateUserTK(r.Context(), database.UpdateUserTKParams{
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
					log.Println("Failed to create new refresh token error: ", err)
					return
				}
			} else {
				log.Println("Failed to update new refresh token error: ", err)
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
