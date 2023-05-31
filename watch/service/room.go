package service

import (
	"github.com/achekslime/flixray/watch/models_json"
	jsoniter "github.com/json-iterator/go"
	"log"
	"time"
)

type room struct {
	// действующие соединения.
	connections map[*connection]bool

	// лог сообщений
	messages []models_json.Event

	// каналы
	addSub    chan *connection
	delSub    chan *connection
	broadcast chan models_json.Event
	doneCh    chan bool

	// доп инфа
	currentVideoUrl string
	// длительность видео в мс.
	duration int
	// время в которое началось воспроизведение.
	timeStart time.Time
	status    string

	timePassed int64
}

func newRoom() *room {
	currentRoom := &room{
		connections: make(map[*connection]bool),
		broadcast:   make(chan models_json.Event),
		addSub:      make(chan *connection),
		delSub:      make(chan *connection),
		messages:    make([]models_json.Event, 0),

		// ToDo заглушка
		currentVideoUrl: "https://app.flixray.ru/srv/streaming/minecraft.m3u8",
		// ToDo заглушка
		duration:   299120,
		timeStart:  time.Now(),
		status:     "PAUSE",
		timePassed: 0,
	}
	return currentRoom
}

func (room *room) eventHandler(msg models_json.Event) {
	switch msg.Type {
	case "ENTER":
		room.enterHandler(msg)
	case "PLAY":
		room.playHandler(msg)
	case "PAUSE":
		room.pauseHandler(msg)
	}
}

func (room *room) enterHandler(msg models_json.Event) {
	response := models_json.Event{
		Author: "",
		Type:   "ENTER",
		Status: room.status,
	}

	var timing int64
	if room.status == "PAUSE" {
		timing = room.timePassed
	} else {
		timing = time.Since(room.timeStart).Milliseconds()
	}
	response.Timing = timing

	// отправляем обратно тайминг.
	for conn := range room.connections {
		if conn.user.Name == msg.Author {
			bytes, err := jsoniter.Marshal(response)
			if err != nil {
				log.Printf("error: %v", err)
			}
			conn.send <- bytes
		}
	}
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

const delta = 3

func (room *room) playHandler(msg models_json.Event) {
	if room.status == "PLAY" {
		if abs(msg.Timing-room.timePassed)/1e3 < delta {
			return
		}
	}

	if room.status == "PAUSE" {
		room.status = "PLAY"
		room.timeStart = time.Now()
	}

	room.sendAll(msg)
}

func (room *room) pauseHandler(msg models_json.Event) {
	if room.status == "PAUSE" {
		return
	}

	room.timePassed = msg.Timing
	room.status = "PAUSE"

	room.sendAll(msg)
}

func (room *room) sendAll(msg models_json.Event) {
	// маршалим сообщение.
	marshalledMessage, err := jsoniter.Marshal(msg)
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
			room.eventHandler(msg)
		}
	}
}
