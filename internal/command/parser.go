package command

import (
	"fmt"
	"slices"
	"strings"
)

// Parser handles command parsing.
type Parser struct {
	//
}

// ParseMoveCommand parses a move command from the input string.
func (p *Parser) ParseMoveCommand(input string) (*Command, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	cmd := Command{}
	parts := strings.Split(input, " ")

	if len(parts) > 2 || parts[0] != string(Move) {
		return nil, fmt.Errorf("invalid move command format")
	}

	dirs := GetMoveDirections()

	found := slices.Contains(dirs, parts[1])
	if !found {
		return nil, fmt.Errorf("invalid move direction: %s", parts[1])
	}

	cmd.Type = Move
	cmd.Value = []string{parts[1]}

	return &cmd, nil
}
