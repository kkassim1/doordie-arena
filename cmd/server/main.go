package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/kkassim1/doordie-arena/internal/game"
	"github.com/kkassim1/doordie-arena/internal/transport/ws"
)

func main() {
	addr := ":8080" // later: read from env/config

	// Create a global match manager for the server.
	matchManager := game.NewMatchManager()

	// Create and start the default match loop.
	defaultMatch := matchManager.GetOrCreateMatch("default")

	// In a real server we'd use a cancellable context on shutdown.
	ctx := context.Background()
	go defaultMatch.Run(ctx.Done(), 50*time.Millisecond) // ~20 ticks/sec

	// Create WebSocket handler that knows about the match manager.
	wsHandler := ws.NewHandler(matchManager)

	// Register WebSocket endpoint
	http.HandleFunc("/ws", wsHandler.HandleWebSocket)

	log.Printf("Do or Die: Arena of Chaos server listening on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
