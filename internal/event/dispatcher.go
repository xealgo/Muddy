package event

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/xealgo/muddy/internal/game"
)

// Simple event data type
type Event struct {
	Type      string      `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// EventDispatcher is responsible for dispatching events to their respective handlers.
type EventDispatcher struct {
	//
}

// SendToRoom sends an event to all players in a specific room.
func (e EventDispatcher) SendToRoom(event Event, sm *game.SessionManager, roomId int) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("unable to send event %s to room %d: %w", event.Type, roomId, err)
	}

	message := string(data)
	active := sm.GetActivePlayers()
	prefixed := "event:" + message

	// This will be slow if we ramp the max player count to 1000+
	// At that point, we'll want to create a slice within the room struct
	// or a shared map roomId -> []playerId.
	for _, ps := range active {
		if ps.CurrentRoomId == roomId {
			err := ps.WriteString(prefixed)
			if err != nil {
				slog.Error("failed to broadcast to player %s: %w", ps.DisplayName, err)
			}
		}
	}

	return nil
}
