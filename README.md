# booking-be-go

![code coverage badge](https://github.com/STaninnat/booking-backend/actions/workflows/ci.yml/badge.svg)

Backend API for the booking system and handling data.

## Features

- **Environment Configuration**: The application loads environment variables from a `.env` file for configuration
- **User Authentication**
- **Room Management**
- **Booking Management**
- **Middlewares**

## Installation and Tools Used

- **[Go](https://golang.org/dl/)**: The primary language for building the API.
- **[SQLC](https://github.com/sqlc-dev/sqlc/)**: A Go package to generate type-safe Go code from SQL queries.
- **[Goose](https://github.com/pressly/goose/)**: A tool to manage database migrations.
- **[Chi](https://github.com/go-chi/chi/)**: A lightweight, idiomatic HTTP router for Go.
- **[CORS](https://github.com/go-chi/cors/)**: Middleware to handle cross-origin resource sharing.
- **[godotenv](https://github.com/joho/godotenv/)**: A Go package used to load environment variables from a `.env` file.
- **[golang-jwt/jwt](https://github.com/golang-jwt/jwt)**: A library for working with JSON Web Tokens (JWT) for authentication.
- **[google/uuid](https://github.com/google/uuid)**: A package to generate and handle UUIDs.
- **[golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto)**: A collection of cryptographic algorithms and utilities for Go.
- **[github.com/lib/pq v1.10.9](https://pkg.go.dev/github.com/lib/pq@v1.10.9)**: A Go driver for PostgreSQL, used to interact with the database.

## Local Development

Clone the repository

```bash
git clone https://github.com/STaninnat/booking-backend
cd booking-backend
```

Configure environment variables. Copy the .env.example file to .env and fill in the values. You'll need to update values in the .env file to match your configuration.

```bash
cp .env.example .env
```

Run .sh file, but if it doesn't work, make sure to run chmod +x first.

```bash
 ./scripts/run_all.sh
 chmod +x ./run_all.sh
```

Run the server:

```bash
./booking
```

## Notes

- Sorry but, this project requests PostgreSQL for the database.
