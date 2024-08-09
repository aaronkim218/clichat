package message

import (
	"backend/utils"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageStore struct {
	Pool *pgxpool.Pool
}

func (ms *MessageStore) SelectMessagesByRoom(rid string) ([]*Message, *utils.HTTPError) {
	query := `
		SELECT index, room_id, timestamp, author, content
		FROM messages
		WHERE room_id = $1
		ORDER BY timestamp ASC;
	`

	rows, err := ms.Pool.Query(context.TODO(), query, rid)
	if err != nil {
		return nil, &utils.HTTPError{
			Code:    500,
			Message: "internal server error",
		}
	}

	var messages []*Message
	for rows.Next() {
		m := new(Message)
		var author *string

		if err := rows.Scan(&m.Index, &m.RoomID, &m.Timestamp, &author, &m.Content); err != nil {
			return nil, &utils.HTTPError{
				Code:    500,
				Message: "internal server error",
			}
		}

		if author == nil {
			m.UserID = "deleted-user"
		} else {
			m.UserID = *author
		}

		messages = append(messages, m)
	}

	if rows.Err() != nil {
		return nil, &utils.HTTPError{
			Code:    500,
			Message: "internal server error",
		}
	}

	return messages, nil
}

func (ms *MessageStore) InsertMessage(m *Message) (*Message, *utils.HTTPError) {
	var index int
	if err := ms.Pool.QueryRow(context.TODO(), "INSERT INTO messages (room_id, timestamp, author, content) VALUES ($1, $2, $3, $4) RETURNING index", m.RoomID, m.Timestamp, m.UserID, m.Content).Scan(&index); err != nil {
		return nil, &utils.HTTPError{
			Code:    500,
			Message: "internal server error",
		}
	}

	m.Index = index

	return m, nil
}
