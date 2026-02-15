package token

import "testing"

func TestContains(t *testing.T) {
	tokenVar := Token{
		Type:    IDENT,
		Literal: "foo",
		Pos: Position{
			StartLine: 4,
			StartCol:  5,
			EndLine:   4,
			EndCol:    7,
		},
	}

	tokenHTML := Token{
		Type:    HTML,
		Literal: "<div>\n    <h1>Hello</h1>\n</div>",
		Pos: Position{
			StartLine: 2,
			StartCol:  0,
			EndLine:   4,
			EndCol:    5,
		},
	}

	testCases := []struct {
		desc   string
		line   uint
		col    uint
		token  Token
		expect bool
	}{
		{
			desc:   "Cursor is out of range",
			line:   3,
			col:    6,
			token:  tokenVar,
			expect: false,
		},
		{
			// (f)oo
			desc: "Cursor is at the start of the file",
			line: 0,
			col:  0,
			token: Token{
				Type:    IDENT,
				Literal: "foo",
				Pos: Position{
					StartLine: 0,
					StartCol:  0,
					EndLine:   0,
					EndCol:    2,
				},
			},
			expect: true,
		},
		{
			// ( )foo
			desc:   "Cursor is one position before token",
			line:   4,
			col:    4,
			token:  tokenVar,
			expect: false,
		},
		{
			// (f)oo
			desc:   "Cursor is exactly at start of token",
			line:   4,
			col:    5, // Assuming token col is 5-7
			token:  tokenVar,
			expect: true,
		},
		{
			// f(o)o
			desc:   "Cursor is exactly at center of token",
			line:   4,
			col:    6, // Assuming token col is 5-7
			token:  tokenVar,
			expect: true,
		},
		{
			// fo(o)
			desc:   "Cursor is exactly at end of token",
			line:   4,
			col:    7, // Assuming token col is 5-7
			token:  tokenVar,
			expect: true,
		},
		{
			// foo( )
			desc:   "Cursor is one position after token",
			line:   4,
			col:    8,
			token:  tokenVar,
			expect: false,
		},
		{
			// (<)div>
			desc:   "Cursor is inside a multi-line token (start line)",
			line:   2, // Assuming token line is 2-4
			col:    0,
			token:  tokenHTML,
			expect: true,
		},
		{
			// ____<h1>(H)ellow</h1>
			desc:   "Cursor is inside a multi-line token (middle line)",
			line:   3, // Assuming token line is 2-4
			col:    9,
			token:  tokenHTML,
			expect: true,
		},
		{
			// </div(>)
			desc:   "Cursor is inside a multi-line token (last line)",
			line:   4, // Assuming token line is 2-4
			col:    5,
			token:  tokenHTML,
			expect: true,
		},
		{
			// </div>( )
			desc:   "Cursor is one column after token",
			line:   4,
			col:    6,
			token:  tokenHTML,
			expect: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			if actual := tc.token.Pos.Contains(tc.line, tc.col); actual != tc.expect {
				t.Errorf("Expected position %v but got %v", tc, actual)
			}
		})
	}
}
