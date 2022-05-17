package http

import (
	"fmt"

	"github.com/spezifisch/silphtelescope/pkg/pogo"
)

func sendGymUpdate(msg *GymMessage) {
	m := pogo.Gym{
		GUID:      msg.GymID,
		TeamColor: pogo.ToTeamColor(msg.TeamID),
		Location: pogo.Location{
			Latitude:  float64(msg.Latitude),
			Longitude: float64(msg.Longitude),
		},
	}
	GymUpdates <- m
}

func sendSpawnUpdate(msg *PokemonMessage) {
	m := pogo.Spawn{
		EncounterID:        msg.EncounterID.String(),
		VerifiedSpawnpoint: msg.KnownDisappearTime,
		Pokemon: pogo.Pokemon{
			ID:     msg.PokemonID,
			Name:   "", // too expensive to look up right now, it's probaby ignored anyway
			Gender: pogo.ToGender(msg.Gender),
		},
		TimestampRange: pogo.TimestampRange{
			StartTime: 0, // MAD doesn't supply that info
			EndTime:   msg.DisappearTimestamp,
		},
		Location: pogo.Location{
			Latitude:  float64(msg.Latitude),
			Longitude: float64(msg.Longitude),
		},
	}
	SpawnUpdates <- m
}

func sendRaidUpdate(msg *RaidMessage) {
	var mon *pogo.Pokemon = nil
	if msg.PokemonID != 0 {
		mon = &pogo.Pokemon{
			ID:     msg.PokemonID,
			Name:   "", // too expensive to look up right now, it's probaby ignored anyway
			Gender: pogo.ToGender(msg.Gender),
		}
	}

	m := pogo.Raid{
		Hash:  fmt.Sprintf("%s:%d", msg.GymID, msg.StartTimestamp),
		GymID: msg.GymID,
		Location: pogo.Location{
			Latitude:  float64(msg.Latitude),
			Longitude: float64(msg.Longitude),
		},
		Pokemon: mon,
		Level:   msg.Level,
		TimestampRange: pogo.TimestampRange{
			StartTime: int64(msg.StartTimestamp),
			EndTime:   int64(msg.EndTimestamp),
		},
	}
	RaidUpdates <- m
}
