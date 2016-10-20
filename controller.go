/*
Controller handles new connections by registering them to correct rooms
*/
package main

import "log"

type Controller struct {
	clients map[*Client]*Room

	rooms map[string]*Room

	register chan *NewConnection
}

func NewController() *Controller {
	c := &Controller{
		clients:  make(map[*Client]*Room),
		rooms:    make(map[string]*Room),
		register: make(chan *NewConnection),
	}
	return c
}

func (c *Controller) run() {
	for {
		select {
		case nc := <-c.register:
			c.addClient(nc)
		}
	}
}

func (c *Controller) addClient(nc *NewConnection) (err error) {
	room, ok := c.rooms[nc.room]
	if !ok {
		room, err = c.roomLookup(nc.room)
		if err != nil {
			return err
		}
	}
	newClient := NewClient(nc.conn, room, nc.hash)
	c.clients[newClient] = room
	room.addClient(newClient)
	newClient.run()
	return
}

// Performs a db lookup for an existing room
func (c *Controller) roomLookup(room string) (*Room, error) {
	r := NewRoom()
	c.rooms[room] = r
	go r.run()
	log.Printf("New room %s", room)
	return r, nil
}
