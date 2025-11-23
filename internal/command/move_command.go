package command

import (
	"fmt"
	"strings"

	"github.com/xealgo/muddy/internal/game"
)

// Directions for movement
const (
	MoveDirNorth string = "north"
	MoveDirSouth string = "south"
	MoveDirEast  string = "east"
	MoveDirWest  string = "west"
)

// MoveCommand type represents a move command with a direction.
type MoveCommand struct {
	Choice string
}

// executeMoveCommand handles the execution of a move command.
func (cmd MoveCommand) Execute(game *game.Game, ps *game.Player) string {
	currentRoom, ok := game.World.GetRoomById(ps.CurrentRoomId)
	if !ok {
		return MessageInvalidMove
	}

	if !currentRoom.IsValidDoorChoice(cmd.Choice) {
		return MessageInvalidMove
	}

	for _, door := range currentRoom.Doors {
		if door.MoveCommand != cmd.Choice {
			continue
		}

		if door.IsLocked {
			return MessageDoorLocked
		}

		ps.CurrentRoomId = door.RoomId
	}

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf(MessageMoveSuccess, cmd.Choice))
	builder.WriteString("\nYou entered the ")

	currentRoom, ok = game.World.GetRoomById(ps.CurrentRoomId)
	if !ok {
		return "The void..no there is a bug here"
	}

	builder.WriteString(currentRoom.GetBasicInfo())
	builder.WriteByte('\n')

	return builder.String()
}
