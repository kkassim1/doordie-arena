package game

import (
	"log"
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

// ApplyInput applies a simple input (move_up, move_down, etc) to a player.
// This keeps all game-state mutation under the match lock.
func (m *Match) ApplyInput(playerID PlayerID, input string) {
	const step = 1.0

	m.mu.Lock()
	defer m.mu.Unlock()

	p, ok := m.Players[playerID]
	if !ok {
		return
	}

	switch input {
	case "move_up":
		p.Y += step
	case "move_down":
		p.Y -= step
	case "move_left":
		p.X -= step
	case "move_right":
		p.X += step
	default:
		// unknown input: ignore for now
	}
}

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

// Run starts the match tick loop at the given tickRate (e.g. 50ms).
// In the future this is where we'll update bots, tasks, physics, etc.
func (m *Match) Run(ctxDone <-chan struct{}, tickRate time.Duration) {
	ticker := time.NewTicker(tickRate)
	defer ticker.Stop()

	log.Printf("match %s loop started with tickRate=%s", m.ID, tickRate)

	for {
		select {
		case <-ctxDone:
			log.Printf("match %s loop stopped", m.ID)
			return
		case t := <-ticker.C:
			m.tick(t)
		}
	}
}

// tick is a single simulation step for the match.
func (m *Match) tick(now time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// TODO: update bots, tasks, timeouts, etc.
	// For now we just keep this as a placeholder.
	// log.Printf("match %s tick at %s", m.ID, now.Format(time.RFC3339Nano))
	_ = now
}
