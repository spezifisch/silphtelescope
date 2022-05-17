package geodex

import "regexp"

var guidRegex *regexp.Regexp

// IsValidGUID returns true if the value might be a valid fort GUID
func IsValidGUID(val string) bool {
	if guidRegex == nil {
		guidRegex = regexp.MustCompile("^[a-fA-F0-9.]+$")
	}

	return guidRegex.Match([]byte(val))
}
