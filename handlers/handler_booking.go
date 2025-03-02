package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/STaninnat/booking-backend/internal/config"
	"github.com/STaninnat/booking-backend/internal/database"
	"github.com/STaninnat/booking-backend/internal/models"
	"github.com/STaninnat/booking-backend/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func HandlerCreateBooking(cfg *config.ApiConfig, w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		CheckIn  string `json:"check_in"`
		CheckOut string `json:"check_out"`
		RoomID   string `json:"room_id"`
	}

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		middlewares.RespondWithError(w, http.StatusBadRequest, "couldn't decode json")
		return
	}

	layout := "2006-01-02 15:04:05"
	checkInAt, err := time.Parse(layout, params.CheckIn+" 14:00:00")
	if err != nil {
		middlewares.RespondWithError(w, http.StatusBadRequest, "invalid check_in format")
		return
	}

	checkOutAt, err := time.Parse(layout, params.CheckOut+" 12:00:00")
	if err != nil {
		middlewares.RespondWithError(w, http.StatusBadRequest, "invalid check_out format")
		return
	}

	if !checkInAt.Before(checkOutAt) {
		http.Error(w, "check_in must be before check_out", http.StatusBadRequest)
		return
	}

	exists, err := cfg.DB.CheckRoomAvailability(r.Context(), database.CheckRoomAvailabilityParams{
		RoomID:   params.RoomID,
		CheckIn:  checkOutAt,
		CheckOut: checkInAt,
	})

	var booking_db database.Booking
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			booking_db, err = cfg.DB.CreateBooking(r.Context(), database.CreateBookingParams{
				ID:        uuid.New().String(),
				CreatedAt: time.Now().Local(),
				UpdatedAt: time.Now().Local(),
				CheckIn:   checkInAt,
				CheckOut:  checkOutAt,
				UserID:    user.ID,
				RoomID:    params.RoomID,
			})
			if err != nil {
				middlewares.RespondWithError(w, http.StatusInternalServerError, "couldn't create booking")
				return
			}
		} else {
			middlewares.RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	if exists != "" {
		middlewares.RespondWithError(w, http.StatusConflict, "room is already booked")
		return
	}

	userResp := map[string]any{
		"message": "Booking created successfully",
		"room":    models.DBBookingToBooking(booking_db),
	}

	middlewares.RespondWithJSON(w, http.StatusCreated, userResp)
}

func HandlerGetBookingsByUserID(cfg *config.ApiConfig, w http.ResponseWriter, r *http.Request, user database.User) {
	bookings, err := cfg.DB.GetBookingsByUserID(r.Context(), user.ID)
	if err != nil {
		middlewares.RespondWithError(w, http.StatusInternalServerError, "couldn't get booking by user id")
		return
	}

	userResp := map[string]any{
		"message": "Got all bookings successfully",
		"rooms":   bookings,
	}

	middlewares.RespondWithJSON(w, http.StatusOK, userResp)
}

func HandlerGetBookingsByRoomID(cfg *config.ApiConfig, w http.ResponseWriter, r *http.Request, user database.User) {
	roomID := chi.URLParam(r, "room_id")
	if roomID == "" {
		middlewares.RespondWithError(w, http.StatusBadRequest, "missing room_id")
	}

	bookings, err := cfg.DB.GetBookingsByRoomID(r.Context(), roomID)
	if err != nil {
		middlewares.RespondWithError(w, http.StatusInternalServerError, "couldn't get booking by room id")
		return
	}

	userResp := map[string]any{
		"message": "Got all bookings successfully",
		"rooms":   bookings,
	}

	middlewares.RespondWithJSON(w, http.StatusOK, userResp)
}

func HandlerDeleteBooking(cfg *config.ApiConfig, w http.ResponseWriter, r *http.Request, user database.User) {
	booking_id := chi.URLParam(r, "id")
	if booking_id == "" {
		middlewares.RespondWithError(w, http.StatusBadRequest, "missing booking_id")
	}

	err := cfg.DB.DeleteBooking(r.Context(), booking_id)
	if err != nil {
		middlewares.RespondWithError(w, http.StatusInternalServerError, "couldn't delete booking")
		return
	}

	userResp := map[string]any{
		"message": "Booking deleted successfully",
	}

	middlewares.RespondWithJSON(w, http.StatusOK, userResp)
}
