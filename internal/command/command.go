package command

// MoveCommand type represents a move command with a direction.
type MoveCommand struct {
	Direction string
}

// SayCommand type represents a say command with a message.
type SayCommand struct {
	Message string
}

// PickupCommand type represents a pickup command.
type PickupCommand struct {
	Identifier string
}
