package db

import (
	"testing"

	"github.com/spezifisch/silphtelescope/pkg/roomservice"
	"github.com/stretchr/testify/assert"
)

func TestRoomConfig(t *testing.T) {
	basePath := "test-data-2342"
	db := NewDB(&basePath)
	defer db.dvRoomConfig.EraseAll()

	readRCs := make(map[string]*roomservice.RoomConfig)
	db.ReadRoomConfigs(readRCs)
	assert.Equal(t, 0, len(readRCs))

	roomA := "!foo:bar.baz"
	rcA := &roomservice.RoomConfig{
		RoomID: roomA,
	}
	roomB := "!pebkac:snafu.bofh"
	rcB := &roomservice.RoomConfig{
		RoomID: roomB,
	}
	rcs := map[string]*roomservice.RoomConfig{
		roomA: rcA,
		roomB: rcB,
	}

	db.SaveRoomConfig(roomA, rcA)
	readRCs = make(map[string]*roomservice.RoomConfig)
	db.ReadRoomConfigs(readRCs)
	assert.Equal(t, 1, len(readRCs))

	db.SaveRoomConfigs(rcs)

	readRCs = make(map[string]*roomservice.RoomConfig)
	db.ReadRoomConfigs(readRCs)
	assert.Equal(t, 2, len(readRCs))
}

func TestRoomState(t *testing.T) {
	basePath := "test-data-2342"
	db := NewDB(&basePath)
	defer db.dvRoomState.EraseAll()

	readRSs := make(map[string]*roomservice.RoomState)
	db.ReadRoomStates(readRSs)
	assert.Equal(t, 0, len(readRSs))

	roomA := "!foo:bar.baz"
	rsA := &roomservice.RoomState{
		Spawns: make(map[string]*roomservice.SpawnState),
	}
	rsA.Spawns["encfoo"] = &roomservice.SpawnState{
		EndTime: 23,
		Posted:  true,
	}
	rsA.Spawns["bar"] = &roomservice.SpawnState{
		EndTime: 1337,
		Posted:  false,
	}

	roomB := "!pebkac:snafu.bofh"
	rsB := &roomservice.RoomState{
		Spawns: make(map[string]*roomservice.SpawnState),
	}
	rsB.Spawns["encfoo"] = &roomservice.SpawnState{
		EndTime: 23,
		Posted:  false,
	}

	rss := map[string]*roomservice.RoomState{
		roomA: rsA,
		roomB: rsB,
	}

	db.SaveRoomState(roomA, rsA)
	readRSs = make(map[string]*roomservice.RoomState)
	db.ReadRoomStates(readRSs)
	assert.Equal(t, 1, len(readRSs))

	db.SaveRoomStates(rss)

	readRSs = make(map[string]*roomservice.RoomState)
	db.ReadRoomStates(readRSs)
	assert.Equal(t, 2, len(readRSs))
}

func TestRCRSSimultaneous(t *testing.T) {
	done := make(chan bool)

	go func() {
		TestRoomConfig(t)
		done <- true
	}()
	go func() {
		TestRoomState(t)
		done <- true
	}()

	<-done
	<-done
}
