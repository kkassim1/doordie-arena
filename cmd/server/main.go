package main

import (
	"log"
	"net/http"

	"github.com/kkassim1/doordie-arena/internal/game"
	"github.com/kkassim1/doordie-arena/internal/transport/ws"
)

func main() {
	addr := ":8080" // later: read from env/config

	// Create a global match manager for the server.
	matchManager := game.NewMatchManager()

	// Create WebSocket handler that knows about the match manager.
	wsHandler := ws.NewHandler(matchManager)

	// Register WebSocket endpoint
	http.HandleFunc("/ws", wsHandler.HandleWebSocket)

	log.Printf("Do or Die: Arena of Chaos server listening on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
