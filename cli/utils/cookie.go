package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func LoadCookie() *http.Cookie {
	file, err := os.Open("cache.json")
	if err != nil {
		fmt.Printf("Error opening cookie file: %v", err)
	}
	defer file.Close()

	var data map[string]http.Cookie
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Printf("Error decoding cookie from JSON: %v", err)
	}

	cookie, exists := data["session_cookie"]
	if !exists {
		fmt.Printf("Cookie not found in the file")
	}

	return &cookie
}

func ResetCookie() {
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

	data["session_cookie"] = struct{}{}

	updatedJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling updated JSON: %v", err)
	}

	err = os.WriteFile("cache.json", updatedJSON, 0644)
	if err != nil {
		fmt.Printf("Error writing updated JSON to cache file: %v", err)
	}
}
