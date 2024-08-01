package commands

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Version() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	v := os.Getenv("VERSION")
	fmt.Println(v)

	return nil
}
