package command

import (
	"fmt"

	"github.com/xealgo/muddy/internal/game"
)

// SellCommand allows a player to sell an item to a merchant NPC in the room.
type SellCommand struct {
	Target string
	ItemID string
}

// Execute allows the player to talk to an NPC in the current room.
func (cmd SellCommand) Execute(g *game.Game, ps *game.Player) string {
	currentRoom, ok := g.World.GetRoomById(ps.CurrentRoomId)
	if !ok {
		return MessageInvalidCmd
	}

	npc, ok := currentRoom.GetNpcByName(cmd.Target)
	if !ok {
		return "There is no such NPC here to sell to."
	}

	merchant, ok := npc.(*game.Merchant)
	if !ok {
		return "You can only sell items to merchants."
	}

	item, ok := ps.Inventory.Sell(cmd.ItemID, merchant)
	if !ok || item == nil {
		return "You don't have that item to sell."
	}

	return fmt.Sprintf("You sold the item %s to %s for $$%d.\n", item.Name, merchant.Name, item.SellingPrice)
}
