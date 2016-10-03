package main

import "github.com/gorilla/websocket"

type NewConnection struct {
	conn *websocket.Conn
	room string
}

type VideoState struct {
	Playing     bool    `json:"playing"`
	CurrentTime float32 `json:"currentTime"`
}

type Message struct {
	MessageType string     `json:"messageType"`
	Message     string     `json:"message"`
	Video       VideoState `json:"videoState"`
}
