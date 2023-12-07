package parser

import (
	"testing"

	"github.com/textwire/textwire/lexer"
)

func TestParseIdentifier(t *testing.T) {
	inp := "{{ myName }}"

	l := lexer.New(inp)
	p := New(l)

	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program must have 1 statement, got %d", len(program.Statements))
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}
