package evaluator

import "testing"

func TestEvalStringFunctions(t *testing.T) {
	tests := []struct {
		inp      string
		expected string
	}{
		{`{{ "anna".len() }}`, "4"},
		{`{{ "".len() }}`, "0"},
		{`{{ "one two three".split() }}`, "one, two, three"},
		{`{{ "one|two|three".split("|") }}`, "one, two, three"},
		{`{{ "one-two".split("-") }}`, "one, two"},
		{`{{ "<h1>nice</h1>".raw()`, "<h1>nice</h1>"},
		{`{{ "cool".raw()`, "cool"},
		{`{{ " 	test		".trim()`, "test"},
		{`{{ "ease".trim("e")`, "as"},
		{`{{ "(no war!)".trim("()")`, "no war!"},
	}

	for _, tt := range tests {
		evaluationExpected(t, tt.inp, tt.expected)
	}
}
