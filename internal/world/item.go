package world

// Item represents an item in the game world
type Item struct {
	Name         string  `json:"name" yaml:"name"`
	Description  string  `json:"description" yaml:"description"`
	SellingPrice float64 `json:"sellingPrice" yaml:"sellingPrice"`
}
