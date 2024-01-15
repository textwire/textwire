package evaluator

import (
	"strings"
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
		t.Errorf("result is not %q, got %q", expect, result)
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
		{"{{ 2 * (5 + 10) }}", "30"},
		{`{{ 3 * 3 * 3 + 10 }}`, "37"},
		{`{{ (5 + 10 * 2 + 15 / 3) * 2 + -10 }}`, "50"},
	}

	for _, tt := range tests {
		evaluationExpected(t, tt.input, tt.expected)
	}
}

func TestEvalFloatExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"{{ 5.11 }}", "5.11"},
		{"{{ -12.3 }}", "-12.3"},
		{`{{ 2.123 + 1.111 }}`, "3.234"},
		{`{{ 2. + 1.2 }}`, "3.2"},
	}

	for _, tt := range tests {
		evaluationExpected(t, tt.input, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"{{ true }}", "1"},
		{"{{ false }}", "0"},
		{"{{ !true }}", "0"},
		{"{{ !false }}", "1"},
		{"{{ !nil }}", "1"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)

		if ok {
			t.Errorf("evaluation failed: %s", errObj.Message)
			return
		}

		result := evaluated.String()

		if result != tt.expected {
			t.Errorf("result is not %s, got %s", tt.expected, result)
		}
	}
}

func TestEvalNilExpression(t *testing.T) {
	inp := "<h1>{{ nil }}</h1>"
	evaluationExpected(t, inp, "<h1></h1>")
}

func TestEvalStringExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`{{ "Hello World" }}`, "Hello World"},
		{`{{ "She \"is\" pretty" }}`, `She "is" pretty`},
		{`{{ "Korotchaeva" + " " + "Anna" }}`, "Korotchaeva Anna"},
		{`{{ "She" + " " + "is" + " " + "nice" }}`, "She is nice"},
		{"{{ `` }}", ""},
	}

	for _, tt := range tests {
		evaluationExpected(t, tt.input, tt.expected)
	}
}

func TestTernaryExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`{{ true ? "Yes" : "No" }}`, "Yes"},
		{`{{ false ? "Yes" : "No" }}`, "No"},
		{`{{ nil ? "Yes" : "No" }}`, "No"},
		{`{{ 1 ? "Yes" : "No" }}`, "Yes"},
		{`{{ 0 ? "Yes" : "No" }}`, "No"},
		{`{{ "" ? "Yes" : "No" }}`, "No"},
		{`{{ !true ? "Yes" : "No" }}`, "No"},
		{`{{ !false ? "Yes" : "No" }}`, "Yes"},
		{`{{ !!true ? 1 : 0 }}`, "1"},
		{`{{ !!false ? 1 : 0 }}`, "0"},
	}

	for _, tt := range tests {
		evaluationExpected(t, tt.input, tt.expected)
	}
}

func TestIfStatement(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{`{{ if true }}Hello{{ end }}`, "Hello"},
		{`{{ if false }}Hello{{ end }}`, ""},
		{`{{ if true }}Anna{{ else if true }}Serhii{{ end }}`, "Anna"},
		{`{{ if false }}Anna{{ else if true }}Serhii{{ end }}`, "Serhii"},
		{`{{ if false }}Anna{{ else if false }}Serhii{{ else }}Great{{ end }}`, "Great"},
		{`{{ if false }}Anna{{ else if false }}Serhii{{ else if true }}Great{{ end }}`, "Great"},
		{`{{ if false }}Anna{{ else if true }}Serhii{{ else if true }}Great{{ end }}`, "Serhii"},
		{`<h2>{{ if true }}Hello{{ end }}</h2>`, "<h2>Hello</h2>"},
		{`<h2>{{ if false }}Hello{{ end }}</h2>`, "<h2></h2>"},
		{`{{if true}}Hello{{end}}`, "Hello"},
		{
			`
				{{ if true }}
					{{ if false }}
					    James
					{{ else if false }}
						John
					{{ else }}
						{{ if true }}{{ "Serhii" }}{{ end }}
					{{ end }}
				{{ else }}
					{{ if true }}Anna{{ end }}
				{{ end }}
			`,
			"Serhii",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)

		if ok {
			t.Errorf("evaluation failed: %s", errObj.Message)
		}

		result := strings.TrimSpace(evaluated.String())

		if result != tt.expect {
			t.Errorf("result is not %q, got %q", tt.expect, result)
		}
	}
}

func TestVariableDeclaration(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`{{ var age = 18 }}{{ age }}`, "18"},
		{`{{ var age = 18 }}{{ age + 2 }}`, "20"},
		{`{{ var age = 18 }}{{ age + age }}`, "36"},
		{`{{ var herName = "Anna" }}{{ herName }}`, "Anna"},
	}

	for _, tt := range tests {
		evaluationExpected(t, tt.input, tt.expected)
	}
}
