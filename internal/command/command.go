package command

// Command represents a parsed command with its type and associated value.
type Command struct {
	Type  CommandType // The type of command such as "move"
	Value []string    // The associated values for commands that require more than 1 value
}
