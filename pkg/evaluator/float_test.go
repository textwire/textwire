package evaluator

import (
	"math"
	"testing"

	"github.com/textwire/textwire/v5/config"
	"github.com/textwire/textwire/v5/pkg/value"
)

// TestFloatArithmeticMatchesGo pins down that float arithmetic follows IEEE 754
// like Go does: overflow goes to infinity, underflow goes to zero, and NaN
// propagates, all without an error.
func TestFloatInfixExp_MatchesGoArithmetic(t *testing.T) {
	e := New(&config.Func{}, nil)
	ctx := NewContext(value.NewScope(), "")

	cases := []struct {
		name string
		op   string
		a, b float64
		want float64
	}{
		{
			name: "positive overflow",
			op:   "+",
			a:    math.MaxFloat64,
			b:    math.MaxFloat64,
			want: math.Inf(1),
		},
		{
			name: "negative overflow",
			op:   "*",
			a:    -math.MaxFloat64,
			b:    2,
			want: math.Inf(-1),
		},
		{
			name: "underflow to zero",
			op:   "/",
			a:    math.SmallestNonzeroFloat64,
			b:    2,
			want: 0,
		},
		{
			name: "infinity minus infinity",
			op:   "-",
			a:    math.Inf(1),
			b:    math.Inf(1),
			want: math.NaN(),
		},
		{
			name: "NaN propagates",
			op:   "+",
			a:    math.NaN(),
			b:    1,
			want: math.NaN(),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := e.floatInfixExp(
				tc.op,
				&value.Float{Val: tc.b},
				&value.Float{Val: tc.a},
				nil,
				ctx,
			)

			result, ok := got.(*value.Float)
			if !ok {
				t.Fatalf("expected *value.Float, got %T", got)
			}

			if !sameFloat(result.Val, tc.want) {
				t.Fatalf("expected %v, got %v", tc.want, result.Val)
			}
		})
	}
}

func sameFloat(a, b float64) bool {
	if math.IsNaN(b) {
		return math.IsNaN(a)
	}
	return a == b
}
