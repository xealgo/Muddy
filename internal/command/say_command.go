package command

import (
	"strings"
	"time"

	"github.com/xealgo/muddy/internal/event"
	"github.com/xealgo/muddy/internal/game"
	"github.com/xealgo/muddy/internal/session"
)

// SayCommand type represents a say command with a message.
type SayCommand struct {
	Message string
}

// executeLookCommand handles the execution of a look command.
func (cmd SayCommand) Execute(game *game.Game, ps *session.PlayerSession) string {
	currentRoom, ok := game.World.GetRoomById(ps.GetData().CurrentRoomId)
	if !ok {
		return MessageInvalidCmd
	}

	e := event.EventDispatcher{}

	m := strings.TrimRight(cmd.Message, "\n")

	event := event.Event{
		Type:      "RoomChat",
		Timestamp: time.Now(),
		Data:      ps.GetData().DisplayName + ": " + m,
	}

	e.SendToRoom(event, game.Sm, currentRoom.ID)

	return ""
}
