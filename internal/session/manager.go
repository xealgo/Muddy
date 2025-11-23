package session

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"unsafe"

	"github.com/quic-go/webtransport-go"
	"github.com/xealgo/muddy/internal/player"
)

type SessionManagerErrorType string

const (
	DefaultMaxSessions = 64 // Default maximum number of player sessions

	// Errors
	ErrorMaxPlayers SessionManagerErrorType = "MAX_PLAYERS_REACHED"
)

// Custom session manager error
type SessionManagerError struct {
	Type    SessionManagerErrorType
	Message string
	Wrapped error
}

// Unwrap returns the underlying error
func (e SessionManagerError) Error() string {
	return fmt.Sprintf("Type: %v, Message: %s, Wrapped: %w", e.Type, e.Message, e.Error())
}

// Unwrap returns the underlying error
func (e SessionManagerError) Unwrap() error {
	return e.Wrapped
}

// SessionManager manages player sessions in the game.
type SessionManager struct {
	Active  []*PlayerSession
	Pending map[string]*player.Player

	maxSessions int
	mutex       *sync.RWMutex
	sessionMap  map[uintptr]string // Session pointer -> player UUID
}

// NewSessionManager creates a new SessionManager with a specified maximum number of sessions.
func NewSessionManager(maxSessions int) *SessionManager {
	return &SessionManager{
		Active:      make([]*PlayerSession, maxSessions),
		Pending:     make(map[string]*player.Player),
		maxSessions: maxSessions,
		mutex:       &sync.RWMutex{},
		sessionMap:  make(map[uintptr]string),
	}
}

// Register adds a new player to the pending list.
func (sm *SessionManager) Register(player *player.Player) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	count := len(sm.sessionMap)
	if count >= sm.maxSessions {
		return &SessionManagerError{Type: ErrorMaxPlayers, Message: "Max player limit reached, please try again", Wrapped: nil}
	}

	sm.Pending[player.GetUUID()] = player
	return nil
}

// Connect adds a new PlayerSession to the manager.
func (sm *SessionManager) Connect(uuid string, session *webtransport.Session, stream *webtransport.Stream) (*PlayerSession, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	playerData, exists := sm.Pending[uuid]
	if !exists {
		return nil, fmt.Errorf("no pending session found")
	}

	for i := 0; i < sm.maxSessions; i++ {
		if sm.Active[i] == nil {
			ps := NewPlayerSession(playerData, session, stream)
			sm.Active[i] = ps
			sessionPtr := uintptr(unsafe.Pointer(ps.GetSession()))

			sm.sessionMap[sessionPtr] = ps.GetData().GetUUID()

			delete(sm.Pending, uuid)

			return ps, nil
		}
	}

	return nil, fmt.Errorf("unable to create player session")
}

// RemovePlayer removes a PlayerSession from the manager.
func (sm *SessionManager) RemovePlayerBySession(session *webtransport.Session) bool {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	sessionPtr := uintptr(unsafe.Pointer(session))
	uuid, exists := sm.sessionMap[sessionPtr]
	if !exists {
		return false
	}

	for i := 0; i < sm.maxSessions; i++ {
		if sm.Active[i] != nil && sm.Active[i].GetData() != nil && sm.Active[i].GetData().GetUUID() == uuid {
			sm.Active[i] = nil
			delete(sm.sessionMap, sessionPtr)
			return true
		}
	}

	return false
}

// RemovePlayer removes a PlayerSession from the manager.
func (sm *SessionManager) RemovePlayer(uuid string) bool {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	for i := 0; i < sm.maxSessions; i++ {
		if sm.Active[i] != nil && sm.Active[i].GetData() != nil && sm.Active[i].GetData().GetUUID() == uuid {
			sm.Active[i] = nil
			return true
		}
	}

	return false
}

// RemovePending removes a player from the pending list.
func (sm *SessionManager) RemovePending(uuid string) bool {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if _, exists := sm.Pending[uuid]; exists {
		delete(sm.Pending, uuid)
		return true
	}

	return false
}

// GetSession retrieves a PlayerSession by player UUID.
func (sm *SessionManager) GetSession(uuid string) (*PlayerSession, bool) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	for i := 0; i < sm.maxSessions; i++ {
		if sm.Active[i] != nil && sm.Active[i].GetData() != nil && sm.Active[i].GetData().GetUUID() == uuid {
			return sm.Active[i], true
		}
	}

	return nil, false
}

// GetActiveSessionCount returns the number of active player sessions.
func (sm *SessionManager) GetActiveSessionCount() int {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	return len(sm.sessionMap)
}

// GetActivePlayers returns a slice of all active PlayerSessions.
func (sm *SessionManager) GetActivePlayers() []*PlayerSession {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	active := []*PlayerSession{}
	for _, ps := range sm.Active {
		if ps != nil {
			active = append(active, ps)
		}
	}

	return active
}

// GetPlayersInRoom returns a slice of players currently in the specified room.
func (sm *SessionManager) GetPlayersInRoom(roomId int, skipPlayerUUID string) []player.Player {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	players := []player.Player{}

	for _, ps := range sm.Active {
		if ps != nil && ps.GetData().CurrentRoomId == roomId && ps.data.GetUUID() != skipPlayerUUID {
			players = append(players, *ps.GetData())
		}
	}

	return players
}

// SendToPlayer sends a message to a specific player by UUID.
func (sm SessionManager) SendToPlayer(playerUuid string, message string) {
	active := sm.GetActivePlayers()
	trimmed := strings.TrimRight(message, "\n") + "\n"

	var session *PlayerSession

	// Eventually we'll want to track players in a map for fast lookup..
	for _, ps := range active {
		if ps.GetData().GetUUID() == playerUuid {
			session = ps
			break
		}
	}

	if session != nil {
		err := session.WriteString(trimmed)
		if err != nil {
			slog.Error("failed to send to player %s: %w", session.data.DisplayName, err)
		}
	}
}
