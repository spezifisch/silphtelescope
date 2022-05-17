package pogo

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"

	log "github.com/sirupsen/logrus"
)

// PokedexEntry associates a Pokedex ID with international names
type PokedexEntry struct {
	ID     int
	NameEN string
	NameDE string
}

// Pokedex holds the whole dex for lookups
type Pokedex struct {
	fileName string
	entries  []*PokedexEntry
}

// NewPokedex creates a ready-to-use Pokedex
func NewPokedex(fileName string) (p *Pokedex, err error) {
	p = &Pokedex{
		fileName: fileName,
	}
	err = p.ReadFile()
	return
}

// ReadFile parses pokedex data from a json file
func (p *Pokedex) ReadFile() (err error) {
	file, err := ioutil.ReadFile(p.fileName)
	if err != nil {
		return
	}

	p.entries = []*PokedexEntry{}
	err = json.Unmarshal(file, &p.entries)
	if err != nil {
		return
	}

	log.Infof("read %d pokedex entries", len(p.entries))
	return
}

// GetNamesByID returns the english and german name of the pokemon with its id from 1-898
func (p *Pokedex) GetNamesByID(id int) (nameEN, nameDE string, err error) {
	arrayIdx := id - 1
	if arrayIdx < 0 || arrayIdx >= len(p.entries) {
		err = &InvalidPokedexIDError{id}
		return
	}

	entry := p.entries[arrayIdx]
	nameEN = entry.NameEN
	nameDE = entry.NameDE
	return
}

// GetIDByName returns ID and name for given english/german name
func (p *Pokedex) GetIDByName(wantedName string) (id int, nameEN, nameDE string, err error) {
	wantedName = strings.ToLower(wantedName)
	for i, entry := range p.entries {
		if strings.ToLower(entry.NameDE) == wantedName || strings.ToLower(entry.NameEN) == wantedName {
			id = 1 + i
			nameEN = entry.NameEN
			nameDE = entry.NameDE
			return
		}
	}
	err = errors.New("Pokemon not found")
	return
}
