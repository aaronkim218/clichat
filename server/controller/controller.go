package controller

import (
	"net/http"
	"server/model/hub"

	"github.com/labstack/echo/v4"
)

func RegisterAPIHandlers(e *echo.Echo) {
	api := e.Group("/api")

	api.GET("/healthcheck", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})
}

func RegisterWSHandlers(e *echo.Echo, h *hub.Hub) {
	ws := e.Group("/ws")

	RegisterRoomsWSHandlers(ws, h)
}
