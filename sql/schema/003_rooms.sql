-- +goose Up
CREATE TABLE
    rooms (
        id TEXT PRIMARY KEY,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        room_name TEXT NOT NULL,
        description TEXT,
        price NUMERIC(10,2) NOT NULL,
        max_guests INT NOT NULL
    );

-- +goose Down
DROP TABLE IF EXISTS rooms;