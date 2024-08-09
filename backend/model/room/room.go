package room

import (
	"backend/model/client"
	"backend/model/message"
	"fmt"
)

type Room struct {
	ID        string                      `json:"room_id"`
	Host      string                      `json:"host"`
	Clients   map[*client.Client]struct{} `json:"-"`
	Broadcast chan *message.Message       `json:"-"`
	Join      chan *client.Client         `json:"-"`
	Leave     chan *client.Client         `json:"-"`
}

func (r *Room) HandleClient(c *client.Client) {
	defer func() {
		fmt.Println("ending handle client")
		r.Leave <- c
		c.CloseConn()
	}()

	readChan := make(chan []byte)
	go readPump(c, readChan)

	for {
		select {
		case m, ok := <-readChan:
			if !ok {
				return
			} else {
				// go ms.InsertMessage()
				// make broadcast a channel of messages and then will have to write conn JSON. add json tags to message struct for writeJSON gorilla websocket
				// make client Write channel message type as well
				// convert incoming msg on read chan to Message struct and then send to broad cast
				// i have user id from client (authenticate client and add user id to struct), room id from r, content from msg and timestamp
				message := &message.Message{
					RoomID:  r.ID,
					UserID:  c.Username,
					Content: string(m),
				}

				r.Broadcast <- message
			}
		case m := <-c.Write:
			if err := c.WriteConn(m); err != nil {
				return
			}
		}
	}
}

func readPump(c *client.Client, readChan chan []byte) {
	for {
		if m, err := c.ReadConn(); err != nil {
			fmt.Println(err.Error())
			close(readChan)
			fmt.Println("readpump ending")
			return
		} else {
			readChan <- m
		}
	}
}
