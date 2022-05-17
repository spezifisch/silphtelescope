package roomservice

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/spezifisch/silphtelescope/internal/helpers"
	"github.com/spezifisch/silphtelescope/pkg/geodex"
	"github.com/spezifisch/silphtelescope/pkg/pogo"
)

// Poster is an agent that posts info from various datasources without user interaction
type Poster struct {
	// data input channels
	GymUpdates   chan pogo.Gym
	SpawnUpdates chan pogo.Spawn
	RaidUpdates  chan pogo.Raid

	// control channels
	Quit             chan bool
	MainControl      MainController
	saveStateAndQuit chan bool

	// pokedex
	Pokedex *pogo.Pokedex

	// geodex
	GeoDex *geodex.GeoDex

	// read RoomConfigs and States from disk before starting the main loop in Run()
	ResumeStateOnStartup bool

	// remove expired spawns and raids from memory
	ExpiryCheckPeriod time.Duration

	// output
	chatter Chatter

	// storage
	db Persister

	// persistent configuration: roomID -> RoomConfig
	roomConfigs map[string]*RoomConfig

	// in-memory state like posted encounter IDs
	roomStates map[string]*RoomState

	// timestamp of last message from MAD
	lastDataTime time.Time

	// timestamp of bot start
	startTime time.Time
}

// NewPoster creates some objects
func NewPoster(chatter Chatter, persister Persister) *Poster {
	return &Poster{
		ExpiryCheckPeriod:    30 * time.Second,
		ResumeStateOnStartup: false,
		roomConfigs:          make(map[string]*RoomConfig),
		roomStates:           make(map[string]*RoomState),
		saveStateAndQuit:     make(chan bool),
		chatter:              chatter,
		db:                   persister,
	}
}

// Run runs the poster main loop blockingly
func (p *Poster) Run() {
	p.startTime = time.Now()

	if p.ResumeStateOnStartup {
		p.readRoomConfigs()
		p.readRoomStates()
	}

	expiryTicker := time.NewTicker(p.ExpiryCheckPeriod)

	for {
		// wait for updates
		select {
		case <-p.GymUpdates:
			//log.Debug("got gym update#")
			p.updateLastData()
		case s := <-p.SpawnUpdates:
			p.updateLastData()
			p.processSpawnUpdate(s)
		case r := <-p.RaidUpdates:
			p.updateLastData()
			p.processRaidUpdate(r)
		case <-expiryTicker.C:
			p.cleanupTick()
		case <-p.saveStateAndQuit:
			expiryTicker.Stop()

			// Save everything from memory that we need to resume our operation after starting again.
			// RoomConfigs are already saved on modification.
			// RoomStates aren't because they're changed very frequently, so save them here.
			p.saveRoomStates()

			if p.MainControl != nil {
				log.Info("poster stopped, shutting down maincontrol")
				defer p.MainControl.Stop()
			}
			return
		case <-p.Quit:
			return
		}
	}
}

// periodical memory cleanup of
// * expired spawns
// * nothing else yet
func (p *Poster) cleanupTick() {
	deleted := 0
	now := time.Now().Unix()
	for _, roomState := range p.roomStates {
		deleted += roomState.removeExpired(now)
	}

	if deleted > 0 {
		log.Debugf("removed %d expired elements", deleted)
	}
}

func (p *Poster) updateLastData() {
	p.lastDataTime = time.Now()
}

func (p *Poster) processRaidUpdate(r pogo.Raid) {
	if r.Pokemon == nil || r.Pokemon.ID == 0 {
		return
	}

	// TODO reduce complexity
	for _, room := range p.roomConfigs {
		roomState := p.getOrCreateRoomState(room.RoomID)
		if roomState.raidIsPosted(r.Hash) {
			continue
		}

		for _, filter := range room.Filter {
			if !filter.ListRaids {
				continue
			}

			if helpers.IntArrayContains(filter.PokemonIDs, r.Pokemon.ID) {
				if filter.Area.Contains(&r.Location) {
					p.postRaid(room, &r)
					roomState.postedRaid(&r, true)
					break
				}
			}
		}
	}
}

func (p *Poster) processSpawnUpdate(s pogo.Spawn) {
	// TODO reduce complexity
	for _, room := range p.roomConfigs {
		roomState := p.getOrCreateRoomState(room.RoomID)
		if roomState.spawnIsPosted(s.EncounterID) {
			continue
		}

		for _, filter := range room.Filter {
			if !filter.ListWanted {
				continue
			}

			if helpers.IntArrayContains(filter.PokemonIDs, s.Pokemon.ID) {
				if filter.Area.Contains(&s.Location) {
					p.postSpawn(room, &s)
					roomState.postedSpawn(&s, true)
					break
				}
			}
		}
	}
}

