package hub

import (
	"backend/model/client"
	"backend/model/room"
	"fmt"
	"sync"
)

type Hub struct {
	Rooms sync.Map
}

// load room from hub
// retrieve room from db and add to hub and return room
// for non-db load room from hub or create new room by default if id not found
func (h *Hub) LoadRoom(rid string) *room.Room {
	r, loaded := h.Rooms.LoadOrStore(rid, &room.Room{
		ID:        rid,
		Clients:   make(map[*client.Client]struct{}),
		Broadcast: make(chan []byte),
		Join:      make(chan *client.Client),
		Leave:     make(chan *client.Client),
	})

	if !loaded {
		go h.Run(r.(*room.Room))
	}

	return r.(*room.Room)
}

func (h *Hub) Run(r *room.Room) {
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
			for c := range r.Clients {
				c.Write <- m
			}
		}
	}
}
