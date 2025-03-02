// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: rooms.sql

package database

import (
	"context"
	"database/sql"
	"time"
)

const createRoom = `-- name: CreateRoom :one
INSERT INTO rooms (id, created_at, updated_at, room_name, description, price, max_guests)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, created_at, updated_at, room_name, description, price, max_guests
`

type CreateRoomParams struct {
	ID          string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	RoomName    string
	Description sql.NullString
	Price       string
	MaxGuests   int32
}

func (q *Queries) CreateRoom(ctx context.Context, arg CreateRoomParams) (Room, error) {
	row := q.db.QueryRowContext(ctx, createRoom,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.RoomName,
		arg.Description,
		arg.Price,
		arg.MaxGuests,
	)
	var i Room
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.RoomName,
		&i.Description,
		&i.Price,
		&i.MaxGuests,
	)
	return i, err
}

const getAllRooms = `-- name: GetAllRooms :many
SELECT id, updated_at, room_name, description, price, max_guests
FROM rooms
`

type GetAllRoomsRow struct {
	ID          string
	UpdatedAt   time.Time
	RoomName    string
	Description sql.NullString
	Price       string
	MaxGuests   int32
}

func (q *Queries) GetAllRooms(ctx context.Context) ([]GetAllRoomsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllRooms)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllRoomsRow
	for rows.Next() {
		var i GetAllRoomsRow
		if err := rows.Scan(
			&i.ID,
			&i.UpdatedAt,
			&i.RoomName,
			&i.Description,
			&i.Price,
			&i.MaxGuests,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRoomByID = `-- name: GetRoomByID :one
SELECT id, updated_at, room_name, description, price, max_guests
FROM rooms
WHERE id = $1
`

type GetRoomByIDRow struct {
	ID          string
	UpdatedAt   time.Time
	RoomName    string
	Description sql.NullString
	Price       string
	MaxGuests   int32
}

func (q *Queries) GetRoomByID(ctx context.Context, id string) (GetRoomByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getRoomByID, id)
	var i GetRoomByIDRow
	err := row.Scan(
		&i.ID,
		&i.UpdatedAt,
		&i.RoomName,
		&i.Description,
		&i.Price,
		&i.MaxGuests,
	)
	return i, err
}
