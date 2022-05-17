package pogo

import (
	"fmt"
	"math"
)

// Location contains the geocoordinates of something
type Location struct {
	Latitude  float64
	Longitude float64
}

// ToLinkOSM return an OpenStreetMap link with a marker
func (l Location) ToLinkOSM() string {
	zoomLevel := 17
	return fmt.Sprintf("https://www.openstreetmap.org/?mlat=%f&mlon=%f#map=%d/%f/%f",
		l.Latitude, l.Longitude, zoomLevel, l.Latitude, l.Longitude)
}

// ToLinkGMaps return a Google Maps link with a marker
func (l Location) ToLinkGMaps() string {
	return fmt.Sprintf("https://maps.google.de/maps?q=%f,%f",
		l.Latitude, l.Longitude)
}

// Based on golang-geo, MIT license: https://github.com/kellydunn/golang-geo/blob/2c6d5d781da2b42e60f1d7507c1055ef625d8652/point.go
const (
	// According to Wikipedia, the Earth's radius is about 6,371km
	earthRadiusKM = 6371
)

// DistanceTo calculates the Haversine distance between two points in meters
// Original Implementation from: http://www.movable-type.co.uk/scripts/latlong.html
// Based on golang-geo, MIT license: https://github.com/kellydunn/golang-geo/blob/2c6d5d781da2b42e60f1d7507c1055ef625d8652/point.go
func (l Location) DistanceTo(p2 *Location) float64 {
	dLat := (p2.Latitude - l.Latitude) * (math.Pi / 180.0)
	dLon := (p2.Longitude - l.Longitude) * (math.Pi / 180.0)

	lat1 := l.Latitude * (math.Pi / 180.0)
	lat2 := p2.Latitude * (math.Pi / 180.0)

	a1 := math.Sin(dLat/2) * math.Sin(dLat/2)
	a2 := math.Sin(dLon/2) * math.Sin(dLon/2) * math.Cos(lat1) * math.Cos(lat2)

	a := a1 + a2

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return 1000 * earthRadiusKM * c
}

// BearingTo Calculates the initial bearing (sometimes referred to as forward azimuth), in degrees
// Original Implementation from: http://www.movable-type.co.uk/scripts/latlong.html
// Based on golang-geo, MIT license: https://github.com/kellydunn/golang-geo/blob/2c6d5d781da2b42e60f1d7507c1055ef625d8652/point.go
func (l Location) BearingTo(p2 *Location) float64 {
	dLon := (p2.Longitude - l.Longitude) * math.Pi / 180.0

	lat1 := l.Latitude * math.Pi / 180.0
	lat2 := p2.Latitude * math.Pi / 180.0

	y := math.Sin(dLon) * math.Cos(lat2)
	x := math.Cos(lat1)*math.Sin(lat2) -
		math.Sin(lat1)*math.Cos(lat2)*math.Cos(dLon)
	brng := math.Atan2(y, x) * 180.0 / math.Pi

	return brng
}

// LocationRadius is a circle around the given coordinates with the given radius
type LocationRadius struct {
	Location
	RadiusM float64
}

// NewLocationRadius conveniently returns the area
func NewLocationRadius(lat, lon, radiusM float64) LocationRadius {
	return LocationRadius{
		Location: Location{
			Latitude:  lat,
			Longitude: lon,
		},
		RadiusM: radiusM,
	}
}

// Contains returns true if this area contains the given location
func (lr LocationRadius) Contains(o *Location) bool {
	distance := lr.Location.DistanceTo(o)
	return distance <= lr.RadiusM
}
