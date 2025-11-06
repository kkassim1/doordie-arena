package ws

import "encoding/json"

// IncomingMessage is what the client sends.
type IncomingMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// JoinPayload is the data for a "join" message.
type JoinPayload struct {
	Name string `json:"name"`
}

// InputPayload is the data for a "input" message (like move, jump, etc).
type InputPayload struct {
	Input string `json:"input"` // e.g. "move_up", "move_down", "move_left", "move_right"
}

// OutgoingMessage is what the server sends.
type OutgoingMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
}

// WelcomePayload is sent back after a successful join.
type WelcomePayload struct {
	PlayerID string `json:"playerId"`
	MatchID  string `json:"matchId"`
}

// PlayerState is what the client sees for each player.
type PlayerState struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
}

// StatePayload is the payload for a "state" message.
type StatePayload struct {
	MatchID string        `json:"matchId"`
	Players []PlayerState `json:"players"`
}
