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
	inbound chan InboundMessage
}

func NewRoom() *Room {
	r := &Room{
		clients:  make(map[*Client]bool),
		inbound:  make(chan InboundMessage),
		outbound: make(chan []byte),
	}
	return r
}

func (r *Room) run() {
	for {
		select {
		case inboundMessage := <-r.inbound:
			if err := r.handleMessage(inboundMessage); err != nil {
				log.Printf("Message not handled\n")
			}
		}
	}
}

func (r *Room) handleMessage(inboundMessage InboundMessage) (err error) {
	message := inboundMessage.RawMessage
	switch m := parse(message); m.MessageType {
	case "debug":
		log.Printf("%+v", m)
	case "INIT":
		client := inboundMessage.Sender
		outbound, _ := json.Marshal(Message{
			MessageType: "INIT",
			Hash:        client.hash,
		})
		client.send <- outbound
	case "SYNC_VIDEO":
		for client := range r.clients {
			if client != inboundMessage.Sender {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(r.clients, client)
				}
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
