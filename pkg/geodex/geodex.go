package geodex

import (
	"errors"

	"github.com/spezifisch/silphtelescope/pkg/pogo"
)

// GeoDex is the wrapper that should be used when using an already initialized DB
type GeoDex struct {
	Disk *DiskDB
	Tile *TDB
}

// NewGeoDex connects to Tile38 DB and sets up diskv ready to supply fort info
func NewGeoDex(ddbBasePath, tdbHostname, tdbPassword string) (gd *GeoDex, err error) {
	d := NewDiskDB(&ddbBasePath)
	t, err := NewTDB(tdbHostname, tdbPassword)
	if err != nil {
		return
	}

	gd = &GeoDex{
		Disk: d,
		Tile: t,
	}
	return
}

// LookupFortNear get the nearest fort within the radius and resolves its name
func (gd *GeoDex) LookupFortNear(point pogo.Location, radiusM float64) (f *Fort, err error) {
	// get nearest fort from tile38
	f, err = gd.Tile.GetNearestFort(point, radiusM)
	if err != nil {
		return
	}
	if f.GUID == nil {
		err = errors.New("nearest fort GUID is nil")
		return
	}

	// get fort name from diskv because tile38 doesn't store the name
	diskFort, err := gd.Disk.GetFort(*f.GUID)
	if err != nil {
		return f, nil // that's ok, just return the fort without name
	}

	f.Name = diskFort.Name
	return
}
