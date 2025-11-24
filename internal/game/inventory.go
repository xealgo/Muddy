package game

import (
	"fmt"
	"strings"
)

// Inventory represents a player's inventory.
type Inventory struct {
	Gold     int             `json:"gold"`
	Items    []*Item         `json:"items"`
	ItemsMap map[string]Item `json:"-"`
}

// NewInventory creates a new empty inventory.
func NewInventory() *Inventory {
	return &Inventory{
		ItemsMap: make(map[string]Item),
	}
}

// Initialize the inventory items
func (inv *Inventory) Initialize() {
	inv.Items = []*Item{}
	for _, item := range inv.ItemsMap {
		inv.Items = append(inv.Items, &item)
	}
}

// Sell removes an item from the inventory and sells it to a merchant.
func (inv *Inventory) Sell(itemId string, merchant *Merchant) (*Item, bool) {
	id := strings.ToLower(itemId)

	fmt.Printf("Attempting to sell item id: %s\n", id)

	item, exists := inv.ItemsMap[id]
	if !exists {
		return nil, false
	}

	if ok := merchant.Sell(item); !ok {
		return &item, false
	}

	inv.Gold += item.SellingPrice
	delete(inv.ItemsMap, id)

	return &item, true
}

// Add adds an item to the inventory.
func (inv *Inventory) Add(item Item) {
	item.ID = fmt.Sprintf("%d", len(inv.ItemsMap)+101)
	inv.ItemsMap[item.ID] = item
}

// List returns a string representation of the inventory contents.
func (inv Inventory) List() string {
	builder := strings.Builder{}

	builder.WriteString("You have ")
	builder.WriteString(fmt.Sprintf("%d gold\n", inv.Gold))

	if len(inv.ItemsMap) == 0 {
		builder.WriteString("Your inventory is empty.\n")
		return builder.String()
	}

	builder.WriteString("Your inventory contains:\n")
	for _, invItem := range inv.ItemsMap {
		builder.WriteString(fmt.Sprintf("- %s (ID: %s)\n", invItem.String(), invItem.ID))
	}

	return builder.String()
}
