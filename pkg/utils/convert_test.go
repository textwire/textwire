package utils

import (
	"math"
	"testing"
)

func TestFloatToStr(t *testing.T) {
	tc := []struct {
		name   string
		inp    float64
		expect string
	}{
		{name: "Positive float", inp: 3.14159, expect: "3.14159"},
		{name: "Negative float", inp: -3.14159, expect: "-3.14159"},
		{name: "Zero as float", inp: 0.0, expect: "0.0"},
		{name: "Negative zero as float", inp: 0.0, expect: "0.0"},
		{name: "Float with fraction", inp: 123.23, expect: "123.23"},
		{name: "Zero at the end is present", inp: 1.0, expect: "1.0"},
		{name: "Zero at the end is present", inp: 1.000, expect: "1.0"},
		{name: "Medium float", inp: 1234567890.1234567, expect: "1234567890.1234567"},
		{name: "Large float", inp: 1.234567890123456e+30, expect: "1.234567890123456e+30"},
		{name: "Very small float", inp: 0.00000123456789, expect: "0.00000123456789"},
		{name: "Negative medium float", inp: -987654321.9876543, expect: "-987654321.9876543"},
		{name: "Float with many trailing zeros", inp: 42.0000000000, expect: "42.0"},
		{name: "Very large decimal", inp: 999999999999999.999, expect: "1000000000000000.0"},
		{name: "Positive infinity", inp: math.Inf(1), expect: "+Inf"},
		{name: "Negative infinity", inp: math.Inf(-1), expect: "-Inf"},
		{name: "NaN", inp: math.NaN(), expect: "NaN"},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			if actual := FloatToStr(tt.inp); actual != tt.expect {
				t.Errorf("expected %q, got %q", tt.expect, actual)
			}
		})
	}
}
