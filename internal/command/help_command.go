package command

import (
	"strings"

	"github.com/xealgo/muddy/internal/game"
)

type HelpCommand struct{}

// Execute provides a list of available commands to the player.
func (h HelpCommand) Execute(g *game.Game, ps *game.Player) string {
	builder := strings.Builder{}

	builder.WriteString("The following commands are available\n")
	builder.WriteString("- look: Describe your surroundings\n")
	builder.WriteString("- move <direction>: Move in a direction (north, south, east, west)\n")
	builder.WriteString("- say <message>: Send a message to other players in the same room\n")
	builder.WriteString("- help: Show this help message\n")
	builder.WriteString("- sell <merchant name> <item name>: Sell an inventory item\n")
	builder.WriteString("- talk <merchant name>: Talk to an NPC\n")

	return builder.String()
}
