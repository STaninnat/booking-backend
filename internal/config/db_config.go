package config

import "github.com/STaninnat/booking-backend/internal/database"

type ApiConfig struct {
	DB            *database.Queries
	JWTSecret     string
	RefreshSecret string
}
