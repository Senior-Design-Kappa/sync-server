package models

import "github.com/gorilla/websocket"

type NewConnection struct {
	Conn *websocket.Conn
	Room string
	Hash string
}

type VideoState struct {
	Playing     bool    `json:"playing"`
	CurrentTime float32 `json:"currentTime"`
}

type CanvasState struct {
}

type LineSegment struct {
  PrevX int `json:"prevX"`
  PrevY int `json:"prevY"`
  CurrX int `json:"currX"`
  CurrY int `json:"currY"`
}

type Message struct {
	MessageType string     `json:"messageType"`
	Message     string     `json:"message"`
	Hash        string     `json:"hash"`
	Video       VideoState `json:"videoState"`

  PrevX       int        `json:"prevX"`
  PrevY       int        `json:"prevY"`
  CurrX       int        `json:"currX"`
  CurrY       int        `json:"currY"`

  Lines       []LineSegment `json:"lines"`
}
