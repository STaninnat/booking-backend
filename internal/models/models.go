package models

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/STaninnat/booking-backend/internal/database"
)

type Room struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	RoomName    string    `json:"room_name"`
	Description *string   `json:"description"`
	Price       float64   `json:"price"`
	MaxGuests   int       `json:"max_guests"`
}

func DBRoomToRoom(room database.Room) Room {
	price, err := strconv.ParseFloat(room.Price, 64)
	if err != nil {
		log.Printf("couldn't parse to float: %v\n", err)
		return Room{}
	}

	return Room{
		ID:          room.ID,
		CreatedAt:   room.CreatedAt,
		UpdatedAt:   room.UpdatedAt,
		RoomName:    room.RoomName,
		Description: nullStringToStringPtr(room.Description),
		Price:       price,
		MaxGuests:   int(room.MaxGuests),
	}
}

type Booking struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CheckIn   time.Time `json:"check_in"`
	CheckOut  time.Time `json:"check_out"`
	UserID    string    `json:"user_id"`
	RoomID    string    `json:"room_id"`
}

func DBBookingToBooking(booking database.Booking) Booking {
	return Booking{
		ID:        booking.ID,
		CreatedAt: booking.CreatedAt,
		UpdatedAt: booking.UpdatedAt,
		CheckIn:   booking.CheckIn,
		CheckOut:  booking.CheckOut,
		UserID:    booking.UserID,
		RoomID:    booking.RoomID,
	}
}

func nullStringToStringPtr(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}
