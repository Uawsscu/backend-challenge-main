package unit

import (
	"testing"

	"github.com/backend-challenge/user-api/pkg/validator"
)

func TestValidator_ValidateEmail(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"test@example.com", true},
		{"user@domain.co", true},
		{"invalid-email", false},
		{"@domain.com", false},
		{"test@", false},
		{"", false},
	}

	for _, tt := range tests {
		if got := validator.ValidateEmail(tt.email); got != tt.valid {
			t.Errorf("ValidateEmail(%q) = %v; want %v", tt.email, got, tt.valid)
		}
	}
}

func TestValidator_ValidatePassword(t *testing.T) {
	tests := []struct {
		password string
		valid    bool
	}{
		{"Password123!", true},
		{"validP@ss1", true},
		{"short", false},
		{"lowercaseonly123!", false},
		{"UPPERCASEONLY123!", false},
		{"NoNumber!", false},
		{"NoSpecial123", false},
		{"", false},
	}

	for _, tt := range tests {
		if got := validator.ValidatePassword(tt.password); got != tt.valid {
			t.Errorf("ValidatePassword(%q) = %v; want %v", tt.password, got, tt.valid)
		}
	}
}

func TestValidator_ValidateRequired(t *testing.T) {
	tests := []struct {
		value string
		valid bool
	}{
		{"some value", true},
		{" ", false},
		{"", false},
	}

	for _, tt := range tests {
		if got := validator.ValidateRequired(tt.value); got != tt.valid {
			t.Errorf("ValidateRequired(%q) = %v; want %v", tt.value, got, tt.valid)
		}
	}
}
