package pogo

import "fmt"

// InvalidPokedexIDError happens when you lookup a non-existent pokemon
type InvalidPokedexIDError struct {
	ID int // the ID that caused the error
}

func (e *InvalidPokedexIDError) Error() string {
	return fmt.Sprintf("pogo: invalid pokedex id %d", e.ID)
}
