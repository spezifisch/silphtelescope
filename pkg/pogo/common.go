package pogo

// TimestampRange contains how long a raid/spawn is active for
type TimestampRange struct {
	StartTime int64
	EndTime   int64
}

// Pokemon describes a Pokemon
type Pokemon struct {
	ID     int
	Name   string // English name
	Gender Gender
}
