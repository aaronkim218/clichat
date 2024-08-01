package controller

import (
	"fmt"
	"server/model/client"
	"server/model/hub"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func RegisterRoomsWSHandlers(ws *echo.Group, h *hub.Hub) {
	rws := ws.Group("/rooms")

	rws.GET("/:rid", func(c echo.Context) error {
		rid := c.Param("rid")

		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}

		r := h.LoadRoom(rid)

		client := client.Client{
			Conn:  conn,
			Write: make(chan []byte),
		}

		r.Join <- &client

		h.Rooms.Range(func(key, value any) bool {
			fmt.Println("room stored in hub with id:", key)
			return true
		})
		return nil
	})
}
