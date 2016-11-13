package room

import (
	"encoding/json"
	"errors"
	"log"
  "time"

	"github.com/Senior-Design-Kappa/sync-server/models"
)

type Room struct {
	// list of clients in this room
	clients map[*Client]bool

	// channel for outbound messages
	outbound chan []byte

	// channel for inbound messages
	inbound chan InboundMessage

  // state of room (canvas + video state)
  state RoomState
}

func NewRoom() *Room {
	r := &Room{
		clients:  make(map[*Client]bool),
		inbound:  make(chan InboundMessage),
		outbound: make(chan []byte),
    state:    *NewRoomState(),
	}
	return r
}

func (r *Room) Run() {
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
  r.state.UpdateStateFromInboundMessage(inboundMessage)

	message := inboundMessage.RawMessage
	switch m := parse(message); m.MessageType {
	case "debug":
		log.Printf("%+v", m)
	case "INIT":
		client := inboundMessage.Sender
    videoTime := r.state.CurrentTime
    if r.state.Playing {
      videoTime += float32(time.Now().Sub(r.state.LastTime).Seconds())
    }
		outbound, _ := json.Marshal(models.Message{
			MessageType: "INIT",
			Hash:        client.hash,

      Video: models.VideoState {
        Playing: r.state.Playing,
        CurrentTime: videoTime,
        Volume: r.state.Volume,
        Muted: r.state.Muted,
      },

      Actions: r.state.Actions,
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
	case "SYNC_CANVAS":
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

func (r *Room) AddClient(client *Client) {
	r.clients[client] = true
}

func parse(message []byte) models.Message {
	var m models.Message
	if err := json.Unmarshal(message, &m); err != nil {
		log.Printf("error unmarshaling message: %+v\n", err)
	}
	return m
}
