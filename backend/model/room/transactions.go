package room

import (
	"backend/utils"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoomStore struct {
	Pool *pgxpool.Pool
}

func checkRoomExistsDB(ctx context.Context, tx pgx.Tx, rid string) (bool, *utils.HTTPError) {
	var exists bool
	if err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM rooms WHERE room_id = $1)", rid).Scan(&exists); err != nil {
		fmt.Println(err.Error())
		return exists, &utils.HTTPError{
			Code:    500,
			Message: "Internal server error a",
		}
	}

	return exists, nil
}

func insertRoomDB(ctx context.Context, tx pgx.Tx, r *Room) *utils.HTTPError {
	if _, err := tx.Exec(ctx, "INSERT INTO rooms (room_id, host) VALUES ($1, $2)", r.ID, r.Host); err != nil {
		fmt.Println(err.Error())
		return &utils.HTTPError{
			Code:    500,
			Message: "Internal server error b",
		}
	}

	return nil
}

// come up with better name for this
func (rs *RoomStore) InsertRoom(r *Room) *utils.HTTPError {
	tx, err := rs.Pool.Begin(context.TODO())
	if err != nil {
		fmt.Println(err.Error())
		return &utils.HTTPError{
			Code:    500,
			Message: "Internal server error c",
		}
	}
	defer tx.Rollback(context.TODO())

	exists, httpErr := checkRoomExistsDB(context.TODO(), tx, r.ID)
	if httpErr != nil {
		return httpErr
	}

	if exists {
		return &utils.HTTPError{
			Code:    400,
			Message: "Room already exists",
		}
	}

	httpErr = insertRoomDB(context.TODO(), tx, r)
	if httpErr != nil {
		return httpErr
	}

	if err := tx.Commit(context.TODO()); err != nil {
		return &utils.HTTPError{
			Code:    500,
			Message: "Internal server error d",
		}
	}

	return nil
}

func (rs *RoomStore) SelectRoom(rid string) (*Room, *utils.HTTPError) {
	r := new(Room)

	if err := rs.Pool.QueryRow(context.TODO(), "SELECT room_id, host FROM rooms WHERE room_id = $1", rid).Scan(&r.ID, &r.Host); err == pgx.ErrNoRows {
		return nil, &utils.HTTPError{
			Code:    400,
			Message: "Room does not exist",
		}
	} else if err != nil {
		return nil, &utils.HTTPError{
			Code:    500,
			Message: "Internal server error",
		}
	}

	return r, nil
}

func (rs *RoomStore) SelectHost(rid string) (string, *utils.HTTPError) {
	var host string

	if err := rs.Pool.QueryRow(context.TODO(), "SELECT host FROM rooms WHERE room_id = $1", rid).Scan(&host); err == pgx.ErrNoRows {
		return "", &utils.HTTPError{
			Code:    400,
			Message: "Room does not exist",
		}
	} else if err != nil {
		return "", &utils.HTTPError{
			Code:    500,
			Message: "internal server error",
		}
	}

	return host, nil
}

// come up with better name for this and helper
func (rs *RoomStore) DeleteRoom(rid string) *utils.HTTPError {
	tx, err := rs.Pool.Begin(context.TODO())
	if err != nil {
		fmt.Println(err.Error())
		return &utils.HTTPError{
			Code:    500,
			Message: "Internal server error c",
		}
	}
	defer tx.Rollback(context.TODO())

	if exists, httpErr := checkRoomExistsDB(context.TODO(), tx, rid); httpErr != nil {
		return httpErr
	} else if !exists {
		return &utils.HTTPError{
			Code:    400,
			Message: "room does not exist",
		}
	}

	if httpErr := deleteRoomDB(context.TODO(), tx, rid); httpErr != nil {
		return httpErr
	}

	if err := tx.Commit(context.TODO()); err != nil {
		return &utils.HTTPError{
			Code:    500,
			Message: "Internal server error d",
		}
	}

	return nil
}

func deleteRoomDB(ctx context.Context, tx pgx.Tx, rid string) *utils.HTTPError {
	if _, err := tx.Exec(ctx, "DELETE FROM rooms WHERE room_id = $1", rid); err != nil {
		fmt.Println(err.Error())
		return &utils.HTTPError{
			Code:    500,
			Message: "Internal server error b",
		}
	}

	return nil
}
