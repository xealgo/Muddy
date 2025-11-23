package game

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"
)

// Room represents a room in the game world
type Room struct {
	ID          int    `yaml:"id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Doors       []Door `yaml:"doors"`
	Items       []Item `yaml:"items"`

	doorMap map[string]*Door
	itemMap map[string]*Item
	mutex   *sync.RWMutex
}

// NewRoom creates a new Room instance
func NewRoom(id int, name string, desc string) *Room {
	room := &Room{
		ID:          id,
		Name:        name,
		Description: desc,
		Items:       []Item{},
		Doors:       []Door{},
		doorMap:     make(map[string]*Door),
		itemMap:     make(map[string]*Item),
		mutex:       &sync.RWMutex{},
	}

	room.Init()

	return room
}

// Init initializes the room's internal structures
func (room *Room) Init() {
	if room.doorMap == nil {
		room.doorMap = make(map[string]*Door)
	}

	if room.itemMap == nil {
		room.itemMap = make(map[string]*Item)
	}

	if room.mutex == nil {
		room.mutex = &sync.RWMutex{}
	}

	for _, door := range room.Doors {
		room.doorMap[door.MoveCommand] = &door
	}

	for _, item := range room.Items {
		room.itemMap[item.Name] = &item
	}
}

// Validate checks if the room has valid attributes
func (room Room) Validate() bool {
	if len(room.Items) > 0 {
		for _, item := range room.Items {
			if !item.Validate() {
				slog.Warn("Invalid item found in room", "roomId", room.ID, "itemName", item.Name)
				return false
			}
		}
	}

	if len(room.Doors) > 0 {
		for _, door := range room.Doors {
			if !door.Validate() {
				slog.Warn("Invalid door found in room", "roomId", room.ID, "doorName", door.Name)
				return false
			}
		}
	}

	if room.ID < 0 || room.Name == "" || room.Description == "" {
		slog.Warn("Invalid room attributes", "roomId", room.ID)
		return false
	}

	return true
}

// GetBasicInfo returns the basic information of the room
func (room Room) GetBasicInfo() string {
	builder := strings.Builder{}

	builder.WriteString(room.Name)
	builder.WriteString(", ")
	builder.WriteString(room.Description)

	return builder.String()
}

// GetDetails returns detailed information about the room, including exits, items, etc.
func (room Room) GetDetails(ps *Player, sm *SessionManager) string {
	builder := strings.Builder{}

	doorStr, count := room.GetDoors()
	if count == 0 {
		builder.WriteString("No exits")
	}

	if count == 1 {
		builder.WriteString(fmt.Sprintf("%d exit:\n", count))
	} else {
		builder.WriteString(fmt.Sprintf("%d exits:\n", count))
	}

	builder.WriteString(doorStr)

	items := room.Items
	if len(items) > 0 {
		builder.WriteString("You see the following items in the room:\n")
		for _, item := range items {
			builder.WriteString("- ")
			builder.WriteString(item.String())
		}
	}

	players := sm.GetPlayersInRoom(room.ID, ps.GetUUID())
	if len(players) > 0 {
		playerCount := 0

		psb := strings.Builder{}

		for _, player := range players {
			if player.CurrentRoomId == room.ID {
				playerCount++
				psb.WriteString(fmt.Sprintf("- %s\n", player.DisplayName))
			}
		}

		builder.WriteString("You see ")
		builder.WriteString(fmt.Sprintf("%d players:\n", playerCount))
		builder.WriteString(psb.String())
	}

	return builder.String()
}

// IsValidDoorChoice checks if the given door name is valid in the room
func (room Room) IsValidDoorChoice(choice string) bool {
	_, exists := room.doorMap[choice]
	return exists
}

// GetDoors returns a formatted string of the room's doors and the count of doors
func (room Room) GetDoors() (string, int) {
	count := len(room.Doors)
	builder := strings.Builder{}

	for _, door := range room.Doors {
		builder.WriteString("> ")
		builder.WriteString(door.Name)

		if door.Description != "" {
			builder.WriteString(": ")
			builder.WriteString(door.Description)
		}

		builder.WriteByte('\n')
	}

	return builder.String(), count
}

// RemoveItem removes an item from the room by its name
func (room *Room) RemoveItem(itemName string) (Item, bool) {
	room.mutex.Lock()
	defer room.mutex.Unlock()

	item, ok := room.itemMap[itemName]
	if !ok {
		return Item{}, false
	}

	copy := *item

	delete(room.itemMap, itemName)

	newItems := []Item{}
	for _, item := range room.Items {
		if item.Name != itemName {
			newItems = append(newItems, item)
		}
	}

	room.Items = newItems
	return copy, true
}

// AddItem adds an item to the room
func (room *Room) AddItem(item Item) bool {
	room.mutex.Lock()
	defer room.mutex.Unlock()

	_, ok := room.itemMap[item.Name]
	if ok {
		return false
	}

	room.itemMap[item.Name] = &item
	room.Items = append(room.Items, item)

	return true
}
