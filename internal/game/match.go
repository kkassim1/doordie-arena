package game

import (
	"sync"
	"time"
)

// PlayerSnapshot is a read-only copy of player state used for broadcasting.
type PlayerSnapshot struct {
	ID   PlayerID
	Name string
	X    float64
	Y    float64
}

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
// Snapshot returns a copy of all players for safe read access.
func (m *Match) Snapshot() []PlayerSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]PlayerSnapshot, 0, len(m.Players))
	for _, p := range m.Players {
		out = append(out, PlayerSnapshot{
			ID:   p.ID,
			Name: p.Name,
			X:    p.X,
			Y:    p.Y,
		})
	}
	return out
}
