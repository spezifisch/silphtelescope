package helpers

// IntArrayContains checks if array s contains value e.
// Source: https://stackoverflow.com/a/10485970
func IntArrayContains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// IntArrayUnorderedRemove removes element i from array s but changes element order.
// Source: https://stackoverflow.com/a/37335777
func IntArrayUnorderedRemove(s []int, i int) []int {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
