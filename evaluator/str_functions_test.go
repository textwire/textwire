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
		{`{{ "<h1>nice</h1>".raw() }}`, "<h1>nice</h1>"},
		{`{{ "cool".raw() }}`, "cool"},
		{`{{ " 	test		".trim() }}`, "test"},
		{`{{ "ease".trim("e") }}`, "as"},
		{`{{ "(no war!)".trim("()") }}`, "no war!"},
		{`{{ "Hello World 你好".upper() }}`, "HELLO WORLD 你好"},
		{`{{ "Hello World 你好".lower() }}`, "hello world 你好"},
	}

	for _, tc := range tests {
		evaluationExpected(t, tc.inp, tc.expected)
	}
}
