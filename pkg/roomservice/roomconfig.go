package roomservice

import (
	"fmt"

	"github.com/spezifisch/silphtelescope/pkg/pogo"
)

// PokemonFilter specifies which pokemon in an area should be posted
type PokemonFilter struct {
	Area       pogo.LocationRadius // area to include
	ListRaids  bool                // true if this filter only matches raids, false for only spawns
	ListWanted bool                // true if only wanted pokemon ids are in the list, false for unwanted pokemon
	PokemonIDs []int               // pokedex numbers
}

// RoomConfig contains settings for a room with one or more people
type RoomConfig struct {
	RoomID         string
	AcceptCommands bool // parse commands from users in this room (admin privileges are checked seperately)
	FormatText     bool
	Filter         []PokemonFilter
}

// ToString converts RoomConfig into a human-readable string
func (r RoomConfig) ToString() string {
	s := fmt.Sprintf("RoomID:%s", r.RoomID)
	return s
}
