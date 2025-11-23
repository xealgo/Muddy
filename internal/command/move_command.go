package command

import (
	"fmt"

	"github.com/xealgo/muddy/internal/game"
	"github.com/xealgo/muddy/internal/session"
	"github.com/xealgo/muddy/internal/world"
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
	Direction string
}

// executeMoveCommand handles the execution of a move command.
func (cmd MoveCommand) Execute(game *game.Game, ps *session.PlayerSession) string {
	currentRoom, ok := game.World.GetRoomById(ps.GetData().CurrentRoomId)
	if !ok {
		return MessageInvalidMove
	}

	switch cmd.Direction {
	case MoveDirNorth:
		return doMove(ps, currentRoom.Exits.North, &cmd)
	case MoveDirSouth:
		return doMove(ps, currentRoom.Exits.South, &cmd)
	case MoveDirEast:
		return doMove(ps, currentRoom.Exits.East, &cmd)
	case MoveDirWest:
		return doMove(ps, currentRoom.Exits.West, &cmd)
	default:
		return MessageInvalidMove
	}
}

// GetMoveDirections returns a list of valid move directions.
func GetMoveDirections() []string {
	return []string{MoveDirNorth, MoveDirSouth, MoveDirEast, MoveDirWest}
}

// doMove performs the move action for the player session.
func doMove(ps *session.PlayerSession, door *world.Door, cmd *MoveCommand) string {
	if door == nil {
		return MessageInvalidMove
	}

	if door.IsLocked {
		return MessageDoorLocked
	}

	ps.GetData().CurrentRoomId = door.RoomId

	return fmt.Sprintf(MessageMoveSuccess, cmd.Direction)
}
