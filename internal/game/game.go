package game

import (
	"log/slog"
	"strings"
)

type Game struct {
	World *World
	Sm    *SessionManager
	state *GameState
}

// NewGame creates a new Game instance.
func NewGame(world *World) *Game {
	g := &Game{
		state: NewGameState(),
		World: world,
	}

	return g
}

// State returns the current game state.
func (g Game) State() *GameState {
	return g.state
}

// GreetPlayer sends a greeting message to the player upon joining the game.
func (g Game) GreetPlayer(ps *Player) {
	startingRoom, ok := g.World.GetRoomById(1)
	if !ok {
		slog.Error("Could not access room 1")
		return
	}

	builder := strings.Builder{}
	builder.WriteString("Greetings ")
	builder.WriteString(ps.DisplayName)
	builder.WriteString("!\nYou seem to find your self in ")
	builder.WriteString(startingRoom.Name)
	builder.WriteString(".\n")
	builder.WriteString(startingRoom.Description)

	ps.WriteString(builder.String())
}
