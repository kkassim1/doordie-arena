package main

import (
	"log"
	"net/http"

	"github.com/kkassim1/doordie-arena/internal/transport/ws"
)

func main() {
	addr := ":8080" // later: read from env/config

	// Register WebSocket endpoint
	http.HandleFunc("/ws", ws.HandleWebSocket)

	log.Printf("Do or Die: Arena of Chaos server listening on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
