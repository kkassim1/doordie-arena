package ws

import (
	"context"
	"log"
	"net/http"
	"time"

	"nhooyr.io/websocket"
)

// For now this is a super simple echo handler so we can verify
// Godot <-> Go connectivity. Later we’ll plug this into the game.Match.
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Accept the WebSocket connection
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"}, // TODO: tighten this in production
	})
	if err != nil {
		log.Printf("websocket accept error: %v", err)
		return
	}
	defer func() {
		_ = conn.Close(websocket.StatusInternalError, "internal error")
	}()

	log.Printf("new client connected from %s", r.RemoteAddr)

	ctx := r.Context()

	for {
		// Read a message from the client
		msgType, data, err := conn.Read(ctx)
		if err != nil {
			log.Printf("read error from %s: %v", r.RemoteAddr, err)
			return
		}

		log.Printf("recv from %s: %s", r.RemoteAddr, string(data))

		// For now: echo it back
		writeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = conn.Write(writeCtx, msgType, data)
		cancel()

		if err != nil {
			log.Printf("write error to %s: %v", r.RemoteAddr, err)
			return
		}
	}
}
