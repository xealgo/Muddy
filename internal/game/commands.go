package game

import (
	"fmt"

	"github.com/xealgo/muddy/internal/command"
	"github.com/xealgo/muddy/internal/session"
	"github.com/xealgo/muddy/internal/world"
)

const (
	MessageInvalidCmd  string = "You can't do that"
	MessageInvalidMove string = "You can't move there"
	MessageDoorLocked  string = "The door seems to be locked"
	MessageMoveSuccess string = "You move to the %s"
)

// ExecuteCommand executes a command based on its type and data.
func (g Game) ExecuteCommand(ps *session.PlayerSession, cmdType command.CommandType, data interface{}) string {
	switch cmdType {
	case command.CommandMove:
		moveCmd, ok := data.(*command.MoveCommand)
		if !ok {
			return MessageInvalidMove
		}
		return g.executeMoveCommand(ps, moveCmd)
	default:
		return MessageInvalidCmd
	}
}

// executeMoveCommand handles the execution of a move command.
func (g Game) executeMoveCommand(ps *session.PlayerSession, cmd *command.MoveCommand) string {
	currentRoom, ok := g.world.GetRoomById(ps.GetData().CurrentRoomId)
	if !ok {
		return MessageInvalidMove
	}

	switch cmd.Direction {
	case command.MoveDirNorth:
		return doMove(ps, currentRoom.Exits.North, cmd)
	case command.MoveDirSouth:
		return doMove(ps, currentRoom.Exits.South, cmd)
	case command.MoveDirEast:
		return doMove(ps, currentRoom.Exits.East, cmd)
	case command.MoveDirWest:
		return doMove(ps, currentRoom.Exits.West, cmd)
	default:
		return MessageInvalidMove
	}
}

// doMove performs the move action for the player session.
func doMove(ps *session.PlayerSession, door *world.Door, cmd *command.MoveCommand) string {
	if door == nil {
		return MessageInvalidMove
	}

	if door.IsLocked {
		return MessageDoorLocked
	}

	ps.GetData().CurrentRoomId = door.RoomId

	return fmt.Sprintf(MessageMoveSuccess, cmd.Direction)
}
