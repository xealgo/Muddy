package game

import (
	"fmt"

	"github.com/xealgo/muddy/internal/command"
	"github.com/xealgo/muddy/internal/session"
	"github.com/xealgo/muddy/internal/world"
)

type Game struct {
	Sm *session.SessionManager

	cp    *command.Parser
	state *GameState
	world *world.World
}

// NewGame creates a new Game instance.
func NewGame(world *world.World) *Game {
	g := &Game{
		state: NewGameState(),
		world: world,
		cp:    command.NewParser(),
	}

	return g
}

// State returns the current game state.
func (g Game) State() *GameState {
	return g.state
}

// ProcessPlayerCommand processes a command input from a player.
func (g Game) ProcessPlayerCommand(ps *session.PlayerSession, input string) string {
	ctype, cmd, err := g.cp.ParseAnyCommand(input)
	if err != nil {
		return fmt.Sprintln(err.Error())
	}

	return g.ExecuteCommand(ps, ctype, cmd)
}
