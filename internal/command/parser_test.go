package command

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAnyCommand(t *testing.T) {
	type CommandTest struct {
		input       string
		expected    any
		ExpectError bool
	}

	tests := []CommandTest{
		{input: "move north", expected: &MoveCommand{Direction: "north"}},
		{input: "say hello world!!", expected: &SayCommand{Message: "hello world!!"}},
		{input: "pickup lizards", expected: &PickupCommand{Identifier: "lizards"}},
	}

	p := Parser{}

	for _, test := range tests {
		typ, cmd, err := p.ParseAnyCommand(test.input)

		if test.expected != nil && test.ExpectError == false {
			assert.Nil(t, err)

			switch typ {
			case CommandMove:
				expectedCmd := test.expected.(*MoveCommand)
				actualCmd := cmd.(*MoveCommand)
				assert.Equal(t, expectedCmd.Direction, actualCmd.Direction)
			case CommandSay:
				expectedCmd := test.expected.(*SayCommand)
				actualCmd := cmd.(*SayCommand)
				assert.Equal(t, expectedCmd.Message, actualCmd.Message)
			case CommandPickup:
				expectedCmd := test.expected.(*PickupCommand)
				actualCmd := cmd.(*PickupCommand)
				assert.Equal(t, expectedCmd.Identifier, actualCmd.Identifier)
			}
		}

		if test.ExpectError {
			assert.NotNil(t, err)
			assert.Nil(t, cmd)
		}
	}
}

func TestParseMoveCommand(t *testing.T) {
	type CommandTest struct {
		input       string
		expected    *MoveCommand
		ExpectError bool
	}

	tests := []CommandTest{
		{input: "move north", expected: &MoveCommand{Direction: "north"}},
		{input: "move  north", expected: &MoveCommand{Direction: "north"}},
		{input: "move     north", expected: &MoveCommand{Direction: "north"}},
		{input: " move north ", expected: &MoveCommand{Direction: "north"}},
		{input: "move\nnorth", expected: &MoveCommand{Direction: "north"}},
		{input: "move south", expected: &MoveCommand{Direction: "south"}},
		{input: "move east", expected: &MoveCommand{Direction: "east"}},
		{input: "move west", expected: &MoveCommand{Direction: "west"}},
		{input: "mehhh north", expected: nil, ExpectError: true},
		{input: "move", expected: nil, ExpectError: true},
		{input: "move out_the_way", expected: nil, ExpectError: true},
		{input: "move south west", expected: nil, ExpectError: true},
	}

	p := Parser{}

	for _, test := range tests {
		cmd, err := p.ParseMoveCommand(test.input)

		if test.expected != nil && test.ExpectError == false {
			assert.Nil(t, err)
			assert.NotNil(t, cmd)
			assert.Equal(t, cmd.Direction, test.expected.Direction)
		}

		if test.ExpectError {
			assert.NotNil(t, err)
			assert.Nil(t, cmd)
		}
	}
}

func TestSayCommand(t *testing.T) {
	type CommandTest struct {
		input       string
		expected    *SayCommand
		ExpectError bool
	}

	longStr := strings.Repeat("a", 129)
	tests := []CommandTest{
		{input: "say hello world!!", expected: &SayCommand{Message: "hello world!!"}},
		{input: "say hey! I just began playing this!", expected: &SayCommand{Message: "hey! I just began playing this!"}},
		{input: "say hello\nworld!!", expected: &SayCommand{Message: "hello world!!"}},
		{input: "say hello\n\rworld!!", expected: &SayCommand{Message: "hello world!!"}},
		{input: "say hello\n\n\n\n\r\r\rworld!!", expected: &SayCommand{Message: "hello world!!"}},
		{input: "sayy hello world!!", expected: nil, ExpectError: true},
		{input: "say " + longStr, expected: nil, ExpectError: true},
	}

	p := Parser{}

	for _, test := range tests {
		cmd, err := p.ParseSayCommand(test.input)

		if test.expected != nil && test.ExpectError == false {
			assert.Nil(t, err)
			assert.NotNil(t, cmd)
			assert.Equal(t, cmd.Message, test.expected.Message)
			// fmt.Printf("%s: %v\n", cmd.Type, cmd.Value)
		}

		if test.ExpectError {
			assert.NotNil(t, err)
			assert.Nil(t, cmd)
		}
	}
}

func TestPickupMoveCommand(t *testing.T) {
	type CommandTest struct {
		input       string
		expected    *PickupCommand
		ExpectError bool
	}

	tests := []CommandTest{
		{input: "pickup lizards", expected: &PickupCommand{Identifier: "lizards"}},
		{input: "pickup  spock", expected: &PickupCommand{Identifier: "spock"}},
		{input: "pickup     rock", expected: &PickupCommand{Identifier: "rock"}},
		{input: " pickup paper ", expected: &PickupCommand{Identifier: "paper"}},
		{input: "pickup\nscissors", expected: &PickupCommand{Identifier: "scissors"}},
		{input: "pickups lizards", expected: nil, ExpectError: true},
		{input: "pickup lizards and spocks", expected: nil, ExpectError: true},
	}

	p := Parser{}

	for _, test := range tests {
		cmd, err := p.ParsePickupCommand(test.input)

		if test.expected != nil && test.ExpectError == false {
			assert.Nil(t, err)
			assert.NotNil(t, cmd)
			assert.Equal(t, cmd.Identifier, test.expected.Identifier)
		}

		if test.ExpectError {
			assert.NotNil(t, err)
			assert.Nil(t, cmd)
		}
	}
}
