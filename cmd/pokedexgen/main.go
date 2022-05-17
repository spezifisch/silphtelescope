package main

import (
	"encoding/json"
	"os"

	"github.com/mtslzr/pokeapi-go"
	"github.com/mtslzr/pokeapi-go/structs"
	log "github.com/sirupsen/logrus"

	"github.com/spezifisch/silphtelescope/pkg/pogo"
)

func main() {
	f := pokedexFetcher{
		OutputFile: "pokedex.json",
	}
	ok := f.Run()
	if ok {
		log.Info("success!")
	} else {
		log.Error("failure!")
	}
}

type pokedexFetcher struct {
	OutputFile string
	pokedex    []*pogo.PokedexEntry
}

func (f *pokedexFetcher) Run() bool {
	if _, err := os.Stat(f.OutputFile); err == nil {
		log.Errorf("output file %s already exists, bye", f.OutputFile)
		return false
	}

	speciesList, err := pokeapi.Resource("pokemon-species", 0, 99999)
	if err != nil {
		log.WithError(err).Error("fetching species failed")
		return false
	}

	f.pokedex = make([]*pogo.PokedexEntry, len(speciesList.Results))
	goodCount := 0

	for i, species := range speciesList.Results {
		log.Infof("fetching %d/%d %s", i, len(speciesList.Results), species.Name)

		var info structs.PokemonSpecies
		info, err = pokeapi.PokemonSpecies(species.Name)
		if err != nil {
			log.WithError(err).Errorf("failed fetching %s species info", species.Name)
			continue
		}

		arrayIdx := info.ID - 1
		if arrayIdx < 0 || arrayIdx >= len(speciesList.Results) {
			log.Errorf("invalid pokedex number #%d ignored", info.ID)
			continue
		}

		entry := &pogo.PokedexEntry{
			ID: info.ID,
		}

		for _, translName := range info.Names {
			switch translName.Language.Name {
			case "en":
				entry.NameEN = translName.Name
			case "de":
				entry.NameDE = translName.Name
			}
		}

		if f.pokedex[arrayIdx] == nil && entry.NameEN != "" && entry.NameDE != "" {
			goodCount++
		}

		f.pokedex[arrayIdx] = entry
	}

	if goodCount != len(f.pokedex) {
		log.Errorf("expected %d pokedex entries, but only got %d", len(f.pokedex), goodCount)
		return false
	}

	log.Infof("processed %d mons", len(f.pokedex))

	file, err := os.OpenFile(f.OutputFile, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	defer file.Close()
	if err != nil {
		log.WithError(err).Errorf("failed opening %s", f.OutputFile)
		return false
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(f.pokedex)
	if err != nil {
		log.WithError(err).Error("serializing pokedex to json failed")
		return false
	}

	return true
}
