package command

import (
	"github.com/xealgo/muddy/internal/game"
	"github.com/xealgo/muddy/internal/session"
)

// PickupCommand type represents a pickup command.
type PickupCommand struct {
	Identifier string
}

// executeLookCommand handles the execution of a look command.
func (cmd PickupCommand) Execute(game *game.Game, ps *session.PlayerSession) string {
	currentRoom, ok := game.World.GetRoomById(ps.GetData().CurrentRoomId)
	if !ok {
		return MessageInvalidCmd
	}

	_ = currentRoom

	return "Pickup command not yet available"
}
