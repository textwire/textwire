package utils

import "testing"

func TestFloatToStr(t *testing.T) {
	tc := []struct {
		name     string
		input    float64
		expected string
	}{
		{name: "Positive float", input: 3.14159, expected: "3.14159"},
		{name: "Negative float", input: -3.14159, expected: "-3.14159"},
		{name: "Zero as float", input: 0.0, expected: "0"},
		{name: "Negative zero as float", input: 0.0, expected: "0"},
		{name: "Float with fraction", input: 123.23, expected: "123.23"},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			actual := FloatToStr(tt.input)

			if actual != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, actual)
			}
		})
	}
}
