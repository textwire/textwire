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

func testIntegerObject(t *testing.T, obj object.Object, expected int64) {
	result, ok := obj.(*object.Integer)

	if !ok {
		t.Errorf("obj is not an Integer, got %T", obj)
	}

	if result.Value != expected {
		t.Errorf("result.Value is not %d, got %d", expected, result.Value)
	}
}

func testStringObject(t *testing.T, obj object.Object, expected string) {
	result, ok := obj.(*object.String)

	if !ok {
		t.Errorf("obj is not a String, got %T", obj)
	}

	if result.Value != expected {
		t.Errorf("result.Value is not %s, got %s", expected, result.Value)
	}
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"{{ 5 }}", 5},
		{"{{ 10 }}", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalStringLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`{{ "Hello World" }}`, "Hello World"},
		{`{{ "She is pretty" }}`, "She is pretty"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}
