package game

// Door represents a door leading to another room
type Door struct {
	Name        string `yaml:"name"`        // Name of the door
	Description string `yaml:"description"` // Description of the door
	MoveCommand string `yaml:"moveCommand"` // Command to move through the door
	IsLocked    bool   `yaml:"isLocked"`    // Is the door locked?
	RoomId      int    `yaml:"roomId"`      // The room this door leads to
}

// String returns the name of the door
func (door Door) String() string {
	return door.Name
}

// Validate checks if the door has valid attributes
func (door Door) Validate() bool {
	if door.Name == "" || door.MoveCommand == "" || door.RoomId < 0 {
		return false
	}

	return true
}
