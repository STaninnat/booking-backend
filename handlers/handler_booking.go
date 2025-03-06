package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/STaninnat/booking-backend/internal/config"
	"github.com/STaninnat/booking-backend/internal/database"
	"github.com/STaninnat/booking-backend/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func HandlerCreateBooking(cfg *config.ApiConfig, w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		CheckIn  string `json:"check_in"`
		CheckOut string `json:"check_out"`
		RoomID   string `json:"room_id"`
		Phone    string `json:"phone"`
	}

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		log.Println("Decode error: ", err)
		return
	}

	layout := "2006-01-02"
	checkInAt, err := time.Parse(layout, params.CheckIn)
	if err != nil {
		log.Println("Invalid check_in format error: ", err)
		return
	}

	checkOutAt, err := time.Parse(layout, params.CheckOut)
	if err != nil {
		log.Println("Invalid check_out format error: ", err)
		return
	}

	if !checkInAt.Before(checkOutAt) {
		log.Println("Check_in must be before check_out")
		return
	}

	exists, err := cfg.DB.CheckRoomAvailability(r.Context(), database.CheckRoomAvailabilityParams{
		RoomID:   params.RoomID,
		CheckIn:  checkOutAt,
		CheckOut: checkInAt,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = cfg.DB.CreateBooking(r.Context(), database.CreateBookingParams{
				ID:        uuid.New().String(),
				CreatedAt: time.Now().Local(),
				UpdatedAt: time.Now().Local(),
				CheckIn:   checkInAt,
				CheckOut:  checkOutAt,
				UserID:    user.ID,
				RoomID:    params.RoomID,
				Phone:     sql.NullString{String: params.Phone, Valid: params.Phone != ""},
			})
			if err != nil {
				log.Println("Couldn't create booking error: ", err)
				return
			}
		} else {
			log.Println("Internal server error: ", err)
			return
		}
	}

	if exists != "" {
		middlewares.RespondWithError(w, http.StatusConflict, "Room is already booked")
		return
	}
	fmt.Println("exists: ", exists)

	userResp := map[string]any{
		"message": "Booking created successfully",
	}

	middlewares.RespondWithJSON(w, http.StatusCreated, userResp)
}

func HandlerGetBookingsByUserID(cfg *config.ApiConfig, w http.ResponseWriter, r *http.Request, user database.User) {
	bookings, err := cfg.DB.GetBookingsByUserID(r.Context(), user.ID)
	if err != nil {
		log.Println("Couldn't get booking by user id error: ", err)
		return
	}

	middlewares.RespondWithJSON(w, http.StatusOK, bookings)
}

func HandlerGetBookingsByRoomID(cfg *config.ApiConfig, w http.ResponseWriter, r *http.Request, user database.User) {
	roomID := chi.URLParam(r, "room_id")
	if roomID == "" {
		log.Println("Missing room id")
	}

	bookings, err := cfg.DB.GetBookingsByRoomID(r.Context(), roomID)
	if err != nil {
		log.Println("Couldn't get booking by room id error: ", err)
		return
	}

	middlewares.RespondWithJSON(w, http.StatusOK, bookings)
}

func HandlerDeleteBooking(cfg *config.ApiConfig, w http.ResponseWriter, r *http.Request, user database.User) {
	booking_id := chi.URLParam(r, "id")
	if booking_id == "" {
		log.Println("Missing booking id")
	}

	err := cfg.DB.DeleteBooking(r.Context(), booking_id)
	if err != nil {
		log.Println("Couldn't delete booking error: ", err)
		return
	}

	userResp := map[string]any{
		"message": "Booking deleted successfully",
	}

	middlewares.RespondWithJSON(w, http.StatusOK, userResp)
}

func HandlerGetAllBookings(cfg *config.ApiConfig, w http.ResponseWriter, r *http.Request, user database.User) {
	type Booking struct {
		ID       string    `json:"id"`
		CheckIn  time.Time `json:"check_in"`
		CheckOut time.Time `json:"check_out"`
		RoomName string    `json:"room_name"`
	}

	bookings, err := cfg.DB.GetAllBookings(r.Context())
	if err != nil {
		if err == sql.ErrNoRows {
			middlewares.RespondWithJSON(w, http.StatusOK, Booking{})
			return
		}
		log.Println("Couldn't get all bookings error: ", err)
		return
	}

	middlewares.RespondWithJSON(w, http.StatusOK, bookings)
}
