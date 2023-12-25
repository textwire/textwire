package evaluator

import (
	"testing"

	"github.com/textwire/textwire/lexer"
	"github.com/textwire/textwire/object"
	"github.com/textwire/textwire/parser"
)

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
