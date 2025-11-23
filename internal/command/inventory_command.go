package command

import (
	"github.com/xealgo/muddy/internal/game"
)

type InventoryCommand struct{}

// executeLookCommand handles the execution of a look command.
func (cmd InventoryCommand) Execute(game *game.Game, ps *game.Player) string {
	return ps.Inventory.List()
}
