package pogo

// TeamColor for a gym
type TeamColor int

// mapping from protos
const (
	Neutral TeamColor = iota
	Blue
	Red
	Yellow
)

// ToTeamColor converts an int from a proto to a TeamColor
func ToTeamColor(val int) TeamColor {
	switch val {
	case 1:
		return Blue
	case 2:
		return Red
	case 3:
		return Yellow
	default:
		return Neutral
	}
}

// ToString converts the TeamColor to a string with the name of the color
func (t TeamColor) ToString() (s string) {
	switch t {
	case Neutral:
		s = "white"
	case Blue:
		s = "blue"
	case Red:
		s = "red"
	case Yellow:
		s = "yellow"
	}
	return
}

// Gym describes a gym, optionally with a raid
type Gym struct {
	TeamColor TeamColor
	GUID      string // Ingress GUID, also used in Pogo
	Name      string // gym name as shown in game
	Location  Location
	Raid      *Raid
}
