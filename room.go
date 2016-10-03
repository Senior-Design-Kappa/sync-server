package main

import (
	"errors"
	"fmt"
)

type Room struct {
	// list of clients in this room
	clients []*Client

	// channel for outbound messages
	outbound chan []byte

	// channel for inbound messages
	inbound chan []byte
}

func NewRoom() *Room {
	h := &Room{
		clients:  make([]*Client, 0),
		inbound:  make(chan []byte),
		outbound: make(chan []byte),
	}
	return h
}

func (r *Room) run() {
	for {
		select {
		case message := <-r.inbound:
			if err := r.handleMessage(message); err != nil {
				fmt.Printf("Message (%s) not handled\n", message)
			}
		}
	}
}

func (r *Room) handleMessage(message []byte) (err error) {
	switch m := parse(message); m.messageType {
	case "test":
		fmt.Printf(m.message)
	default:
		return errors.New("Message not handled")
	}

	return
}

func (r *Room) addClient(client *Client) {
	r.clients = append(r.clients, client)
}

func parse(message []byte) Message {
	m := Message{
		messageType: "test",
		message:     "this is working",
	}
	return m
}
