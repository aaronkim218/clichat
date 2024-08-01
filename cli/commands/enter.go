package commands

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/gorilla/websocket"
	"github.com/rivo/tview"
)

func Enter(rid string) {
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial("ws://localhost:42069/ws/rooms/"+rid, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	app := tview.NewApplication()

	textView := tview.NewTextView().
		SetText("").
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
	msgs := make(chan []byte)
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

func manageView(app *tview.Application, textView *tview.TextView, msgs chan []byte, modeChan chan string) {
	mode := ""

	for {
		select {
		case m := <-modeChan:
			mode = m
		case msg := <-msgs:

			app.QueueUpdateDraw(func() {
				textView.Write(append(msg, '\n'))
				if mode != "r" {
					textView.ScrollToEnd()
				}
			})
		}
	}
}

func readPump(conn *websocket.Conn, msgs chan []byte) {
	defer close(msgs)

	for {
		_, m, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("readPump error:", err.Error())
			return
		}
		msgs <- m
	}
}
