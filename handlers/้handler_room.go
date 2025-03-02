package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/STaninnat/booking-backend/internal/config"
	"github.com/STaninnat/booking-backend/internal/database"
	"github.com/STaninnat/booking-backend/internal/models"
	"github.com/STaninnat/booking-backend/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func HandlerCreateRoom(cfg *config.ApiConfig, w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		RoomName    string  `json:"room_name"`
		Description *string `json:"description"`
		Price       float64 `json:"price"`
		MaxGuests   int32   `json:"max_guests"`
	}

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		middlewares.RespondWithError(w, http.StatusBadRequest, "couldn't decode json")
		return
	}

	description := sql.NullString{
		String: "",
		Valid:  false,
	}
	if params.Description != nil {
		description.String = *params.Description
		description.Valid = true
	}

	room_db, err := cfg.DB.CreateRoom(r.Context(), database.CreateRoomParams{
		ID:          uuid.New().String(),
		CreatedAt:   time.Now().Local(),
		UpdatedAt:   time.Now().Local(),
		RoomName:    params.RoomName,
		Description: description,
		Price:       fmt.Sprintf("%.2f", params.Price),
		MaxGuests:   int32(params.MaxGuests),
	})
	if err != nil {
		middlewares.RespondWithError(w, http.StatusInternalServerError, "couldn't create room")
		return
	}

	userResp := map[string]any{
		"message": "Room created successfully",
		"room":    models.DBRoomToRoom(room_db),
	}

	middlewares.RespondWithJSON(w, http.StatusCreated, userResp)
}

func HandlerGetAllRooms(cfg *config.ApiConfig, w http.ResponseWriter, r *http.Request, user database.User) {
	rooms, err := cfg.DB.GetAllRooms(r.Context())
	if err != nil {
		middlewares.RespondWithError(w, http.StatusInternalServerError, "couldn't get all rooms")
		return
	}

	userResp := map[string]any{
		"message": "Got all rooms successfully",
		"rooms":   rooms,
	}

	middlewares.RespondWithJSON(w, http.StatusOK, userResp)
}

func HandlerGetRoom(cfg *config.ApiConfig, w http.ResponseWriter, r *http.Request, user database.User) {
	roomID := chi.URLParam(r, "id")
	room, err := cfg.DB.GetRoomByID(r.Context(), roomID)
	if err != nil {
		middlewares.RespondWithError(w, http.StatusNotFound, "couldn't find room")
		return
	}

	userResp := map[string]any{
		"message": "Room was found",
		"rooms":   room,
	}

	middlewares.RespondWithJSON(w, http.StatusOK, userResp)
}
