package game

import "time"

// GameState represents the current state of the game.
type GameState struct {
	uptime    uint32
	startTime time.Time
}

// NewGameState creates a new GameState instance.
func NewGameState() *GameState {
	return &GameState{
		uptime:    0,
		startTime: time.Now(),
	}
}

// Uptime returns the duration since the game started.
func (gs GameState) Uptime() time.Duration {
	return time.Since(gs.startTime)
}
