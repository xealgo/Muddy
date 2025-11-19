package session

import (
	"sync"
	"unsafe"

	"github.com/quic-go/webtransport-go"
	"github.com/xealgo/muddy/internal/player"
)

const (
	DefaultMaxSessions = 64 // Default maximum number of player sessions
)

// SessionManager manages player sessions in the game.
type SessionManager struct {
	maxSessions int
	Active      []PlayerSession
	Pending     map[string]*player.Player

	mutex      *sync.RWMutex
	sessionMap map[uintptr]string // Session pointer -> player UUID
}

// NewSessionManager creates a new SessionManager with a specified maximum number of sessions.
func NewSessionManager(maxSessions int) *SessionManager {
	return &SessionManager{
		maxSessions: maxSessions,
		Active:      make([]PlayerSession, maxSessions),
		mutex:       &sync.RWMutex{},
		sessionMap:  make(map[uintptr]string),
		Pending:     make(map[string]*player.Player),
	}
}

// Register adds a new player to the pending list.
func (sm *SessionManager) Register(player *player.Player) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	sm.Pending[player.GetUUID()] = player
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
