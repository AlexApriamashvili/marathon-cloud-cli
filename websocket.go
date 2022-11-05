package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type RuntimeState struct {
	TotalEmulators   int    `json:"total_emulators"`
	WorkingEmulators int    `json:"working_emulators"`
	State            string `json:"state"`
}

func Subscribe(token string, runId string) {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "runtime.testwise.pro:1005", Path: "/hello", RawQuery: "token=" + token + "&run_id=" + runId}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("Dial fatal")
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	func() {
		defer close(done)
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			var message RuntimeState
			err = json.Unmarshal(data, &message)
			if err != nil {
				fmt.Println("Error reading runtime")
				continue
			}
			if len(message.State) > 0 {
				fmt.Println(time.Now().Format(time.Stamp), message.State)
				continue
			}
			fmt.Printf("%s Running %d emulators of %d total\n", time.Now().Format(time.Stamp), message.WorkingEmulators, message.TotalEmulators)

		}
	}()

}