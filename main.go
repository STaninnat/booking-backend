package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/STaninnat/booking-backend/handlers"
	"github.com/STaninnat/booking-backend/internal/config"
	"github.com/STaninnat/booking-backend/internal/database"
	"github.com/STaninnat/booking-backend/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(".env.development"); err != nil {
		log.Printf("warning: assuming default configuration. '.env.development' unreadable: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Println("warning: PORT environment variable is not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("warning: JWT_SECRET environment variable is not set")
	}

	refreshSecret := os.Getenv("REFRESH_SECRET")
	if refreshSecret == "" {
		log.Fatal("warning: REFRESH_SECRET environment variable is not set")
	}

	apicfg := config.ApiConfig{
		JWTSecret:     jwtSecret,
		RefreshSecret: refreshSecret,
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("warning: DATABASE_URL environment variable is not set")
		log.Println("running without CRUD endpoints")
	} else {
		db, err := sql.Open("postgres", dbURL)
		if err != nil {
			log.Fatalf("warning: can't connect to database: %v", err)
		}

		if err := db.Ping(); err != nil {
			log.Fatalf("failed to ping database: %v", err)
		}

		dbQueries := database.New(db)
		apicfg.DB = dbQueries
		apicfg.DBConn = db
		log.Println("Connected to database successfully!")
	}

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	if apicfg.DB != nil {
		v1Router.Get("/healthz", handlers.HandlerReadiness)
		v1Router.Get("/error", handlers.HandlerError)

		v1Router.Post("/user/signup", handlers.HandlerCreateUser(&apicfg))
		v1Router.Post("/user/signin", handlers.HandlerSignin(&apicfg))
		v1Router.Post("/user/signout", middlewares.MiddlewareAuth(&apicfg, handlers.HandlerSignout))
		v1Router.Post("/user/refresh-key", handlers.HandlerRefreshKey(&apicfg))

		v1Router.Post("/rooms", middlewares.MiddlewareAuth(&apicfg, handlers.HandlerCreateRoom))
		v1Router.Get("/rooms", middlewares.MiddlewareAuth(&apicfg, handlers.HandlerGetAllRooms))
		v1Router.Get("/rooms/{id}", middlewares.MiddlewareAuth(&apicfg, handlers.HandlerGetRoom))

		v1Router.Post("/bookings", middlewares.MiddlewareAuth(&apicfg, handlers.HandlerCreateBooking))
		v1Router.Get("/bookings/user/{user_id}", middlewares.MiddlewareAuth(&apicfg, handlers.HandlerGetBookingsByUserID))
		v1Router.Get("/bookings/room/{room_id}", middlewares.MiddlewareAuth(&apicfg, handlers.HandlerGetBookingsByRoomID))
		v1Router.Delete("/bookings/{id}", middlewares.MiddlewareAuth(&apicfg, handlers.HandlerDeleteBooking))

	}

	router.Mount("/v1", v1Router)
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Serving on port: %s\n", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
