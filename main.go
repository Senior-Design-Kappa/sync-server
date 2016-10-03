package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	addr          = "localhost:8000"
	indexTemplate = template.Must(template.ParseFiles("index.html"))
	upgrader      = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	c *Controller
)

func main() {
	c = NewController()
	go c.run()

	r := mux.NewRouter()

	r.HandleFunc("/", serveRoot)
	r.HandleFunc("/health", health)
	r.HandleFunc("/connect/{roomID}", handleConnection)

	s := http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}

func serveRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	indexTemplate.Execute(w, r.Host)
}

// health reports 200 if services is up and running
func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

// handleConnection handles websocket requests from client
func handleConnection(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	roomID := vals.Get("roomID")
	if roomID != "" {
		// handle invalid roomID
		fmt.Printf("Invalid roomID\n")
		http.Error(w, "Room not found", 404)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	nc := &NewConnection{conn, roomID}
	c.register <- nc
}
