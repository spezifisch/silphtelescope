package http

import (
	"encoding/json"
	"math/big"
)

// Location in geo coordinates
type Location struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

// GymMessage contains raid/egg and defender info
type GymMessage struct {
	Location
	GymID            string `json:"gym_id"`
	TeamID           int    `json:"team_id"`
	Name             string `json:"name"`
	SlotsAvailable   int    `json:"slots_available"`
	URL              string `json:"url,omitempty"`
	IsExRaidEligible int    `json:"is_ex_raid_eligible"`
}

// BasePokemon for raids and spawns
type BasePokemon struct {
	PokemonID int `json:"pokemon_id"`
	Form      int `json:"form,omitempty"`
	Gender    int `json:"gender,omitempty"`
	Costume   int `json:"costume,omitempty"`
}

// PokemonMessage describes a spawned pokemon
type PokemonMessage struct {
	Location
	BasePokemon
	EncounterID        big.Int `json:"encounter_id"`
	SpawnpointID       big.Int `json:"spawnpoint_id"`
	DisappearTimestamp int64   `json:"disappear_time"`
	KnownDisappearTime bool    `json:"verified"`
	Rarity             int     `json:"rarity"`
}

// RaidPokemon describes an egg or spawned raid boss
type RaidPokemon struct {
	CP        int `json:"cp"`
	Evolution int `json:"evolution"`
	Move1     int `json:"move_1"`
	Move2     int `json:"move_2"`

	// when spawned:
	BasePokemon
}

// RaidMessage contains raid info for a gym
type RaidMessage struct {
	Location
	Level            int    `json:"level"`
	TeamID           int    `json:"team_id"`
	StartTimestamp   int    `json:"start"`
	EndTimestamp     int    `json:"end"`
	Name             string `json:"name"`
	GymID            string `json:"gym_id"`
	URL              string `json:"url,omitempty"`
	IsExRaidEligible bool   `json:"is_ex_raid_eligible"`
	IsExclusive      bool   `json:"is_exclusive"`

	BasePokemon
}

// Envelope is the thing that MAD posts to the configures webhooks
type Envelope struct {
	Type    string          `json:"type"` // gym, pokemon, raid, pokestop, weather
	Message json.RawMessage `json:"message"`
}
