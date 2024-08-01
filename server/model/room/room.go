package room

import (
	"fmt"
	"server/model/client"
)

type Room struct {
	ID        string
	Clients   map[*client.Client]struct{}
	Broadcast chan []byte
	Join      chan *client.Client
	Leave     chan *client.Client
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
				r.Broadcast <- m
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
