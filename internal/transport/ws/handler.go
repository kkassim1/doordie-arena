package ws

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/kkassim1/doordie-arena/internal/game"
	"nhooyr.io/websocket"
)

// Handler holds dependencies for WebSocket connections.
type Handler struct {
	matchManager *game.MatchManager
}

// NewHandler creates a new WebSocket handler.
func NewHandler(mm *game.MatchManager) *Handler {
	return &Handler{
		matchManager: mm,
	}
}

// HandleWebSocket is the HTTP → WebSocket upgrade entrypoint.
func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
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

	// Connection-scoped state: the player and match for this client
	var player *game.Player
	var match *game.Match

	for {
		msgType, data, err := conn.Read(ctx)
		if err != nil {
			log.Printf("read error from %s: %v", r.RemoteAddr, err)
			return
		}

		if msgType != websocket.MessageText {
			log.Printf("ignoring non-text message from %s", r.RemoteAddr)
			continue
		}

		// Decode the incoming JSON message
		var incoming IncomingMessage
		if err := json.Unmarshal(data, &incoming); err != nil {
			log.Printf("invalid JSON from %s: %v", r.RemoteAddr, err)
			continue
		}

		switch incoming.Type {

		case "join":
			// Parse join payload
			var jp JoinPayload
			if err := json.Unmarshal(incoming.Payload, &jp); err != nil {
				log.Printf("invalid join payload from %s: %v", r.RemoteAddr, err)
				continue
			}

			if jp.Name == "" {
				jp.Name = "Player"
			}

			// Create a new player and add to the default match
			player = &game.Player{
				ID:   game.NewPlayerID(),
				Name: jp.Name,
				Type: game.PlayerTypeHuman,
				X:    0,
				Y:    0,
			}

			match = h.matchManager.GetOrCreateMatch("default")
			match.AddPlayer(player)

			log.Printf("player %s joined match %s", player.Name, match.ID)

			// Send welcome message back
			out := OutgoingMessage{
				Type: "welcome",
				Payload: WelcomePayload{
					PlayerID: string(player.ID),
					MatchID:  match.ID,
				},
			}

			if err := writeJSON(ctx, conn, out); err != nil {
				log.Printf("write welcome error to %s: %v", r.RemoteAddr, err)
				return
			}

		case "input":
			// Make sure they joined first
			if player == nil || match == nil {
				log.Printf("received input before join from %s", r.RemoteAddr)
				continue
			}

			// Parse input payload
			var ip InputPayload
			if err := json.Unmarshal(incoming.Payload, &ip); err != nil {
				log.Printf("invalid input payload from %s: %v", r.RemoteAddr, err)
				continue
			}

			// Let the game layer apply the input under lock
			match.ApplyInput(player.ID, ip.Input)

			// Build a state snapshot and send back
			snap := match.Snapshot()
			playerStates := make([]PlayerState, 0, len(snap))
			for _, ps := range snap {
				playerStates = append(playerStates, PlayerState{
					ID:   string(ps.ID),
					Name: ps.Name,
					X:    ps.X,
					Y:    ps.Y,
				})
			}

			stateMsg := OutgoingMessage{
				Type: "state",
				Payload: StatePayload{
					MatchID: match.ID,
					Players: playerStates,
				},
			}

			if err := writeJSON(ctx, conn, stateMsg); err != nil {
				log.Printf("write state error to %s: %v", r.RemoteAddr, err)
				return
			}

		default:
			// For now, log unknown message types.
			log.Printf("unknown message type %q from %s", incoming.Type, r.RemoteAddr)
		}
	}
}

// writeJSON writes a JSON message to the WebSocket with a timeout.
func writeJSON(ctx context.Context, conn *websocket.Conn, msg OutgoingMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	writeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return conn.Write(writeCtx, websocket.MessageText, data)
}
