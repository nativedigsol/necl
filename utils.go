package necl

import "strings"

// ContainsMany will compare a target string with multiple strings
func ContainsMany(target string, compareTo []string) bool {
	for _, compare := range compareTo {
		if strings.Contains(target, compare) {
			return true
		}
	}

	return false
}
