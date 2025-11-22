package command

type CommandType string

// Supported command types
const (
	Move      CommandType = "move"      // move north - move in one of 4 directions
	Look      CommandType = "look"      // Tells the player what they can see in the room
	Pickup    CommandType = "pickup"    // pickup {item-name} - adds an item to the player's inventory
	Inventory CommandType = "inventory" // reports what's in the player's inventory
	Say       CommandType = "say"       // say hello everyone! broadcasts a chat message to everyone in the room
)

// Directions for movement
const (
	MoveNorth string = "north"
	MoveSouth string = "south"
	MoveEast  string = "east"
	MoveWest  string = "west"
)

// GetMoveDirections returns a list of valid move directions.
func GetMoveDirections() []string {
	return []string{MoveNorth, MoveSouth, MoveEast, MoveWest}
}
