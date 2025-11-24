package game

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// World container for all rooms in the game.
type World struct {
	rooms   []*Room
	roomMap map[int]*Room
}

// NewWorld creates a new World instance.
func NewWorld() *World {
	w := &World{
		rooms:   []*Room{},
		roomMap: make(map[int]*Room),
	}
	return w
}

// GetRoomById retrieves a room by its ID.
func (w World) GetRoomById(roomId int) (*Room, bool) {
	room, exists := w.roomMap[roomId]
	return room, exists
}

// LoadRoomsFromYaml loads rooms from a YAML file.
func (w *World) LoadRoomsFromYaml(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to load file %s: %w", file, err)
	}

	rooms := []*Room{}

	err = yaml.Unmarshal(data, &rooms)
	if err != nil {
		return fmt.Errorf("failed to parse rooms from file %s: %w", file, err)
	}

	for _, room := range rooms {
		if !room.Validate() {
			return fmt.Errorf("invalid room data found in file %s", file)
		}

		newRoom := NewRoom(room.ID, room.Name, room.Description)
		newRoom.Copy(room)

		w.rooms = append(w.rooms, newRoom)
		w.roomMap[room.ID] = newRoom
	}

	return nil
}
