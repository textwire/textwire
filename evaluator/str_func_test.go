package evaluator

import "testing"

func TestEvalStringFunctions(t *testing.T) {
	cases := []struct {
		inp    string
		expect string
	}{
		// len
		{`{{ "anna".len() }}`, "4"},
		{`{{ "".len() }}`, "0"},
		{`{{ "中国很漂亮".len() }}`, "5"},
		// split
		{`{{ "one two three".split() }}`, "one, two, three"},
		{`{{ "one|two|three".split("|") }}`, "one, two, three"},
		{`{{ "one-two".split("-") }}`, "one, two"},
		{`{{ "我喜欢中文".split("欢") }}`, "我喜, 中文"},
		// raw
		{`{{ "<h1>nice</h1>".raw() }}`, "<h1>nice</h1>"},
		{`{{ "cool".raw() }}`, "cool"},
		{`{{ "<b>中国很大</b>".raw() }}`, "<b>中国很大</b>"},
		// trim
		{`{{ " 	test		".trim() }}`, "test"},
		{`{{ "ease".trim("e") }}`, "as"},
		{`{{ "(no war!)".trim("()") }}`, "no war!"},
		{`{{ " 中国很大   ".trim("中 大") }}`, "国很"},
		{`{{ "😡🤣🤣🤣😤".trim("😡😤") }}`, "🤣🤣🤣"},
		// trimRight
		{`{{ " 	test		".trimRight() }}`, " 	test"},
		{`{{ "ease".trimRight("e") }}`, "eas"},
		{`{{ "(no war!)".trimRight("()") }}`, "(no war!"},
		{`{{ " 中国很大   ".trimRight("中 大") }}`, " 中国很"},
		// trimLeft
		{`{{ " 	test		".trimLeft() }}`, "test		"},
		{`{{ "Textwire".trimLeft('t') }}`, "Textwire"},
		{`{{ "Textwire".trimLeft('T') }}`, "extwire"},
		{`{{ "ease".trimLeft("e") }}`, "ase"},
		{`{{ "(no war!)".trimLeft("()") }}`, "no war!)"},
		{`{{ " 中国很大   ".trimLeft("中 大") }}`, "国很大   "},
		// repeat
		{`{{ "a".repeat(3) }}`, "aaa"},
		{`{{ "a".repeat(0) }}`, ""},
		{`{{ "a".repeat(1) }}`, "a"},
		{`{{ "b".repeat(10) }}`, "bbbbbbbbbb"},
		{`{{ "".repeat(10) }}`, ""},
		{`{{ " ".repeat(4) }}`, "    "},
		{`{{ "nice ".repeat(4) }}`, "nice nice nice nice "},
		{`{{ "中国 ".repeat(4) }}`, "中国 中国 中国 中国 "},
		{`{{ "просто ".repeat(2) }}`, "просто просто "},
		{`{{ '🤣'.repeat(5) }}`, "🤣🤣🤣🤣🤣"},
		// upper
		{`{{ "Hello World".upper() }}`, "HELLO WORLD"},
		{`{{ "upper_-1234567890!@#$%^*()=+".upper() }}`, "UPPER_-1234567890!@#$%^*()=+"},
		{`{{ "".upper() }}`, ""},
		{`{{ "中国很大".upper() }}`, "中国很大"},
		{`{{ "😡🤣😤".upper() }}`, "😡🤣😤"},
		// lower
		{`{{ "Hello World".lower() }}`, "hello world"},
		{`{{ "LOWER_-1234567890!@#$%^*()=+".lower() }}`, "lower_-1234567890!@#$%^*()=+"},
		{`{{ "".lower() }}`, ""},
		{`{{ "中国很大".lower() }}`, "中国很大"},
		{`{{ "😡🤣😤".lower() }}`, "😡🤣😤"},
		// reverse
		{`{{ "Hello World".reverse() }}`, "dlroW olleH"},
		{`{{ "reverse_-1234567890!@#$%^*()=+".reverse() }}`, "+=)(*^%$#@!0987654321-_esrever"},
		{`{{ "".reverse() }}`, ""},
		{`{{ "T".reverse() }}`, "T"},
		{`{{ "我爱中文".reverse() }}`, "文中爱我"},
		{`{{ "😡🤣😤".reverse() }}`, "😤🤣😡"},
		// contains
		{`{{ "Hello World".contains("World") }}`, "1"},
		{`{{ "Hello World".contains("world") }}`, "0"},
		{`{{ "Hello World 你好".contains("你好") }}`, "1"},
		{`{{ "Hello World 你好".contains("你") }}`, "1"},
		{`{{ "Hello World 你好".contains("你好 ") }}`, "0"},
		{`{{ "".contains("") }}`, "1"},
		{`{{ "some".contains("") }}`, "1"},
		{`{{ "Hello, World!".lower().contains("world") }}`, "1"},
		{`{{ !"aaa".contains("a") }}`, "0"},
		{`{{ !"aaa".contains("b") }}`, "1"},
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
		{`{{ "-900".decimal(',') }}`, "-900,00"},
		// at
		{`{{ "Textwire is awesome".at() }}`, "T"},
		{`{{ "Textwire is awesome".at(0) }}`, "T"},
		{`{{ "Textwire is awesome".at(1) }}`, "e"},
		{`{{ "Textwire is awesome".at(5) }}`, "i"},
		{`{{ "Textwire is awesome".at(8) }}`, " "},
		{`{{ "我爱你".at(2) }}`, "你"},
		{`{{ "привет".at(2) }}`, "и"},
		{`{{ "".at(0) }}`, ""},
		{`{{ "".at(99) }}`, ""},
		{`{{ "cho".at(-1) }}`, "o"},
		{`{{ "Hello World".at(-1) }}`, "d"},
		{`{{ "cho".at(-3) }}`, "c"},
		{`{{ "我爱中国".at(-2) }}`, "中"},
		// first
		{`{{ "Textwire is awesome".first() }}`, "T"},
		{`{{ "我爱你".first() }}`, "我"},
		{`{{ "привет".first() }}`, "п"},
		{`{{ "".first() }}`, ""},
		// last
		{`{{ "Textwire is awesome".last() }}`, "e"},
		{`{{ "我爱你".last() }}`, "你"},
		{`{{ "привет".last() }}`, "т"},
		{`{{ "".last() }}`, ""},
	}

	for i, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, i)
	}
}
