package game

import (
	"fmt"
	"strings"
)

type InventoryItem struct {
	ID   string
	Item Item
}

// Inventory represents a player's inventory.
type Inventory struct {
	Items []InventoryItem
}

// NewInventory creates a new empty inventory.
func NewInventory() *Inventory {
	return &Inventory{
		Items: []InventoryItem{},
	}
}

// Sell removes an item from the inventory and sells it to a merchant.
func (inv *Inventory) Sell(id string, merchant *Merchant) bool {
	fmt.Println("Selling items is not yet implemented")

	for i, item := range inv.Items {
		if item.ID == id {
			inv.Items = append(inv.Items[:i], inv.Items[i+1:]...)
			return true
		}
	}
	return false
}

// Add adds an item to the inventory.
func (inv *Inventory) Add(item Item) {
	inv.Items = append(inv.Items, InventoryItem{
		ID:   fmt.Sprintf("%s-%d", item.Name, len(inv.Items)+1),
		Item: item,
	})
}

// List returns a string representation of the inventory contents.
func (inv Inventory) List() string {
	builder := strings.Builder{}
	if len(inv.Items) == 0 {
		builder.WriteString("Your inventory is empty.\n")
		return builder.String()
	}

	builder.WriteString("Your inventory contains:\n")
	for _, invItem := range inv.Items {
		builder.WriteString(fmt.Sprintf("- %s (ID: %s)\n", invItem.Item.String(), invItem.ID))
	}

	return builder.String()
}
