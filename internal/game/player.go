package game

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/quic-go/webtransport-go"
)

// Player represents a player in the game.
type Player struct {
	uuid          string
	Username      string
	DisplayName   string
	CurrentRoomId int
	Inventory     *Inventory

	session *webtransport.Session
	stream  *webtransport.Stream
}

// NewPlayer creates a new player with a unique UUID.
func NewPlayer(username string, displayName string) *Player {
	p := &Player{
		uuid:          uuid.NewString(),
		Username:      username,
		DisplayName:   displayName,
		CurrentRoomId: 1,
		Inventory:     NewInventory(),
	}

	return p
}

// GetUUID returns the UUID of the player.
func (p Player) GetUUID() string {
	return p.uuid
}

// SetSession sets the webtransport session for the player.
func (p *Player) SetSession(session *webtransport.Session) {
	p.session = session
}

// SetStream sets the webtransport stream for the player.
func (p *Player) SetStream(stream *webtransport.Stream) {
	p.stream = stream
}

// GetSession returns the webtransport session.
func (p Player) GetSession() *webtransport.Session {
	return p.session
}

// GetStream returns the webtransport stream.
func (p Player) GetStream() *webtransport.Stream {
	return p.stream
}

// WriteString writes a string message to the player's stream.
func (p Player) WriteString(message string) error {
	if p.stream == nil {
		return fmt.Errorf("player session (%s) stream is nil", p.uuid)
	}

	_, err := p.stream.Write([]byte(message))
	return err
}
