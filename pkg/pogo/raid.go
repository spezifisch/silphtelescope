package pogo

// Raid contains raid info
type Raid struct {
	Hash     string
	GymID    string
	Location Location
	Pokemon  *Pokemon // nil if Spawned=false
	Level    int
	TimestampRange
}
