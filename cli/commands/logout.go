package commands

import (
	"cli/utils"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func Logout() error {
	url := "http://localhost:42069/api/auth/logout"

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		fmt.Println("error creating request")
		return err
	}

	req.Header.Add("Cookie", utils.LoadCookie().String())

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making req")
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading body")
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		fmt.Println("unexpected response", string(respBody))
		return errors.New("unexpected response")
	}

	utils.ResetCookie()
	fmt.Println("cookie removed from cache successfully")

	return nil
}
