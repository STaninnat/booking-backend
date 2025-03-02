-- +goose Up
CREATE TABLE
    users (
        id TEXT PRIMARY KEY,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        full_name TEXT NOT NULL,
        email TEXT NOT NULL,
        phone TEXT,
        username TEXT NOT NULL,
        password TEXT NOT NULL,
        api_key TEXT UNIQUE NOT NULL,
        api_key_expires_at TIMESTAMP  NOT NULL
    );

-- +goose Down
DROP TABLE IF EXISTS users;