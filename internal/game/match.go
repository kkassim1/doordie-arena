package game

import (
	"sync"
	"time"
)

// Match represents a running arena instance.
type Match struct {
	ID        string
	Players   map[PlayerID]*Player
	Tasks     map[TaskID]*Task
	StartTime time.Time

	mu sync.RWMutex
	// Later: tick rate, world state, bots, etc.
}

func NewMatch(id string) *Match {
	return &Match{
		ID:        id,
		Players:   make(map[PlayerID]*Player),
		Tasks:     make(map[TaskID]*Task),
		StartTime: time.Now(),
	}
}

func (m *Match) AddPlayer(p *Player) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Players[p.ID] = p
}

// Next steps: RemovePlayer, AssignRandomTask, Tick(), etc.
