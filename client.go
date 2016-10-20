package main

import (
	"bytes"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 8192

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

// Client is the interface to the websocket connection.
type Client struct {
	// The room this client is connected to
	room *Room

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// Generated unique identifier
	hash string
}

func NewClient(conn *websocket.Conn, room *Room, hash string) *Client {
	c := &Client{
		conn: conn,
		send: make(chan []byte, 256),
		room: room,
		hash: hash,
	}
	return c
}

func (c *Client) run() {
	go c.readFrom()
	go c.writeTo()
}

// write writes a message with the given message type and payload.
func (c *Client) write(mt int, payload []byte) error {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteMessage(mt, payload)
}

func (c *Client) readFrom() {
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(c.pongHandler)
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
				delete(c.room.clients, c)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.room.inbound <- InboundMessage{c, message}
	}
}

func (c *Client) pongHandler(str string) error {
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	return nil
}

func (c *Client) writeTo() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// The room closed the channel.
				c.write(websocket.CloseMessage, []byte{})
				return
			}

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
