package controller

import (
	"backend/model/user"

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
		u := new(user.User)

		if err := c.Bind(u); err != nil {
			return echo.NewHTTPError(400, "Invalid user")
		}

		if err := us.InsertUser(u); err != nil {
			return echo.NewHTTPError(err.Code, err.Message)
		}

		return c.JSON(200, u)
	})
}
