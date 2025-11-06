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

	// Very simple 2D position for now
	X float64
	Y float64

	// Later: position, velocity, health, score, etc.
}
