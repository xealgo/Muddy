package command

// MoveCommand type represents a move command with a direction.
type MoveCommand struct {
	Direction string
}

// SayCommand type represents a say command with a message.
type SayCommand struct {
	Message string
}

// LookCommand type represents a look command.
type LookCommand struct {
	//
}

// PickupCommand type represents a pickup command.
type PickupCommand struct {
	Identifier string
}
