package command

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

type CommandParseFunc = func(input string) (any, error)

// Parser handles command parsing.
type Parser struct {
	// Internal cached parser funcs
	parseFuncs []struct {
		typ CommandType
		fn  CommandParseFunc
	}
}

// NewParser creates a new Parser instance.
func NewParser() *Parser {
	p := &Parser{}
	p.parseFuncs = []struct {
		typ CommandType
		fn  CommandParseFunc
	}{
		{CommandMove, func(input string) (any, error) { return p.ParseMoveCommand(input) }},
		{CommandSay, func(input string) (any, error) { return p.ParseSayCommand(input) }},
		{CommandPickup, func(input string) (any, error) { return p.ParsePickupCommand(input) }},
	}

	return p
}

// ParseAnyCommand parses any command from the input string.
func (p *Parser) ParseAnyCommand(input string) (CommandType, interface{}, error) {
	if p.parseFuncs == nil {
		p = NewParser()
	}

	for _, pf := range p.parseFuncs {
		cmd, err := pf.fn(input)
		if err == nil {
			return pf.typ, cmd, nil
		}
	}

	return "", nil, fmt.Errorf("no valid command found")
}

// ParseMoveCommand parses a move command from the input string.
func (p *Parser) ParseMoveCommand(input string) (*MoveCommand, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	input = replaceNewlines(strings.TrimSpace(input))
	parts := strings.Split(input, " ")

	if len(parts) != 2 || parts[0] != string(CommandMove) {
		return nil, fmt.Errorf("invalid move command format")
	}

	dirs := GetMoveDirections()

	found := slices.Contains(dirs, parts[1])
	if !found {
		return nil, fmt.Errorf("invalid move direction: %s", parts[1])
	}

	cmd := MoveCommand{
		Direction: parts[1],
	}

	return &cmd, nil
}

// ParseSayCommand parses a say command from the input string.
func (p *Parser) ParseSayCommand(input string) (*SayCommand, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	input = replaceNewlines(input)
	parts := strings.SplitN(input, " ", 2)

	if len(parts) != 2 || parts[0] != string(CommandSay) {
		return nil, fmt.Errorf("invalid say command format")
	}

	if len(parts[1]) > 128 {
		return nil, fmt.Errorf("message too long: %d characters (max 128)", len(parts[1]))
	}

	cmd := SayCommand{
		Message: parts[1],
	}

	return &cmd, nil
}

// ParsePickupCommand parses a pickup command from the input string.
func (p *Parser) ParsePickupCommand(input string) (*PickupCommand, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	input = replaceNewlines(strings.TrimSpace(input))
	parts := strings.Split(input, " ")

	if len(parts) != 2 || parts[0] != string(CommandPickup) {
		return nil, fmt.Errorf("invalid pickup command format")
	}

	if len(parts[1]) > 32 {
		return nil, fmt.Errorf("invalid item identifier: %s", parts[1])
	}

	cmd := PickupCommand{
		Identifier: parts[1],
	}

	return &cmd, nil
}

// replaceNewlines replaces newline characters with spaces in the input string.
func replaceNewlines(input string) string {
	re := regexp.MustCompile(`(\r\n|\r|\n)+| +`)
	return re.ReplaceAllString(input, " ")
}
