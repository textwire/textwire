package evaluator

import "testing"

func TestEvalStringFunctions(t *testing.T) {
	cases := []struct {
		id     uint
		inp    string
		expect string
	}{
		// len
		{10, `{{ "anna".len() }}`, "4"},
		{20, `{{ "".len() }}`, "0"},
		{30, `{{ "中国很漂亮".len() }}`, "5"},
		{40, `{{ "👋🏽🌍".len() }}`, "3"}, // 👋 + 🏽 skin tone modifier give length 2
		// split
		{50, `{{ "one two three".split() }}`, "one, two, three"},
		{60, `{{ "".split() }}`, ""},
		{70, `{{ "".split("|") }}`, ""},
		{80, `{{ "one|two|three".split("|") }}`, "one, two, three"},
		{90, `{{ "one-two".split("-") }}`, "one, two"},
		{100, `{{ "我喜欢中文".split("欢") }}`, "我喜, 中文"},
		{110, `{{ "abc".split("") }}`, "a, b, c"},
		{120, `{{ "".split("") }}`, ""},
		{130, `{{ "hello".split("x") }}`, "hello"},
		{140, `{{ "test".split("xyz") }}`, "test"},
		{150, `{{ "line1\nline2".split("\n") }}`, "line1, line2"},
		{160, `{{ "col1\tcol2".split("\t") }}`, "col1, col2"},
		// raw
		{170, `{{ "<h1>nice</h1>" }}`, "&lt;h1&gt;nice&lt;/h1&gt;"},
		{180, `{{ "\"\"" }}`, "&#34;&#34;"},
		{190, `{{ "<h1>nice</h1>".raw() }}`, "<h1>nice</h1>"},
		{200, `{{ "cool".raw() }}`, "cool"},
		{210, `{{ "<b>中国很大</b>".raw() }}`, "<b>中国很大</b>"},
		// trim
		{220, `{{ " 	test		".trim() }}`, "test"},
		{230, `{{ "ease".trim("e") }}`, "as"},
		{240, `{{ "(no war!)".trim("()") }}`, "no war!"},
		{250, `{{ " 中国很大   ".trim("中 大") }}`, "国很"},
		{260, `{{ "😡 Elton 😂 Elton".trim("😡😤") }}`, " Elton 😂 Elton"},
		{270, `{{ "  test  ".trim("") }}`, "  test  "},
		{280, `{{ "test\n".trim() }}`, "test"},
		{290, `{{ "  \t\n  ".trim() }}`, ""},
		{300, `{{ "hello\r\n".trim() }}`, "hello"},
		// trimRight
		{310, `{{ " 	test		".trimRight() }}`, " 	test"},
		{320, `{{ "ease".trimRight("e") }}`, "eas"},
		{330, `{{ "(no war!)".trimRight("()") }}`, "(no war!"},
		{340, `{{ " 中国很大   ".trimRight("中 大") }}`, " 中国很"},
		{350, `{{ "nice\n".trimRight() }}`, "nice"},
		// trimLeft
		{360, `{{ " 	test		".trimLeft() }}`, "test		"},
		{370, `{{ "Textwire".trimLeft('t') }}`, "Textwire"},
		{380, `{{ "Textwire".trimLeft('T') }}`, "extwire"},
		{390, `{{ "ease".trimLeft("e") }}`, "ase"},
		{400, `{{ "(no war!)".trimLeft("()") }}`, "no war!)"},
		{410, `{{ " 中国很大   ".trimLeft("中 大") }}`, "国很大   "},
		// repeat
		{420, `{{ "a".repeat(3) }}`, "aaa"},
		{430, `{{ "a".repeat(0) }}`, ""},
		{440, `{{ "a".repeat(1) }}`, "a"},
		{450, `{{ "b".repeat(10) }}`, "bbbbbbbbbb"},
		{460, `{{ "".repeat(10) }}`, ""},
		{470, `{{ " ".repeat(4) }}`, "    "},
		{480, `{{ "nice ".repeat(4) }}`, "nice nice nice nice "},
		{490, `{{ "中国 ".repeat(4) }}`, "中国 中国 中国 中国 "},
		{500, `{{ "просто ".repeat(2) }}`, "просто просто "},
		{510, `{{ '🤣'.repeat(5) }}`, "🤣🤣🤣🤣🤣"},
		{520, `{{ "a".repeat(-1) }}`, ""},
		{
			530,
			`{{ "a".repeat(100) }}`,
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		},
		// upper
		{540, `{{ "Hello World".upper() }}`, "HELLO WORLD"},
		{550, `{{ "upper_-1234567890!@#$%^*()=+".upper() }}`, "UPPER_-1234567890!@#$%^*()=+"},
		{560, `{{ "".upper() }}`, ""},
		{570, `{{ "中国很大".upper() }}`, "中国很大"},
		{580, `{{ "😡🤣😤".upper() }}`, "😡🤣😤"},
		// lower
		{590, `{{ "Hello World".lower() }}`, "hello world"},
		{600, `{{ "LOWER_-1234567890!@#$%^*()=+".lower() }}`, "lower_-1234567890!@#$%^*()=+"},
		{610, `{{ "".lower() }}`, ""},
		{620, `{{ "中国很大".lower() }}`, "中国很大"},
		{630, `{{ "😡🤣😤".lower() }}`, "😡🤣😤"},
		// reverse
		{640, `{{ "Hello World".reverse() }}`, "dlroW olleH"},
		{650, `{{ "reverse_-1234567890!@#$%^*()=+".reverse() }}`, "+=)(*^%$#@!0987654321-_esrever"},
		{660, `{{ "".reverse() }}`, ""},
		{670, `{{ "T".reverse() }}`, "T"},
		{680, `{{ "我爱中文".reverse() }}`, "文中爱我"},
		{690, `{{ "😡🤣😤".reverse() }}`, "😤🤣😡"},
		// contains
		{700, `{{ "Hello World".contains("World") }}`, "1"},
		{710, `{{ "Hello World".contains("world") }}`, "0"},
		{720, `{{ "Hello World 你好".contains("你好") }}`, "1"},
		{730, `{{ "Hello World 你好".contains("你") }}`, "1"},
		{740, `{{ "Hello World 你好".contains("你好 ") }}`, "0"},
		{750, `{{ "".contains("") }}`, "1"},
		{760, `{{ "some".contains("") }}`, "1"},
		{770, `{{ "".contains("test") }}`, "0"},
		{780, `{{ "".contains("") }}`, "1"},
		{790, `{{ "test".contains("") }}`, "1"},
		// truncate
		{800, `{{ "Hello World".truncate(5) }}`, "Hello..."},
		{810, `{{ "谢尔盖".truncate(3) }}`, "谢尔盖"},
		{820, `{{ "anna".truncate(4) }}`, "anna"},
		{830, `{{ "anna".truncate(4, "!!!") }}`, "anna"},
		{840, `{{ "Hello World".truncate(5, "!!!") }}`, "Hello!!!"},
		{850, `{{ "".truncate(0, "") }}`, ""},
		{860, `{{ "1234567890".truncate(4, "~") }}`, "1234~"},
		{870, `{{ "Hello World".truncate(0) }}`, "..."},
		{880, `{{ "Hello World".truncate(0, "---") }}`, "---"},
		{890, `{{ "test".truncate(-1) }}`, "..."},
		{900, `{{ "test".truncate(100) }}`, "test"},
		{910, `{{ "".truncate(5) }}`, ""},
		// decimal
		{920, `{{ "".decimal() }}`, ""},
		{930, `{{ "0".decimal() }}`, "0.00"},
		{940, `{{ "100".decimal() }}`, "100.00"},
		{950, `{{ "2352".decimal() }}`, "2352.00"},
		{960, `{{ "1000".decimal('_') }}`, "1000_00"},
		{970, `{{ "9000".decimal('_', 10) }}`, "9000_0000000000"},
		{980, `{{ "100".decimal('|', 0) }}`, "100"},
		{990, `{{ "100".decimal('|', 1) }}`, "100|0"},
		{1000, `{{ "hello".decimal() }}`, "hello"},
		{1010, `{{ "nice".decimal('|', 10) }}`, "nice"},
		{1020, `{{ "12.02".decimal() }}`, "12.02"},
		{1030, `{{ "10,10".decimal() }}`, "10,10"},
		{1040, `{{ "-900".decimal(',') }}`, "-900,00"},
		{1050, `{{ "0.001".decimal() }}`, "0.001"},
		{1060, `{{ "1000000".decimal() }}`, "1000000.00"},
		{1070, `{{ "-0.5".decimal() }}`, "-0.50"},
		// at
		{1080, `{{ "Textwire is awesome".at() }}`, "T"},
		{1090, `{{ "Textwire is awesome".at(0) }}`, "T"},
		{1100, `{{ "Textwire is awesome".at(1) }}`, "e"},
		{1110, `{{ "Textwire is awesome".at(5) }}`, "i"},
		{1120, `{{ "Textwire is awesome".at(8) }}`, " "},
		{1130, `{{ "我爱你".at(2) }}`, "你"},
		{1140, `{{ "привет".at(2) }}`, "и"},
		{1150, `{{ "".at(0) }}`, ""},
		{1160, `{{ "".at(99) }}`, ""},
		{1170, `{{ "cho".at(-1) }}`, "o"},
		{1180, `{{ "Hello World".at(-1) }}`, "d"},
		{1190, `{{ "cho".at(-3) }}`, "c"},
		{1200, `{{ "我爱中国".at(-2) }}`, "中"},
		// first
		{1210, `{{ "Textwire is awesome".first() }}`, "T"},
		{1220, `{{ "我爱你".first() }}`, "我"},
		{1230, `{{ "привет".first() }}`, "п"},
		{1240, `{{ "".first() }}`, ""},
		// last
		{1250, `{{ "Textwire is awesome".last() }}`, "e"},
		{1260, `{{ "我爱你".last() }}`, "你"},
		{1270, `{{ "привет".last() }}`, "т"},
		{1280, `{{ "".last() }}`, ""},
		// format
		{1290, `{{ "He has %s apples".format(2) }}`, "He has 2 apples"},
		{1300, `{{ "First: %s. Last: %s".format('Amy', 'Adams') }}`, "First: Amy. Last: Adams"},
		{1310, `{{ "%s-%s-%s".format(0.1, false, true) }}`, "0.1-0-1"},
		{1320, `{{ "[%s]".format([1, 2]) }}`, "[1, 2]"},
		{1330, `{{ "%s-%d".format("nice") }}`, "nice-%d"},
		{1340, `{{ "%%s".format("Sydney") }}`, "%Sydney"},
		{1350, `{{ "|%s and %s|".format("Anna") }}`, "|Anna and %s|"},
		{1360, `{{ "".format("ignored") }}`, ""},
		{1370, `{{ "Hello World".format("extra") }}`, "Hello World"},
		{1380, `{{ "Only one: %s".format("first", "second", "third") }}`, "Only one: first"},
		{1390, `{{ "%s%s%s".format("a", "b", "c") }}`, "abc"},
		{1400, `{{ "50%% complete".format("ignored") }}`, "50%% complete"},
		{1410, `{{ "This %s is nice".format('%s') }}`, "This %s is nice"},
		{1420, `{{ "%s".format(42) }}`, "42"},
		{1430, `{{ "%s".format(-99) }}`, "-99"},
		{1440, `{{ "Value: %s".format(true) }}`, "Value: 1"},
		{1450, `{{ "Value: %s".format(false) }}`, "Value: 0"},
		{1460, `{{ "Empty: [%s]".format("") }}`, "Empty: []"},
		{1470, `{{ "你好%s".format("世界") }}`, "你好世界"},
		{1480, `{{ "End: %s".format("here") }}`, "End: here"},
		{
			1490,
			`{{ "%s".format("a very long string with many characters") }}`,
			"a very long string with many characters",
		},
		{1500, `{{ "%.4f".format(3.14159) }}`, "%.4f"},
		{1510, `{{ "%%%s%%".format("middle") }}`, "%%middle%%"},
		{1520, `{{ "%s	%s".format("a", "b") }}`, "a	b"},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, tc.id)
	}
}

func TestStringMethodChaining(t *testing.T) {
	cases := []struct {
		id     uint
		inp    string
		expect string
	}{
		{1530, `{{ "Hello, World!".lower().contains("world") }}`, "1"},
		{1540, `{{ "HELLO".lower().upper() }}`, "HELLO"},
		{1550, `{{ "test".upper().lower().upper() }}`, "TEST"},
		{1560, `{{ "anna".reverse().upper() }}`, "ANNA"},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, tc.id)
	}
}
