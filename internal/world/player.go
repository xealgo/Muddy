package world

// RoomPlayer represents a player in a room
type RoomPlayer interface {
	GetID() int
	GetDisplayName() string
}
