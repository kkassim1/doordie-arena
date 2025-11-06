package game

import "sync"

// MatchManager holds all running matches.
// For now we'll just use one match called "default".
type MatchManager struct {
	mu      sync.RWMutex
	matches map[string]*Match
}

func NewMatchManager() *MatchManager {
	return &MatchManager{
		matches: make(map[string]*Match),
	}
}

// GetOrCreateMatch returns an existing match or creates it if missing.
func (mm *MatchManager) GetOrCreateMatch(id string) *Match {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	if match, ok := mm.matches[id]; ok {
		return match
	}

	match := NewMatch(id)
	mm.matches[id] = match
	return match
}
