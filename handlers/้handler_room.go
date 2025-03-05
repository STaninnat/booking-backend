package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/STaninnat/booking-backend/internal/config"
	"github.com/STaninnat/booking-backend/internal/database"
	"github.com/STaninnat/booking-backend/internal/models"
	"github.com/STaninnat/booking-backend/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CalendarResponse struct {
	RoomID      string   `json:"room_id"`
	BookedDates []string `json:"booked_dates"`
}

type BookedDate struct {
	CheckIn  time.Time `json:"check_in"`
	CheckOut time.Time `json:"check_out"`
}

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
		log.Println("Decode error: ", err)
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
		middlewares.RespondWithError(w, http.StatusInternalServerError, "Couldn't create room")
		return
	}

	middlewares.RespondWithJSON(w, http.StatusCreated, models.DBRoomToRoom(room_db))
}

func HandlerGetAllRooms(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rooms, err := cfg.DB.GetAllRooms(r.Context())
		if err != nil {
			log.Println("Couldn't get all rooms error: ", err)
			return
		}

		middlewares.RespondWithJSON(w, http.StatusOK, rooms)
	}
}

func HandlerGetRoom(cfg *config.ApiConfig, w http.ResponseWriter, r *http.Request, user database.User) {
	roomID := chi.URLParam(r, "id")
	if roomID == "" {
		log.Println("Missing room id")
	}

	room, err := cfg.DB.GetRoomByID(r.Context(), roomID)
	if err != nil {
		middlewares.RespondWithError(w, http.StatusNotFound, "Couldn't find room")
		return
	}

	middlewares.RespondWithJSON(w, http.StatusOK, room)
}

func HandlerGetRoomCalendar(cfg *config.ApiConfig, w http.ResponseWriter, r *http.Request, user database.User) {
	roomID := chi.URLParam(r, "room_id")
	if roomID == "" {
		log.Println("Missing room id")
	}

	bookings, err := cfg.DB.GetBookedDatesByRoomID(r.Context(), roomID)
	if err != nil {
		if err == sql.ErrNoRows {
			middlewares.RespondWithJSON(w, http.StatusOK, CalendarResponse{RoomID: roomID, BookedDates: []string{}})
			return
		}
		log.Println("Failed to get booking error: ")
		return
	}

	var bookedDatesInput []BookedDate
	for _, b := range bookings {
		bookedDatesInput = append(bookedDatesInput, BookedDate{
			CheckIn:  b.CheckIn,
			CheckOut: b.CheckOut,
		})
	}
	bookedDates := generateBookedDates(bookedDatesInput)

	response := CalendarResponse{
		RoomID:      roomID,
		BookedDates: bookedDates,
	}

	middlewares.RespondWithJSON(w, http.StatusOK, response)
}

func generateBookedDates(bookings []BookedDate) []string {
	dateMap := make(map[string]bool)

	for _, b := range bookings {
		for i := b.CheckIn; !i.After(b.CheckOut); i = i.AddDate(0, 0, 1) {
			dateMap[i.Format("2006-01-02")] = true
		}
		dateMap[b.CheckOut.Format("2006-01-02")] = true
	}

	var bookedDates []string
	for date := range dateMap {
		bookedDates = append(bookedDates, date)
	}

	sort.Strings(bookedDates)
	return bookedDates
}
