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
	evaluated := testEval(input)

	errObj, ok := evaluated.(*object.Error)

	if ok {
		t.Errorf("evaluation failed: %s", errObj.Message)
	}

	result := evaluated.String()

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
		{"{{ -123 }}", "-123"},
		{`{{ 5 + 5 }}`, "10"},
		{`{{ 11 + 13 - 1 }}`, "23"},
		{`{{ 3 * 3 * 3 + 10 }}`, "37"},
		{`{{ (5 + 10 * 2 + 15 / 3) * 2 + -10 }}`, "50"},
	}

	for _, tt := range tests {
		evaluationExpected(t, tt.input, tt.expected)
	}
}

func TestEvalNilLiteral(t *testing.T) {
	input := "<h1>{{ nil }}</h1>"
	evaluationExpected(t, input, "<h1></h1>")
}

func TestEvalStringLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`{{ "Hello World" }}`, "Hello World"},
		{`{{ "She \"is\" pretty" }}`, `She "is" pretty`},
		{`{{ "Korotchaeva" + " " + "Anna" }}`, "Korotchaeva Anna"},
		{`{{ "She" + " " + "is" + " " + "smart" }}`, "She is smart"},
	}

	for _, tt := range tests {
		evaluationExpected(t, tt.input, tt.expected)
	}
}
