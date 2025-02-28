package handlers

import (
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

func HandlerCreateUser(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			FullName string `json:"fullname"`
			LastName string `json:"lastname"`
			UserName string `json:"username"`
			Password string `json:"password"`
		}

		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		if err := decoder.Decode(&params); err != nil {
			middlewares.RespondWithError(w, http.StatusBadRequest, "couldn't decode parameters")
			return
		}

		if params.UserName == "" || params.Password == "" {
			middlewares.RespondWithError(w, http.StatusBadRequest, "invalid input")
			return
		}

		if !security.IsValidUserName(params.UserName) {
			middlewares.RespondWithError(w, http.StatusBadRequest, "invalid username format")
			return
		}

		_, err := cfg.DB.GetUserByName(r.Context(), params.UserName)
		if err == nil {
			middlewares.RespondWithError(w, http.StatusBadRequest, "username already exists")
			return
		}

		if len(params.Password) < 8 {
			middlewares.RespondWithError(w, http.StatusBadRequest, "password must be at least 8 ")
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
		if err != nil {
			middlewares.RespondWithError(w, http.StatusInternalServerError, "couldn't hash password")
			return
		}

		_, hashedApiKey, err := security.GenerateAndHashAPIKey()
		if err != nil {
			middlewares.RespondWithError(w, http.StatusInternalServerError, "couldn't generate apikey")
			return
		}

		apiKeyExpiresAt := time.Now().Local().Add(30 * 24 * time.Hour)

		err = cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
			ID:              uuid.New().String(),
			CreatedAt:       time.Now().Local(),
			UpdatedAt:       time.Now().Local(),
			FullName:        params.FullName,
			LastName:        params.LastName,
			Username:        params.UserName,
			Password:        string(hashedPassword),
			ApiKey:          hashedApiKey,
			ApiKeyExpiresAt: apiKeyExpiresAt,
		})
		if err != nil {
			log.Printf("error while creating user: %v", err)
			middlewares.RespondWithError(w, http.StatusInternalServerError, "couldn't create user")
			return
		}

		jwtExpiresAt := time.Now().Local().Add(15 * time.Minute)

		user, err := cfg.DB.GetUser(r.Context(), hashedApiKey)
		if err != nil {
			log.Printf("error while getting user: %v", err)
			middlewares.RespondWithError(w, http.StatusInternalServerError, "couldn't get user")
			return
		}

		userID, err := uuid.Parse(user.ID)
		if err != nil {
			log.Printf("error parsing user ID: %v", err)
			middlewares.RespondWithError(w, http.StatusInternalServerError, "invalid user ID")
			return
		}

		tokenString, err := security.GenerateJWTToken(userID, cfg.JWTSecret, jwtExpiresAt)
		if err != nil {
			middlewares.RespondWithError(w, http.StatusInternalServerError, "couldn't generate access token")
			return
		}

		refreshExpiresAt := time.Now().Local().Add(30 * 24 * time.Hour)
		refreshToken, err := security.GenerateJWTToken(userID, cfg.RefreshSecret, refreshExpiresAt)
		if err != nil {
			middlewares.RespondWithError(w, http.StatusInternalServerError, "couldn't generate refresh token")
			return
		}

		err = cfg.DB.CreateUserRfKey(r.Context(), database.CreateUserRfKeyParams{
			ID:                    uuid.New().String(),
			CreatedAt:             time.Now().Local(),
			UpdatedAt:             time.Now().Local(),
			AccessTokenExpiresAt:  jwtExpiresAt,
			RefreshToken:          refreshToken,
			RefreshTokenExpiresAt: refreshExpiresAt,
			UserID:                user.ID,
		})
		if err != nil {
			middlewares.RespondWithError(w, http.StatusInternalServerError, "failed to create new refresh token")
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    tokenString,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			Expires:  jwtExpiresAt,
			SameSite: http.SameSiteStrictMode,
			// SameSite: http.SameSiteLaxMode,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			Expires:  refreshExpiresAt,
			SameSite: http.SameSiteStrictMode,
			// SameSite: http.SameSiteLaxMode,
		})

		userResp := map[string]string{
			"message": "User created successfully",
		}

		middlewares.RespondWithJSON(w, http.StatusCreated, userResp)
	}
}
