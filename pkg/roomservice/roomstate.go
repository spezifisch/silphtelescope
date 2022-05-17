package roomservice

import "github.com/spezifisch/silphtelescope/pkg/pogo"

// RoomState keeps track of posted mons
type RoomState struct {
	// EncounterID -> SpawnState
	Spawns map[string]*SpawnState
	// Raid hash -> RaidState
	Raids map[string]*RaidState
}

// SpawnState tracks if a spawn was already posted as long as it has not elapsed
type SpawnState struct {
	EndTime int64
	Posted  bool
}

// RaidState tracks if a raid was already posted as long as it has not elapsed
type RaidState struct {
	EndTime int64
	Posted  bool
}

// NewRoomState creates a RoomState object
func NewRoomState() *RoomState {
	return &RoomState{
		Spawns: make(map[string]*SpawnState),
		Raids:  make(map[string]*RaidState),
	}
}

func (r *RoomState) clear(spawns, raids bool) {
	if spawns {
		r.Spawns = make(map[string]*SpawnState)
	}
	if raids {
		r.Raids = make(map[string]*RaidState)
	}
}

func (r *RoomState) raidIsPosted(hash string) bool {
	s, found := r.Raids[hash]
	return found && s.Posted
}

func (r *RoomState) spawnIsPosted(encounterID string) bool {
	s, found := r.Spawns[encounterID]
	return found && s.Posted
}

func (r *RoomState) postedRaid(s *pogo.Raid, posted bool) {
	r.Raids[s.Hash] = &RaidState{
		EndTime: s.EndTime,
		Posted:  posted,
	}
}

func (r *RoomState) postedSpawn(s *pogo.Spawn, posted bool) {
	r.Spawns[s.EncounterID] = &SpawnState{
		EndTime: s.EndTime,
		Posted:  posted,
	}
}

func (r *RoomState) removeExpired(before int64) (deleted int) {
	deleted = 0
	for k, st := range r.Raids {
		if st.EndTime < before {
			delete(r.Raids, k)
			deleted++
		}
	}
	for k, st := range r.Spawns {
		if st.EndTime < before {
			delete(r.Spawns, k)
			deleted++
		}
	}
	return
}
