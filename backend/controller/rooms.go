package controller

import (
	"backend/model/client"
	"backend/model/hub"
	"backend/model/message"
	"backend/model/room"
	"fmt"
	"net/http"

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

func RegisterRoomsHandlers(api *echo.Group, p *pgxpool.Pool, h *hub.Hub) {
	rooms := api.Group("/rooms")

	rs := room.RoomStore{
		Pool: p,
	}
	ms := message.MessageStore{
		Pool: p,
	}

	rooms.POST("", func(c echo.Context) error {
		username, ok := c.Get("username").(string)
		if !ok {
			return echo.NewHTTPError(400, "bad username")
		}

		r := new(room.Room)

		type RequestBody struct {
			RoomID string `json:"room_id"`
		}

		var reqBody RequestBody
		if err := c.Bind(&reqBody); err != nil {
			return echo.NewHTTPError(400, "bad request")
		}

		r.ID = reqBody.RoomID
		r.Host = username

		if err := rs.InsertRoom(r); err != nil {
			return echo.NewHTTPError(err.Code, err.Message)
		}

		return c.JSON(200, r)
	})

	rooms.GET("/:rid/ws", func(c echo.Context) error {
		rid := c.Param("rid")

		r, httpErr := h.LoadRoom(rid, &rs, &ms)
		if httpErr != nil {
			return echo.NewHTTPError(httpErr.Code, httpErr.Message)
		}

		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}

		// username should be available because of middleware
		username, _ := c.Get("username").(string)

		client := client.Client{
			Username: username,
			Conn:     conn,
			Write:    make(chan *message.Message),
		}

		r.Join <- &client

		fmt.Println("sent client to join")

		h.Rooms.Range(func(key, value any) bool {
			fmt.Println("room stored in hub with id:", key)
			return true
		})
		return nil
	})

	rooms.GET("/:rid/messages", func(c echo.Context) error {
		rid := c.Param("rid")

		// check room exists

		msgs, httpErr := ms.SelectMessagesByRoom(rid)
		if httpErr != nil {
			return echo.NewHTTPError(httpErr.Code, httpErr.Message)
		}

		return c.JSON(200, msgs)
	})

	rooms.DELETE("/:rid", func(c echo.Context) error {
		rid := c.Param("rid")

		if err := rs.DeleteRoom(rid); err != nil {
			return echo.NewHTTPError(err.Code, err.Message)
		}

		return c.NoContent(http.StatusNoContent)

	}, hostMiddleware(&rs))
}

// TODO: middleware for routes that modify rooms and is applied at route level
// e.g. check that user who is trying to update room/ delete room is host
// check that user who is trying to join room has been invited (potentially)
// check that user who is trying to leave room is a member

func hostMiddleware(rs *room.RoomStore) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			fmt.Println("starting host middleware")

			rid := c.Param("rid")
			username, _ := c.Get("username").(string)

			host, err := rs.SelectHost(rid)
			if err != nil {
				return echo.NewHTTPError(err.Code, err.Message)
			}

			if host != username {
				return echo.NewHTTPError(400, "not allowed. user is not host of room")
			}

			fmt.Println("here")
			fmt.Println("host", host)
			fmt.Println("username", username)

			return next(c)
		}
	}
}
