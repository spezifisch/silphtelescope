package geodex

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/spezifisch/silphtelescope/pkg/pogo"

	t38c "github.com/axvq/tile38-client"
)

// TDB is the silphtelescope geodex connection
type TDB struct {
	db *t38c.Client
}

// NewTDB returns a usable silpht db object
func NewTDB(hostname, password string) (db *TDB, err error) {
	db = &TDB{}

	// connect to tile38
	if password != "" {
		db.db, err = t38c.New(hostname, t38c.WithPassword(password))
	} else {
		db.db, err = t38c.New(hostname)
	}
	return
}

// Close closes the database connection
func (tdb *TDB) Close() {
	if tdb.db != nil {
		tdb.db.Close()
	}
}

// Drop deletes the whole fort database
func (tdb *TDB) Drop() (err error) {
	err = tdb.db.Keys.Drop("fort")
	return
}

// InsertFort adds the fort to the db
func (tdb *TDB) InsertFort(f *Fort) (err error) {
	if f.GUID == nil {
		log.Warn("tried to insert fort without GUID")
		return
	}

	err = tdb.db.Keys.Set("fort", *f.GUID).Point(f.Latitude, f.Longitude).
		Field("type", float64(f.Type)).
		Do()

	return
}

func toFortType(typeFieldVal float64) (t FortType) {
	switch int(typeFieldVal) {
	case 1:
		t = FortTypeGym
	case 2:
		t = FortTypeStop
	}
	return t
}

// ToString returns the fort type in human-readable form
func (t FortType) ToString() (s string) {
	switch t {
	case 0:
		s = "Portal"
	case 1:
		s = "Gym"
	case 2:
		s = "Stop"
	default:
		s = "Invalid"
	}

	return
}

// GetNearestFort looks in the given radius (in meters) around the point for the nearest Fort
func (tdb *TDB) GetNearestFort(point pogo.Location, radiusM float64) (f *Fort, err error) {
	response, err := tdb.db.Search.Nearby("fort",
		float64(point.Latitude), float64(point.Longitude), radiusM).
		Format(t38c.FormatPoints).
		Do()
	if err != nil {
		return
	}
	if len(response.Points) < 1 {
		err = errors.New("no fort found")
		return
	}

	pt := response.Points[0]
	f = &Fort{
		GUID:      &pt.ID,
		Latitude:  pt.Point.Lat,
		Longitude: pt.Point.Lon,
		Type:      toFortType(pt.Fields[0]),
	}
	return
}
