package controller

import (
	"backend/model/user"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func RegisterUsersHandlers(api *echo.Group, p *pgxpool.Pool) {
	users := api.Group("/users")
	// usersWS := ws.Group("/users")

	us := user.UserStore{
		Pool: p,
	}

	users.POST("", func(c echo.Context) error {
		type response struct {
			Username string `json:"username"`
		}

		u := new(user.User)

		if err := c.Bind(u); err != nil {
			return echo.NewHTTPError(400, "Invalid user")
		}

		if err := us.InsertUser(u); err != nil {
			return echo.NewHTTPError(err.Code, err.Message)
		}

		return c.JSON(200, response{
			Username: u.Username,
		})
	})

	// need to add uid to path
	users.DELETE("", func(c echo.Context) error {
		username := c.Get("username").(string)
		if httpErr := us.DeleteUser(username); httpErr != nil {
			return echo.NewHTTPError(httpErr.Code, httpErr.Message)
		}

		return c.NoContent(http.StatusNoContent)
	})
}
