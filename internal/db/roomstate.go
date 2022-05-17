package db

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/spezifisch/silphtelescope/pkg/roomservice"
)

// SaveRoomState is called by roomservice to persist them, cached
func (db *DB) SaveRoomState(roomID string, roomState *roomservice.RoomState) {
	data, err := json.Marshal(roomState)
	if err == nil {
		db.dvRoomState.Write(roomID, data)
	}
}

// SaveRoomStates is called by roomservice to persist them, cached
func (db *DB) SaveRoomStates(roomStates map[string]*roomservice.RoomState) {
	for roomID, rs := range roomStates {
		db.SaveRoomState(roomID, rs)
	}
}

// ReadRoomStates returns all saved room states from disk
func (db *DB) ReadRoomStates(roomStates map[string]*roomservice.RoomState) {
	for roomID := range db.dvRoomState.Keys(nil) {
		data, err := db.dvRoomState.Read(roomID)
		if err == nil {
			roomStates[roomID] = roomservice.NewRoomState()
			err = json.Unmarshal(data, roomStates[roomID])
			if err != nil {
				log.WithError(err).Warnf("failed deserializing RoomState of %s", roomID)
				delete(roomStates, roomID)
			}
		} else {
			log.WithError(err).Warnf("failed reading diskv key for RoomState of %s", roomID)
		}
	}
}
