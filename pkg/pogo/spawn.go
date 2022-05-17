package pogo

// Spawn describes a spawned Pokemon
type Spawn struct {
	EncounterID        string
	VerifiedSpawnpoint bool
	Pokemon
	Location
	TimestampRange
}
