package main

import (
	"encoding/json"
	"errors"
	"log"
)

type Room struct {
	// list of clients in this room
	clients map[*Client]bool

	// channel for outbound messages
	outbound chan []byte

	// channel for inbound messages
	inbound chan []byte
}

func NewRoom() *Room {
	r := &Room{
		clients:  make(map[*Client]bool),
		inbound:  make(chan []byte),
		outbound: make(chan []byte),
	}
	return r
}

func (r *Room) run() {
	for {
		select {
		case message := <-r.inbound:
			if err := r.handleMessage(message); err != nil {
				log.Printf("Message (%s) not handled\n", message)
			}
		}
	}
}

func (r *Room) handleMessage(message []byte) (err error) {
	switch m := parse(message); m.MessageType {
	case "test":
		for client := range r.clients {
			// m, err := json.Marshal(message)
			// if err != nil {
			// 	return err
			// }
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(r.clients, client)
			}
		}
	default:
		return errors.New("")
	}

	return
}

func (r *Room) addClient(client *Client) {
	r.clients[client] = true
}

func parse(message []byte) Message {
	var m Message
	if err := json.Unmarshal(message, &m); err != nil {
		log.Printf("error unmarshaling message: %+v\n", err)
	}
	return m
}
