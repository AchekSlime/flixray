package service

import "time"

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to readHandler the next pong message from the peer.
	pongWait = 60 * time.Second

	// send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)
