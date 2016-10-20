package main

import "github.com/gorilla/websocket"

type NewConnection struct {
	conn *websocket.Conn
	room string
	hash string
}

type VideoState struct {
	Playing     bool    `json:"playing"`
	CurrentTime float32 `json:"currentTime"`
}

type CanvasState struct {
}

type InboundMessage struct {
	Sender     *Client
	RawMessage []byte
}

type Message struct {
	MessageType string     `json:"messageType"`
	Message     string     `json:"message"`
	Hash        string     `json:"hash"`
	Video       VideoState `json:"videoState"`
}
