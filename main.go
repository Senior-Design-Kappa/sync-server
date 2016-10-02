package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

var (
	addr          = "localhost:8000"
	indexTemplate = template.Must(template.ParseFiles("index.html"))
)

func main() {
	// hub := newHub()
	// go hub.run()
	r := mux.NewRouter()

	r.HandleFunc("/", serveRoot)
	r.HandleFunc("/health", health)
	r.HandleFunc("/connect", serveWs)

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

// serveWs handles websocket requests from client
func serveWs(w http.ResponseWriter, r *http.Request) {
	// conn, err := upgrader.Upgrade(w, r, nil)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// fmt.Println("client connected")
	// client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	// client.hub.register <- client
	// go client.writePump()
	// client.readPump()
}
