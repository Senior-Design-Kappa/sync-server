package main

import (
	"github.com/gorilla/websocket"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type SocketHandler interface {
}

// SocketHandler is an middleman between the websocket connection and the hub.
type socketHandler struct {
	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

func NewSocketHandler() SocketHandler {
	s := &socketHandler{}
	return s
}
