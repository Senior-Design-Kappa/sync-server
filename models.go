package main

import "github.com/gorilla/websocket"

type NewConnection struct {
	conn *websocket.Conn
	room string
}

type Message struct {
	messageType string
	message     string
}
