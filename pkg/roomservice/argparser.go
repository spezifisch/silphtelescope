package roomservice

import (
	"errors"
	"strconv"
	"strings"

	"github.com/spezifisch/silphtelescope/pkg/pogo"
)

// ArgParser parses elements of an arg list as different types
type ArgParser struct {
	args []string
}

// NewArgParser takes an arg list like a string split at every space.
func NewArgParser(args []string) *ArgParser {
	return &ArgParser{args}
}

func (a *ArgParser) isValidIndex(index int) bool {
	return index >= 0 && index < len(a.args)
}

// Count returns the number of args
func (a *ArgParser) Count() int {
	return len(a.args)
}

// AsFloat parses the arg as float64
func (a *ArgParser) AsFloat(index int) (float64, error) {
	if !a.isValidIndex(index) {
		return 0, errors.New("invalid index")
	}
	return strconv.ParseFloat(a.args[index], 64)
}

// AsInt parses the arg as an int
func (a *ArgParser) AsInt(index int) (int, error) {
	if !a.isValidIndex(index) {
		return 0, errors.New("invalid index")
	}
	v, err := strconv.ParseInt(a.args[index], 10, 64)
	return int(v), err
}

// AsString returns the arg as is, but with index bounds check
func (a *ArgParser) AsString(index int) (string, error) {
	if !a.isValidIndex(index) {
		return "", errors.New("invalid index")
	}
	return a.args[index], nil
}

// AsIntArray returns the single value or comma separated values as an array
func (a *ArgParser) AsIntArray(index int) (arr []int, err error) {
	val, _ := a.AsString(index)
	gotMultiple := strings.Contains(val, ",")

	if !gotMultiple {
		// only got a single int
		var singleVal int
		singleVal, err = a.AsInt(index)
		if err != nil {
			return
		}
		arr = []int{singleVal}
		return
	}

	// split the numbers
	parts := strings.Split(val, ",")
	arr = make([]int, len(parts))

	// use an ArgParse to parse the numbers, why not?
	arrayParser := NewArgParser(parts)
	for i := 0; i < arrayParser.Count(); i++ {
		arr[i], err = arrayParser.AsInt(i)
		if err != nil {
			arr = nil
			return
		}
	}

	return
}

// AsLocation returns a pogo.Location from the given indices
func (a *ArgParser) AsLocation(latIdx, lonIdx int) (l pogo.Location, err error) {
	lat, err := a.AsFloat(latIdx)
	if err != nil {
		return
	}
	lon, err := a.AsFloat(lonIdx)
	if err != nil {
		return
	}
	// TODO check for inf, nan and stuff
	return pogo.Location{Latitude: lat, Longitude: lon}, nil
}

// AsLocationRadius returns a pogo.LocationRadius from the given indices
func (a *ArgParser) AsLocationRadius(latIdx, lonIdx, radiusIdx int) (lr pogo.LocationRadius, err error) {
	lat, err := a.AsFloat(latIdx)
	if err != nil {
		return
	}
	lon, err := a.AsFloat(lonIdx)
	if err != nil {
		return
	}
	radius, err := a.AsFloat(radiusIdx)
	if err != nil {
		return
	}
	// TODO check for inf, nan and stuff
	return pogo.NewLocationRadius(lat, lon, radius), nil
}
