package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type Streamer struct {
	upgrader websocket.Upgrader
	clients  map[*websocket.Conn]bool
	mutex    *sync.Mutex
}

func NewStreamer() Streamer {
	return Streamer{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		clients: make(map[*websocket.Conn]bool, 0),
		mutex:   new(sync.Mutex),
	}
}

func (s *Streamer) Broadcast(msg []byte) {
	s.mutex.Lock()
	for c := range s.clients {
		err := c.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("Client disconnected:", err)
			delete(s.clients, c)
		}
	}
	s.mutex.Unlock()
}

func (s Streamer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	s.mutex.Lock()
	s.clients[c] = true
	s.mutex.Unlock()

	log.Println("Client connected:", r.RemoteAddr)
}