func (p *Poster) getOrCreateRoomState(roomID string) *RoomState {
	rs, ok := p.roomStates[roomID]
	if !ok {
		// create
		rs = NewRoomState()
		p.roomStates[roomID] = rs
	}

	return rs
}

func (p *Poster) postRaid(room *RoomConfig, r *pogo.Raid) {
	startTime := time.Unix(r.StartTime, 0)
	endTime := time.Unix(r.EndTime, 0)
	startTimeStr := startTime.Format("15:04:05")
	endTimeStr := endTime.Format("15:04:05")

	pokemonStr := fmt.Sprintf("Pokemon #%d", r.Pokemon.ID)
	if p.Pokedex != nil {
		nameEN, nameDE, err := p.Pokedex.GetNamesByID(r.Pokemon.ID)
		if err == nil {
			if nameEN != nameDE {
				pokemonStr = fmt.Sprintf("%s (de: %s)", nameEN, nameDE)
			} else {
				pokemonStr = nameEN
			}
		}
	}

	raidLocation := r.Location
	fortName := r.GymID
	if p.GeoDex != nil {
		fort, err := p.GeoDex.Disk.GetFort(r.GymID)
		if err == nil {
			fortName = fort.GetName()
		}
	}

	text := fmt.Sprintf("Raid %s %s-%s at %s (Level %d)",
		pokemonStr, startTimeStr, endTimeStr, fortName, r.Level)
	if room.FormatText {
		fortStr := fmt.Sprintf("<a href=\"%s\">%s</a>", raidLocation.ToLinkGMaps(), fortName)
		fText := fmt.Sprintf("Raid %s %s-%s at %s (Level %d)",
			pokemonStr, startTimeStr, endTimeStr, fortStr, r.Level)
		p.chatter.SendFormattedText(room.RoomID, text, fText)
	} else {
		p.chatter.SendText(room.RoomID, text)
	}
}

func (p *Poster) postSpawn(room *RoomConfig, s *pogo.Spawn) {
	endTime := time.Unix(s.EndTime, 0)
	timeLeft := endTime.Sub(time.Now().Round(time.Second))

	endTimeStr := endTime.Format("15:04:05")

	pokemonStr := fmt.Sprintf("Pokemon #%d", s.Pokemon.ID)
	if p.Pokedex != nil {
		nameEN, nameDE, err := p.Pokedex.GetNamesByID(s.Pokemon.ID)
		if err == nil {
			if nameEN != nameDE {
				pokemonStr = fmt.Sprintf("%s (de: %s)", nameEN, nameDE)
			} else {
				pokemonStr = nameEN
			}
		}
	}

	gmapsLink := s.Location.ToLinkGMaps()

	// lookup nearest stop or gym
	nearStr := ""
	fmtNearStr := ""
	if p.GeoDex != nil {
		radiusM := 500.0
		if nearestFort, err := p.GeoDex.LookupFortNear(s.Location, radiusM); err == nil {
			// get distance and bearing from fort to spawn point
			fortLocation := nearestFort.Location()
			distanceM := fortLocation.DistanceTo(&s.Location)
			bearingDeg := fortLocation.BearingTo(&s.Location)
			nearStr = fmt.Sprintf(" near %s (%dm, %d°)", nearestFort.GetName(), int(distanceM), int(bearingDeg))
			fmtNearStr = fmt.Sprintf(" near <a href=\"%s\">%s (%dm, %d°)</a>", gmapsLink, nearestFort.GetName(), int(distanceM), int(bearingDeg))
		}
	}

	text := fmt.Sprintf("%s until %s (%s left)%s at %s",
		pokemonStr, endTimeStr, timeLeft, nearStr, gmapsLink)
	if room.FormatText {
		if fmtNearStr == "" {
			fmtNearStr = fmt.Sprintf(" at <a href=\"%s\">(%f,%f)</a>", gmapsLink, s.Longitude, s.Latitude)
		}

		fText := fmt.Sprintf("%s until %s (%s left)%s",
			pokemonStr, endTimeStr, timeLeft, fmtNearStr)
		p.chatter.SendFormattedText(room.RoomID, text, fText)
	} else {
		p.chatter.SendText(room.RoomID, text)
	}
}
