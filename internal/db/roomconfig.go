package db

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/spezifisch/silphtelescope/pkg/roomservice"
)

// SaveRoomConfig saves a single room's RoomConfig
func (db *DB) SaveRoomConfig(roomID string, roomConfig *roomservice.RoomConfig) {
	data, err := json.Marshal(roomConfig)
	if err == nil {
		db.dvRoomConfig.Write(roomID, data)
	} else {
		log.WithError(err).Warnf("failed serializing RoomConfig of %s", roomID)
	}
}

// SaveRoomConfigs is called by roomservice to persist them
func (db *DB) SaveRoomConfigs(roomConfigs map[string]*roomservice.RoomConfig) {
	for roomID, rc := range roomConfigs {
		db.SaveRoomConfig(roomID, rc)
	}
}

// ReadRoomConfigs returns all saved room configs from disk
func (db *DB) ReadRoomConfigs(roomConfigs map[string]*roomservice.RoomConfig) {
	for roomID := range db.dvRoomConfig.Keys(nil) {
		data, err := db.dvRoomConfig.Read(roomID)
		if err == nil {
			roomConfigs[roomID] = &roomservice.RoomConfig{}
			err = json.Unmarshal(data, roomConfigs[roomID])
			if err != nil {
				log.WithError(err).Warnf("failed deserializing RoomConfig of %s", roomID)
				delete(roomConfigs, roomID)
			}
		} else {
			log.WithError(err).Warnf("failed reading diskv key for RoomConfig of %s", roomID)
		}
	}
}
