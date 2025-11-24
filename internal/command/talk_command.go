package command

import "github.com/xealgo/muddy/internal/game"

type TalkCommand struct {
	Target string
}

// Execute allows the player to talk to an NPC in the current room.
func (cmd TalkCommand) Execute(game *game.Game, ps *game.Player) string {
	currentRoom, ok := game.World.GetRoomById(ps.CurrentRoomId)
	if !ok {
		return MessageInvalidCmd
	}

	npc, ok := currentRoom.GetNpcByName(cmd.Target)
	if ok {
		return npc.Greet(ps)
	}

	return ""
}
