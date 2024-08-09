package hub

import (
	"backend/model/client"
	"backend/model/message"
	"backend/model/room"
	"backend/utils"
	"fmt"
	"sync"
	"time"
)

type Hub struct {
	Rooms sync.Map
}

// load room from hub
// retrieve room from db and add to hub and return room
// for non-db load room from hub or create new room by default if id not found
func (h *Hub) LoadRoom(rid string, rs *room.RoomStore, ms *message.MessageStore) (*room.Room, *utils.HTTPError) {
	if r, ok := h.Rooms.Load(rid); ok {
		go h.Run(r.(*room.Room), ms)
		return r.(*room.Room), nil
	}

	if r, httpErr := rs.SelectRoom(rid); httpErr != nil {
		return nil, httpErr
	} else {
		r.Clients = make(map[*client.Client]struct{})
		r.Broadcast = make(chan *message.Message)
		r.Join = make(chan *client.Client)
		r.Leave = make(chan *client.Client)
		h.Rooms.Store(rid, r)
		go h.Run(r, ms)
		return r, nil
	}
}

func (h *Hub) Run(r *room.Room, ms *message.MessageStore) {
	for {
		select {
		case c := <-r.Join:
			r.Clients[c] = struct{}{}
			go r.HandleClient(c)
		case c := <-r.Leave:
			fmt.Println("removing client")
			delete(r.Clients, c)
			if len(r.Clients) == 0 {
				fmt.Println("deleting room", r.ID)
				h.Rooms.Delete(r.ID)
				return
			}
		case m := <-r.Broadcast:
			m.Timestamp = time.Now()
			msg, err := ms.InsertMessage(m)
			if err != nil {
				fmt.Println("error inserting msgs into db. shutting down")
				return
			}
			for c := range r.Clients {
				c.Write <- msg
			}
		}
	}
}
