-- name: CreateRoom :exec
INSERT INTO rooms (id, created_at, updated_at, room_name, description, price, max_guests)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;
--

-- name: GetAllRooms :one
SELECT id, updated_at, room_name, description, price, max_guests
FROM rooms;
--

-- name: GetRoomByID :one
SELECT id, updated_at, room_name, description, price, max_guests
FROM rooms
WHERE id = $1;
--