package pogo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	err := InvalidPokedexIDError{23}
	assert.Contains(t, err.Error(), "id 23")
}

func TestGender(t *testing.T) {
	assert.Equal(t, Unset, ToGender(0))
	assert.Equal(t, Male, ToGender(1))
	assert.Equal(t, Female, ToGender(2))
	assert.Equal(t, Genderless, ToGender(3))
	assert.Equal(t, Unset, ToGender(4))

	assert.Equal(t, "Unset", ToGender(0).ToString())
	assert.Equal(t, "Male", ToGender(1).ToString())
	assert.Equal(t, "Female", ToGender(2).ToString())
	assert.Equal(t, "Genderless", ToGender(3).ToString())
	assert.Equal(t, "Unset", ToGender(4).ToString())

	assert.Equal(t, "", ToGender(0).ToSymbol())
	assert.Equal(t, "♂", ToGender(1).ToSymbol())
	assert.Equal(t, "♀", ToGender(2).ToSymbol())
	assert.Equal(t, "", ToGender(3).ToSymbol())
	assert.Equal(t, "", ToGender(4).ToSymbol())
}

func TestGym(t *testing.T) {
	assert.Equal(t, Neutral, ToTeamColor(0))
	assert.Equal(t, Blue, ToTeamColor(1))
	assert.Equal(t, Red, ToTeamColor(2))
	assert.Equal(t, Yellow, ToTeamColor(3))
	assert.Equal(t, Neutral, ToTeamColor(4))

	assert.Equal(t, "white", ToTeamColor(0).ToString())
	assert.Equal(t, "blue", ToTeamColor(1).ToString())
	assert.Equal(t, "red", ToTeamColor(2).ToString())
	assert.Equal(t, "yellow", ToTeamColor(3).ToString())
	assert.Equal(t, "white", ToTeamColor(4).ToString())
}

func TestLocation(t *testing.T) {
	lat := 12.345
	lon := 67.8910
	latStr := fmt.Sprintf("%f", lat)
	lonStr := fmt.Sprintf("%f", lon)

	l := Location{lat, lon}

	osm := l.ToLinkOSM()
	assert.Contains(t, osm, "openstreetmap")
	assert.Contains(t, osm, latStr)
	assert.Contains(t, osm, lonStr)

	gm := l.ToLinkGMaps()
	assert.Contains(t, gm, "maps.google")
	assert.Contains(t, gm, latStr)
	assert.Contains(t, gm, lonStr)
}

func TestPokedex(t *testing.T) {
	testPokedexFile := "../../data/pokedex.json"

	_, err := NewPokedex("non/existent/329384723894")
	assert.Error(t, err)

	_, err = NewPokedex("../../test/data/invalid-pokedex.json")
	assert.Error(t, err)

	dex, err := NewPokedex(testPokedexFile)
	if assert.NoError(t, err) {
		assert.NotEqual(t, nil, dex)
	}

	ne, nd, err := dex.GetNamesByID(1)
	if assert.NoError(t, err) {
		assert.Equal(t, "Bulbasaur", ne)
		assert.Equal(t, "Bisasam", nd)
	}

	_, _, err = dex.GetNamesByID(0)
	assert.Error(t, err)
	_, _, err = dex.GetNamesByID(1 + len(dex.entries))
	assert.Error(t, err)
	_, _, err = dex.GetNamesByID(-1)
	assert.Error(t, err)
}
