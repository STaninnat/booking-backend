-- +goose Up
CREATE TABLE 
    users_key (
        id TEXT PRIMARY KEY,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        access_token_expires_at TIMESTAMP  NOT NULL,
        refresh_token TEXT UNIQUE NOT NULL,
        refresh_token_expires_at TIMESTAMP  NOT NULL,
        user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE
    );

-- +goose Down
DROP TABLE IF EXISTS users_key;