package controller

import (
	"backend/model/client"
	"backend/model/hub"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func RegisterRoomsHandlers(api *echo.Group, ws *echo.Group, p *pgxpool.Pool, h *hub.Hub) {
	// roomsAPI := api.Group("/rooms")
	roomsWS := ws.Group("/rooms")

	roomsWS.GET("/:rid", func(c echo.Context) error {
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
