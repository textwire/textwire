package evaluator

import (
	"strings"
	"testing"

	"github.com/textwire/textwire/v4/config"
	"github.com/textwire/textwire/v4/pkg/fail"
	"github.com/textwire/textwire/v4/pkg/file"
	"github.com/textwire/textwire/v4/pkg/lexer"
	"github.com/textwire/textwire/v4/pkg/parser"
	"github.com/textwire/textwire/v4/pkg/value"
)

func testEval(inp string) (value.Value, *fail.Error) {
	l := lexer.New(inp)
	p := parser.New(l, file.New("file", "to/file", "/path/to/file", nil))
	prog := p.ParseProgram()
	scope := value.NewScope()

	if p.HasErrors() {
		return nil, p.Errors()[0]
	}

	e := New(&config.Func{}, nil)
	ctx := NewContext(scope, prog.AbsPath)

	return e.Eval(prog, ctx), nil
}

func evaluationExpected(t *testing.T, inp, expect string, idx uint) {
	evaluated, failure := testEval(inp)
	if failure != nil {
		t.Fatalf("Case: %d. evaluation failed: %s", idx, failure)
	}

	errObj, ok := evaluated.(*value.Error)
	if ok {
		t.Fatalf("Case: %d. evaluation failed: %s", idx, errObj)
	}

	res := evaluated.String()
	if res != expect {
		t.Fatalf("Case: %d. Result is not '%s', got '%s'", idx, expect, res)
	}
}

func TestEvalText(t *testing.T) {
	cases := []struct {
		id     uint
		inp    string
		expect string
	}{
		{10, "<h1>Hello World</h1>", "<h1>Hello World</h1>"},
		{20, "<ul><li><span>Email: anna@protonmail.com</span></li></ul>",
			"<ul><li><span>Email: anna@protonmail.com</span></li></ul>"},
		{30, "<b>Nice</b>@foo", "<b>Nice</b>@foo"},
		{40, `<h1>\@continue</h1>`, "<h1>@continue</h1>"},
		{50, `<h1>@\@break</h1>`, "<h1>@@break</h1>"},
		{60, `<h1>@@@\@break</h1>`, "<h1>@@@@break</h1>"},
		{70, `\@`, `\@`},
		{80, `\\@`, `\\@`},
		{90, `\@if(true)`, `@if(true)`},
		{100, `\\@if(true)`, `\@if(true)`},
		{110, `\{{ 5 }}`, `{{ 5 }}`},
		{120, `\\{{ "nice" }}`, `\{{ "nice" }}`},
		{130, `\\\{{ x }}`, `\\{{ x }}`},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, tc.id)
	}
}

func TestEvalNumericExp(t *testing.T) {
	cases := []struct {
		id     uint
		inp    string
		expect string
	}{
		// Ints
		{10, "{{ 5; 5 }}", "55"},
		{20, "{{ 5 }}", "5"},
		{30, "{{ 10 }}", "10"},
		{40, "{{ -123 }}", "-123"},
		{50, `{{ 5 + 5 }}`, "10"},
		{60, `{{ 5 - 5 }}`, "0"},
		{70, `{{ 20 / 2 }}`, "10"},
		{80, `{{ 23 * 2 }}`, "46"},
		{90, `{{ 11 + 13 - 1 }}`, "23"},
		{100, "{{ 2 * (5 + 10) }}", "30"},
		{110, `{{ (3 + 5) * 2 }}`, "16"},
		{120, `{{ 3 * 3 * 3 + 10 }}`, "37"},
		{130, `{{ (5 + 10 * 2 + 15 / 3) * 2 + -10 }}`, "50"},
		{140, `{{ ((5 + 10) * ((2 + 15) / 3) + 2) }}`, "77"},
		// Floats
		{220, "{{ 5.11 }}", "5.11"},
		{230, "{{ -12.3 }}", "-12.3"},
		{240, `{{ 2.123 + 1.111 }}`, "3.234"},
		{250, `{{ 2.0 + 1.2 }}`, "3.2"},
		{260, `{{ 0.0 + 0.0 }}`, "0.0"},
		{270, `{{ 0.0 }}`, "0.0"},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, tc.id)
	}
}

