package command

import (
	"github.com/xealgo/muddy/internal/game"
)

type InventoryCommand struct{}

// Execute lists the items in the player's inventory.
func (cmd InventoryCommand) Execute(game *game.Game, ps *game.Player) string {
	return ps.Inventory.List()
}
