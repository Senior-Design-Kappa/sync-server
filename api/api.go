package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sync-server/config"
	"github.com/sync-server/logic"
)

type Server struct {
	*http.Server
	logic  logic.Logic
	Config config.Config
}

func NewServer(conf config.Config, logic logic.Logic) *Server {
	r := mux.NewRouter()

	s := &Server{
		Server: &http.Server{
			Handler:      r,
			Addr:         conf.Addr,
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		},
		logic:  logic,
		Config: conf,
	}
	r.HandleFunc("/health", health)
	return s
}

// health reports 200 if services is up and running
func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

// serveWs handles websocket requests from client
func serveWs(w http.ResponseWriter, r *http.Request) {
	//conn, err := upgrader.Upgrade(w, r, nil)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	// client.hub.register <- client
	// go client.writePump()
	// client.readPump()
}
