package command

import (
	"github.com/xealgo/muddy/internal/game"
)

type CommandType string

// Supported command types
const (
	CommandUnknown   CommandType = "unknown"   // unknown command
	CommandHelp      CommandType = "help"      // provides help information about available commands
	CommandMove      CommandType = "move"      // move north - move in one of 4 directions
	CommandLook      CommandType = "look"      // Tells the player what they can see in the room
	CommandPickup    CommandType = "pickup"    // pickup {item-name} - adds an item to the player's inventory
	CommandInventory CommandType = "inventory" // reports what's in the player's inventory
	CommandSay       CommandType = "say"       // say hello everyone! broadcasts a chat message to everyone in the room
)

// Command interface for executing commands
type Command interface {
	Execute(game *game.Game, ps *game.Player) string
}
