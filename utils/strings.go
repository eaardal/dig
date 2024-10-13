package utils

import "regexp"

// ContainsSpecialChars checks if a string contains special characters.
// Only letters, numbers, underscores and hyphens are allowed.
func ContainsSpecialChars(value string) bool {
	regex := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	return regex.MatchString(value)
}