func TestEvalBooleanExp(t *testing.T) {
	cases := []struct {
		id     uint
		inp    string
		expect string
	}{
		// --------------------------- Booleans ---------------------------
		{10, "{{ true }}", "1"},
		{20, "{{ false }}", "0"},
		{30, "{{ !true }}", "0"},
		{40, "{{ !false }}", "1"},
		{50, "{{ !nil }}", "1"},
		{60, "{{ !!true }}", "1"},
		{70, "{{ !!false }}", "0"},
		// --------------------------- Logical ---------------------------
		// Logical OR
		{80, "{{ true || false }}", "1"},
		{90, "{{ false || true }}", "1"},
		{100, "{{ false || false }}", "0"},
		{110, "{{ false || false || true }}", "1"},
		{120, "{{ '' || '3' }}", "1"},
		{130, "{{ 3 || 0 }}", "1"},
		{140, "{{ [] || [] }}", "0"},
		{150, "{{ {} || {} }}", "0"},
		{160, "{{ 0.2 || 2.3 }}", "1"},
		{170, "{{ 'a' || 'b' }}", "1"},
		{180, "{{ nil || nil }}", "0"},
		// Logical AND
		{190, "{{ false && false || true }}", "1"},
		{200, "{{ true && false || false }}", "0"},
		{210, "{{ false && (false || true) }}", "0"},
		{220, "{{ 3 && 0 }}", "0"},
		{230, "{{ [] && [] }}", "0"},
		{240, "{{ {} && {} }}", "0"},
		{250, "{{ 0.2 && 2.3 }}", "1"},
		{260, "{{ 'a' && 'b' }}", "1"},
		{270, "{{ '' && '' }}", "0"},
		{280, "{{ nil && nil }}", "0"},
		{290, "{{ true && true }}", "1"},
		{300, "{{ !false && !false }}", "1"},
		{310, "{{ false && false }}", "0"},
		{320, "{{ false && !false }}", "0"},
		// --------------------------- Integers ---------------------------
		// Integers - equality operators
		{330, "{{ 1 == 1 }}", "1"},
		{340, "{{ 0 == 0 }}", "1"},
		{350, "{{ 1 == 2 }}", "0"},
		{360, "{{ -1 == -1 }}", "1"},
		{370, "{{ -1 == 1 }}", "0"},
		{380, "{{ 100 == 100 }}", "1"},
		{390, "{{ 100 == 99 }}", "0"},
		// Integers - inequality operators
		{400, "{{ 1 != 1 }}", "0"},
		{410, "{{ 0 != 0 }}", "0"},
		{420, "{{ 1 != 2 }}", "1"},
		{430, "{{ -1 != 1 }}", "1"},
		{440, "{{ 100 != 100 }}", "0"},
		{450, "{{ 100 != 99 }}", "1"},
		// Integers - less than
		{460, "{{ 1 < 2 }}", "1"},
		{470, "{{ 0 < 1 }}", "1"},
		{480, "{{ -1 < 0 }}", "1"},
		{490, "{{ 1 < 1 }}", "0"},
		{500, "{{ 2 < 1 }}", "0"},
		{510, "{{ -2 < -3 }}", "0"},
		{520, "{{ 100 < 200 }}", "1"},
		// Integers - greater than
		{530, "{{ 1 > 2 }}", "0"},
		{540, "{{ 0 > 1 }}", "0"},
		{550, "{{ -1 > 0 }}", "0"},
		{560, "{{ 2 > 1 }}", "1"},
		{570, "{{ 0 > -1 }}", "1"},
		{580, "{{ 1 > 1 }}", "0"},
		{590, "{{ 200 > 100 }}", "1"},
		// Integers - less than or equal
		{600, "{{ 1 <= 2 }}", "1"},
		{610, "{{ 1 <= 1 }}", "1"},
		{620, "{{ 0 <= 0 }}", "1"},
		{630, "{{ 2 <= 1 }}", "0"},
		{640, "{{ -1 <= 0 }}", "1"},
		{650, "{{ -1 <= -1 }}", "1"},
		{660, "{{ 100 <= 50 }}", "0"},
		// Integers - greater than or equal
		{670, "{{ 1 >= 2 }}", "0"},
		{680, "{{ 1 >= 1 }}", "1"},
		{690, "{{ 0 >= 0 }}", "1"},
		{700, "{{ 2 >= 1 }}", "1"},
		{710, "{{ 0 >= -1 }}", "1"},
		{720, "{{ -1 >= -1 }}", "1"},
		{730, "{{ 50 >= 100 }}", "0"},
		// Integers with negative numbers
		{740, "{{ -5 == -5 }}", "1"},
		{750, "{{ -5 != -3 }}", "1"},
		{760, "{{ -10 < -5 }}", "1"},
		{770, "{{ -10 > -5 }}", "0"},
		{780, "{{ -5 <= -5 }}", "1"},
		{790, "{{ -5 >= -10 }}", "1"},
		// --------------------------- Floats ---------------------------
		// Floats - equality operators
		{800, "{{ 1.1 == 1.1 }}", "1"},
		{810, "{{ 0.0 == 0.0 }}", "1"},
		{820, "{{ 1.1 == 2.1 }}", "0"},
		{830, "{{ -1.5 == -1.5 }}", "1"},
		{840, "{{ -1.5 == 1.5 }}", "0"},
		{850, "{{ 3.14159 == 3.14159 }}", "1"},
		{860, "{{ 3.14 == 3.15 }}", "0"},
		// Floats - inequality operators
		{870, "{{ 1.1 != 1.1 }}", "0"},
		{880, "{{ 0.0 != 0.0 }}", "0"},
		{890, "{{ 1.1 != 2.1 }}", "1"},
		{900, "{{ -1.5 != 1.5 }}", "1"},
		{910, "{{ 3.14 != 3.141 }}", "1"},
		// Floats - less than
		{920, "{{ 1.1 < 2.1 }}", "1"},
		{930, "{{ 0.0 < 0.1 }}", "1"},
		{940, "{{ -1.5 < 0.0 }}", "1"},
		{950, "{{ 1.1 < 1.1 }}", "0"},
		{960, "{{ 2.5 < 1.5 }}", "0"},
		{970, "{{ -2.5 < -3.5 }}", "0"},
		// Floats - greater than
		{980, "{{ 1.1 > 2.1 }}", "0"},
		{990, "{{ 0.0 > 1.0 }}", "0"},
		{1000, "{{ -1.0 > 0.0 }}", "0"},
		{1010, "{{ 2.1 > 1.1 }}", "1"},
		{1020, "{{ 0.5 > -0.5 }}", "1"},
		{1030, "{{ 1.1 > 1.1 }}", "0"},
		// Floats - less than or equal
		{1040, "{{ 1.1 <= 2.1 }}", "1"},
		{1050, "{{ 1.1 <= 1.1 }}", "1"},
		{1060, "{{ 0.0 <= 0.0 }}", "1"},
		{1070, "{{ 2.1 <= 1.1 }}", "0"},
		{1080, "{{ -1.0 <= 0.0 }}", "1"},
		{1090, "{{ -1.0 <= -1.0 }}", "1"},
		// Floats - greater than or equal
		{1100, "{{ 1.1 >= 2.1 }}", "0"},
		{1110, "{{ 1.1 >= 1.1 }}", "1"},
		{1120, "{{ 0.0 >= 0.0 }}", "1"},
		{1130, "{{ 2.1 >= 1.1 }}", "1"},
		{1140, "{{ 0.0 >= -1.0 }}", "1"},
		{1150, "{{ -1.0 >= -1.0 }}", "1"},
		// Floats with negative numbers
		{1160, "{{ -5.5 == -5.5 }}", "1"},
		{1170, "{{ -5.5 != -3.5 }}", "1"},
		{1180, "{{ -10.5 < -5.5 }}", "1"},
		{1190, "{{ -10.5 > -5.5 }}", "0"},
		{1200, "{{ -5.5 <= -5.5 }}", "1"},
		{1210, "{{ -5.5 >= -10.5 }}", "1"},
		// --------------------------- Strings ---------------------------
		// Strings - equality operators
		{1220, "{{ 'hello' == 'hello' }}", "1"},
		{1230, "{{ '' == '' }}", "1"},
		{1240, "{{ 'hello' == 'world' }}", "0"},
		{1250, "{{ 'abc' == 'ABC' }}", "0"},
		{1260, "{{ ' test ' == ' test ' }}", "1"},
		{1270, "{{ 'a' == 'ab' }}", "0"},
		// Strings - inequality operators
		{1280, "{{ 'hello' != 'hello' }}", "0"},
		{1290, "{{ '' != '' }}", "0"},
		{1300, "{{ 'hello' != 'world' }}", "1"},
		{1310, "{{ 'abc' != 'ABC' }}", "1"},
		{1320, "{{ 'a' != 'ab' }}", "1"},
		// Strings - empty vs non-empty
		{1330, "{{ '' == 'a' }}", "0"},
		{1340, "{{ '' != 'a' }}", "1"},
		// Strings - numbers as strings
		{1350, "{{ '10' == '10' }}", "1"},
		{1360, "{{ '10' == '2' }}", "0"},
		// Strings - special characters
		{1370, "{{ 'hello world' == 'hello world' }}", "1"},
		{1380, "{{ 'test\nline' == 'test\nline' }}", "1"},
		// Booleans - equality operators
		{1390, "{{ true == true }}", "1"},
		{1400, "{{ false == false }}", "1"},
		{1410, "{{ true == false }}", "0"},
		{1420, "{{ false == true }}", "0"},
		// Booleans - inequality operators
		{1430, "{{ true != true }}", "0"},
		{1440, "{{ false != false }}", "0"},
		{1450, "{{ true != false }}", "1"},
		{1460, "{{ false != true }}", "1"},
		// --------------------------- Nils ---------------------------
		// Nil with boolean
		{1470, "{{ true == nil }}", "0"},
		{1480, "{{ nil == true }}", "0"},
		{1490, "{{ false == nil }}", "0"},
		{1500, "{{ nil == false }}", "0"},
		// Nil with integers
		{1510, "{{ nil == 0 }}", "0"},
		{1520, "{{ 0 == nil }}", "0"},
		{1530, "{{ nil == 1 }}", "0"},
		{1540, "{{ 1 == nil }}", "0"},
		{1550, "{{ nil != 0 }}", "1"},
		{1560, "{{ 0 != nil }}", "1"},
		{1570, "{{ nil != 5 }}", "1"},
		{1580, "{{ 5 != nil }}", "1"},
		// Nil with floats
		{1590, "{{ nil == 0.0 }}", "0"},
		{1600, "{{ 0.0 == nil }}", "0"},
		{1610, "{{ nil == 1.5 }}", "0"},
		{1620, "{{ 1.5 == nil }}", "0"},
		{1630, "{{ nil != 0.0 }}", "1"},
		{1640, "{{ 0.0 != nil }}", "1"},
		{1650, "{{ nil != 3.14 }}", "1"},
		{1660, "{{ 3.14 != nil }}", "1"},
		// Nil with strings
		{1670, "{{ nil == '' }}", "0"},
		{1680, "{{ '' == nil }}", "0"},
		{1690, "{{ nil == 'test' }}", "0"},
		{1700, "{{ 'test' == nil }}", "0"},
		{1710, "{{ nil != '' }}", "1"},
		{1720, "{{ '' != nil }}", "1"},
		{1730, "{{ nil != 'hello' }}", "1"},
		{1740, "{{ 'hello' != nil }}", "1"},
		// Nil with arrays
		{1750, "{{ nil == [] }}", "0"},
		{1760, "{{ [] == nil }}", "0"},
		{1770, "{{ nil == [1, 2] }}", "0"},
		{1780, "{{ [1, 2] == nil }}", "0"},
		{1790, "{{ nil != [] }}", "1"},
		{1800, "{{ [] != nil }}", "1"},
		{1810, "{{ nil != [1] }}", "1"},
		{1820, "{{ [1] != nil }}", "1"},
		// Nil with objects
		{1830, "{{ nil == {} }}", "0"},
		{1840, "{{ {} == nil }}", "0"},
		{1850, "{{ nil == {name: 'test'} }}", "0"},
		{1860, "{{ {name: 'test'} == nil }}", "0"},
		{1870, "{{ nil != {} }}", "1"},
		{1880, "{{ {} != nil }}", "1"},
		{1890, "{{ nil != {x: 1} }}", "1"},
		{1900, "{{ {x: 1} != nil }}", "1"},
		// Nil with nil
		{1910, "{{ nil == nil }}", "1"},
		{1920, "{{ nil != nil }}", "0"},
		// --------------------------- Arrays ---------------------------
		// Arrays with arrays
		{1930, "{{ [] == [] }}", "1"},
		{1940, "{{ [] != [] }}", "0"},
		{1950, "{{ [1, 2] == [1] }}", "0"},
		{1960, "{{ [1, 2] == [1, 2] }}", "1"},
		{1970, "{{ [1, 2] != [1, 2] }}", "0"},
		{1980, "{{ [1] == [2] }}", "0"},
		{1990, "{{ [1] != [2] }}", "1"},
		{2000, "{{ [1] == [1] }}", "1"},
		{2010, "{{ [1] != [1] }}", "0"},
		{2020, "{{ [[[1]]] == [[[1]]] }}", "1"},
		{2030, "{{ [[[1]]] != [[[1]]] }}", "0"},
		{2040, "{{ [[[1]]] == [[[2]]] }}", "0"},
		{2050, "{{ [[[1]]] != [[[2]]] }}", "1"},
		// --------------------------- Objects ---------------------------
		// Objects with objects
		{2060, "{{ {} == {} }}", "1"},
		{2070, "{{ {} != {} }}", "0"},
		{2080, "{{ {name: 'test'} == {name: 'test'} }}", "1"},
		{2090, "{{ {name: 'test'} != {name: 'test'} }}", "0"},
		{2100, "{{ {x: 1} == {x: 1} }}", "1"},
		{2110, "{{ {x: 1} != {x: 1} }}", "0"},
		{2120, "{{ {x: {y: 2}} == {x: {y: 2}} }}", "1"},
		{2130, "{{ {x: {y: 2}} != {x: {y: 2}} }}", "0"},
		{2140, "{{ {x: {y: {z: 'Anna'}}} == {x: {y: {z: 'Anna'}}} }}", "1"},
		{2150, "{{ {x: {y: {z: 'Anna'}}} != {x: {y: {z: 'Anna'}}} }}", "0"},
		// Mixed type comparisons with == and !=
		{2160, "{{ 1 == '1' }}", "0"},
		{2170, "{{ '1' == 1 }}", "0"},
		{2180, "{{ true == 1 }}", "0"},
		{2190, "{{ 1 == true }}", "0"},
		{2200, "{{ false == 0 }}", "0"},
		{2210, "{{ 0 == false }}", "0"},
		{2220, "{{ 1.0 == 1 }}", "0"},
		{2230, "{{ 1 == 1.0 }}", "0"},
		{2240, "{{ {} == 1 }}", "0"},
		{2250, "{{ [] == 1 }}", "0"},
	}

	for _, tc := range cases {
		evaluated, failure := testEval(tc.inp)
		if failure != nil {
			t.Fatalf("Case: %d. evaluation failed: %s", tc.id, failure)
		}

		err, ok := evaluated.(*value.Error)
		if ok {
			t.Errorf("Case: %d. Evaluation failed: %s", tc.id, err)
			return
		}

		if res := evaluated.String(); res != tc.expect {
			t.Errorf("Case: %d. Result is not %q, got %q", tc.id, tc.expect, res)
		}
	}
}

