package db

import (
	"path/filepath"
	"strings"

	"github.com/peterbourgon/diskv"
)

// DB stores dynamic state so that the bot can be restarted seamlessly
type DB struct {
	dvRoomConfig, dvRoomState *diskv.Diskv
}

// NewDB returns a ready-to-use DB object
func NewDB(basePath *string) (db *DB) {
	// !room@home.server --> /home.server/room
	// Room ID spec: https://matrix.org/docs/spec/appendices#room-ids-and-event-ids
	matrixRoomIDTransform := func(s string) []string {
		s = strings.TrimLeft(s, "!")
		parts := strings.SplitN(s, ":", 2)
		room := parts[0]
		homeserver := parts[1]
		return []string{homeserver, room}
	}

	pathRoomConfig := filepath.Join(*basePath, "roomconfig")
	pathRoomState := filepath.Join(*basePath, "roomstate")

	return &DB{
		dvRoomConfig: diskv.New(diskv.Options{
			BasePath:  pathRoomConfig,
			Transform: matrixRoomIDTransform,
		}),
		dvRoomState: diskv.New(diskv.Options{
			BasePath:  pathRoomState,
			Transform: matrixRoomIDTransform,
		}),
	}
}
