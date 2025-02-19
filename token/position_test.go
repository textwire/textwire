package token

import "testing"

func TestContains(t *testing.T) {
	testCases := []struct {
		desc   string
		line   uint
		col    uint
		expect bool
	}{
		{
			desc:   "line is out of range",
			line:   3,
			col:    6,
			expect: false,
		},
		{
			desc:   "column is one position before the token",
			line:   4,
			col:    4,
			expect: false,
		},
		{
			desc:   "column is one position after the token",
			line:   4,
			col:    8,
			expect: false,
		},
		{
			desc:   "column and line are in range at start of token",
			line:   4,
			col:    5,
			expect: true,
		},
		{
			desc:   "column and line are in range at center of token",
			line:   4,
			col:    6,
			expect: true,
		},
		{
			desc:   "column and line are in range at end of token",
			line:   4,
			col:    7,
			expect: true,
		},
	}

	token := Token{
		Type:    IDENT,
		Literal: "foo",
		Pos: Position{
			StartLine: 4,
			StartCol:  5,
			EndLine:   4,
			EndCol:    7,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actual := token.Pos.Contains(tc.line, tc.col)

			if actual != tc.expect {
				t.Errorf("Got: %v, Expect: %v", actual, tc)
			}
		})
	}
}
