package world

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoomBasicInfo(t *testing.T) {
	room := NewRoom(1, "Test Room", "A room for testing.")
	assert.Equal(t, 1, room.ID)
	assert.Equal(t, "Test Room", room.Name)
	assert.Equal(t, "A room for testing.", room.Description)

	info := room.GetBasicInfo()
	expectedInfo := "Test Room\nA room for testing."

	assert.Equal(t, expectedInfo, info)
}