func TestEvalNilExp(t *testing.T) {
	inp := "<h1>{{ nil }}</h1>"
	evaluationExpected(t, inp, "<h1></h1>", 10)
}

func TestEvalStringExp(t *testing.T) {
	cases := []struct {
		id     uint
		inp    string
		expect string
	}{
		// Basic string output
		{10, `{{ "Hello World" }}`, "Hello World"},
		{20, `{{ 'Hello World 2' }}`, "Hello World 2"},
		// String in HTML attributes
		{30, `<div {{ 'data-attr="Test"' }}></div>`, "<div data-attr=&#34;Test&#34;></div>"},
		{40, `<div {{ "data-attr='Test'" }}></div>`, "<div data-attr=&#39;Test&#39;></div>"},
		// String with escaped characters
		{50, `{{ "She \"is\" pretty" }}`, "She &#34;is&#34; pretty"},
		// String concatenation
		{60, `{{ "Korotchaeva" + " " + "Anna" }}`, "Korotchaeva Anna"},
		{70, `{{ "She" + " " + "is" + " " + "nice" }}`, "She is nice"},
		// Empty string
		{80, "{{ '' }}", ""},
		{90, `{{ "" }}`, ""},
		// String with HTML escaping
		{100, `{{ "<h1>Test</h1>" }}`, "&lt;h1&gt;Test&lt;/h1&gt;"},
		{110, `{{ "<div>Hello</div>" }}`, "&lt;div&gt;Hello&lt;/div&gt;"},
		{
			120,
			`{{ "<script>alert('xss')</script>" }}`,
			"&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;",
		},
		{130, `{{ "<img src='test.jpg'>" }}`, "&lt;img src=&#39;test.jpg&#39;&gt;"},
		{140, `{{ "&amp;" }}`, "&amp;amp;"},
		{150, `{{ "<br/>" }}`, "&lt;br/&gt;"},
		{160, `{{ "<a href='#'>Link</a>" }}`, "&lt;a href=&#39;#&#39;&gt;Link&lt;/a&gt;"},
		{
			170,
			`{{ "<p class='text'>Content</p>" }}`,
			"&lt;p class=&#39;text&#39;&gt;Content&lt;/p&gt;",
		},
		{180, `{{ "<ul><li>Item</li></ul>" }}`, "&lt;ul&gt;&lt;li&gt;Item&lt;/li&gt;&lt;/ul&gt;"},
		{
			190,
			`{{ "<input type='text' value='test'>" }}`,
			"&lt;input type=&#39;text&#39; value=&#39;test&#39;&gt;",
		},
		// String concatenation with variables
		{200, `{{ name = "Anna"; "Hello " + name }}`, "Hello Anna"},
		{210, `{{ a = "Hello"; b = "World"; a + " " + b }}`, "Hello World"},
		// String with numbers
		{220, `{{ "Count: " + 5.str() }}`, "Count: 5"},
		{230, `{{ "Pi: " + 3.14.str() }}`, "Pi: 3.14"},
		// Empty string concatenation
		{240, `{{ "" + "test" }}`, "test"},
		{250, `{{ "test" + "" }}`, "test"},
		{260, `{{ "" + "" }}`, ""},
		{270, `{{ '' + '' }}`, ""},
		{280, `{{ '' + "" }}`, ""},
		{290, `{{ "" + '' }}`, ""},
		// String with special characters
		{300, `{{ "test\nline" }}`, "test\nline"},
		{310, `{{ "tab\there" }}`, "tab\there"},
		// Long string concatenation
		{320, `{{ "a" + "b" + "c" + "d" + "e" }}`, "abcde"},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, tc.id)
	}
}

