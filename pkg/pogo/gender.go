package pogo

// Gender for a Pokemon
type Gender int

// Gender enum
const (
	Unset Gender = iota
	Male
	Female
	Genderless
)

// ToGender converts an int from proto to a Gender
func ToGender(val int) Gender {
	switch val {
	case 1:
		return Male
	case 2:
		return Female
	case 3:
		return Genderless
	default:
		return Unset
	}
}

// ToString converts the Gender to a string with the name of the Gender
func (g Gender) ToString() (s string) {
	switch g {
	case 1:
		s = "Male"
	case 2:
		s = "Female"
	case 3:
		s = "Genderless"
	default:
		s = "Unset"
	}
	return
}

// ToSymbol converts the Gender to a symbol for short display
func (g Gender) ToSymbol() (s string) {
	switch g {
	case 1:
		s = "♂"
	case 2:
		s = "♀"
	default:
		s = ""
	}
	return
}
