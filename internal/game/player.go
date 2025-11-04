package game

import "github.com/google/uuid"

// PlayerID is a unique identifier for each player (human or bot).
type PlayerID string

func NewPlayerID() PlayerID {
	return PlayerID(uuid.NewString())
}

type PlayerType int

const (
	PlayerTypeHuman PlayerType = iota
	PlayerTypeBot
)

type Player struct {
	ID   PlayerID
	Name string
	Type PlayerType
	// Later: position, velocity, health, score, etc.
}
