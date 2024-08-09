package message

import "time"

type Message struct {
	Index     int       `json:"index"`
	RoomID    string    `json:"room_id"`
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
}
