package game

const (
	NpcMerchant string = "merchant"
)

// Npc interface represents a non-player character in the game.
type Npc interface {
	GetData() *NpcData
	Greet(*Player) string
	Description() string
}

// NpcData holds basic information about an NPC.
type NpcData struct {
	ID          string `yaml:"-"`
	Type        string `yaml:"type"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Greeting    string `yaml:"greeting"`
}
