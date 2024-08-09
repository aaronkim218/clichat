package commands

import (
	"cli/utils"
	"fmt"
	"io"
	"net/http"
)

func Delete() error {
	url := "http://localhost:42069/api/users"

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Cookie", utils.LoadCookie().String())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		fmt.Println("unexpected status code:", resp.StatusCode, "body:", string(body))
		return nil
	}

	fmt.Println("user deleted successfully")

	return nil
}
