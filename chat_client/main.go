package main

import (
	"github.com/achekslime/flixray/chat_client/clients"
	_ "github.com/gogo/protobuf/proto"
	"os"
	"time"
)

func main() {
	interrupt := make(chan os.Signal, 1)

	go clients.WriteClient(interrupt)
	time.Sleep(time.Second * 12)
	//go clients.ReadClient(interrupt)

	for {
		time.Sleep(time.Hour)
	}
}
