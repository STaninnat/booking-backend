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
			FirstName string `json:"firstname"`
			LastName  string `json:"lastname"`
			Email     string `json:"email"`
			UserName  string `json:"username"`
			Password  string `json:"password"`
		}

		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		if err := decoder.Decode(&params); err != nil {
			log.Println("Decode error: ", err)
			return
		}

		if params.FirstName == "" || params.LastName == "" || params.Email == "" || params.UserName == "" || params.Password == "" {
			middlewares.RespondWithError(w, http.StatusBadRequest, "Invalid input")
			return
		}

		if !security.IsValidUserNameFormat(params.UserName) {
			middlewares.RespondWithError(w, http.StatusBadRequest, "Invalid username format")
			return
		}
		exists, err := cfg.DB.CheckUserExistsByUsername(r.Context(), params.UserName)
		if err != nil {
			log.Println("Checking username error: ", err)
			return
		}
		if exists {
			middlewares.RespondWithError(w, http.StatusBadRequest, "Username already exists")
			return
		}

		if !security.IsValidateEmailFormat(params.Email) {
			middlewares.RespondWithError(w, http.StatusBadRequest, "Invalid email format")
			return
		}

		exists, err = cfg.DB.CheckUserExistsByEmail(r.Context(), params.Email)
		if err != nil {
			log.Println("Checking email error: ", err)
			return
		}
		if exists {
			middlewares.RespondWithError(w, http.StatusBadRequest, "An account with this email already exists")
			return
		}

		fullName := params.FirstName + " " + params.LastName
		exists, err = cfg.DB.CheckUserExistsByFullname(r.Context(), fullName)
		if err != nil {
			log.Println("Checking full name error: ", err)
			return
		}
		if exists {
			middlewares.RespondWithError(w, http.StatusBadRequest, "An account with this name already exists")
			return
		}

		if len(params.Password) < 8 {
			middlewares.RespondWithError(w, http.StatusBadRequest, "Password must be at least 8 ")
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Couldn't hash password error: ", err)
			return
		}

		_, hashedApiKey, err := security.GenerateAndHashAPIKey()
		if err != nil {
			log.Println("Couldn't generate apikey error: ", err)
			return
		}

		apiKeyExpiresAt := time.Now().Local().Add(30 * 24 * time.Hour)

		err = cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
			ID:              uuid.New().String(),
			CreatedAt:       time.Now().Local(),
			UpdatedAt:       time.Now().Local(),
			FullName:        fullName,
			Email:           params.Email,
			Username:        params.UserName,
			Password:        string(hashedPassword),
			ApiKey:          hashedApiKey,
			ApiKeyExpiresAt: apiKeyExpiresAt,
		})
		if err != nil {
			log.Printf("Error while creating user: %v\n", err)
			middlewares.RespondWithError(w, http.StatusInternalServerError, "Couldn't create user")
			return
		}

		jwtExpiresAt := time.Now().Local().Add(15 * time.Minute)

		user, err := cfg.DB.GetUserByKey(r.Context(), hashedApiKey)
		if err != nil {
			log.Printf("Error while getting user: %v\n", err)
			return
		}

		userID, err := uuid.Parse(user.ID)
		if err != nil {
			log.Printf("Error parsing user ID: %v\n", err)
			return
		}

		tokenString, err := security.GenerateJWTToken(userID, cfg.JWTSecret, jwtExpiresAt)
		if err != nil {
			log.Println("Couldn't generate access token error: ", err)
			return
		}

		refreshExpiresAt := time.Now().Local().Add(30 * 24 * time.Hour)
		refreshToken, err := security.GenerateJWTToken(userID, cfg.RefreshSecret, refreshExpiresAt)
		if err != nil {
			log.Println("Couldn't generate refresh token error: ", err)
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
			log.Println("Failed to create new refresh token error: ", err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    tokenString,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			Expires:  jwtExpiresAt,
			// SameSite: http.SameSiteStrictMode,
			SameSite: http.SameSiteLaxMode,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			Expires:  refreshExpiresAt,
			// SameSite: http.SameSiteStrictMode,
			SameSite: http.SameSiteLaxMode,
		})

		userResp := map[string]string{
			"message": "User created successfully",
		}

		middlewares.RespondWithJSON(w, http.StatusCreated, userResp)
	}
}
