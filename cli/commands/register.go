package commands

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Username string `json:"username"`
}

func Register() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter username: ")
	scanner.Scan()
	username := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter password: ")
	scanner.Scan()
	password := strings.TrimSpace(scanner.Text())

	if err := scanner.Err(); err != nil {
		fmt.Println("scanner err")
		os.Exit(1)
	}

	u := RegisterRequest{
		Username: username,
		Password: password,
	}

	jsonData, err := json.Marshal(u)
	if err != nil {
		fmt.Println("marshalling err")
		return
	}

	url := "http://localhost:42069/api/users"

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("error creating new request")
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error making request", err.Error())
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("bad response")
		return
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading response body")
		return
	}

	var response RegisterResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		fmt.Println("error unmarshalling data")
	}

	fmt.Println("printing username field: ", response.Username)

}
