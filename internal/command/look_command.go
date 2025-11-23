package command

import (
	"strings"

	"github.com/xealgo/muddy/internal/game"
)

// LookCommand type represents a look command.
type LookCommand struct {
	//
}

// executeLookCommand handles the execution of a look command.
func (cmd LookCommand) Execute(game *game.Game, ps *game.Player) string {
	currentRoom, ok := game.World.GetRoomById(ps.CurrentRoomId)
	if !ok {
		return MessageInvalidCmd
	}

	builder := strings.Builder{}

	builder.WriteString("You look around the room\n")
	builder.WriteString("You see ")
	builder.WriteString(currentRoom.GetDetails(ps, game.Sm))

	return builder.String()
}
