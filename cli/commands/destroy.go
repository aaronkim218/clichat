package commands

import (
	"cli/utils"
	"fmt"
	"io"
	"net/http"
)

func Destroy(rid string) error {
	url := "http://localhost:42069/api/rooms/" + rid

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Cookie", utils.LoadCookie().String())

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

	fmt.Println("room deleted successfully")

	return nil
}