func TestEvalTernaryExp(t *testing.T) {
	cases := []struct {
		id     uint
		inp    string
		expect string
	}{
		// Basic boolean conditions
		{10, `{{ true ? "Yes" : "No" }}`, "Yes"},
		{20, `{{ false ? "Yes" : "No" }}`, "No"},
		// Falsy values
		{30, `{{ nil ? "Yes" : "No" }}`, "No"},
		{40, `{{ 1 ? "Yes" : "No" }}`, "Yes"},
		{50, `{{ 0 ? "Yes" : "No" }}`, "No"},
		{60, `{{ "" ? "Yes" : "No" }}`, "No"},
		// Negation
		{70, `{{ !true ? "Yes" : "No" }}`, "No"},
		{80, `{{ !false ? "Yes" : "No" }}`, "Yes"},
		// Double negation
		{90, `{{ !!true ? 1 : 0 }}`, "1"},
		{100, `{{ !!false ? 1 : 0 }}`, "0"},
		// Comparison operators
		{110, `{{ 1 == 1 ? "Yes" : "No" }}`, "Yes"},
		{120, `{{ 1 == 2 ? "Yes" : "No" }}`, "No"},
		{130, `{{ 5 > 3 ? "Yes" : "No" }}`, "Yes"},
		{140, `{{ 3 > 5 ? "Yes" : "No" }}`, "No"},
		// String comparison
		{150, `{{ "a" == "a" ? "Yes" : "No" }}`, "Yes"},
		{160, `{{ "a" == "b" ? "Yes" : "No" }}`, "No"},
		// Array truthiness
		{170, `{{ [] ? "Yes" : "No" }}`, "No"},
		{180, `{{ [1] ? "Yes" : "No" }}`, "Yes"},
		// Object truthiness
		{190, `{{ {} ? "Yes" : "No" }}`, "No"},
		{200, `{{ {x: 1} ? "Yes" : "No" }}`, "Yes"},
		// Nested ternary
		{210, `{{ true ? (true ? "A" : "B") : "C" }}`, "A"},
		{220, `{{ true ? (false ? "A" : "B") : "C" }}`, "B"},
		{230, `{{ false ? (true ? "A" : "B") : "C" }}`, "C"},
		// Arithmetic expressions
		{240, `{{ 1 + 1 == 2 ? "Yes" : "No" }}`, "Yes"},
		{250, `{{ 10 - 5 > 3 ? "Yes" : "No" }}`, "Yes"},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, tc.id)
	}
}

func TestEvalIfDir(t *testing.T) {
	cases := []struct {
		id     uint
		inp    string
		expect string
	}{
		// Basic if statements
		{10, `@if(true)Hello@end`, "Hello"},
		{20, `@if(false)Hello@end`, ""},
		{30, `@if(true)X@else@end`, "X"},
		{40, `@if(false)@elseX@end`, "X"},
		{50, `@if(true)@else@end`, ""},
		{60, `@if(false)1@else@end`, ""},
		// HTML wrapper
		{70, `<h2>@if(true)Hello@end</h2>`, "<h2>Hello</h2>"},
		{80, `<h2>@if(false)Hello@end</h2>`, "<h2></h2>"},
		// @elseif statements
		{90, `@if(true)Anna@elseif(true)Lili@end`, "Anna"},
		{100, `@if(false)Alan@elseif(true)Serhii@end`, "Serhii"},
		{110, `@if(false)Ana De Armaz@elseif(false)David@elseVladimir@end`, "Vladimir"},
		{120, `@if(false)Will@elseif(false)Daria@elseif(true)Poll@end`, "Poll"},
		{130, `@if(false)Lara@elseif(true)Susan@elseif(true)Smith@end`, "Susan"},
		{140, `@if(false)@elseif(false)@elsemy@mail.com@end`, "my@mail.com"},
		{150, `@if(false)A@elseif(false)B@end`, ""},
		// Functions
		{160, `@if(true.binary())Hello@end`, "Hello"},
		{170, `@if(false.binary())Hello@end`, ""},
		{180, `@if("".len() > 0)Non empty@elseEmpty@end`, "Empty"},
		{190, `@if("x".len() > 0)Non empty@elseEmpty@end`, "Non empty"},
		// Truthy/falsy values
		{200, `@if(1)Yes@end`, "Yes"},
		{210, `@if(0)Yes@end`, ""},
		{220, `@if("")Yes@end`, ""},
		{230, `@if([])Yes@end`, ""},
		{240, `@if({})Yes@end`, ""},
		{250, `@if(nil)Yes@end`, ""},
		// Logical operators
		{260, `@if(true && true)Yes@end`, "Yes"},
		{270, `@if(true && false)Yes@end`, ""},
		{280, `@if(false || true)Yes@end`, "Yes"},
		{290, `@if(false || false)Yes@end`, ""},
		// Boolean negation
		{300, `@if(!true)Yes@end`, ""},
		{310, `@if(!false)Yes@end`, "Yes"},
		{320, `@if(!!true)Yes@end`, "Yes"},
		{330, `@if(!!false)Yes@end`, ""},
		// Comparison operators
		{340, `@if(1 == 1)Yes@end`, "Yes"},
		{350, `@if(1 == 2)Yes@end`, ""},
		{360, `@if(1 != 2)Yes@end`, "Yes"},
		{370, `@if(1 < 2)Yes@end`, "Yes"},
		{380, `@if(2 > 1)Yes@end`, "Yes"},
		{390, `@if(2 >= 2)Yes@end`, "Yes"},
		{400, `@if(2 <= 2)Yes@end`, "Yes"},
		// Expression results as conditions
		{410, `@if(1 + 1 == 2)Yes@end`, "Yes"},
		{420, `@if(5 - 3 == 1)Yes@end`, ""},
		{430, `@if(2 - 2)No@elseYes@end`, "Yes"},
		{440, `@if(-1 + 2)Yes@elseNo@end`, "Yes"},
		// Nested if statements
		{
			390,
			`
				@if(true)
					@if(false)
					    James
					@elseif(false)
						John
					@else
						@if(true){{ "Marry" }}@end
					@end
				@else
					@if(true)Anna@end
				@end
			`,
			"Marry",
		},
	}

	for _, tc := range cases {
		evaluated, failure := testEval(tc.inp)
		if failure != nil {
			t.Fatalf("Case: %d. evaluation failed: %s", tc.id, failure)
		}

		err, ok := evaluated.(*value.Error)
		if ok {
			t.Errorf("Case: %d. Evaluation failed: %s", tc.id, err)
		}

		if res := strings.TrimSpace(evaluated.String()); res != tc.expect {
			t.Errorf("Case: %d. Result is not %q, got %q", tc.id, tc.expect, res)
		}
	}
}

func TestEvalArr(t *testing.T) {
	cases := []struct {
		id     uint
		inp    string
		expect string
	}{
		// Empty arrays
		{10, `{{ [] }}`, ""},
		{20, `{{ [[]] }}`, ""},
		{30, `{{ [[[[[]]]]] }}`, ""},
		// Integer arrays
		{40, `{{ [1, 2, 3] }}`, "1, 2, 3"},
		{50, `{{ [0] }}`, "0"},
		{60, `{{ [-1, -2, -3] }}`, "-1, -2, -3"},
		// Float arrays
		{70, `{{ [1.5, 2.5, 3.5] }}`, "1.5, 2.5, 3.5"},
		{80, `{{ [0.0] }}`, "0.0"},
		// String arrays
		{90, `{{ ["Anna", "Serhii" ] }}`, "Anna, Serhii"},
		{100, `{{ ["hello"] }}`, "hello"},
		{110, `{{ ["with space", "another"] }}`, "with space, another"},
		// Boolean arrays
		{120, `{{ [true, false] }}`, "1, 0"},
		{130, `{{ [true] }}`, "1"},
		{140, `{{ [false] }}`, "0"},
		// Mixed type arrays
		{150, `{{ [1, "hello", true] }}`, "1, hello, 1"},
		{160, `{{ [0, "", false] }}`, "0, , 0"},
		// Nested arrays
		{170, `{{ [[1, 2], [3, 4]] }}`, "1, 2, 3, 4"},
		{180, `{{ [[1, [2]], 3] }}`, "1, 2, 3"},
		{190, `{{ [[[11]]] }}`, "11"},
		// Arrays with expressions
		{200, `{{ [1 + 2, 3 * 2] }}`, "3, 6"},
		{210, `{{ [5 - 3, 10 / 2] }}`, "2, 5"},
		// Arrays with variables
		{220, `{{ a = 1; b = 2; [a, b] }}`, "1, 2"},
		{230, `{{ x = "test"; [x, x] }}`, "test, test"},
		// Arrays with ternary
		{240, `{{ [true ? 1 : 0, false ? "yes" : "no"] }}`, "1, no"},
		// Trailing comma
		{250, `{{ [1, 2,] }}`, "1, 2"},
		{260, `{{ ["a", "b",] }}`, "a, b"},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, tc.id)
	}
}

