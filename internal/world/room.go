package world

import (
	"fmt"
	"strings"
)

// Door represents a door leading to another room
type Door struct {
	IsLocked bool `json:"isLocked" yaml:"isLocked"` // Is the door locked?
	RoomId   int  `json:"roomId" yaml:"roomId"`     // The room this door leads to
}

// RoomExits holds the exits for a room
type RoomExits struct {
	North *Door `json:"north" yaml:"north"`
	South *Door `json:"south" yaml:"south"`
	West  *Door `json:"west" yaml:"east"`
	East  *Door `json:"east" yaml:"west"`
}

// GetExits returns a string listing the available exits and their count
func (exits RoomExits) GetExits() (string, int) {
	count := 0

	builder := strings.Builder{}
	doors := []*Door{exits.North, exits.South, exits.East, exits.West}
	labels := []string{"North", "South", "East", "West"}

	for direction, door := range doors {
		if door == nil {
			continue
		}

		count++
		label := labels[direction]
		builder.WriteString(label + "\n")
	}

	return builder.String(), count
}

// Room represents a room in the game world
type Room struct {
	ID          int       `json:"id" yaml:"id"`
	Name        string    `json:"name" yaml:"name"`
	Description string    `json:"description" yaml:"description"`
	Exits       RoomExits `json:"exits" yaml:"exits"`
	Items       []Item    `json:"items" yaml:"items"` // Shared items
	// Players     []RoomPlayer `json:"-"`     // Will be mapped in the service layer
}

// NewRoom creates a new Room instance
func NewRoom(id int, name string, desc string) *Room {
	return &Room{
		ID:          id,
		Name:        name,
		Description: desc,
		Items:       make([]Item, 0),
		// Players:     make([]RoomPlayer, 0),
	}
}

// GetBasicInfo returns the basic information of the room
func (room Room) GetBasicInfo() string {
	builder := strings.Builder{}

	builder.WriteString(room.Name)
	builder.WriteString("\n")
	builder.WriteString(room.Description)

	return builder.String()
}

// GetExits returns a string listing the available exits from the room
func (room Room) GetExits() string {
	builder := strings.Builder{}

	exitsStr, count := room.Exits.GetExits()

	if count == 0 {
		builder.WriteString("No exits available.\n")
	}

	builder.WriteString(fmt.Sprintf("You see %d exits\n", count))
	builder.WriteString(exitsStr)

	return builder.String()
}
