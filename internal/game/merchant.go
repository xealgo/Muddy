package game

import (
	"encoding/json"
	"fmt"
	"sync"
)

// Merchant represents a merchant in the game world.
type Merchant struct {
	NpcData
	Inventory *Inventory `json:"inventory"`

	mutex *sync.Mutex
}

// NewMerchant creates a new Merchant instance.
func NewMerchant(id string) *Merchant {
	m := &Merchant{
		Inventory: NewInventory(),
		mutex:     &sync.Mutex{},
	}

	m.ID = id
	m.Type = NpcMerchant

	return m
}

// GetData returns the NPC data of the merchant.
func (m *Merchant) GetData() *NpcData {
	return &m.NpcData
}

// Greet sends a greeting message to the player.
func (m *Merchant) Greet(player *Player) string {
	return m.Greeting
}

// Description returns a description of the merchant.
func (m *Merchant) Description() string {
	return fmt.Sprintf("%s the merchant", m.Name)
}

// Sell allows the merchant to buy an item from a player.
func (m *Merchant) Sell(item Item) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.Inventory.Add(item)
	return true
}

// Convert converts raw NPC data into a Merchant instance.
func (m *Merchant) Convert(rawNpc map[string]any) error {
	if rawNpc["type"] != NpcMerchant {
		return fmt.Errorf("invalid NPC data format")
	}

	// Marshal the map to a JSON byte slice
	jsonBytes, err := json.Marshal(rawNpc)
	if err != nil {
		return fmt.Errorf("error marshaling NPC data: %w", err)
	}

	// Unmarshal the JSON byte slice into your struct
	npcData := NewMerchant("")

	if err = json.Unmarshal(jsonBytes, npcData); err != nil {
		return fmt.Errorf("error unmarshaling NPC data: %w", err)
	}

	if npcData.Type != NpcMerchant {
		return fmt.Errorf("invalid NPC type for merchant: %s", npcData.Type)
	}

	if npcData.Inventory != nil {
		m.Inventory = npcData.Inventory
	}

	m.Name = npcData.Name
	m.Greeting = npcData.Greeting
	m.Inventory.Initialize()

	return nil
}
