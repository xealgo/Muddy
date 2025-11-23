package command

import (
	"strings"

	"github.com/xealgo/muddy/internal/game"
	"github.com/xealgo/muddy/internal/session"
)

type HelpCommand struct{}

// Execute provides a list of available commands to the player.
func (h HelpCommand) Execute(g *game.Game, ps *session.PlayerSession) string {
	builder := strings.Builder{}

	builder.WriteString("The following commands are available\n")
	builder.WriteString("- look: Describe your surroundings\n")
	builder.WriteString("- move <direction>: Move in a direction (north, south, east, west)\n")
	builder.WriteString("- say <message>: Send a message to other players in the same room\n")
	builder.WriteString("- help: Show this help message\n")
	return builder.String()
}
