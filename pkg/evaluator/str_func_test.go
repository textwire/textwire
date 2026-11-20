package evaluator

import "testing"

func TestEvalStringFunctions(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		// len
		{"len of ascii string", `{{ "anna".len() }}`, "4"},
		{"len of empty string", `{{ "".len() }}`, "0"},
		{"len of chinese string", `{{ "中国很漂亮".len() }}`, "5"},
		{
			"len of emoji with skin tone",
			`{{ "👋🏽🌍".len() }}`,
			"2",
		}, // 👋🏽 skin tone modifier is one character
		{"len of flag emoji", `{{ "🇨🇭".len() }}`, "1"},
		{"len of family emoji", `{{ "👨‍👩‍👧".len() }}`, "1"},
		// split
		{"split on whitespace", `{{ "one two three".split() }}`, "one, two, three"},
		{"split empty string default", `{{ "".split() }}`, ""},
		{"split empty string by separator", `{{ "".split("|") }}`, ""},
		{"split on pipe", `{{ "one|two|three".split("|") }}`, "one, two, three"},
		{"split on dash", `{{ "one-two".split("-") }}`, "one, two"},
		{"split chinese on character", `{{ "我喜欢中文".split("欢") }}`, "我喜, 中文"},
		{"split on empty separator", `{{ "abc".split("") }}`, "a, b, c"},
		{"split empty string on empty separator", `{{ "".split("") }}`, ""},
		{"split when separator absent", `{{ "hello".split("x") }}`, "hello"},
		{"split on multi-char absent separator", `{{ "test".split("xyz") }}`, "test"},
		{"split on newline", `{{ "line1\nline2".split("\n") }}`, "line1, line2"},
		{"split on tab", `{{ "col1\tcol2".split("\t") }}`, "col1, col2"},
		// trim
		{"trim whitespace", `{{ " 	test		".trim() }}`, "test"},
		{"trim custom character", `{{ "ease".trim("e") }}`, "as"},
		{"trim character set", `{{ "(no war!)".trim("()") }}`, "no war!"},
		{"trim chinese and spaces", `{{ " 中国很大   ".trim("中 大") }}`, "国很"},
		{
			"trim emoji set",
			`{{ "😡 Elton 😂 Elton".trim("😡😤") }}`,
			" Elton 😂 Elton",
		},
		{"trim with empty cutset", `{{ "  test  ".trim("") }}`, "  test  "},
		{"trim trailing newline", `{{ "test\n".trim() }}`, "test"},
		{"trim mixed whitespace", `{{ "  \t\n  ".trim() }}`, ""},
		{"trim crlf", `{{ "hello\r\n".trim() }}`, "hello"},
		// trimRight
		{"trim right whitespace", `{{ " 	test		".trimRight() }}`, " 	test"},
		{"trim right character", `{{ "ease".trimRight("e") }}`, "eas"},
		{"trim right character set", `{{ "(no war!)".trimRight("()") }}`, "(no war!"},
		{"trim right chinese and spaces", `{{ " 中国很大   ".trimRight("中 大") }}`, " 中国很"},
		{"trim right newline", `{{ "nice\n".trimRight() }}`, "nice"},
		// trimLeft
		{"trim left whitespace", `{{ " 	test		".trimLeft() }}`, "test		"},
		{"trim left lowercase t does not match", `{{ "Textwire".trimLeft('t') }}`, "Textwire"},
		{"trim left capital t", `{{ "Textwire".trimLeft('T') }}`, "extwire"},
		{"trim left character", `{{ "ease".trimLeft("e") }}`, "ase"},
		{"trim left character set", `{{ "(no war!)".trimLeft("()") }}`, "no war!)"},
		{"trim left chinese and spaces", `{{ " 中国很大   ".trimLeft("中 大") }}`, "国很大   "},
		// repeat
		{"repeat three times", `{{ "a".repeat(3) }}`, "aaa"},
		{"repeat zero times", `{{ "a".repeat(0) }}`, ""},
		{"repeat once", `{{ "a".repeat(1) }}`, "a"},
		{"repeat ten times", `{{ "b".repeat(10) }}`, "bbbbbbbbbb"},
		{"repeat empty string", `{{ "".repeat(10) }}`, ""},
		{"repeat space", `{{ " ".repeat(4) }}`, "    "},
		{"repeat string with trailing space", `{{ "nice ".repeat(4) }}`, "nice nice nice nice "},
		{"repeat chinese with space", `{{ "中国 ".repeat(4) }}`, "中国 中国 中国 中国 "},
		{"repeat cyrillic with space", `{{ "просто ".repeat(2) }}`, "просто просто "},
		{"repeat emoji", `{{ '🤣'.repeat(5) }}`, "🤣🤣🤣🤣🤣"},
		{"repeat negative count", `{{ "a".repeat(-1) }}`, ""},
		{
			"repeat hundred times",
			`{{ "a".repeat(100) }}`,
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		},
		// upper
		{"upper ascii", `{{ "Hello World".upper() }}`, "HELLO WORLD"},
		{
			"upper leaves symbols",
			`{{ "upper_-1234567890!@#$%^*()=+".upper() }}`,
			"UPPER_-1234567890!@#$%^*()=+",
		},
		{"upper empty", `{{ "".upper() }}`, ""},
		{"upper chinese", `{{ "中国很大".upper() }}`, "中国很大"},
		{"upper emoji", `{{ "😡🤣😤".upper() }}`, "😡🤣😤"},
		// lower
		{"lower ascii", `{{ "Hello World".lower() }}`, "hello world"},
		{
			"lower leaves symbols",
			`{{ "LOWER_-1234567890!@#$%^*()=+".lower() }}`,
			"lower_-1234567890!@#$%^*()=+",
		},
		{"lower empty", `{{ "".lower() }}`, ""},
		{"lower chinese", `{{ "中国很大".lower() }}`, "中国很大"},
		{"lower emoji", `{{ "😡🤣😤".lower() }}`, "😡🤣😤"},
		// capitalize
		{"capitalize cyrillic", `{{ "їжак".capitalize() }}`, "Їжак"},
		{"capitalize accented latin", `{{ "élan".capitalize() }}`, "Élan"},
		{"capitalize ascii", `{{ "skirk".capitalize() }}`, "Skirk"},
		{"capitalize invalid utf8", "{{ \"\xffabc\".capitalize() }}", "\xffabc"},
		{"capitalize titlecase digraph dz", `{{ "ǆemal".capitalize() }}`, "ǅemal"},
		{"capitalize titlecase digraph lj", `{{ "ǉubav".capitalize() }}`, "ǈubav"},
		{"capitalize titlecase digraph nj", `{{ "ǌiva".capitalize() }}`, "ǋiva"},
		// reverse
		{"reverse ascii sentence", `{{ "Hello World".reverse() }}`, "dlroW olleH"},
		{
			"reverse symbols and digits",
			`{{ "reverse_-1234567890!@#$%^*()=+".reverse() }}`,
			"+=)(*^%$#@!0987654321-_esrever",
		},
		{"reverse empty", `{{ "".reverse() }}`, ""},
		{"reverse single char", `{{ "T".reverse() }}`, "T"},
		{"reverse chinese", `{{ "我爱中文".reverse() }}`, "文中爱我"},
		{"reverse emoji", `{{ "😡🤣😤".reverse() }}`, "😤🤣😡"},
		{"reverse combining accent", `{{ "éx".reverse() }}`, "xé"},
		{"reverse invalid utf8", "{{ \"a\xffb\".reverse() }}", "b\xffa"},
		{"reverse skin tone emoji", `{{ "👍🏽".reverse() }}`, "👍🏽"},
		{"reverse leading combining accent", "{{ \"\u0301a\".reverse() }}", "a\u0301"},
		{
			"reverse emoji followed by ascii",
			"{{ \"\U0001F44D\U0001F3FDx\".reverse() }}",
			"x\U0001F44D\U0001F3FD",
		},
		{"reverse flags", `{{ "🇨🇳🇵🇱".reverse() }}`, "🇵🇱🇨🇳"},
		{"reverse family emoji with text", `{{ "👨‍👩‍👧 hi".reverse() }}`, "ih 👨‍👩‍👧"},
		{"reverse rainbow flag with ascii", `{{ "🏳️‍🌈x".reverse() }}`, "x🏳️‍🌈"},
		// contains
		{"contains matching word", `{{ "Hello World".contains("World") }}`, "1"},
		{"contains is case sensitive", `{{ "Hello World".contains("world") }}`, "0"},
		{"contains chinese substring", `{{ "Hello World 你好".contains("你好") }}`, "1"},
		{"contains single chinese char", `{{ "Hello World 你好".contains("你") }}`, "1"},
		{"contains chinese with trailing space", `{{ "Hello World 你好".contains("你好 ") }}`, "0"},
		{"empty contains empty", `{{ "".contains("") }}`, "1"},
		{"non-empty contains empty", `{{ "some".contains("") }}`, "1"},
		{"empty contains non-empty", `{{ "".contains("test") }}`, "0"},
		{"empty contains empty duplicate", `{{ "".contains("") }}`, "1"},
		{"string contains empty", `{{ "test".contains("") }}`, "1"},
		// truncate
		{"truncate ascii", `{{ "Hello World".truncate(5) }}`, "Hello..."},
		{"truncate chinese shorter than limit", `{{ "谢尔盖".truncate(3) }}`, "谢尔盖"},
		{"truncate exact length", `{{ "anna".truncate(4) }}`, "anna"},
		{"truncate exact length with suffix", `{{ "anna".truncate(4, "!!!") }}`, "anna"},
		{
			"truncate ascii with suffix",
			`{{ "Hello World".truncate(5, "!!!") }}`,
			"Hello!!!",
		},
		{"truncate empty with empty suffix", `{{ "".truncate(0, "") }}`, ""},
		{"truncate digits with tilde suffix", `{{ "1234567890".truncate(4, "~") }}`, "1234~"},
		{"truncate zero length", `{{ "Hello World".truncate(0) }}`, "..."},
		{"truncate zero length with suffix", `{{ "Hello World".truncate(0, "---") }}`, "---"},
		{"truncate longer than string", `{{ "test".truncate(100) }}`, "test"},
		{"truncate empty string", `{{ "".truncate(5) }}`, ""},
		{"truncate cyrillic", `{{ "привіт".truncate(2) }}`, "пр..."},
		{"truncate combining accent", `{{ "éx".truncate(1) }}`, "é..."},
		{"truncate cyrillic without suffix", `{{ "привіт".truncate(3, "") }}`, "при"},
		{
			"truncate combining accents exact",
			"{{ \"e\u0301e\u0301\".truncate(2) }}",
			"e\u0301e\u0301",
		},
		{"truncate combining accent at zero", "{{ \"e\u0301x\".truncate(0) }}", "..."},
		{
			"truncate emoji with skin tone",
			"{{ \"\U0001F44D\U0001F3FDx\".truncate(1) }}",
			"\U0001F44D\U0001F3FD...",
		},
		{"truncate flag", `{{ "🇨🇳🇵🇱".truncate(1) }}`, "🇨🇳..."},
		{"truncate family emoji", `{{ "👨‍👩‍👧 hi".truncate(1) }}`, "👨‍👩‍👧..."},
		// decimal
		{"decimal of empty string", `{{ "".decimal() }}`, ""},
		{"decimal of zero", `{{ "0".decimal() }}`, "0.00"},
		{"decimal of hundred", `{{ "100".decimal() }}`, "100.00"},
		{"decimal of 2352", `{{ "2352".decimal() }}`, "2352.00"},
		{"decimal with underscore", `{{ "1000".decimal('_') }}`, "1000_00"},
		{
			"decimal with underscore and ten precision",
			`{{ "9000".decimal('_', 10) }}`,
			"9000_0000000000",
		},
		{"decimal with zero precision", `{{ "100".decimal('|', 0) }}`, "100"},
		{"decimal with one precision", `{{ "100".decimal('|', 1) }}`, "100|0"},
		{"decimal of non-numeric", `{{ "hello".decimal() }}`, "hello"},
		{"decimal of non-numeric with options", `{{ "nice".decimal('|', 10) }}`, "nice"},
		{"decimal of existing decimal", `{{ "12.02".decimal() }}`, "12.02"},
		{"decimal of comma number", `{{ "10,10".decimal() }}`, "10,10"},
		{"decimal of negative with comma", `{{ "-900".decimal(',') }}`, "-900,00"},
		{"decimal of small fraction", `{{ "0.001".decimal() }}`, "0.001"},
		{"decimal of million", `{{ "1000000".decimal() }}`, "1000000.00"},
		{"decimal of negative fraction", `{{ "-0.5".decimal() }}`, "-0.50"},
		{"decimal with negative precision", `{{ "1".decimal('.', -1) }}`, "1"},
		{"decimal of three decimals with comma", `{{ "1.555".decimal(",") }}`, "1,555"},
		{"decimal of one decimal with comma", `{{ "1.5".decimal(",") }}`, "1,50"},
		{
			"decimal of large integer",
			`{{ "12345678901234567890".decimal() }}`,
			"12345678901234567890.00",
		},
		{"decimal of nan", `{{ "nan".decimal() }}`, "nan"},
		{"decimal of inf", `{{ "inf".decimal() }}`, "inf"},
		{"decimal of negative inf", `{{ "-inf".decimal(',') }}`, "-inf"},
		{"decimal of exponent notation", `{{ "1e5".decimal() }}`, "1e5"},
		{"decimal of malformed number", `{{ "1.2.3".decimal(',') }}`, "1.2.3"},
		{"decimal of negative one decimal", `{{ "-1.5".decimal(',') }}`, "-1,50"},
		{
			"decimal of large number with fraction",
			`{{ "123456789012345678901.5".decimal() }}`,
			"123456789012345678901.50",
		},
		// at
		{"at without argument", `{{ "Textwire is awesome".at() }}`, "T"},
		{"at index zero", `{{ "Textwire is awesome".at(0) }}`, "T"},
		{"at index one", `{{ "Textwire is awesome".at(1) }}`, "e"},
		{"at index five", `{{ "Textwire is awesome".at(5) }}`, "i"},
		{"at index eight", `{{ "Textwire is awesome".at(8) }}`, " "},
		{"at chinese index", `{{ "我爱你".at(2) }}`, "你"},
		{"at cyrillic index", `{{ "привет".at(2) }}`, "и"},
		{"at empty string index zero", `{{ "".at(0) }}`, ""},
		{"at out of bounds empty", `{{ "".at(99) }}`, ""},
		{"at negative last char", `{{ "cho".at(-1) }}`, "o"},
		{"at negative last char longer string", `{{ "Hello World".at(-1) }}`, "d"},
		{"at negative first char", `{{ "cho".at(-3) }}`, "c"},
		{"at negative chinese index", `{{ "我爱中国".at(-2) }}`, "中"},
		{"at negative out of bounds", `{{ "cho".at(-4) }}`, ""},
		{"at large negative out of bounds", `{{ "abc".at(-100) }}`, ""},
		{"at skin tone emoji", `{{ "👍🏽".at(0) }}`, "👍🏽"},
		{"at flag emoji", `{{ "🇨🇭".at(0) }}`, "🇨🇭"},
		{"at second flag emoji", `{{ "🇨🇭🇵🇱".at(1) }}`, "🇵🇱"},
		{"at last flag emoji", `{{ "🇨🇭🇵🇱".at(-1) }}`, "🇵🇱"},
		{"at family emoji", `{{ "👨‍👩‍👧".at(0) }}`, "👨‍👩‍👧"},
		// first
		{"first ascii char", `{{ "Textwire is awesome".first() }}`, "T"},
		{"first chinese char", `{{ "我爱你".first() }}`, "我"},
		{"first cyrillic char", `{{ "привет".first() }}`, "п"},
		{"first from empty string", `{{ "".first() }}`, ""},
		{"first skin tone emoji", `{{ "👍🏽".first() }}`, "👍🏽"},
		{"first flag emoji", `{{ "🇨🇭".first() }}`, "🇨🇭"},
		{"first family emoji", `{{ "👨‍👩‍👧".first() }}`, "👨‍👩‍👧"},
		// last
		{"last ascii char", `{{ "Textwire is awesome".last() }}`, "e"},
		{"last chinese char", `{{ "我爱你".last() }}`, "你"},
		{"last cyrillic char", `{{ "привет".last() }}`, "т"},
		{"last from empty string", `{{ "".last() }}`, ""},
		{"last skin tone emoji", `{{ "👍🏽".last() }}`, "👍🏽"},
		{"last flag emoji", `{{ "🇨🇭".last() }}`, "🇨🇭"},
		{"last family emoji", `{{ "👨‍👩‍👧".last() }}`, "👨‍👩‍👧"},
		// format
		{"format single int placeholder", `{{ "He has %s apples".format(2) }}`, "He has 2 apples"},
		{
			"format two string placeholders",
			`{{ "First: %s. Last: %s".format('Amy', 'Adams') }}`,
			"First: Amy. Last: Adams",
		},
		{
			"format mixed types",
			`{{ "%s-%s-%s".format(0.1, false, true) }}`,
			"0.1-0-1",
		},
		{"format array argument", `{{ "[%s]".format([1, 2]) }}`, "[1, 2]"},
		{"format leaves unsupported verb", `{{ "%s-%d".format("nice") }}`, "nice-%d"},
		{"format escaped percent before s", `{{ "%%s".format("Sydney") }}`, "%Sydney"},
		{
			"format missing argument keeps placeholder",
			`{{ "|%s and %s|".format("Anna") }}`,
			"|Anna and %s|",
		},
		{"format empty template", `{{ "".format("ignored") }}`, ""},
		{"format without placeholders", `{{ "Hello World".format("extra") }}`, "Hello World"},
		{
			"format ignores extra arguments",
			`{{ "Only one: %s".format("first", "second", "third") }}`,
			"Only one: first",
		},
		{"format three placeholders", `{{ "%s%s%s".format("a", "b", "c") }}`, "abc"},
		{
			"format double percent stays literal",
			`{{ "50%% complete".format("ignored") }}`,
			"50%% complete",
		},
		{
			"format placeholder as argument",
			`{{ "This %s is nice".format('%s') }}`,
			"This %s is nice",
		},
		{"format positive int", `{{ "%s".format(42) }}`, "42"},
		{"format negative int", `{{ "%s".format(-99) }}`, "-99"},
		{"format true", `{{ "Value: %s".format(true) }}`, "Value: 1"},
		{"format false", `{{ "Value: %s".format(false) }}`, "Value: 0"},
		{"format empty string", `{{ "Empty: [%s]".format("") }}`, "Empty: []"},
		{"format chinese", `{{ "你好%s".format("世界") }}`, "你好世界"},
		{"format trailing placeholder", `{{ "End: %s".format("here") }}`, "End: here"},
		{
			"format long string",
			`{{ "%s".format("a very long string with many characters") }}`,
			"a very long string with many characters",
		},
		{"format unsupported float verb", `{{ "%.4f".format(3.14159) }}`, "%.4f"},
		{
			"format escaped percents around placeholder",
			`{{ "%%%s%%".format("middle") }}`,
			"%%middle%%",
		},
		{"format tab separated", `{{ "%s	%s".format("a", "b") }}`, "a	b"},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, nil, tc.name)
	}
}

func TestString_MethodChaining(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		{
			"lower then contains",
			`{{ "Hello, World!".lower().contains("world") }}`,
			"1",
		},
		{"lower then upper", `{{ "HELLO".lower().upper() }}`, "HELLO"},
		{"upper lower upper", `{{ "test".upper().lower().upper() }}`, "TEST"},
		{"reverse then upper", `{{ "anna".reverse().upper() }}`, "ANNA"},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, nil, tc.name)
	}
}