func TestEvalIndexExp(t *testing.T) {
	cases := []struct {
		id     uint
		inp    string
		expect string
	}{
		// Basic indexing
		{10, `{{ [1, 2, 3][0] }}`, "1"},
		{20, `{{ [1, 2, 3][1] }}`, "2"},
		{30, `{{ [1, 2, 3][2] }}`, "3"},
		// Negative indexing (not allowed)
		{40, `{{ [1, 2, 3][-1] }}`, ""},
		{50, `{{ [1, 2, 3][-2] }}`, ""},
		{60, `{{ ["a", "b", "c"][-3] }}`, ""},
		// Out of bounds indexing
		{70, `{{ [][2] }}`, ""},
		{80, `{{ [1, 2][5] }}`, ""},
		{90, `{{ [1, 2][-5] }}`, ""},
		// String indexing
		{100, `{{ ["Some string"][0] }}`, "Some string"},
		{110, `{{ ["hello", "world"][1] }}`, "world"},
		// Nested array indexing
		{120, `{{ [[1, 2], [3, 4]][0][1] }}`, "2"},
		{130, `{{ [[1, 2], [3, 4]][1][0] }}`, "3"},
		{140, `{{ [[[11]]][0][0][0] }}`, "11"},
		// Index with expression
		{150, `{{ arr = [1, 2, 3]; arr[1 + 1] }}`, "3"},
		{160, `{{ arr = [1, 2, 3]; arr[3 - 1] }}`, "3"},
		// Variable as index
		{170, `{{ arr = [1, 2, 3]; i = 2; arr[i] }}`, "3"},
		{180, `{{ arr = ["x", "y", "z"]; i = 0; arr[i] }}`, "x"},
		// Ternary as index
		{190, `{{ ["a", "b", "c"][true ? 0 : 1] }}`, "a"},
		{200, `{{ ["a", "b", "c"][false ? 0 : 1] }}`, "b"},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, tc.id)
	}
}

func TestEvalAssign(t *testing.T) {
	cases := []struct {
		id     uint
		inp    string
		expect string
	}{
		// Integer assignment
		{10, `{{ age = 18 }}`, ""},
		{20, `{{ age = 18; age }}`, "18"},
		{30, `{{ age = -5; age }}`, "-5"},
		{40, `{{ age = 0; age }}`, "0"},
		// Float assignment
		{50, `{{ pi = 3.14; pi }}`, "3.14"},
		{60, `{{ price = 0.0; price }}`, "0.0"},
		{70, `{{ negative = -2.5; negative }}`, "-2.5"},
		{80, `{{ f = 0.0; f = f + 0.5; f }}`, "0.5"},
		// String assignment
		{90, `{{ name = "Anna"; name }}`, "Anna"},
		{100, `{{ empty = ""; empty }}`, ""},
		{110, `{{ quote = "He said \"Hello\""; quote }}`, "He said &#34;Hello&#34;"},
		// Boolean assignment
		{120, `{{ flag = true; flag }}`, "1"},
		{130, `{{ flag = false; flag }}`, "0"},
		// Nil assignment
		{140, `{{ nothing = nil; nothing }}`, ""},
		// Multiple assignments
		{150, `{{ a = 1; b = 2; c = 3; a + b + c }}`, "6"},
		{160, `{{ x = "Hello"; y = "World"; x + " " + y }}`, "Hello World"},
		// Reassignment
		{170, `{{ age = 18; age = 25; age }}`, "25"},
		{180, `{{ name = "Anna"; name = "Maria"; name }}`, "Maria"},
		{190, `{{ x = 1; x = x + 1; x = x + 1; x }}`, "3"},
		// Assignment with expression
		{200, `{{ sum = 5 + 3; sum }}`, "8"},
		{210, `{{ calc = 10 * 2 - 5; calc }}`, "15"},
		{220, `{{ result = (2 + 3) * 4; result }}`, "20"},
		// Assignment with string concatenation
		{230, `{{ full = "John" + " " + "Doe"; full }}`, "John Doe"},
		{240, `{{ msg = "Count: " + 5.str(); msg }}`, "Count: 5"},
		// Assignment with ternary
		{250, `{{ val = true ? 1 : 0; val }}`, "1"},
		{260, `{{ val = false ? "yes" : "no"; val }}`, "no"},
		// Assignment with method call
		{270, `{{ upper = "hello".upper(); upper }}`, "HELLO"},
		{280, `{{ len = [1, 2, 3].len(); len }}`, "3"},
		// Assignment with index
		{290, `{{ arr = [10, 20, 30]; val = arr[1]; val }}`, "20"},
		{300, `{{ str = "abc"; val = str.at(0); val }}`, "a"},
		// Assignment in conditional
		{310, `@if(true){{ x = 5 }}{{ x }}@end`, "5"},
		{320, `@if(false){{ x = 5 }}{{ x }}@end`, ""},
		// Assignment with loop variable
		{330, `@each(n in [1, 2, 3]){{ x = n }}{{ x }}@end`, "123"},
		// Object assignment
		{340, `{{ user = {name: "John"}; user.name }}`, "John"},
		{350, `{{ data = {"age": 25}; data.age }}`, "25"},
		{360, `{{ user = {"name": "Ann"}; user = {"name": "Anna"}; user.name }}`, "Anna"},
		{370, `{{ user = {}; user.name = "Anna"; user.name }}`, "Anna"},
		{380, `{{ user = {}; user.name = "Anna"; user.name }}`, "Anna"},
		{
			id:     420,
			inp:    `{{ user = {}; user.address = { street: "x" }; user.address.street = 'y'; user.address.street }}`,
			expect: "y",
		},
		{
			id:     430,
			inp:    `{{ user = {}; user.address = {}; user.address.street = "Pushkina"; user.address.street }}`,
			expect: "Pushkina",
		},
		// Array assignment
		{390, `{{ names = ["Anna", "Serhii"]; names }}`, "Anna, Serhii"},
		{400, `{{ empty = []; empty }}`, ""},
		{410, `{{ nested = [[1, 2], [3, 4]]; nested }}`, "1, 2, 3, 4"},
		{420, `{{ names = ['Serhii', 'Nastya']; names[1] = 'Anna'; names }}`, "Serhii, Anna"},
		{430, `{{ nums = [10, 20, 30]; nums[0] = 1; nums[1] = 2; nums }}`, "1, 2, 30"},
		{440, `{{ x = [[[20]]]; x[0][0][0] = 30; x }}`, "30"},
		{450, `{{ x = ['1', ['2', ['3', ['4']]]]; x[1][1][1][0] = '5'; x }}`, "1, 2, 3, 5"},
		{460, `{{ x = [0]; newVal = 10; x[0] = newVal; x[0] }}`, "10"},
		// Mixed assignment
		{470, `{{ x = [{ name: 'Chiori' }]; x[0].name = 'Mavuika'; x[0].name }}`, "Mavuika"},
		{
			480,
			`{{ name = 'Mavuika'; x = [{ name: 'Chiori' }]; x[0].name = name; x[0].name }}`,
			"Mavuika",
		},
		// Index assignment edge cases
		{490, `{{ arr = [1, 2]; arr[0] = {name: 'x'}; arr[0].name }}`, "x"},
		{500, `{{ arr = [1, 2]; arr[0] = [3, 4]; arr[0][0] }}`, "3"},
		{510, `{{ arr = [1, 2, 3]; arr[1 + 1] = 5; arr[2] }}`, "5"},
		{520, `{{ arr = [[1], [2]]; arr[0][0] = 10; arr[0][0] }}`, "10"},
		{530, `{{ x = 5; arr = [1, 2]; arr[0] = x; arr[0] }}`, "5"},
		{540, `{{ arr = [1, 2]; arr[0] = arr[1]; arr[0] }}`, "2"},
		{550, `{{ arr = [1, 2, 3]; arr[2] = 10; arr[2] }}`, "10"},
		{560, `{{ arr = ['a', 'b', 'c']; arr[0] = arr[0] + 'x'; arr[0] }}`, "ax"},
		{570, `{{ arr = [true, false]; arr[0] = false; arr[0] }}`, "0"},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, tc.id)
	}
}

func TestIncDecDir(t *testing.T) {
	cases := []struct {
		id     uint
		inp    string
		expect string
	}{
		// Integer
		{10, `{{ i = 10; i++; i }}`, "11"},
		{20, `{{ x = 0; x++; x }}`, "1"},
		// Float
		{30, `{{ x = 4.4; x--; x }}`, "3.4"},
		{48, `{{ x = 0.0; x--; x }}`, "-1.0"},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, tc.id)
	}
}

