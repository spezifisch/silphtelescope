package roomservice

import (
	"errors"

	"github.com/jinzhu/copier"
	log "github.com/sirupsen/logrus"

	"github.com/spezifisch/silphtelescope/internal/helpers"
)

// GetRoomConfig returns a copy of the room configuration if there is one
func (p *Poster) GetRoomConfig(roomID string) (*RoomConfig, bool) {
	rc, ok := p.roomConfigs[roomID]
	return rc, ok
}

// UpdateRoomConfig overwrites or adds the RoomConfig and appends the Filter
func (p *Poster) UpdateRoomConfig(rcUpdate *RoomConfig) (created bool) {
	defer p.commitRoomConfig(rcUpdate.RoomID)

	if rc, ok := p.roomConfigs[rcUpdate.RoomID]; ok {
		rc.AcceptCommands = rcUpdate.AcceptCommands
		rc.FormatText = rcUpdate.FormatText
		rc.Filter = append(rc.Filter, rcUpdate.Filter...)

		created = false
		return
	}

	// room config doesn't exist yet
	p.roomConfigs[rcUpdate.RoomID] = &RoomConfig{}
	copier.Copy(p.roomConfigs[rcUpdate.RoomID], rcUpdate)
	created = true
	return
}

// DeleteFilters deletes all pokemon filters for a room if a RoomConfig exists
func (p *Poster) DeleteFilters(roomID string) (deleted bool) {
	if rc, ok := p.roomConfigs[roomID]; ok {
		rc.Filter = nil
		p.commitRoomConfig(roomID)
		deleted = true
	}
	return
}

// RoomConfigChange describes an update operation for an existing RoomConfig
type RoomConfigChange struct {
	ChangeAcceptCommands bool // update RC with value from given RoomConfig
	ChangeFormatText     bool // same as above
	Operation            RoomConfigOperation
	FilterIndex          int          // only when UpdateFilter=true
	FilterChange         FilterChange // only when UpdateFilter=true
}

// RoomConfigOperation describes what should be done in a RoomConfigChange
type RoomConfigOperation int

// this flag says what you want to do with the roomconfig
const (
	RoomConfigOperationNone RoomConfigOperation = iota
	RoomConfigOperationAppendFilter
	RoomConfigOperationUpdateFilter
	RoomConfigOperationRemoveFilter
)

// FilterChange is used to describe what should be done with the existing Filter
type FilterChange int

const (
	// FilterChangeAddPokemon adds pokemon to list
	FilterChangeAddPokemon FilterChange = iota
	// FilterChangeRemovePokemon removes pokemon from list
	FilterChangeRemovePokemon
	// FilterChangeArea replaces the area
	FilterChangeArea
)

// ChangeRoomConfig edits an existing RoomConfig with the given changeset
func (p *Poster) ChangeRoomConfig(roomID string, rcChange *RoomConfigChange, newValues *RoomConfig) (err error) {
	// get RC
	rc, ok := p.roomConfigs[roomID]
	if !ok {
		return errors.New("roomconfig doesn't exist")
	}
	if rcChange == nil {
		return errors.New("no rcChange supplied")
	}

	// sanity check
	if rcChange.Operation == RoomConfigOperationAppendFilter {
		if len(newValues.Filter) == 0 {
			return errors.New("there are no filters to append")
		}
	}
	if rcChange.Operation == RoomConfigOperationUpdateFilter {
		if rcChange.FilterIndex < 0 || rcChange.FilterIndex >= len(rc.Filter) {
			return errors.New("invalid filter id")
		}
		if len(newValues.Filter) != 1 {
			return errors.New("supplied Filter count must be 1 for UpdateFilter operation")
		}
	}
	if rcChange.Operation == RoomConfigOperationRemoveFilter {
		if rcChange.FilterIndex < 0 || rcChange.FilterIndex >= len(rc.Filter) {
			return errors.New("invalid filter id")
		}
	}

	// everything needs newValues except RemoveFilter
	if rcChange.Operation != RoomConfigOperationRemoveFilter {
		if newValues == nil {
			return errors.New("no newValues supplied")
		}
	}

	// do the changes
	defer p.commitRoomConfig(roomID)
	if rcChange.ChangeAcceptCommands {
		rc.AcceptCommands = newValues.AcceptCommands
	}
	if rcChange.ChangeFormatText {
		rc.FormatText = newValues.FormatText
	}
	switch rcChange.Operation {
	case RoomConfigOperationAppendFilter:
		rc.Filter = append(rc.Filter, newValues.Filter...)
	case RoomConfigOperationUpdateFilter:
		f := &rc.Filter[rcChange.FilterIndex]
		p.changeRoomConfigFilter(f, rcChange.FilterChange, &newValues.Filter[0])
	case RoomConfigOperationRemoveFilter:
		rc.Filter = append(rc.Filter[:rcChange.FilterIndex], rc.Filter[rcChange.FilterIndex+1:]...)
	}

	return
}

// Change a single Filter. All input is assumed to be benign and checked by ChangeRoomConfig().
func (p *Poster) changeRoomConfigFilter(f *PokemonFilter, op FilterChange, newFilter *PokemonFilter) {
	switch op {
	case FilterChangeAddPokemon:
		// TODO check for duplicate values and array length
		f.PokemonIDs = append(f.PokemonIDs, newFilter.PokemonIDs...)
	case FilterChangeRemovePokemon:
		newPokemonList := []int{}
		// copy all pokemon from old list who should be kept
		for _, pokemonID := range f.PokemonIDs {
			if !helpers.IntArrayContains(newFilter.PokemonIDs, pokemonID) {
				// append if not in exclusion list
				if !helpers.IntArrayContains(newPokemonList, pokemonID) {
					// and if not already in list
					newPokemonList = append(newPokemonList, pokemonID)
				}
			}
		}
		f.PokemonIDs = newPokemonList
	case FilterChangeArea:
		f.Area = newFilter.Area
	}
}

// Write RoomConfig changes to disk.
// This needs to be called after every modification of a RoomConfig object!
func (p *Poster) commitRoomConfig(roomID string) {
	if p.db == nil {
		return
	}

	p.db.SaveRoomConfig(roomID, p.roomConfigs[roomID])
}

func (p *Poster) readRoomConfigs() {
	if p.db == nil {
		return
	}

	p.db.ReadRoomConfigs(p.roomConfigs)
	log.Infof("read RoomConfig for %d rooms", len(p.roomConfigs))
}
