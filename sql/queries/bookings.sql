-- name: CreateBooking :one
INSERT INTO bookings (id, created_at, updated_at, check_in, check_out, user_id, room_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: CheckRoomAvailability :one
SELECT id FROM bookings
WHERE room_id = $1
AND (
    (check_in < $2 AND check_out >= $3)
    OR (check_in >= $3 AND check_out <= $2)
)
LIMIT 1;

-- name: GetBookingsByUserID :many
SELECT b.id, b.updated_at, b.check_in, b.check_out, b.user_id, b.room_id, r.room_name
FROM bookings b
JOIN rooms r ON b.room_id = r.id
WHERE b.user_id = $1;

-- name: GetBookingsByRoomID :many
SELECT b.id, b.updated_at, b.check_in, b.check_out, b.user_id, u.email AS user_email
FROM bookings b
JOIN users u ON b.user_id = u.id
WHERE b.room_id = $1;

-- name: GetBookedDatesByRoomID :many
SELECT check_in, check_out
FROM bookings
WHERE room_id = $1;

-- name: DeleteBooking :exec
DELETE FROM bookings
WHERE id = $1;