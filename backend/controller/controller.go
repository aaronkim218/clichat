package controller

import (
	"backend/model/hub"
	"fmt"
	"os"

	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(e *echo.Echo, p *pgxpool.Pool, h *hub.Hub) {

	store := sessions.NewCookieStore([]byte(os.Getenv("AUTHENTICATION_KEY")))
	store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
	}

	api := e.Group("/api")

	api.Use(SessionMiddleware(store))
	api.Use(AuthenticationMiddleware())

	api.GET("/healthcheck", func(c echo.Context) error {
		return c.NoContent(204)
	})

	RegisterUsersHandlers(api, p)
	RegisterRoomsHandlers(api, p, h)
	RegisterAuthHandlers(api, p, store)
}

func SessionMiddleware(store *sessions.CookieStore) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			fmt.Println("starting session middle")
			session, err := store.Get(c.Request(), "clichat-session")
			if err != nil {
				return echo.NewHTTPError(400, "no cookie present")
			}

			if username, ok := session.Values["username"].(string); ok {
				c.Set("username", username)
			}

			if authenticated, ok := session.Values["authenticated"].(bool); ok {
				c.Set("authenticated", authenticated)
			}

			fmt.Println("passed session middleware")

			return next(c)
		}
	}
}

func AuthenticationMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var excluded = map[string]map[string]struct{}{
				"/api/auth/login": {
					"POST": {},
				},
				"/api/users": {
					"POST": {},
				},
			}

			if methods, ok := excluded[c.Path()]; ok {
				if _, ok := methods[c.Request().Method]; ok {
					return next(c)
				}
			}

			if authenticated, ok := c.Get("authenticated").(bool); !ok || !authenticated {
				return echo.NewHTTPError(401, "middleware Unauthorized")
			}

			fmt.Println("passed authentication")

			return next(c)
		}
	}
}
