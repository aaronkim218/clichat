package commands

import (
	"fmt"
	"os"
)

func Version() error {
	data, err := os.ReadFile("VERSION")
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}
