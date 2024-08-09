package commands

import (
	"bytes"
	"cli/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type request struct {
	RoomID string `json:"room_id"`
}

type response struct {
	RoomID string `json:"room_id"`
	Host   string `json:"host"`
}

func Create(rid string) error {

	req := request{
		RoomID: rid,
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		fmt.Println("error marshaling req")
		return err
	}

	url := "http://localhost:42069/api/rooms"

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonReq))
	if err != nil {
		fmt.Println("error creating http req")
		return err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Add("Cookie", utils.LoadCookie().String())

	client := new(http.Client)
	resp, err := client.Do(httpReq)
	if err != nil {
		fmt.Println("error making req" + err.Error())
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading body")
		return err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("unexpected status code:", resp.StatusCode, "body:", string(respBody))
		return nil
	}

	var parsedResp response
	err = json.Unmarshal(respBody, &parsedResp)
	if err != nil {
		fmt.Println("error unmarshalling")
		return err
	}

	fmt.Println("Created room with host:" + parsedResp.Host)
	return nil
}
