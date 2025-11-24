package command

import (
	"strings"
	"time"

	"github.com/xealgo/muddy/internal/event"
	"github.com/xealgo/muddy/internal/game"
)

// SayCommand type represents a say command with a message.
type SayCommand struct {
	Message string
}

// Execute allows the player to say a message in the current room.
func (cmd SayCommand) Execute(game *game.Game, ps *game.Player) string {
	currentRoom, ok := game.World.GetRoomById(ps.CurrentRoomId)
	if !ok {
		return MessageInvalidCmd
	}

	e := event.EventDispatcher{}

	m := strings.TrimRight(cmd.Message, "\n")

	event := event.Event{
		Type:      "RoomChat",
		Timestamp: time.Now(),
		Data:      ps.DisplayName + ": " + m,
	}

	e.SendToRoom(event, game.Sm, currentRoom.ID)

	return ""
}