func TestEvalForDir(t *testing.T) {
	cases := []struct {
		id     uint
		inp    string
		expect string
	}{
		// Basic for loops
		{10, `@for(i = 0; i < 2; i++){{ i }}@end`, "01"},
		{20, `@for(i = 1; i <= 3; i++){{ i }}@end`, "123"},
		{30, `@for(i = 5; i > 0; i--){{ i }}@end`, "54321"},
		// Empty loop header parts
		{50, `@for(; false;)Here@end`, ""},
		{60, `@for(; true;)x@break@end`, "x"},
		{70, `@for(c = 1; false; c++){{ c }}@end`, ""},
		{80, `@for(i = 0; i < 0; i++){{ i }}@end`, ""},
		{90, `@for(;;){{ 1 }}@break@end`, "1"},
		// Single iteration
		{100, `@for(c = 1; c == 1; c++){{ c }}@end`, "1"},
		{110, `@for(i = 0; i < 1; i++){{ i }}@end`, "0"},
		// @else directive
		{120, `@for(c = 1; false; c++){{ c }}@else@end`, ""},
		{130, `@for(c = 1; false; c++){{ c }}@else<b>Empty</b>@end`, "<b>Empty</b>"},
		{140, `@for(c = 0; c < 0; c++){{ c }}@elseEmpty@end`, "Empty"},
		// @break directive
		{150, `@for(i = 1; i <= 3; i++){{ i }}@break@end`, "1"},
		{160, `@for(i = 1; i <= 3; i++)@break{{ i }}@end`, ""},
		{170, `@for(i = 1; i <= 3; i++)@if(i == 3)@break@end{{ i }}@end`, "12"},
		{180, `@for(i = 0; i < 10; i++)@break@end`, ""},
		// @continue directive
		{190, `@for(i = 1; i <= 3; i++)@continue{{ i }}@end`, ""},
		{200, `@for(i = 1; i <= 3; i++){{ i }}@continue@end`, "123"},
		{210, `@for(i = 1; i <= 3; i++)@if(i == 2)@continue@end{{ i }}@end`, "13"},
		{220, `@for(i = 1; i <= 5; i++)@if(i % 2 == 0)@continue@end{{ i }}@end`, "135"},
		// @breakif directive
		{230, `@for(i = 1; i <= 3; i++)@breakif(i == 3){{ i }}@end`, "12"},
		{240, `@for(i = 1; i <= 3; i++)@breakif(i == 2){{ i }}@end`, "1"},
		{250, `@for(i = 1; i <= 10; i++)@breakif(i > 5){{ i }}@end`, "12345"},
		// @continueif directive
		{260, `@for(i = 1; i <= 3; i++)@continueif(i == 3){{ i }}@end`, "12"},
		{270, `@for(i = 1; i <= 3; i++)@continueif(i == 2){{ i }}@end`, "13"},
		{280, `@for(i = 1; i <= 5; i++)@continueif(i % 2 == 0){{ i }}@end`, "135"},
		// Nested for loops
		{290, `@for(i = 0; i < 2; i++)@for(j = 0; j < 2; j++){{ i }}{{ j }}@end@end`, "00011011"},
		{300, `@for(i = 1; i <= 2; i++)@for(j = 1; j <= 2; j++){{ i * j }}@end@end`, "1224"},
		// For loop with HTML
		{
			300,
			`<ul>@for(i = 1; i <= 3; i++)<li>{{ i }}</li>@end</ul>`,
			"<ul><li>1</li><li>2</li><li>3</li></ul>",
		},
		// Variable modification from outside scope
		{310, `{{ sum = 0 }}@for(i = 1; i <= 5; i++){{ sum = sum + i }}@end{{ sum }}`, "15"},
		{311, `{{ x = 0 }}@for(; x < 2; x++){{ x }}@end`, "01"},
		{320, `{{ count = 0 }}@for(i = 0; i < 3; i++){{ count = count + 1 }}@end{{ count }}`, "3"},
		{330, `{{ n = 0 }}@for(; true; n++){{ n }}@breakif(n == 2)@end`, "012"},
		{340, `{{ i = 0 }}@for(; i < 3; i++){{ i }}@end`, "012"},
		// Multiple statements in loop body
		{350, `@for(i = 0; i < 3; i++){{ i }};{{ i * 2 }}@end`, "0;01;22;4"},
		{
			id:     360,
			inp:    `@for(i = 1; i <= 3; i++){{ i }} -> {{ i * i }}|@end`,
			expect: "1 -> 1|2 -> 4|3 -> 9|",
		},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, tc.id)
	}
}

func TestEvalEachDir(t *testing.T) {
	cases := []struct {
		id     uint
		inp    string
		expect string
	}{
		{10, `@each(name in ["anna", "serhii"]){{ name }} @end`, "anna serhii "},
		{20, `@each(num in [1, 2, 3]){{ num }}@end`, "123"},
		{30, `@each(num in []){{ num }}@end`, ""},
		// test loop variable
		{40, `@each(num in [43, 12, 53]){{ loop.index }}@end`, "012"},
		{50, `@each(num in [100]){{ loop.index }}@end`, "0"},
		{60, `@each(num in [1, 2, 3, 4]){{ loop.first }}@end`, "1000"},
		{70, `@each(num in [1, 2, 3, 4]){{ loop.last }}@end`, "0001"},
		{80, `@each(num in [4, 2, 8]){{ loop.iter }}@end`, "123"},
		{90, `@each(num in [9, 3, 44, 24, 1, 3]){{ loop.iter }}@end`, "123456"},
		// test @else directive
		{100, `@each(v in []){{ v }}@else<b>Empty array</b>@end`, "<b>Empty array</b>"},
		{110, `@each(n in []){{ n }}@else@end`, ""},
		{120, `@each(n in []){{ n }}@elsetest@end`, "test"},
		{130, `@each(n in [1, 2, 3, 4, 5]){{ n }}@end`, "12345"},
		// test @break directive
		{140, `@each(n in [1, 2, 3, 4, 5])@break{{ n }}@end`, ""},
		{150, `@each(n in [1, 2, 3, 4, 5]){{ n }}@break@end`, "1"},
		{160, `@each(n in [1, 2, 3, 4, 5])@if(n == 3)@break@end{{ n }}@end`, "12"},
		// test @continue directive
		{170, `@each(n in [1, 2, 3, 4, 5])@continue{{ n }}@end`, ""},
		{180, `@each(n in [1, 2, 3, 4, 5]){{ n }}@continue@end`, "12345"},
		{190, `@each(n in [1, 2, 3, 4, 5])@if(n == 3)@continue@end{{ n }}@end`, "1245"},
		// test @breakif directive
		{200, `@each(n in [1, 2, 3, 4, 5])@breakif(n == 3){{ n }}@end`, "12"},
		{
			210,
			`@each(n in ["ann", "serhii", "sam"])@breakif(n == 'sam'){{ n }} @end`,
			"ann serhii ",
		},
		// test @continueif directive
		{210, `@each(n in [1, 2, 3, 4, 5])@continueif(n == 3){{ n }}@end`, "1245"},
		{
			230,
			`@each(n in ["ann", "serhii", "sam"])@continueif(n == 'sam'){{ n }} @end`,
			"ann serhii ",
		},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, tc.id)
	}
}

