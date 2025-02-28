package config

import (
	"database/sql"

	"github.com/STaninnat/booking-backend/internal/database"
)

type ApiConfig struct {
	DB            *database.Queries
	DBConn        *sql.DB
	JWTSecret     string
	RefreshSecret string
}
