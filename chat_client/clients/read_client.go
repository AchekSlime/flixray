package clients

import (
	"github.com/achekslime/flixray/chat_client/models"
	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"time"
)

func ReadClient(interrupt chan os.Signal) {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws/test"}
	log.Printf("[read client] connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("[read client] dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("[read client] read:", err)
				return
			}
			var m models.Message
			err = proto.Unmarshal(message, &m)
			if err != nil {
				log.Printf("error: %v", err)
			}
			log.Printf("[read client] receive message: %v", m)
		}
	}()

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("[read client] interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("[read client] write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
