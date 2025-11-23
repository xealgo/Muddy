package player

import "github.com/google/uuid"

// Player represents a player in the game.
type Player struct {
	uuid          string
	Username      string
	DisplayName   string
	CurrentRoomId int
}

// NewPlayer creates a new player with a unique UUID.
func NewPlayer(username string, displayName string) *Player {
	p := &Player{
		uuid:          uuid.NewString(),
		Username:      username,
		DisplayName:   displayName,
		CurrentRoomId: 1,
	}

	return p
}

// GetUUID returns the UUID of the player.
func (p Player) GetUUID() string {
	return p.uuid
}
