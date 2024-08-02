package controller

import (
	"backend/model/hub"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(e *echo.Echo, p *pgxpool.Pool, h *hub.Hub) {

	api := e.Group("/api")
	ws := e.Group("/ws")

	api.GET("/healthcheck", func(c echo.Context) error {
		return c.NoContent(204)
	})

	RegisterUsersHandlers(api, p)
	RegisterRoomsHandlers(api, ws, p, h)
	RegisterAuthHandlers(api, p)
}
