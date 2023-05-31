package clients

import (
	"github.com/achekslime/flixray/chat_client/models"
	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func genTestMessage() *models.Message {
	return &models.Message{
		Author:   "achek",
		RoomName: "test",
		Message:  "test message",
	}
}

func WriteClient(interrupt chan os.Signal) {
	u := url.URL{Scheme: "ws", Host: "localhost:8086", Path: "/ws/test"}
	log.Printf("[write client] connecting to %s", u.String())

	header := http.Header{}
	header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImltcHNmYWNlQG1haWwucnUiLCJleHAiOjE2ODU0MTkyMzF9.yRJw3Ym0rOM_PfRijV3RVg7F58hgti7fU2PEy-E8tCE")
	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Fatal("[write client] dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("[write client] read:", err)
				return
			}
			var m models.Message
			err = proto.Unmarshal(message, &m)
			if err != nil {
				log.Printf("error: %v", err)
			}
			log.Printf("[write client] receive message: %v", m)
		}
	}()

	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()

	// test message.
	testMessage := genTestMessage()
	messageBytes, err := testMessage.Marshal()
	if err != nil {
		log.Printf("[write client] error: %v", err)
	}

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			err := c.WriteMessage(websocket.BinaryMessage, messageBytes)
			if err != nil {
				log.Println("[write client] write err:", err)
				return
			}
			log.Println("[write client] write message")
		case <-interrupt:
			log.Println("[write client] interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("[write client] write close:", err)
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
