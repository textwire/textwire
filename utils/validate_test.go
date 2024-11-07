package utils

import "testing"

func TestStrIsInt(t *testing.T) {
	tc := []struct {
		inp      string
		expected bool
	}{
		{"anna", false},
		{"123", true},
		{"-123", true},
		{"0", true},
		{"-1", true},
		{"123.23", false},
		{"123.0", false},
	}

	for _, tt := range tc {
		t.Run("Test case: "+tt.inp, func(t *testing.T) {
			got := StrIsInt(tt.inp)

			if got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}