func TestEvalObjLit(t *testing.T) {
	cases := []struct {
		id     uint
		inp    string
		expect string
	}{
		// Bracket and dot access
		{10, `{{ {"name": "Ann"}['name'] }}`, "Ann"},
		{20, `{{ {"name": "Ann"}.name }}`, "Ann"},
		// Basic property access
		{30, `{{ obj = {name: "Ann"}; obj.name }}`, "Ann"},
		{40, `{{ o = {"name": "Ann", "age": 22}; o.age }}`, "22"},
		// Nested objects
		{50, `{{ user = {"father": {"name": "Ann"}}; user.father.name }}`, "Ann"},
		{
			60,
			`{{ user = {"father": {"name": {"first": "Serhii"}}}; user.father.name.first }}`,
			"Serhii",
		},
		{
			70,
			`{{ u = {"father": {name: {"first": "Serhii",},},}; u['father']['name'].first }}`,
			"Serhii",
		},
		// Shorthand properties
		{80, `{{ name = "Serhii"; age = 12; obj = { name, age }; obj.name }}`, "Serhii"},
		{90, `{{ name = "Serhii"; age = 12; obj = { name, age }; obj.age }}`, "12"},
		// Case-insensitive first character access
		{100, `{{ {"Name": "Ann"}.name }}`, "Ann"},
		{110, `{{ {"name": "Ann"}.Name }}`, ""},
		// Non-existent keys
		{120, `{{ obj = {"name": "Ann"}; obj.age }}`, ""},
		{130, `{{ obj = {"name": "Ann"}; obj['missing'] }}`, ""},
		// Empty object
		{140, `{{ {}.name }}`, ""},
		{150, `{{ obj = {}; obj.name }}`, ""},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, tc.id)
	}
}

func TestEvalComments(t *testing.T) {
	cases := []struct {
		id     uint
		inp    string
		expect string
	}{
		// Basic comments
		{10, "{{-- This is a comment --}}", ""},
		{20, "{{-- --}}", ""},
		{30, "{{--a--}}", ""},
		// Comments in HTML
		{40, "<section>{{-- This is a comment --}}</section>", "<section></section>"},
		{50, "<div>{{-- comment --}}</div>", "<div></div>"},
		// Comments with text around
		{60, "Some {{-- --}}text", "Some text"},
		{70, "Hello{{-- comment --}}World", "HelloWorld"},
		{80, "Start{{-- middle --}}End", "StartEnd"},
		// Empty or minimal comments
		{90, "{{----}}", ""},
		{100, "{{--     --}}", ""},
		{110, "{{--\n--}}", ""},
		// Comments with directives
		{120, "{{-- @each(u in users){{ u }}@end --}}", ""},
		{130, "{{-- @if(true)Hello@end --}}", ""},
		{140, "{{-- @for(i=0;i<10;i++){{i}}@end --}}", ""},
		// Comments with special characters
		{150, "{{-- <html> &amp; @#$%^&*() --}}", ""},
		{160, "{{-- Japanese: こんにちは --}}", ""},
		{170, "{{-- Emoji: 😀🎉🚀 --}}", ""},
		// Comments inside directives
		{180, "@if(true){{-- comment --}}Yes@end", "Yes"},
		{190, "@if(true)Yes{{-- comment --}}@end", "Yes"},
		{200, "@each(n in [1,2]){{-- c --}}{{n}}@end", "12"},
		// Multiple comments
		{210, "{{-- a --}}Text{{-- b --}}", "Text"},
		{220, "{{-- x --}}{{-- y --}}Result", "Result"},
		{230, "A{{-- 1 --}}B{{-- 2 --}}C", "ABC"},
		// Comments at boundaries
		{240, "{{-- start --}}End", "End"},
		{250, "Start{{-- end --}}", "Start"},
		// Multi-line comments
		{260, "{{-- Line 1\nLine 2\nLine 3 --}}", ""},
		{270, "Before{{--\nmulti\nline\n--}}After", "BeforeAfter"},
		// Nested comment-like text (should be treated as part of comment)
		{280, "{{-- Contains {{-- and --}} inside --}}", ""},
		{290, "{{-- {{ 1 + 1 }} inside --}}", ""},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, tc.id)
	}
}

func TestCannotUseOperatorError(t *testing.T) {
	cases := []struct {
		id    uint
		inp   string
		left  value.ValueType
		op    string
		right value.ValueType
	}{
		// Int + Float (both directions)
		{10, "{{ 3 + 2.0 }}", value.INT_VAL, "+", value.FLOAT_VAL},
		{20, "{{ 2.0 + 3 }}", value.FLOAT_VAL, "+", value.INT_VAL},
		{30, "{{ 5 - 1.5 }}", value.INT_VAL, "-", value.FLOAT_VAL},
		{40, "{{ 4.5 - 2 }}", value.FLOAT_VAL, "-", value.INT_VAL},
		{50, "{{ 3 * 2.5 }}", value.INT_VAL, "*", value.FLOAT_VAL},
		{60, "{{ 2.5 * 4 }}", value.FLOAT_VAL, "*", value.INT_VAL},
		{70, "{{ 10 / 2.5 }}", value.INT_VAL, "/", value.FLOAT_VAL},
		{80, "{{ 7.5 / 3 }}", value.FLOAT_VAL, "/", value.INT_VAL},
		// Arithmetic with strings
		{90, "{{ 5 * 'x' }}", value.INT_VAL, "*", value.STR_VAL},
		{100, "{{ 'x' * 3 }}", value.STR_VAL, "*", value.INT_VAL},
		{110, "{{ 'x' / 2 }}", value.STR_VAL, "/", value.INT_VAL},
		{120, "{{ 2 / 'x' }}", value.INT_VAL, "/", value.STR_VAL},
		{130, "{{ 'a' - 2 }}", value.STR_VAL, "-", value.INT_VAL},
		{140, "{{ 'a' / 3.0 }}", value.STR_VAL, "/", value.FLOAT_VAL},
		// Arithmetic with booleans
		{150, "{{ true + 5 }}", value.BOOL_VAL, "+", value.INT_VAL},
		{160, "{{ 5 - true }}", value.INT_VAL, "-", value.BOOL_VAL},
		{170, "{{ false - 2.0 }}", value.BOOL_VAL, "-", value.FLOAT_VAL},
		{180, "{{ 2.0 + true }}", value.FLOAT_VAL, "+", value.BOOL_VAL},
		{190, "{{ true * 1 }}", value.BOOL_VAL, "*", value.INT_VAL},
		{200, "{{ true / 2 }}", value.BOOL_VAL, "/", value.INT_VAL},
		// Boolean with strings
		{210, "{{ true + 'str' }}", value.BOOL_VAL, "+", value.STR_VAL},
		{220, "{{ 'str' - false }}", value.STR_VAL, "-", value.BOOL_VAL},
		// Modulo operator
		{230, "{{ 5 % 2.0 }}", value.INT_VAL, "%", value.FLOAT_VAL},
		{240, "{{ 5.0 % 2 }}", value.FLOAT_VAL, "%", value.INT_VAL},
		{250, "{{ 'a' % 2 }}", value.STR_VAL, "%", value.INT_VAL},
		{260, "{{ true % 2 }}", value.BOOL_VAL, "%", value.INT_VAL},
		// Array/Object operations
		{270, "{{ [] * {} }}", value.ARR_VAL, "*", value.OBJ_VAL},
		{280, "{{ {} / [] }}", value.OBJ_VAL, "/", value.ARR_VAL},
		{290, "{{ [] + {} }}", value.ARR_VAL, "+", value.OBJ_VAL},
		{300, "{{ {} - [] }}", value.OBJ_VAL, "-", value.ARR_VAL},
		{310, "{{ [] * 1 }}", value.ARR_VAL, "*", value.INT_VAL},
		{320, "{{ {} / 1 }}", value.OBJ_VAL, "/", value.INT_VAL},
		// Int/Object operations
		{330, "{{ 3 + {} }}", value.INT_VAL, "+", value.OBJ_VAL},
		{340, "{{ {} - 3 }}", value.OBJ_VAL, "-", value.INT_VAL},
		{350, "{{ 5 * {} }}", value.INT_VAL, "*", value.OBJ_VAL},
		// Array with arithmetic
		{360, "{{ [] + 5 }}", value.ARR_VAL, "+", value.INT_VAL},
		{370, "{{ 10 - [] }}", value.INT_VAL, "-", value.ARR_VAL},
		{380, "{{ [] * 3 }}", value.ARR_VAL, "*", value.INT_VAL},
		// Float with arrays/objects
		{390, "{{ 3.14 + [] }}", value.FLOAT_VAL, "+", value.ARR_VAL},
		{400, "{{ {} / 2.5 }}", value.OBJ_VAL, "/", value.FLOAT_VAL},
		// String with arrays/objects
		{410, "{{ 'x' - [] }}", value.STR_VAL, "-", value.ARR_VAL},
		{420, "{{ {} - 'x' }}", value.OBJ_VAL, "-", value.STR_VAL},
		{430, "{{ 'str' + [] }}", value.STR_VAL, "+", value.ARR_VAL},
		{440, "{{ {} + 'str' }}", value.OBJ_VAL, "+", value.STR_VAL},
		// Strings
		{450, "{{ 'a' - 'b' }}", value.STR_VAL, "-", value.STR_VAL},
		{460, "{{ 'a' * 'b' }}", value.STR_VAL, "*", value.STR_VAL},
		{470, "{{ 'a' / 'b' }}", value.STR_VAL, "/", value.STR_VAL},
		{480, "{{ 'a' < 'b' }}", value.STR_VAL, "<", value.STR_VAL},
		{490, "{{ 'a' > 'b' }}", value.STR_VAL, ">", value.STR_VAL},
		{500, "{{ 'a' <= 'b' }}", value.STR_VAL, "<=", value.STR_VAL},
		{510, "{{ 'a' >= 'b' }}", value.STR_VAL, ">=", value.STR_VAL},
		{520, "{{ 'a' % 'b' }}", value.STR_VAL, "%", value.STR_VAL},
		// String + number (addition not supported for strings)
		{530, "{{ 'test' + 5 }}", value.STR_VAL, "+", value.INT_VAL},
		{540, "{{ 5 + 'test' }}", value.INT_VAL, "+", value.STR_VAL},
		{550, "{{ 'test' + 5.5 }}", value.STR_VAL, "+", value.FLOAT_VAL},
		{560, "{{ 5.5 + 'test' }}", value.FLOAT_VAL, "+", value.STR_VAL},
		// Boolean - string
		{570, "{{ true - 'str' }}", value.BOOL_VAL, "-", value.STR_VAL},
		{580, "{{ 'str' - true }}", value.STR_VAL, "-", value.BOOL_VAL},
		// Float modulo operations
		{590, "{{ 3.14 % 2.0 }}", value.FLOAT_VAL, "%", value.FLOAT_VAL},
		{600, "{{ 3.14 % 2 }}", value.FLOAT_VAL, "%", value.INT_VAL},
		// Array / array
		{610, "{{ [] / [] }}", value.ARR_VAL, "/", value.ARR_VAL},
		// Object / object
		{620, "{{ {} / {} }}", value.OBJ_VAL, "/", value.OBJ_VAL},
		// Nil with arithmetic operators
		{630, "{{ nil + 1 }}", value.NIL_VAL, "+", value.INT_VAL},
		{640, "{{ nil - 1 }}", value.NIL_VAL, "-", value.INT_VAL},
		{650, "{{ nil * 1 }}", value.NIL_VAL, "*", value.INT_VAL},
		{660, "{{ nil / 1 }}", value.NIL_VAL, "/", value.INT_VAL},
	}

	for _, tc := range cases {
		evaluated, failure := testEval(tc.inp)
		if failure != nil {
			t.Fatalf("Case: %d. evaluation failed: %s", tc.id, failure)
		}

		err, ok := evaluated.(*value.Error)
		if !ok {
			t.Fatalf("Case: %d. Evaluation failed, got error %q", tc.id, err)
		}

		expect := fail.New(
			nil,
			"/path/to/file",
			fail.OriginEval,
			fail.ErrCannotUseOperator,
			tc.op,
			tc.left,
			tc.op,
			tc.right,
		)
		if err.String() != expect.String() {
			t.Fatalf("Case: %d. Error message must be:\n%q\ngot:\n%q", tc.id, expect, err)
		}
	}
}

