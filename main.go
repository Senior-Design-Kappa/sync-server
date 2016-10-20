package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	addr     = "localhost:8000"
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	c *Controller
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	hashLength    = 64
)

func main() {
	c = NewController()
	go c.run()

	r := mux.NewRouter()

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

// health reports 200 if services is up and running
func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func generateHash(n int) string {
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

// handleConnection handles websocket requests from client
func handleConnection(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	vals := mux.Vars(r)
	roomID, ok := vals["roomID"]
	if !ok {
		// handle invalid roomID
		log.Printf("Invalid roomID\n")
		http.Error(w, "Room not found", 404)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	hash := generateHash(hashLength)
	nc := &NewConnection{conn, roomID, hash}
	c.register <- nc
}
