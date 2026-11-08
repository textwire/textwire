package lexer

import (
	"testing"

	"codeberg.org/textwire/textwire/v5/pkg/token"
)

func TestClosingFor(t *testing.T) {
	cases := []struct {
		name   string
		open   token.TokenType
		expect token.TokenType
	}{
		{name: "{{ closes with }}", open: token.LBRACES, expect: token.RBRACES},
		{name: "{!! closes with !!}", open: token.LBRACESRAW, expect: token.RBRACESRAW},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := closingFor(tc.open); got != tc.expect {
				t.Fatalf("expect %s, got %s", token.String(tc.expect), token.String(got))
			}
		})
	}
}

func TestClosingFor_PanicsOnUnexpectedToken(t *testing.T) {
	cases := []struct {
		name string
		open token.TokenType
	}{
		{name: "closing }} is not an opening token", open: token.RBRACES},
		{name: "closing !!} is not an opening token", open: token.RBRACESRAW},
		{name: "directive is not an opening token", open: token.IF},
		{name: "EOF is not an opening token", open: token.EOF},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatalf("expect panic for %s", token.String(tc.open))
				}
			}()

			closingFor(tc.open)
		})
	}
}
