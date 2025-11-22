package session

import (
	"fmt"
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
	Active  []PlayerSession
	Pending map[string]*player.Player

	maxSessions int
	mutex       *sync.RWMutex
	sessionMap  map[uintptr]string // Session pointer -> player UUID
}

// NewSessionManager creates a new SessionManager with a specified maximum number of sessions.
func NewSessionManager(maxSessions int) *SessionManager {
	return &SessionManager{
		Active:      make([]PlayerSession, maxSessions),
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
func (sm *SessionManager) Connect(uuid string, session *webtransport.Session, stream *webtransport.Stream) (*PlayerSession, bool) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	playerData, exists := sm.Pending[uuid]
	if !exists {
		return nil, false
	}

	for i := 0; i < sm.maxSessions; i++ {
		var data = sm.Active[i].GetData()
		if data == nil {
			ps := NewPlayerSession(playerData, session, stream)
			sm.Active[i] = *ps
			sessionPtr := uintptr(unsafe.Pointer(ps.GetSession()))
			sm.sessionMap[sessionPtr] = data.GetUUID()

			delete(sm.Pending, uuid)

			return ps, true
		}
	}

	return nil, false
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
		if sm.Active[i].GetData() != nil && sm.Active[i].GetData().GetUUID() == uuid {
			sm.Active[i] = PlayerSession{}
			delete(sm.sessionMap, sessionPtr)
			return true
		}
	}

	return false
}

// RemovePlayer removes a PlayerSession from the manager.
func (sm *SessionManager) RemovePlayer(ps PlayerSession) bool {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	for i := 0; i < sm.maxSessions; i++ {
		if sm.Active[i].GetData() != nil && sm.Active[i].GetData().GetUUID() == ps.GetData().GetUUID() {
			sm.Active[i] = PlayerSession{}
			return true
		}
	}

	return false
}

// GetSession retrieves a PlayerSession by player UUID.
func (sm *SessionManager) GetSession(uuid string) (*PlayerSession, bool) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	for i := 0; i < sm.maxSessions; i++ {
		if sm.Active[i].GetData() != nil && sm.Active[i].GetData().GetUUID() == uuid {
			return &sm.Active[i], true
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
