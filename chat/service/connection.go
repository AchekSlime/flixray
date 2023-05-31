package service

import (
	"fmt"
	"github.com/achekslime/core/models"
	"github.com/achekslime/flixray/chat/models_json"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"log"
	"time"
)

// connection обернутое ws соединение пользователя.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn
	// Канал куда с broadcast комнаты пишутся сообщения.
	send chan []byte
	// Имя комнаты.
	roomName string
	// Пользователь.
	user *models.User
}

func newConnection(ws *websocket.Conn, user *models.User, roomName string) *connection {
	return &connection{
		ws:       ws,
		send:     make(chan []byte, 256),
		roomName: roomName,
		user:     user,
	}
}

// readHandler чтение сообщений из соединения в broadcast комнаты.
func (conn *connection) readHandler(broadcast chan models_json.Message, delSub chan *connection) {
	defer func() {
		delSub <- conn
		conn.ws.Close()
	}()
	conn.ws.SetReadLimit(maxMessageSize)
	conn.ws.SetReadDeadline(time.Now().Add(pongWait))
	conn.ws.SetPongHandler(func(string) error { conn.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// readHandler
	for {
		_, msg, err := conn.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		// анмаршалим proto.
		var m models_json.Message
		err = jsoniter.Unmarshal(msg, &m)
		// proto.
		//err = proto.Unmarshal(msg, &m)
		if err != nil {
			log.Printf("error: %v", err)
		}
		log.Printf(fmt.Sprintf("message: %s", m))

		// кладем сообщение в общий канал.
		broadcast <- m
	}
}

// writeHandler чтение сообщений из broadcast комнаты и запись их в ws соединение.
func (conn *connection) writeHandler() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		conn.ws.Close()
	}()
	for {
		select {
		case msg, ok := <-conn.send:
			if !ok {
				conn.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := conn.write(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			if err := conn.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// write writes a message with the given message type and payload.
func (conn *connection) write(mt int, payload []byte) error {
	conn.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return conn.ws.WriteMessage(mt, payload)
}
