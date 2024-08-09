package commands

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Username string `json:"username"`
}

func Login() {
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

	u := LoginRequest{
		Username: username,
		Password: password,
	}

	jsonData, err := json.Marshal(u)
	if err != nil {
		fmt.Println("marshalling err")
		return
	}

	url := "http://localhost:42069/api/auth/login"

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

	if resp.StatusCode != http.StatusNoContent {
		fmt.Println("bad response")
		return
	}

	cookies := resp.Cookies()
	var clichatSessionCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "clichat-session" {
			clichatSessionCookie = cookie
			break
		}
	}

	if clichatSessionCookie == nil {
		fmt.Println("cookie not found")
		return
	}

	var data map[string]interface{}

	file, err := os.ReadFile("cache.json")
	if err != nil && !os.IsNotExist(err) {
		fmt.Printf("Error reading cache file: %v", err)
	}

	if len(file) > 0 {
		err = json.Unmarshal(file, &data)
		if err != nil {
			fmt.Printf("Error unmarshalling existing cache file: %v", err)
		}
	} else {
		data = make(map[string]interface{})
	}

	data["session_cookie"] = clichatSessionCookie
	data["username"] = username

	updatedJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling updated JSON: %v", err)
	}

	err = os.WriteFile("cache.json", updatedJSON, 0644)
	if err != nil {
		fmt.Printf("Error writing updated JSON to cache file: %v", err)
	}

	fmt.Println("cookie saved successfully")
}
