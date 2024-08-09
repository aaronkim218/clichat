package commands

import (
	"cli/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gorilla/websocket"
	"github.com/rivo/tview"
)

type Message struct {
	Index     int       `json:"index"`
	RoomID    string    `json:"room_id"`
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
}

func Enter(rid string) {
	cookie := utils.LoadCookie()
	headers := http.Header{}
	headers.Add("Cookie", cookie.String())

	dialer := websocket.DefaultDialer
	conn, resp, err := dialer.Dial("ws://localhost:42069/api/rooms/"+rid+"/ws", headers)
	if err != nil {
		fmt.Println(err.Error())
		if resp != nil {
			// catch room does not exist err and display message
			fmt.Println("TODO")
		}
		return
	}
	defer conn.Close()

	app := tview.NewApplication()

	history, _ := getHistory(rid, &headers)

	textView := tview.NewTextView().
		SetText(history).
		SetDynamicColors(true).
		SetWordWrap(true).
		SetScrollable(true)

	input := tview.NewInputField().
		SetLabel(">").
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorWhite)

	input.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			text := input.GetText()
			err := conn.WriteMessage(websocket.TextMessage, []byte(text))
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			input.SetText("")
		}
	})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(textView, 0, 1, true).
		AddItem(input, 1, 0, true)

	modeChan := make(chan string)
	msgs := make(chan *Message)
	go manageView(app, textView, msgs, modeChan)
	go readPump(conn, msgs)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyESC:
			app.Stop()
			return nil
		case tcell.KeyCtrlN:
			app.SetFocus(textView)
			modeChan <- "n"
			return nil
		case tcell.KeyCtrlR:
			app.SetFocus(textView)
			modeChan <- "r"
			return nil
		case tcell.KeyCtrlW:
			app.SetFocus(input)
			modeChan <- "w"
			return nil
		default:
			return event
		}
	})

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}

func manageView(app *tview.Application, textView *tview.TextView, msgs chan *Message, modeChan chan string) {
	mode := ""

	for {
		select {
		case m := <-modeChan:
			mode = m
		case msg := <-msgs:

			app.QueueUpdateDraw(func() {
				textView.Write(append([]byte(formatMessage((msg))), '\n'))
				if mode != "r" {
					textView.ScrollToEnd()
				}
			})
		}
	}
}

func readPump(conn *websocket.Conn, msgs chan *Message) {
	defer close(msgs)

	// need to change this to ReadJSON
	for {
		m := new(Message)
		err := conn.ReadJSON(m)
		if err != nil {
			fmt.Println("readPump error:", err.Error())
			return
		}
		msgs <- m
	}
}

func formatMessage(m *Message) string {
	return fmt.Sprintf("%s: %s", m.UserID, m.Content)
}

func getHistory(rid string, headers *http.Header) (string, error) {
	url := "http://localhost:42069/api/rooms/" + rid + "/messages"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("error creating new request")
		return "", err
	}

	req.Header = *headers

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error making request", err.Error())
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("unexpected status code")
		return "", err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading body")
		return "", err
	}

	var response []Message
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		fmt.Println("error unmarshalling data")
		return "", err
	}

	history := ""
	for _, msg := range response {
		fmt.Println("user:" + msg.UserID)
		history = history + formatMessage(&msg) + "\n"
	}

	return history, nil
}
