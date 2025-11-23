package command

import (
	"fmt"
	"regexp"
	"strings"
)

type CommandParseFunc = func(input string) (Command, error)

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
		{CommandMove, func(input string) (Command, error) { return p.ParseMoveCommand(input) }},
		{CommandSay, func(input string) (Command, error) { return p.ParseSayCommand(input) }},
		{CommandPickup, func(input string) (Command, error) { return p.ParsePickupCommand(input) }},
		{CommandLook, func(input string) (Command, error) { return p.ParseLookCommand(input) }},
		{CommandHelp, func(input string) (Command, error) { return p.ParseHelpCommand(input) }},
		{CommandInventory, func(input string) (Command, error) { return p.ParseInventoryCommand(input) }},
	}

	return p
}

// ParseAnyCommand parses any command from the input string.
func (p *Parser) ParseAnyCommand(input string) (CommandType, Command, error) {
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

// ParseHelpCommand parses a help command from the input string.
func (p Parser) ParseHelpCommand(input string) (*HelpCommand, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	input = replaceNewlines(strings.TrimSpace(input))
	parts := strings.Split(input, " ")

	if len(parts) != 1 || parts[0] != string(CommandHelp) {
		return nil, fmt.Errorf("invalid help command format")
	}

	cmd := HelpCommand{}

	return &cmd, nil
}

// ParseMoveCommand parses a move command from the input string.
func (p Parser) ParseMoveCommand(input string) (*MoveCommand, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	input = replaceNewlines(strings.ToLower(strings.TrimSpace(input)))
	parts := strings.SplitN(input, " ", 2)

	if len(parts) != 2 || parts[0] != string(CommandMove) {
		return nil, fmt.Errorf("invalid move command format")
	}

	cmd := MoveCommand{
		Choice: strings.TrimSpace(parts[1]),
	}

	return &cmd, nil
}

// ParseSayCommand parses a say command from the input string.
func (p Parser) ParseSayCommand(input string) (*SayCommand, error) {
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
func (p Parser) ParsePickupCommand(input string) (*PickupCommand, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	input = replaceNewlines(strings.TrimSpace(input))
	parts := strings.SplitN(input, " ", 2)

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

// ParseLookCommand parses a look command from the input string.
func (p Parser) ParseLookCommand(input string) (*LookCommand, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	input = replaceNewlines(strings.TrimSpace(input))
	parts := strings.Split(input, " ")

	if len(parts) != 1 || parts[0] != string(CommandLook) {
		return nil, fmt.Errorf("invalid look command format")
	}

	cmd := LookCommand{}

	return &cmd, nil
}

// ParseInventoryCommand parses an inventory command from the input string.
func (p Parser) ParseInventoryCommand(input string) (*InventoryCommand, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	input = replaceNewlines(strings.TrimSpace(input))
	parts := strings.Split(input, " ")

	if len(parts) != 1 || parts[0] != string(CommandInventory) {
		return nil, fmt.Errorf("invalid inventory command format")
	}

	cmd := InventoryCommand{}

	return &cmd, nil
}

// replaceNewlines replaces newline characters with spaces in the input string.
func replaceNewlines(input string) string {
	re := regexp.MustCompile(`(\r\n|\r|\n)+| +`)
	return re.ReplaceAllString(input, " ")
}
