package command

import (
	"fmt"

	"github.com/xealgo/muddy/internal/game"
)

// Runner executes commands for the game.
type Runner struct {
	game   *game.Game
	parser *Parser
}

// NewRunner creates a new command runner.
func NewRunner(game *game.Game) *Runner {
	return &Runner{
		game:   game,
		parser: NewParser(),
	}
}

// Execute processes and executes commands based on input from the players.
func (r Runner) Execute(ps *game.Player, input string) string {
	if len(input) == 0 || input == "\n" || input == "\r" {
		return ""
	}

	_, cmd, err := r.parser.ParseAnyCommand(input)
	if err != nil {
		return fmt.Sprintln(err.Error())
	}

	return cmd.Execute(r.game, ps)
}
