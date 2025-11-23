package game

import "fmt"

type ItemType string

const (
	Key     ItemType = "key"
	Trinket ItemType = "trinket"
)

var validItemTypes = []ItemType{Key, Trinket}

// Item represents an item in the game world
type Item struct {
	Type         ItemType `yaml:"type"`
	Name         string   `yaml:"name"`
	Description  string   `yaml:"description"`
	SellingPrice int      `yaml:"sellingPrice"`
}

// String returns a formatted string representation of the item
func (item Item) String() string {
	return fmt.Sprintf("%s, %s\n", item.Name, item.Description)
}

// Validate checks if the item has valid attributes
func (item Item) Validate() bool {
	isValidType := false

	for _, validType := range validItemTypes {
		if item.Type == validType {
			isValidType = true
			break
		}
	}

	if item.Name == "" || item.Description == "" || item.SellingPrice < 0 {
		return false
	}

	return isValidType
}