func TestEvalFormatDateGlobalFunc(t *testing.T) {
	cases := []struct {
		id     uint
		inp    string
		expect string
	}{
		// Test date time
		{10, "{{ d = '2026-03-24 13:03:25'; formatDate(d, '15:04:05') }}", "13:03:25"},
		{
			20,
			"{{ d = '2027-12-13 10:00:09'; formatDate(d, '2006-01-02 15:04:05') }}",
			"2027-12-13 10:00:09",
		},
		{30, "{{ d = '2023-10-12 01:10:25'; formatDate(d, '2006-01-02') }}", "2023-10-12"},
		// Test date
		{40, "{{ d = '2026-03-24'; formatDate(d, '2006-01-02') }}", "2026-03-24"},
		{50, "{{ d = '2026-03-24'; formatDate(d, '15:04:05') }}", "00:00:00"},
		{60, "{{ d = '2026-03-24'; formatDate(d, '2006-01-02 15:04:05') }}", "2026-03-24 00:00:00"},
		// Test time
		{70, "{{ t = '10:00:09'; formatDate(t, '2006-01-02') }}", "0000-01-01"},
		{80, "{{ t = '10:00:09'; formatDate(t, '15:04:05') }}", "10:00:09"},
		{90, "{{ t = '10:00:09'; formatDate(t, '2006-01-02 15:04:05') }}", "0000-01-01 10:00:09"},
		// Edge cases empty string
		{100, "{{ t = ''; formatDate(t, '2006-01-02 15:04:05') }}", ""},
		// Edge cases midnight
		{110, "{{ d = '2026-01-01 00:00:00'; formatDate(d, '15:04:05') }}", "00:00:00"},
		{120, "{{ d = '2026-01-01 00:00:00'; formatDate(d, '2006-01-02 15:04:05') }}", "2026-01-01 00:00:00"},
		{130, "{{ d = '2026-01-01 00:00:00'; formatDate(d, '2006-01-02') }}", "2026-01-01"},
		// Edge cases end of day
		{140, "{{ d = '2026-12-31 23:59:59'; formatDate(d, '15:04:05') }}", "23:59:59"},
		{150, "{{ d = '2026-12-31 23:59:59'; formatDate(d, '2006-01-02 15:04:05') }}", "2026-12-31 23:59:59"},
		{160, "{{ d = '2026-12-31 23:59:59'; formatDate(d, '2006-01-02') }}", "2026-12-31"},
		// Leap year date
		{170, "{{ d = '2024-02-29 12:30:45'; formatDate(d, '2006-01-02') }}", "2024-02-29"},
		{180, "{{ d = '2024-02-29 12:30:45'; formatDate(d, '15:04:05') }}", "12:30:45"},
		{190, "{{ d = '2024-02-29 12:30:45'; formatDate(d, '2006-01-02 15:04:05') }}", "2024-02-29 12:30:45"},
		// Single-digit month/day padding
		{200, "{{ d = '2026-03-05 01:02:03'; formatDate(d, '2006-01-02') }}", "2026-03-05"},
		{210, "{{ d = '2026-03-05 01:02:03'; formatDate(d, '15:04:05') }}", "01:02:03"},
		{220, "{{ d = '2026-03-05 01:02:03'; formatDate(d, '2006-01-02 15:04:05') }}", "2026-03-05 01:02:03"},
		// Early morning times
		{230, "{{ d = '2026-06-15 01:23:45'; formatDate(d, '15:04:05') }}", "01:23:45"},
		{240, "{{ d = '2026-06-15 02:00:00'; formatDate(d, '2006-01-02 15:04:05') }}", "2026-06-15 02:00:00"},
		// Afternoon/PM times
		{250, "{{ d = '2026-08-20 13:45:30'; formatDate(d, '15:04:05') }}", "13:45:30"},
		{260, "{{ d = '2026-08-20 15:30:00'; formatDate(d, '2006-01-02 15:04:05') }}", "2026-08-20 15:30:00"},
		{270, "{{ d = '2026-08-20 23:00:00'; formatDate(d, '15:04:05') }}", "23:00:00"},
		// Date-only inputs with various formats
		{280, "{{ d = '2000-01-01'; formatDate(d, '2006-01-02') }}", "2000-01-01"},
		{290, "{{ d = '1999-12-31'; formatDate(d, '2006-01-02 15:04:05') }}", "1999-12-31 00:00:00"},
		// Time-only inputs with various formats
		{300, "{{ t = '00:00:00'; formatDate(t, '15:04:05') }}", "00:00:00"},
		{310, "{{ t = '12:34:56'; formatDate(t, '2006-01-02 15:04:05') }}", "0000-01-01 12:34:56"},
		{320, "{{ t = '23:59:59'; formatDate(t, '15:04:05') }}", "23:59:59"},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, tc.id)
	}
}
