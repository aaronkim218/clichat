package main

import (
	"context"
	"fmt"
	"os"
	"server/controller"
	"server/model/hub"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func main() {
	urlExample := "postgres://admin:admin@localhost:5432/clichat"
	conn, err := pgx.Connect(context.Background(), urlExample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var name string
	var weight int64
	err = conn.QueryRow(context.Background(), "select name, weight from widgets where id=$1", 42).Scan(&name, &weight)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(name, weight)
	e := echo.New()

	h := &hub.Hub{}

	controller.RegisterAPIHandlers(e)
	controller.RegisterWSHandlers(e, h)

	e.Start(":42069")
}
