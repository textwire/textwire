package value

import (
	"math"
	"testing"
)

func TestFloatJSON(t *testing.T) {
	cases := []struct {
		name   string
		val    float64
		expect string
	}{
		{name: "regular float", val: 3.14, expect: "3.14"},
		{name: "whole float", val: 2, expect: "2.0"},
		{name: "NaN", val: math.NaN(), expect: "null"},
		{name: "positive infinity", val: math.Inf(1), expect: "null"},
		{name: "negative infinity", val: math.Inf(-1), expect: "null"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := (&Float{Val: tc.val}).JSON()
			if err != nil {
				t.Fatal(err)
			}

			if got != tc.expect {
				t.Fatalf("expect %q, got %q", tc.expect, got)
			}
		})
	}
}
