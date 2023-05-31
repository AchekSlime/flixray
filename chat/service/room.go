package service

import (
	"github.com/achekslime/flixray/chat/models_json"
	jsoniter "github.com/json-iterator/go"
	"log"
)

type room struct {
	// действующие соединения.
	connections map[*connection]bool

	// лог сообщений
	messages []models_json.Message

	// каналы
	addSub    chan *connection
	delSub    chan *connection
	broadcast chan models_json.Message
	doneCh    chan bool
}

func newRoom() *room {
	return &room{
		connections: make(map[*connection]bool),
		broadcast:   make(chan models_json.Message),
		addSub:      make(chan *connection),
		delSub:      make(chan *connection),
		messages:    make([]models_json.Message, 0),
	}
}

func (room *room) Run() {
	for {
		select {
		case conn := <-room.addSub:
			if room.connections == nil {
				room.connections = make(map[*connection]bool)
			}
			room.connections[conn] = true
		case conn := <-room.delSub:
			if room.connections != nil {
				if _, ok := room.connections[conn]; ok {
					delete(room.connections, conn)
					close(conn.send)
					// удалить комнату.
					//if len(room.connections) == 0 {
					//	delete(h.rooms, s.room)
					//}
				}
			}
		case msg := <-room.broadcast:
			// сохраняем сообщение в лог
			room.messages = append(room.messages, msg)
			// маршалим сообщение.
			marshalledMessage, err := jsoniter.Marshal(msg)
			// proto.
			//marshalledMessage, err := msg.Marshal()
			if err != nil {
				log.Println(err.Error())
			}
			// проходим по всем соединениям.

			for conn := range room.connections {
				if conn.user.Name == msg.Author {
					continue
				}
				select {
				// отправляем сообщения всем соединениям.
				case conn.send <- marshalledMessage:
				default:
					close(conn.send)
					delete(room.connections, conn)
					// удалить комнату.
					//if len(room.connections) == 0 {
					//	delete(h.rooms, s.room)
					//}
				}
			}
		}
	}
}
