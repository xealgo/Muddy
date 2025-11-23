package command

import (
	"github.com/xealgo/muddy/internal/game"
)

const (
	MessageItemNotFound = "There is no such item here."
)

// PickupCommand type represents a pickup command.
type PickupCommand struct {
	Identifier string
}

// executeLookCommand handles the execution of a look command.
func (cmd PickupCommand) Execute(game *game.Game, ps *game.Player) string {
	currentRoom, ok := game.World.GetRoomById(ps.CurrentRoomId)
	if !ok {
		return MessageInvalidCmd
	}

	item, ok := currentRoom.RemoveItem(cmd.Identifier)
	if !ok {
		return MessageItemNotFound
	}

	ps.Inventory.Add(item)

	return "You picked up the " + item.Name + "."
}
