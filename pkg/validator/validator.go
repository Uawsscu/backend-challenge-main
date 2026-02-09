package validator

import (
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// ValidateEmail checks if the email format is valid
func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// ValidateRequired checks if a string field is not empty
func ValidateRequired(value string) bool {
	return strings.TrimSpace(value) != ""
}

// ValidatePassword checks if password meets minimum requirements
func ValidatePassword(password string) bool {
	// At least 6 characters
	return len(password) >= 6
}
