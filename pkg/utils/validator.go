package utils

import (
	"regexp"
	"strings"
	"unicode"
)

var usernameRe = regexp.MustCompile(`^[a-zA-Z0-9_]{3,32}$`)

// ValidateUsername checks length and allowed characters.
func ValidateUsername(s string) bool {
	return usernameRe.MatchString(s)
}

// ValidatePassword enforces minimum strength: ≥8 chars, mixed case, at least one digit.
func ValidatePassword(s string) bool {
	if len(s) < 8 {
		return false
	}
	var hasUpper, hasLower, hasDigit bool
	for _, r := range s {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}
	return hasUpper && hasLower && hasDigit
}

// SanitizeString trims whitespace and returns empty string for pure-whitespace input.
func SanitizeString(s string) string {
	return strings.TrimSpace(s)
}
