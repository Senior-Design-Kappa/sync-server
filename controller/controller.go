/*
Controller handles new connections by registering them to correct rooms
*/
package controller

import (
	"log"

	"github.com/Senior-Design-Kappa/sync-server/models"
	"github.com/Senior-Design-Kappa/sync-server/room"
)

type Controller struct {
	clients map[*room.Client]*room.Room

	rooms map[string]*room.Room

	Register chan *models.NewConnection
}

func NewController() *Controller {
	c := &Controller{
		clients:  make(map[*room.Client]*room.Room),
		rooms:    make(map[string]*room.Room),
		Register: make(chan *models.NewConnection),
	}
	return c
}

func (c *Controller) Run() {
	for {
		select {
		case nc := <-c.Register:
			c.addClient(nc)
		}
	}
}

func (c *Controller) addClient(nc *models.NewConnection) (err error) {
	r, ok := c.rooms[nc.Room]
	if !ok {
		r, err = c.roomLookup(nc.Room)
		if err != nil {
			return err
		}
	}
	newClient := room.NewClient(nc.Conn, r, nc.Hash)
	c.clients[newClient] = r
	r.AddClient(newClient)
	newClient.Run()
	return
}

// Performs a db lookup for an existing room
func (c *Controller) roomLookup(roomID string) (*room.Room, error) {
	r := room.NewRoom()
	c.rooms[roomID] = r
	go r.Run()
	log.Printf("New room %s", roomID)
	return r, nil
}
