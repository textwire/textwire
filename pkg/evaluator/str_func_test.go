package evaluator

import "testing"

func TestEvalStringFunctions(t *testing.T) {
	cases := []struct {
		id     int
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
		{60, `{{ "one|two|three".split("|") }}`, "one, two, three"},
		{70, `{{ "one-two".split("-") }}`, "one, two"},
		{80, `{{ "我喜欢中文".split("欢") }}`, "我喜, 中文"},
		// raw
		{90, `{{ "<h1>nice</h1>".raw() }}`, "<h1>nice</h1>"},
		{100, `{{ "cool".raw() }}`, "cool"},
		{110, `{{ "<b>中国很大</b>".raw() }}`, "<b>中国很大</b>"},
		// trim
		{120, `{{ " 	test		".trim() }}`, "test"},
		{130, `{{ "ease".trim("e") }}`, "as"},
		{140, `{{ "(no war!)".trim("()") }}`, "no war!"},
		{150, `{{ " 中国很大   ".trim("中 大") }}`, "国很"},
		{160, `{{ "😡 Elton 😂 Elton".trim("😡😤") }}`, " Elton 😂 Elton"},
		// trimRight
		{170, `{{ " 	test		".trimRight() }}`, " 	test"},
		{180, `{{ "ease".trimRight("e") }}`, "eas"},
		{190, `{{ "(no war!)".trimRight("()") }}`, "(no war!"},
		{200, `{{ " 中国很大   ".trimRight("中 大") }}`, " 中国很"},
		// trimLeft
		{210, `{{ " 	test		".trimLeft() }}`, "test		"},
		{220, `{{ "Textwire".trimLeft('t') }}`, "Textwire"},
		{230, `{{ "Textwire".trimLeft('T') }}`, "extwire"},
		{240, `{{ "ease".trimLeft("e") }}`, "ase"},
		{250, `{{ "(no war!)".trimLeft("()") }}`, "no war!)"},
		{260, `{{ " 中国很大   ".trimLeft("中 大") }}`, "国很大   "},
		// repeat
		{270, `{{ "a".repeat(3) }}`, "aaa"},
		{280, `{{ "a".repeat(0) }}`, ""},
		{290, `{{ "a".repeat(1) }}`, "a"},
		{300, `{{ "b".repeat(10) }}`, "bbbbbbbbbb"},
		{310, `{{ "".repeat(10) }}`, ""},
		{320, `{{ " ".repeat(4) }}`, "    "},
		{330, `{{ "nice ".repeat(4) }}`, "nice nice nice nice "},
		{340, `{{ "中国 ".repeat(4) }}`, "中国 中国 中国 中国 "},
		{350, `{{ "просто ".repeat(2) }}`, "просто просто "},
		{360, `{{ '🤣'.repeat(5) }}`, "🤣🤣🤣🤣🤣"},
		// upper
		{370, `{{ "Hello World".upper() }}`, "HELLO WORLD"},
		{380, `{{ "upper_-1234567890!@#$%^*()=+".upper() }}`, "UPPER_-1234567890!@#$%^*()=+"},
		{390, `{{ "".upper() }}`, ""},
		{400, `{{ "中国很大".upper() }}`, "中国很大"},
		{410, `{{ "😡🤣😤".upper() }}`, "😡🤣😤"},
		// lower
		{420, `{{ "Hello World".lower() }}`, "hello world"},
		{430, `{{ "LOWER_-1234567890!@#$%^*()=+".lower() }}`, "lower_-1234567890!@#$%^*()=+"},
		{440, `{{ "".lower() }}`, ""},
		{450, `{{ "中国很大".lower() }}`, "中国很大"},
		{460, `{{ "😡🤣😤".lower() }}`, "😡🤣😤"},
		// reverse
		{470, `{{ "Hello World".reverse() }}`, "dlroW olleH"},
		{480, `{{ "reverse_-1234567890!@#$%^*()=+".reverse() }}`, "+=)(*^%$#@!0987654321-_esrever"},
		{490, `{{ "".reverse() }}`, ""},
		{500, `{{ "T".reverse() }}`, "T"},
		{510, `{{ "我爱中文".reverse() }}`, "文中爱我"},
		{520, `{{ "😡🤣😤".reverse() }}`, "😤🤣😡"},
		// contains
		{530, `{{ "Hello World".contains("World") }}`, "1"},
		{540, `{{ "Hello World".contains("world") }}`, "0"},
		{550, `{{ "Hello World 你好".contains("你好") }}`, "1"},
		{560, `{{ "Hello World 你好".contains("你") }}`, "1"},
		{570, `{{ "Hello World 你好".contains("你好 ") }}`, "0"},
		{580, `{{ "".contains("") }}`, "1"},
		{590, `{{ "some".contains("") }}`, "1"},
		{600, `{{ "Hello, World!".lower().contains("world") }}`, "1"},
		{610, `{{ !"aaa".contains("a") }}`, "0"},
		{620, `{{ !"aaa".contains("b") }}`, "1"},
		// truncate
		{630, `{{ "Hello World".truncate(5) }}`, "Hello..."},
		{640, `{{ "谢尔盖".truncate(3) }}`, "谢尔盖"},
		{650, `{{ "anna".truncate(4) }}`, "anna"},
		{660, `{{ "anna".truncate(4, "!!!") }}`, "anna"},
		{670, `{{ "Hello World".truncate(5, "!!!") }}`, "Hello!!!"},
		{680, `{{ "".truncate(0, "") }}`, ""},
		{690, `{{ "1234567890".truncate(4, "~") }}`, "1234~"},
		{700, `{{ "Hello World".truncate(0) }}`, "..."},
		{710, `{{ "Hello World".truncate(0, "---") }}`, "---"},
		// decimal
		{720, `{{ "".decimal() }}`, ""},
		{730, `{{ "0".decimal() }}`, "0.00"},
		{740, `{{ "100".decimal() }}`, "100.00"},
		{750, `{{ "2352".decimal() }}`, "2352.00"},
		{760, `{{ "1000".decimal('_') }}`, "1000_00"},
		{770, `{{ "9000".decimal('_', 10) }}`, "9000_0000000000"},
		{780, `{{ "100".decimal('|', 0) }}`, "100"},
		{790, `{{ "100".decimal('|', 1) }}`, "100|0"},
		{800, `{{ "hello".decimal() }}`, "hello"},
		{810, `{{ "nice".decimal('|', 10) }}`, "nice"},
		{820, `{{ "12.02".decimal() }}`, "12.02"},
		{830, `{{ "10,10".decimal() }}`, "10,10"},
		{840, `{{ "-900".decimal(',') }}`, "-900,00"},
		// at
		{850, `{{ "Textwire is awesome".at() }}`, "T"},
		{860, `{{ "Textwire is awesome".at(0) }}`, "T"},
		{870, `{{ "Textwire is awesome".at(1) }}`, "e"},
		{880, `{{ "Textwire is awesome".at(5) }}`, "i"},
		{890, `{{ "Textwire is awesome".at(8) }}`, " "},
		{900, `{{ "我爱你".at(2) }}`, "你"},
		{910, `{{ "привет".at(2) }}`, "и"},
		{920, `{{ "".at(0) }}`, ""},
		{930, `{{ "".at(99) }}`, ""},
		{940, `{{ "cho".at(-1) }}`, "o"},
		{950, `{{ "Hello World".at(-1) }}`, "d"},
		{960, `{{ "cho".at(-3) }}`, "c"},
		{970, `{{ "我爱中国".at(-2) }}`, "中"},
		// first
		{980, `{{ "Textwire is awesome".first() }}`, "T"},
		{990, `{{ "我爱你".first() }}`, "我"},
		{1000, `{{ "привет".first() }}`, "п"},
		{1010, `{{ "".first() }}`, ""},
		// last
		{1020, `{{ "Textwire is awesome".last() }}`, "e"},
		{1030, `{{ "我爱你".last() }}`, "你"},
		{1040, `{{ "привет".last() }}`, "т"},
		{1050, `{{ "".last() }}`, ""},
		// format
		{1060, `{{ "He has %s apples".format(2) }}`, "He has 2 apples"},
		{1070, `{{ "First: %s. Last: %s".format('Amy', 'Adams') }}`, "First: Amy. Last: Adams"},
		{1080, `{{ "%s-%s-%s".format(0.1, false, true) }}`, "0.1-0-1"},
		{1090, `{{ "[%s]".format([1, 2]) }}`, "[1, 2]"},
		{1100, `{{ "%s-%d".format("nice") }}`, "nice-%d"},
		{1110, `{{ "%%s".format("Sydney") }}`, "%Sydney"},
		{1120, `{{ "|%s and %s|".format("Anna") }}`, "|Anna and %s|"},
		{1130, `{{ "".format("ignored") }}`, ""},
		{1140, `{{ "Hello World".format("extra") }}`, "Hello World"},
		{1150, `{{ "Only one: %s".format("first", "second", "third") }}`, "Only one: first"},
		{1160, `{{ "%s%s%s".format("a", "b", "c") }}`, "abc"},
		{1170, `{{ "50%% complete".format("ignored") }}`, "50%% complete"},
		{1180, `{{ "This %s is nice".format('%s') }}`, "This %s is nice"},
		{1190, `{{ "%s".format(42) }}`, "42"},
		{1200, `{{ "%s".format(-99) }}`, "-99"},
		{1210, `{{ "Value: %s".format(true) }}`, "Value: 1"},
		{1220, `{{ "Value: %s".format(false) }}`, "Value: 0"},
		{1230, `{{ "Empty: [%s]".format("") }}`, "Empty: []"},
		{1240, `{{ "你好%s".format("世界") }}`, "你好世界"},
		{1250, `{{ "End: %s".format("here") }}`, "End: here"},
		{
			1260,
			`{{ "%s".format("a very long string with many characters") }}`,
			"a very long string with many characters",
		},
		{1270, `{{ "%.4f".format(3.14159) }}`, "%.4f"},
		{1280, `{{ "%%%s%%".format("middle") }}`, "%%middle%%"},
		{1290, `{{ "%s	%s".format("a", "b") }}`, "a	b"},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, tc.id)
	}
}
