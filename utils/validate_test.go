package utils

import "testing"

func TestStrIsInt(t *testing.T) {
	tc := []struct {
		name     string
		inp      string
		expected bool
	}{
		{"Non-integer string", "anna", false},
		{"Positive integer", "123", true},
		{"Negative integer", "-123", true},
		{"Zero as integer", "0", true},
		{"Negative one", "-1", true},
		{"Decimal number with fraction", "123.23", false},
		{"Decimal number ending with zero", "123.0", false},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			got := StrIsInt(tt.inp)

			if got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}
