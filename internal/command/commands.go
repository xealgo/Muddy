package command

type CommandType string

// Supported command types
const (
	CommandUnknown   CommandType = "unknown"   // unknown command
	CommandMove      CommandType = "move"      // move north - move in one of 4 directions
	CommandLook      CommandType = "look"      // Tells the player what they can see in the room
	CommandPickup    CommandType = "pickup"    // pickup {item-name} - adds an item to the player's inventory
	CommandInventory CommandType = "inventory" // reports what's in the player's inventory
	CommandSay       CommandType = "say"       // say hello everyone! broadcasts a chat message to everyone in the room
)

// Directions for movement
const (
	MoveDirNorth string = "north"
	MoveDirSouth string = "south"
	MoveDirEast  string = "east"
	MoveDirWest  string = "west"
)

// GetMoveDirections returns a list of valid move directions.
func GetMoveDirections() []string {
	return []string{MoveDirNorth, MoveDirSouth, MoveDirEast, MoveDirWest}
}
