package world

import "fmt"

// Item represents an item in the game world
type Item struct {
	Name         string `json:"name" yaml:"name"`
	Description  string `json:"description" yaml:"description"`
	SellingPrice int    `json:"sellingPrice" yaml:"sellingPrice"`
}

func (item Item) String() string {
	return fmt.Sprintf("%s: %s\n- Selling price: $$%d\n", item.Name, item.Description, item.SellingPrice)
}
