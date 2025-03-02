-- +goose Up
CREATE TABLE
    bookings (
        id TEXT PRIMARY KEY,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        check_in TIMESTAMP NOT NULL,
        check_out TIMESTAMP NOT NULL,
        user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        room_id TEXT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
        UNIQUE (room_id, check_in, check_out)
    );

-- +goose Down
DROP TABLE IF EXISTS bookings;