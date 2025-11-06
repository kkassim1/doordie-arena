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
