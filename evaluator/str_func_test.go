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
		// contains
		{`{{ "Hello World 你好".contains("World") }}`, "1"},
		{`{{ "Hello World 你好".contains("world") }}`, "0"},
		{`{{ "Hello World 你好".contains("你好") }}`, "1"},
		{`{{ "Hello World 你好".contains("你") }}`, "1"},
		{`{{ "Hello World 你好".contains("你好 ") }}`, "0"},
		{`{{ "".contains("") }}`, "1"},
		{`{{ "some".contains("") }}`, "1"},
		// truncate
		{`{{ "Hello World".truncate(5) }}`, "Hello..."},
		{`{{ "谢尔盖".truncate(3) }}`, "谢尔盖"},
		{`{{ "anna".truncate(4) }}`, "anna"},
		{`{{ "anna".truncate(4, "!!!") }}`, "anna"},
		{`{{ "Hello World".truncate(5, "!!!") }}`, "Hello!!!"},
		{`{{ "".truncate(0, "") }}`, ""},
		{`{{ "1234567890".truncate(4, "~") }}`, "1234~"},
		{`{{ "Hello World".truncate(0) }}`, "..."},
		{`{{ "Hello World".truncate(0, "---") }}`, "---"},
		// decimal
		{`{{ "".decimal() }}`, ""},
		{`{{ "0".decimal() }}`, "0.00"},
		{`{{ "100".decimal() }}`, "100.00"},
		{`{{ "2352".decimal() }}`, "2352.00"},
		{`{{ "1000".decimal('_') }}`, "1000_00"},
		{`{{ "9000".decimal('_', 10) }}`, "9000_0000000000"},
		{`{{ "100".decimal('|', 0) }}`, "100"},
		{`{{ "100".decimal('|', 1) }}`, "100|0"},
		{`{{ "hello".decimal() }}`, "hello"},
		{`{{ "nice".decimal('|', 10) }}`, "nice"},
		{`{{ "12.02".decimal() }}`, "12.02"},
		{`{{ "10,10".decimal() }}`, "10,10"},
	}

	for _, tc := range tests {
		evaluationExpected(t, tc.inp, tc.expected)
	}
}
