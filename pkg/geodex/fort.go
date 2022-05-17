package geodex

import (
	"fmt"

	"github.com/spezifisch/silphtelescope/pkg/pogo"
)

// Fort is the data structure for our tile38 table
type Fort struct {
	GUID      *string  `json:"guid"` // mandatory
	Latitude  float64  `json:"latitude"`
	Longitude float64  `json:"longitude"`
	Name      *string  `json:"name,omitempty"` // optional
	Type      FortType `json:"type"`
}

// FortType says if a fort is a portal, gym, or pokestop
type FortType int

// yup
const (
	FortTypePortal FortType = iota
	FortTypeGym
	FortTypeStop
)

// ToString returns all the fort's fields in a string
func (f *Fort) ToString() string {
	name := "nil"
	if f.Name != nil {
		name = f.GetName()
	}
	return fmt.Sprintf("GUID=%s Type=%s (%f,%f) Name: %s",
		*f.GUID, f.Type.ToString(), f.Latitude, f.Longitude, name)
}

// GetName returns the fort's name if set, "<Type>:<GUID>" string otherwise
func (f *Fort) GetName() string {
	if f.Name == nil {
		return fmt.Sprintf("%s:%s", f.Type.ToString(), *f.GUID)
	}
	return *f.Name
}

// Location returns the fort's position as a Location object
func (f *Fort) Location() *pogo.Location {
	return &pogo.Location{
		Latitude:  f.Latitude,
		Longitude: f.Longitude,
	}
}
