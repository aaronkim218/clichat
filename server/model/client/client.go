package client

import (
	"errors"
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn  *websocket.Conn
	Write chan []byte
}

func (c *Client) CloseConn() error {
	// close connection logic
	c.Conn.Close()
	return nil
}

func (c *Client) ReadConn() ([]byte, error) {
	if _, m, err := c.Conn.ReadMessage(); err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("error reading from connection")
	} else {
		return m, nil
	}
}

func (c *Client) WriteConn(m []byte) error {
	if err := c.Conn.WriteMessage(websocket.TextMessage, m); err != nil {
		fmt.Println(err.Error())
		return errors.New("error writing to connection")
	} else {
		return nil
	}
}
