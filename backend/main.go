package main

import (
	"backend/controller"
	"backend/model/hub"
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbURI := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	p, err := pgxpool.New(context.Background(), dbURI)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	h := new(hub.Hub)

	e := echo.New()

	controller.RegisterHandlers(e, p, h)

	if err := e.Start(":42069"); err != nil {
		fmt.Println(err.Error())
		return
	}
}
