package geodex

import (
	"encoding/json"
	"errors"
	"path/filepath"

	"github.com/peterbourgon/diskv"
)

// DiskDB stores fort info like names
type DiskDB struct {
	forts *diskv.Diskv
}

// NewDiskDB returns a ready-to-use DiskDB object
func NewDiskDB(basePath *string) (db *DiskDB) {
	pathFortDex := filepath.Join(*basePath, "fort")

	return &DiskDB{
		forts: diskv.New(diskv.Options{
			BasePath:     pathFortDex,
			Transform:    blockTransform,
			CacheSizeMax: 1 * 1024 * 1024,
		}),
	}
}

// Drop deletes the whole db
func (db *DiskDB) Drop() error {
	return db.forts.EraseAll()
}

// SaveFort saves the fort's info in the db
func (db *DiskDB) SaveFort(f *Fort) (err error) {
	if f.GUID == nil {
		return errors.New("cannot save fort with nil GUID")
	}

	data, err := json.Marshal(f)
	if err == nil {
		db.forts.Write(*f.GUID, data)
	}
	return
}

// GetFort returns the fort's info from the db
func (db *DiskDB) GetFort(GUID string) (f *Fort, err error) {
	data, err := db.forts.Read(GUID)
	if err != nil {
		return
	}

	f = &Fort{}
	err = json.Unmarshal(data, f)
	return
}

// MergeFort copies new values to an existing fort, or created a fort if it doesn't exist
func (db *DiskDB) MergeFort(f *Fort) (err error) {
	if f.GUID == nil {
		return errors.New("cannot save fort with nil GUID")
	}

	guid := *f.GUID
	data, err := db.GetFort(guid)
	if err != nil {
		// doesn't exist
		return db.SaveFort(f)
	}

	// never update GUID. the GUID is already the file path and name.
	data.Latitude = f.Latitude
	data.Longitude = f.Longitude
	data.Type = f.Type
	if data.Name == nil || *data.Name == "" {
		// Update the name if it isn't already set.
		// That's the whole reason for this function.
		data.Name = f.Name
	}
	return db.SaveFort(data)
}

// limit directory levels
const maxSliceSize = 3

// blockTransform based on: https://github.com/peterbourgon/diskv/blob/fc0553497cbfcf78f101d0bf8e82c6e627f4bbb0/examples/content-addressable-store/cas.go
const transformBlockSize = 2 // grouping of chars per directory depth

func blockTransform(s string) []string {
	sliceSize := len(s) / transformBlockSize
	if sliceSize > maxSliceSize {
		sliceSize = maxSliceSize
	}
	pathSlice := make([]string, sliceSize)

	for i := 0; i < sliceSize; i++ {
		from, to := i*transformBlockSize, (i*transformBlockSize)+transformBlockSize
		pathSlice[i] = s[from:to]
	}
	return pathSlice
}
