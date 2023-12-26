package evaluator

import (
	"testing"

	"github.com/textwire/textwire/lexer"
	"github.com/textwire/textwire/object"
	"github.com/textwire/textwire/parser"
)

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnv()

	return Eval(program, env)
}

func evaluationExpected(t *testing.T, input, expect string) {
	result := testEval(input).String()

	if result != expect {
		t.Errorf("result is not %s, got %s", expect, result)
	}
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"{{ 5 }}", "5"},
		{"{{ 10 }}", "10"},
	}

	for _, tt := range tests {
		evaluationExpected(t, tt.input, tt.expected)
	}
}

func TestEvalStringLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`{{ "Hello World" }}`, "Hello World"},
		{`{{ "She \"is\" pretty" }}`, `She "is" pretty`},
	}

	for _, tt := range tests {
		evaluationExpected(t, tt.input, tt.expected)
	}
}

func TestEvalReturnStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`{{ return }}`, ""},
		{`{{ return 23 }}`, "23"},
	}

	for _, tt := range tests {
		evaluationExpected(t, tt.input, tt.expected)
	}
}
