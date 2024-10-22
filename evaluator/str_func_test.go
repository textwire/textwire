package evaluator

import "testing"

func TestEvalStringFunctions(t *testing.T) {
	tests := []struct {
		inp      string
		expected string
	}{
		// len
		{`{{ "anna".len() }}`, "4"},
		{`{{ "".len() }}`, "0"},
		// split
		{`{{ "one two three".split() }}`, "one, two, three"},
		{`{{ "one|two|three".split("|") }}`, "one, two, three"},
		{`{{ "one-two".split("-") }}`, "one, two"},
		// raw
		{`{{ "<h1>nice</h1>".raw() }}`, "<h1>nice</h1>"},
		{`{{ "cool".raw() }}`, "cool"},
		// trim
		{`{{ " 	test		".trim() }}`, "test"},
		{`{{ "ease".trim("e") }}`, "as"},
		{`{{ "(no war!)".trim("()") }}`, "no war!"},
		// upper
		{`{{ "Hello World 你好".upper() }}`, "HELLO WORLD 你好"},
		{`{{ "upper_-1234567890!@#$%^*()=+".upper() }}`, "UPPER_-1234567890!@#$%^*()=+"},
		{`{{ "".upper() }}`, ""},
		// lower
		{`{{ "Hello World 你好".lower() }}`, "hello world 你好"},
		{`{{ "LOWER_-1234567890!@#$%^*()=+".lower() }}`, "lower_-1234567890!@#$%^*()=+"},
		{`{{ "".lower() }}`, ""},
		// reverse
		{`{{ "Hello World 你好".reverse() }}`, "好你 dlroW olleH"},
		{`{{ "reverse_-1234567890!@#$%^*()=+".reverse() }}`, "+=)(*^%$#@!0987654321-_esrever"},
		{`{{ "".reverse() }}`, ""},
		{`{{ "T".reverse() }}`, "T"},
	}

	for _, tc := range tests {
		evaluationExpected(t, tc.inp, tc.expected)
	}
}
