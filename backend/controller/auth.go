package controller

import (
	"backend/model/user"
	"os"

	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func RegisterAuthHandlers(api *echo.Group, p *pgxpool.Pool) {
	auth := api.Group("/auth")

	us := user.UserStore{
		Pool: p,
	}

	store := sessions.NewCookieStore([]byte(os.Getenv("AUTHENTICATION_KEY")))
	store.Options = &sessions.Options{
		HttpOnly: true,
	}

	auth.POST("/login", func(c echo.Context) error {
		reqUser := new(user.User)

		if err := c.Bind(reqUser); err != nil {
			return echo.NewHTTPError(400, "Invalid user")
		}

		storedUser, httpErr := us.SelectUser(reqUser.Username)
		if httpErr != nil {
			return echo.NewHTTPError(httpErr.Code, httpErr.Message)
		}

		if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(reqUser.Password)); err != nil {
			return echo.NewHTTPError(401, "Incorrect password")
		}

		session, err := store.Get(c.Request(), "clichat-session")
		if err != nil {
			return echo.NewHTTPError(500, "Internal server error")
		}

		session.Values["authenticated"] = true
		session.Values["username"] = reqUser.Username

		err = session.Save(c.Request(), c.Response().Writer)
		if err != nil {
			return echo.NewHTTPError(500, "Internal server error")
		}

		return c.NoContent(204)
	})

	auth.POST("/logout", func(c echo.Context) error {
		session, err := store.Get(c.Request(), "clichat-session")
		if err != nil {
			return echo.NewHTTPError(500, "Internal server error")
		}

		session.Options.MaxAge = -1

		err = sessions.Save(c.Request(), c.Response().Writer)
		if err != nil {
			return echo.NewHTTPError(500, "Internal server error")
		}

		return c.NoContent(204)
	})
}
