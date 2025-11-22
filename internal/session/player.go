package session

import (
	"fmt"

	"github.com/quic-go/webtransport-go"
	"github.com/xealgo/muddy/internal/player"
)

// PlayerSession represents a player's session in the game.
type PlayerSession struct {
	data    *player.Player
	session *webtransport.Session
	stream  *webtransport.Stream
}

// NewPlayerSession creates a new PlayerSession.
func NewPlayerSession(data *player.Player, session *webtransport.Session, stream *webtransport.Stream) *PlayerSession {
	return &PlayerSession{
		data:    data,
		session: session,
		stream:  stream,
	}
}

// GetData returns the player data associated with the session.
func (ps PlayerSession) GetData() *player.Player {
	return ps.data
}

// GetSession returns the webtransport session.
func (ps PlayerSession) GetSession() *webtransport.Session {
	return ps.session
}

// GetStream returns the webtransport stream.
func (ps PlayerSession) GetStream() *webtransport.Stream {
	return ps.stream
}

// WriteString writes a string message to the player's stream.
func (ps PlayerSession) WriteString(message string) error {
	if ps.stream == nil {
		return fmt.Errorf("player session (%s) stream is nil", ps.data.GetUUID())
	}

	_, err := ps.stream.Write([]byte(message))
	return err
}
