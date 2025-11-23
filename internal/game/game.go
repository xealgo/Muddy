package game

import (
	"log/slog"
	"strings"

	"github.com/xealgo/muddy/internal/session"
	"github.com/xealgo/muddy/internal/world"
)

type Game struct {
	World *world.World
	Sm    *session.SessionManager
	state *GameState
}

// NewGame creates a new Game instance.
func NewGame(world *world.World) *Game {
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
func (g Game) GreetPlayer(ps *session.PlayerSession) {
	startingRoom, ok := g.World.GetRoomById(1)
	if !ok {
		slog.Error("Could not access room 1")
		return
	}

	builder := strings.Builder{}
	builder.WriteString("Greetings ")
	builder.WriteString(ps.GetData().DisplayName)
	builder.WriteString("!\nYou seem to find your self in ")
	builder.WriteString(startingRoom.Name)
	builder.WriteString(".\n")
	builder.WriteString(startingRoom.Description)

	ps.WriteString(builder.String())
}
