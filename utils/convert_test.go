package utils

import "testing"

func TestFloatToStr(t *testing.T) {
	tc := []struct {
		name   string
		input  float64
		expect string
	}{
		{name: "Positive float", input: 3.14159, expect: "3.14159"},
		{name: "Negative float", input: -3.14159, expect: "-3.14159"},
		{name: "Zero as float", input: 0.0, expect: "0.0"},
		{name: "Negative zero as float", input: 0.0, expect: "0.0"},
		{name: "Float with fraction", input: 123.23, expect: "123.23"},
		{name: "Zero at the end is present", input: 1.0, expect: "1.0"},
		{name: "Zero at the end is present", input: 1.000, expect: "1.0"},
		{name: "Medium float", input: 1234567890.1234567, expect: "1234567890.1234567"},
		{name: "Large float", input: 1.234567890123456e+30, expect: "1.234567890123456e+30"},
		{name: "Very small float", input: 0.00000123456789, expect: "0.00000123456789"},
		{name: "Negative medium float", input: -987654321.9876543, expect: "-987654321.9876543"},
		{name: "Float with many trailing zeros", input: 42.0000000000, expect: "42.0"},
		{name: "Very large decimal", input: 999999999999999.999, expect: "1000000000000000.0"},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			actual := FloatToStr(tt.input)

			if actual != tt.expect {
				t.Errorf("expect %s, got %s", tt.expect, actual)
			}
		})
	}
}
