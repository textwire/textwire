package evaluator

import (
	"strings"
	"testing"

	"github.com/textwire/textwire/v3/config"
	"github.com/textwire/textwire/v3/pkg/fail"
	"github.com/textwire/textwire/v3/pkg/file"
	"github.com/textwire/textwire/v3/pkg/lexer"
	"github.com/textwire/textwire/v3/pkg/object"
	"github.com/textwire/textwire/v3/pkg/parser"
)

func testEval(inp string) (object.Object, *fail.Error) {
	l := lexer.New(inp)
	p := parser.New(l, file.New("file", "to/file", "/path/to/file", nil))
	prog := p.ParseProgram()
	scope := object.NewScope()

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

	errObj, ok := evaluated.(*object.Error)
	if ok {
		t.Fatalf("Case: %d. evaluation failed: %s", idx, errObj)
	}

	res := evaluated.String()
	if res != expect {
		t.Fatalf("Case: %d. Result is not '%s', got '%s'", idx, expect, res)
	}
}

func TestEvalHTML(t *testing.T) {
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
		{150, `{{ 10++ }}`, "11"},
		{160, `{{ 10-- }}`, "9"},
		{170, `{{ 3++ + 2-- }}`, "5"},
		{180, `{{ 3-- + 2-- * 3++ + (4--) }}`, "9"},
		// Float
		{190, `{{ 4.4++ }}`, "5.4"},
		{200, `{{ 4.4-- }}`, "3.4"},
		{210, `{{ 4.0-- }}`, "3.0"},
		{220, "{{ 5.11 }}", "5.11"},
		{230, "{{ -12.3 }}", "-12.3"},
		{240, `{{ 2.123 + 1.111 }}`, "3.234"},
		{250, `{{ 2.0 + 1.2 }}`, "3.2"},
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
		// Booleans
		{10, "{{ true }}", "1"},
		{20, "{{ false }}", "0"},
		{30, "{{ !true }}", "0"},
		{40, "{{ !false }}", "1"},
		{50, "{{ !nil }}", "1"},
		{60, "{{ !!true }}", "1"},
		{70, "{{ !!false }}", "0"},
		{80, "{{ true && true }}", "1"},
		{90, "{{ !false && !false }}", "1"},
		{100, "{{ false && false }}", "0"},
		{110, "{{ false && !false }}", "0"},
		// Logical OR
		{120, "{{ true || false }}", "1"},
		{130, "{{ false || true }}", "1"},
		{140, "{{ false || false }}", "0"},
		{150, "{{ false || false || true }}", "1"},
		{160, "{{ '' || '3' }}", "1"},
		{170, "{{ 3 || 0 }}", "1"},
		{180, "{{ [] || [] }}", "0"},
		{190, "{{ {} || {} }}", "0"},
		{200, "{{ 0.2 || 2.3 }}", "1"},
		{210, "{{ 'a' || 'b' }}", "1"},
		{220, "{{ nil || nil }}", "0"},
		// Logical AND
		{230, "{{ false && false || true }}", "1"},
		{240, "{{ true && false || false }}", "0"},
		{250, "{{ false && (false || true) }}", "0"},
		{260, "{{ 3 && 0 }}", "0"},
		{270, "{{ [] && [] }}", "0"},
		{280, "{{ {} && {} }}", "0"},
		{290, "{{ 0.2 && 2.3 }}", "1"},
		{300, "{{ 'a' && 'b' }}", "1"},
		{310, "{{ '' && '' }}", "0"},
		{320, "{{ nil && nil }}", "0"},
		// Integers - equality operators
		{330, "{{ 1 == 1 }}", "1"},
		{335, "{{ 0 == 0 }}", "1"},
		{340, "{{ 1 == 2 }}", "0"},
		{345, "{{ -1 == -1 }}", "1"},
		{346, "{{ -1 == 1 }}", "0"},
		{347, "{{ 100 == 100 }}", "1"},
		{348, "{{ 100 == 99 }}", "0"},
		// Integers - inequality operators
		{350, "{{ 1 != 1 }}", "0"},
		{355, "{{ 0 != 0 }}", "0"},
		{360, "{{ 1 != 2 }}", "1"},
		{365, "{{ -1 != 1 }}", "1"},
		{366, "{{ 100 != 100 }}", "0"},
		{367, "{{ 100 != 99 }}", "1"},
		// Integers - less than
		{370, "{{ 1 < 2 }}", "1"},
		{375, "{{ 0 < 1 }}", "1"},
		{380, "{{ -1 < 0 }}", "1"},
		{385, "{{ 1 < 1 }}", "0"},
		{386, "{{ 2 < 1 }}", "0"},
		{387, "{{ -2 < -3 }}", "0"},
		{388, "{{ 100 < 200 }}", "1"},
		// Integers - greater than
		{390, "{{ 1 > 2 }}", "0"},
		{395, "{{ 0 > 1 }}", "0"},
		{400, "{{ -1 > 0 }}", "0"},
		{405, "{{ 2 > 1 }}", "1"},
		{406, "{{ 0 > -1 }}", "1"},
		{407, "{{ 1 > 1 }}", "0"},
		{408, "{{ 200 > 100 }}", "1"},
		// Integers - less than or equal
		{410, "{{ 1 <= 2 }}", "1"},
		{415, "{{ 1 <= 1 }}", "1"},
		{420, "{{ 0 <= 0 }}", "1"},
		{425, "{{ 2 <= 1 }}", "0"},
		{426, "{{ -1 <= 0 }}", "1"},
		{427, "{{ -1 <= -1 }}", "1"},
		{428, "{{ 100 <= 50 }}", "0"},
		// Integers - greater than or equal
		{430, "{{ 1 >= 2 }}", "0"},
		{435, "{{ 1 >= 1 }}", "1"},
		{440, "{{ 0 >= 0 }}", "1"},
		{445, "{{ 2 >= 1 }}", "1"},
		{446, "{{ 0 >= -1 }}", "1"},
		{447, "{{ -1 >= -1 }}", "1"},
		{448, "{{ 50 >= 100 }}", "0"},
		// Integers with negative numbers
		{450, "{{ -5 == -5 }}", "1"},
		{455, "{{ -5 != -3 }}", "1"},
		{460, "{{ -10 < -5 }}", "1"},
		{465, "{{ -10 > -5 }}", "0"},
		{470, "{{ -5 <= -5 }}", "1"},
		{475, "{{ -5 >= -10 }}", "1"},
		// Floats - equality operators
		{480, "{{ 1.1 == 1.1 }}", "1"},
		{485, "{{ 0.0 == 0.0 }}", "1"},
		{490, "{{ 1.1 == 2.1 }}", "0"},
		{495, "{{ -1.5 == -1.5 }}", "1"},
		{496, "{{ -1.5 == 1.5 }}", "0"},
		{497, "{{ 3.14159 == 3.14159 }}", "1"},
		{498, "{{ 3.14 == 3.15 }}", "0"},
		// Floats - inequality operators
		{500, "{{ 1.1 != 1.1 }}", "0"},
		{505, "{{ 0.0 != 0.0 }}", "0"},
		{510, "{{ 1.1 != 2.1 }}", "1"},
		{515, "{{ -1.5 != 1.5 }}", "1"},
		{516, "{{ 3.14 != 3.141 }}", "1"},
		// Floats - less than
		{520, "{{ 1.1 < 2.1 }}", "1"},
		{525, "{{ 0.0 < 0.1 }}", "1"},
		{530, "{{ -1.5 < 0.0 }}", "1"},
		{535, "{{ 1.1 < 1.1 }}", "0"},
		{536, "{{ 2.5 < 1.5 }}", "0"},
		{537, "{{ -2.5 < -3.5 }}", "0"},
		// Floats - greater than
		{540, "{{ 1.1 > 2.1 }}", "0"},
		{545, "{{ 0.0 > 1.0 }}", "0"},
		{550, "{{ -1.0 > 0.0 }}", "0"},
		{555, "{{ 2.1 > 1.1 }}", "1"},
		{556, "{{ 0.5 > -0.5 }}", "1"},
		{557, "{{ 1.1 > 1.1 }}", "0"},
		// Floats - less than or equal
		{560, "{{ 1.1 <= 2.1 }}", "1"},
		{565, "{{ 1.1 <= 1.1 }}", "1"},
		{570, "{{ 0.0 <= 0.0 }}", "1"},
		{575, "{{ 2.1 <= 1.1 }}", "0"},
		{576, "{{ -1.0 <= 0.0 }}", "1"},
		{577, "{{ -1.0 <= -1.0 }}", "1"},
		// Floats - greater than or equal
		{580, "{{ 1.1 >= 2.1 }}", "0"},
		{585, "{{ 1.1 >= 1.1 }}", "1"},
		{590, "{{ 0.0 >= 0.0 }}", "1"},
		{595, "{{ 2.1 >= 1.1 }}", "1"},
		{596, "{{ 0.0 >= -1.0 }}", "1"},
		{597, "{{ -1.0 >= -1.0 }}", "1"},
		// Floats with negative numbers
		{600, "{{ -5.5 == -5.5 }}", "1"},
		{605, "{{ -5.5 != -3.5 }}", "1"},
		{610, "{{ -10.5 < -5.5 }}", "1"},
		{615, "{{ -10.5 > -5.5 }}", "0"},
		{620, "{{ -5.5 <= -5.5 }}", "1"},
		{625, "{{ -5.5 >= -10.5 }}", "1"},
		// Strings - equality operators
		{630, "{{ 'hello' == 'hello' }}", "1"},
		{635, "{{ '' == '' }}", "1"},
		{640, "{{ 'hello' == 'world' }}", "0"},
		{645, "{{ 'abc' == 'ABC' }}", "0"},
		{650, "{{ ' test ' == ' test ' }}", "1"},
		{655, "{{ 'a' == 'ab' }}", "0"},
		// Strings - inequality operators
		{660, "{{ 'hello' != 'hello' }}", "0"},
		{665, "{{ '' != '' }}", "0"},
		{670, "{{ 'hello' != 'world' }}", "1"},
		{675, "{{ 'abc' != 'ABC' }}", "1"},
		{680, "{{ 'a' != 'ab' }}", "1"},
		// Strings - empty vs non-empty
		{805, "{{ '' == 'a' }}", "0"},
		{810, "{{ '' != 'a' }}", "1"},
		// Strings - numbers as strings
		{815, "{{ '10' == '10' }}", "1"},
		{820, "{{ '10' == '2' }}", "0"},
		// Strings - special characters
		{835, "{{ 'hello world' == 'hello world' }}", "1"},
		{840, "{{ 'test\nline' == 'test\nline' }}", "1"},
		// Booleans - equality operators
		{850, "{{ true == true }}", "1"},
		{855, "{{ false == false }}", "1"},
		{860, "{{ true == false }}", "0"},
		{865, "{{ false == true }}", "0"},
		// Booleans - inequality operators
		{870, "{{ true != true }}", "0"},
		{875, "{{ false != false }}", "0"},
		{880, "{{ true != false }}", "1"},
		{885, "{{ false != true }}", "1"},
		// Nils
		{980, "{{ true == nil }}", "0"},
		{985, "{{ nil == true }}", "0"},
		{990, "{{ false == nil }}", "0"},
		{995, "{{ nil == false }}", "0"},
		{1000, "{{ nil == 0 }}", "0"},
		{1005, "{{ 0 == nil }}", "0"},
		// Nil with integers
		{1010, "{{ nil == 1 }}", "0"},
		{1015, "{{ 1 == nil }}", "0"},
		{1020, "{{ nil != 0 }}", "1"},
		{1025, "{{ 0 != nil }}", "1"},
		{1030, "{{ nil != 5 }}", "1"},
		{1035, "{{ 5 != nil }}", "1"},
		// Nil with floats
		{1040, "{{ nil == 0.0 }}", "0"},
		{1045, "{{ 0.0 == nil }}", "0"},
		{1050, "{{ nil == 1.5 }}", "0"},
		{1055, "{{ 1.5 == nil }}", "0"},
		{1060, "{{ nil != 0.0 }}", "1"},
		{1065, "{{ 0.0 != nil }}", "1"},
		{1070, "{{ nil != 3.14 }}", "1"},
		{1075, "{{ 3.14 != nil }}", "1"},
		// Nil with strings
		{1080, "{{ nil == '' }}", "0"},
		{1085, "{{ '' == nil }}", "0"},
		{1090, "{{ nil == 'test' }}", "0"},
		{1095, "{{ 'test' == nil }}", "0"},
		{1100, "{{ nil != '' }}", "1"},
		{1105, "{{ '' != nil }}", "1"},
		{1110, "{{ nil != 'hello' }}", "1"},
		{1115, "{{ 'hello' != nil }}", "1"},
		// Nil with arrays
		{1120, "{{ nil == [] }}", "0"},
		{1125, "{{ [] == nil }}", "0"},
		{1130, "{{ nil == [1, 2] }}", "0"},
		{1135, "{{ [1, 2] == nil }}", "0"},
		{1140, "{{ nil != [] }}", "1"},
		{1145, "{{ [] != nil }}", "1"},
		{1150, "{{ nil != [1] }}", "1"},
		{1155, "{{ [1] != nil }}", "1"},
		// Nil with objects
		{1160, "{{ nil == {} }}", "0"},
		{1165, "{{ {} == nil }}", "0"},
		{1170, "{{ nil == {name: 'test'} }}", "0"},
		{1175, "{{ {name: 'test'} == nil }}", "0"},
		{1180, "{{ nil != {} }}", "1"},
		{1185, "{{ {} != nil }}", "1"},
		{1190, "{{ nil != {x: 1} }}", "1"},
		{1195, "{{ {x: 1} != nil }}", "1"},
		// Nil with nil
		{1200, "{{ nil == nil }}", "1"},
		{1205, "{{ nil != nil }}", "0"},
	}

	for _, tc := range cases {
		evaluated, failure := testEval(tc.inp)
		if failure != nil {
			t.Fatalf("Case: %d. evaluation failed: %s", tc.id, failure)
		}

		err, ok := evaluated.(*object.Error)
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
		{10, `{{ "Hello World" }}`, "Hello World"},
		{20, `<div {{ 'data-attr="Test"' }}></div>`, `<div data-attr="Test"></div>`},
		{30, `<div {{ "data-attr='Test'" }}></div>`, `<div data-attr='Test'></div>`},
		{40, `{{ "She \"is\" pretty" }}`, `She "is" pretty`},
		{50, `{{ "Korotchaeva" + " " + "Anna" }}`, "Korotchaeva Anna"},
		{60, `{{ "She" + " " + "is" + " " + "nice" }}`, "She is nice"},
		{70, "{{ '' }}", ""},
		{80, `{{ "<h1>Test</h1>" }}`, "&lt;h1&gt;Test&lt;/h1&gt;"},
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
		{10, `{{ true ? "Yes" : "No" }}`, "Yes"},
		{20, `{{ false ? "Yes" : "No" }}`, "No"},
		{30, `{{ nil ? "Yes" : "No" }}`, "No"},
		{40, `{{ 1 ? "Yes" : "No" }}`, "Yes"},
		{50, `{{ 0 ? "Yes" : "No" }}`, "No"},
		{60, `{{ "" ? "Yes" : "No" }}`, "No"},
		{70, `{{ !true ? "Yes" : "No" }}`, "No"},
		{80, `{{ !false ? "Yes" : "No" }}`, "Yes"},
		{90, `{{ !!true ? 1 : 0 }}`, "1"},
		{100, `{{ !!false ? 1 : 0 }}`, "0"},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, tc.id)
	}
}

func TestEvalIfStmt(t *testing.T) {
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
		{171, `@if("".len() > 0)Non empty@elseEmpty@end`, "Empty"},
		{172, `@if("x".len() > 0)Non empty@elseEmpty@end`, "Non empty"},
		// Truthy/falsy values
		{180, `@if(1)Yes@end`, "Yes"},
		{190, `@if(0)Yes@end`, ""},
		{200, `@if("")Yes@end`, ""},
		{210, `@if([])Yes@end`, ""},
		{220, `@if({})Yes@end`, ""},
		{230, `@if(nil)Yes@end`, ""},
		// Logical operators
		{240, `@if(true && true)Yes@end`, "Yes"},
		{250, `@if(true && false)Yes@end`, ""},
		{260, `@if(false || true)Yes@end`, "Yes"},
		{270, `@if(false || false)Yes@end`, ""},
		// Boolean negation
		{280, `@if(!true)Yes@end`, ""},
		{290, `@if(!false)Yes@end`, "Yes"},
		{300, `@if(!!true)Yes@end`, "Yes"},
		{310, `@if(!!false)Yes@end`, ""},
		// Comparison operators
		{320, `@if(1 == 1)Yes@end`, "Yes"},
		{330, `@if(1 == 2)Yes@end`, ""},
		{340, `@if(1 != 2)Yes@end`, "Yes"},
		{350, `@if(1 < 2)Yes@end`, "Yes"},
		{360, `@if(2 > 1)Yes@end`, "Yes"},
		{361, `@if(2 >= 2)Yes@end`, "Yes"},
		{362, `@if(2 <= 2)Yes@end`, "Yes"},
		// Expression results as conditions
		{370, `@if(1 + 1 == 2)Yes@end`, "Yes"},
		{380, `@if(5 - 3 == 1)Yes@end`, ""},
		{381, `@if(2 - 2)No@elseYes@end`, "Yes"},
		{382, `@if(-1 + 2)Yes@elseNo@end`, "Yes"},
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

		err, ok := evaluated.(*object.Error)
		if ok {
			t.Errorf("Case: %d. Evaluation failed: %s", tc.id, err)
		}

		if res := strings.TrimSpace(evaluated.String()); res != tc.expect {
			t.Errorf("Case: %d. Result is not %q, got %q", tc.id, tc.expect, res)
		}
	}
}

func TestEvalArray(t *testing.T) {
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
		{20, `{{ age = 18 }}`, ""},
		{30, `{{ age = 18; age }}`, "18"},
		{40, `{{ age = -5; age }}`, "-5"},
		{50, `{{ age = 0; age }}`, "0"},
		// Float assignment
		{60, `{{ pi = 3.14; pi }}`, "3.14"},
		{70, `{{ price = 0.0; price }}`, "0.0"},
		{80, `{{ negative = -2.5; negative }}`, "-2.5"},
		{81, `{{ f = 0.0; f = f + 0.5; f }}`, "0.5"},
		// String assignment
		{90, `{{ name = "Anna"; name }}`, "Anna"},
		{100, `{{ empty = ""; empty }}`, ""},
		{110, `{{ quote = "He said \"Hello\""; quote }}`, `He said "Hello"`},
		// Boolean assignment
		{120, `{{ flag = true; flag }}`, "1"},
		{130, `{{ flag = false; flag }}`, "0"},
		// Nil assignment
		{140, `{{ nothing = nil; nothing }}`, ""},
		// Multiple assignments
		{200, `{{ a = 1; b = 2; c = 3; a + b + c }}`, "6"},
		{210, `{{ x = "Hello"; y = "World"; x + " " + y }}`, "Hello World"},
		// Reassignment
		{220, `{{ age = 18; age = 25; age }}`, "25"},
		{230, `{{ name = "Anna"; name = "Maria"; name }}`, "Maria"},
		{240, `{{ x = 1; x = x + 1; x = x + 1; x }}`, "3"},
		// Assignment with expression
		{250, `{{ sum = 5 + 3; sum }}`, "8"},
		{260, `{{ calc = 10 * 2 - 5; calc }}`, "15"},
		{270, `{{ result = (2 + 3) * 4; result }}`, "20"},
		// Assignment with string concatenation
		{280, `{{ full = "John" + " " + "Doe"; full }}`, "John Doe"},
		{290, `{{ msg = "Count: " + 5.str(); msg }}`, "Count: 5"},
		// Assignment with ternary
		{300, `{{ val = true ? 1 : 0; val }}`, "1"},
		{310, `{{ val = false ? "yes" : "no"; val }}`, "no"},
		// Assignment with method call
		{320, `{{ upper = "hello".upper(); upper }}`, "HELLO"},
		{330, `{{ len = [1, 2, 3].len(); len }}`, "3"},
		// Assignment with index
		{340, `{{ arr = [10, 20, 30]; val = arr[1]; val }}`, "20"},
		{350, `{{ str = "abc"; val = str.at(0); val }}`, "a"},
		// Assignment in conditional
		{360, `@if(true){{ x = 5 }}{{ x }}@end`, "5"},
		{370, `@if(false){{ x = 5 }}{{ x }}@end`, ""},
		// Assignment with loop variable
		{380, `@each(n in [1, 2, 3]){{ x = n }}{{ x }}@end`, "123"},
		// Object assignment
		{381, `{{ user = {name: "John"}; user.name }}`, "John"},
		{382, `{{ data = {"age": 25}; data.age }}`, "25"},
		{390, `{{ user = {"name": "Ann"}; user = {"name": "Anna"}; user.name }}`, "Anna"},
		{400, `{{ user = {}; user.name = "Anna"; user.name }}`, "Anna"},
		{410, `{{ user = {}; user.name = "Anna"; user.name }}`, "Anna"},
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
		{440, `{{ names = ["Anna", "Serhii"]; names }}`, "Anna, Serhii"},
		{450, `{{ empty = []; empty }}`, ""},
		{460, `{{ nested = [[1, 2], [3, 4]]; nested }}`, "1, 2, 3, 4"},
		{470, `{{ names = ['Serhii', 'Nastya']; names[1] = 'Anna'; names }}`, "Serhii, Anna"},
		{480, `{{ nums = [10, 20, 30]; nums[0] = 1; nums[1] = 2; nums }}`, "1, 2, 30"},
		{490, `{{ x = [[[20]]]; x[0][0][0] = 30; x }}`, "30"},
		{500, `{{ x = ['1', ['2', ['3', ['4']]]]; x[1][1][1][0] = '5'; x }}`, "1, 2, 3, 5"},
		{501, `{{ x = [0]; newVal = 10; x[0] = newVal; x[0] }}`, "10"},
		// Mixed assignment
		{510, `{{ x = [{ name: 'Chiori' }]; x[0].name = 'Mavuika'; x[0].name }}`, "Mavuika"},
		{520, `{{ name = 'Mavuika'; x = [{ name: 'Chiori' }]; x[0].name = name; x[0].name }}`, "Mavuika"},
		// Index assignment edge cases
		{540, `{{ arr = [1, 2]; arr[0] = {name: 'x'}; arr[0].name }}`, "x"},
		{550, `{{ arr = [1, 2]; arr[0] = [3, 4]; arr[0][0] }}`, "3"},
		{560, `{{ arr = [1, 2, 3]; arr[1 + 1] = 5; arr[2] }}`, "5"},
		{570, `{{ arr = [[1], [2]]; arr[0][0] = 10; arr[0][0] }}`, "10"},
		{580, `{{ x = 5; arr = [1, 2]; arr[0] = x; arr[0] }}`, "5"},
		{590, `{{ arr = [1, 2]; arr[0] = arr[1]; arr[0] }}`, "2"},
		{600, `{{ arr = [1, 2, 3]; arr[2] = 10; arr[2] }}`, "10"},
		{610, `{{ arr = ['a', 'b', 'c']; arr[0] = arr[0] + 'x'; arr[0] }}`, "ax"},
		{620, `{{ arr = [true, false]; arr[0] = false; arr[0] }}`, "0"},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, tc.id)
	}
}

func TestEvalForStmt(t *testing.T) {
	cases := []struct {
		id     uint
		inp    string
		expect string
	}{
		// Basic for loops
		{10, `@for(i = 0; i < 2; i++){{ i }}@end`, "01"},
		{20, `@for(i = 1; i <= 3; i++){{ i }}@end`, "123"},
		{30, `@for(i = 5; i > 0; i--){{ i }}@end`, "54321"},
		{40, `@for(i = 0; i < 4; i+2){{ i }}@end`, "02"},
		// Empty loop header parts
		{50, `@for(; false;)Here@end`, ""},
		{51, `@for(; true;)x@break@end`, "x"},
		{60, `@for(c = 1; false; c++){{ c }}@end`, ""},
		{70, `@for(i = 0; i < 0; i++){{ i }}@end`, ""},
		{71, `@for(;;){{ 1 }}@break@end`, "1"},
		// Single iteration
		{80, `@for(c = 1; c == 1; c++){{ c }}@end`, "1"},
		{90, `@for(i = 0; i < 1; i++){{ i }}@end`, "0"},
		// @else directive
		{110, `@for(c = 1; false; c++){{ c }}@else@end`, ""},
		{120, `@for(c = 1; false; c++){{ c }}@else<b>Empty</b>@end`, "<b>Empty</b>"},
		{130, `@for(c = 0; c < 0; c++){{ c }}@elseEmpty@end`, "Empty"},
		// @break directive
		{140, `@for(i = 1; i <= 3; i++){{ i }}@break@end`, "1"},
		{150, `@for(i = 1; i <= 3; i++)@break{{ i }}@end`, ""},
		{160, `@for(i = 1; i <= 3; i++)@if(i == 3)@break@end{{ i }}@end`, "12"},
		{170, `@for(i = 0; i < 10; i++)@break@end`, ""},
		// @continue directive
		{180, `@for(i = 1; i <= 3; i++)@continue{{ i }}@end`, ""},
		{190, `@for(i = 1; i <= 3; i++){{ i }}@continue@end`, "123"},
		{200, `@for(i = 1; i <= 3; i++)@if(i == 2)@continue@end{{ i }}@end`, "13"},
		{210, `@for(i = 1; i <= 5; i++)@if(i % 2 == 0)@continue@end{{ i }}@end`, "135"},
		// @breakif directive
		{220, `@for(i = 1; i <= 3; i++)@breakif(i == 3){{ i }}@end`, "12"},
		{230, `@for(i = 1; i <= 3; i++)@breakif(i == 2){{ i }}@end`, "1"},
		{240, `@for(i = 1; i <= 10; i++)@breakif(i > 5){{ i }}@end`, "12345"},
		// @continueif directive
		{250, `@for(i = 1; i <= 3; i++)@continueif(i == 3){{ i }}@end`, "12"},
		{260, `@for(i = 1; i <= 3; i++)@continueif(i == 2){{ i }}@end`, "13"},
		{270, `@for(i = 1; i <= 5; i++)@continueif(i % 2 == 0){{ i }}@end`, "135"},
		// Nested for loops
		{280, `@for(i = 0; i < 2; i++)@for(j = 0; j < 2; j++){{ i }}{{ j }}@end@end`, "00011011"},
		{290, `@for(i = 1; i <= 2; i++)@for(j = 1; j <= 2; j++){{ i * j }}@end@end`, "1224"},
		// For loop with HTML
		{
			300,
			`<ul>@for(i = 1; i <= 3; i++)<li>{{ i }}</li>@end</ul>`,
			"<ul><li>1</li><li>2</li><li>3</li></ul>",
		},
		// Variable modification in loop
		{310, `{{ sum = 0 }}@for(i = 1; i <= 5; i++){{ sum = sum + i }}@end{{ sum }}`, "15"},
		{320, `{{ count = 0 }}@for(i = 0; i < 3; i++){{ count = count + 1 }}@end{{ count }}`, "3"},
		// Float iteration
		{330, `@for(f = 0.0; f < 1.0; f + 0.5){{ f }}@end`, "0.00.5"},
		{340, `@for(f = 0.0; f < 2.0; f + 1.0){{ f }}@end`, "0.01.0"},
		// Multiple statements in loop body
		{350, `@for(i = 0; i < 3; i++){{ i }};{{ i * 2 }}@end`, "0;01;22;4"},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, tc.id)
	}
}

func TestEvalEachStmt(t *testing.T) {
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
		{220, `@each(n in [1, 2, 3, 4, 5])@continueif(n == 3){{ n }}@end`, "1245"},
		{
			230,
			`@each(n in ["ann", "serhii", "sam"])@continueif(n == 'sam'){{ n }} @end`,
			"ann serhii ",
		},
		// support continueIf and breakIf
		{240, `@each(n in [1, 2])@continueIf(n == 2){{ n }}@end`, "1"},
		{250, `@each(n in [1, 2, 3])@breakIf(n == 2){{ n }}@end`, "1"},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, tc.id)
	}
}

func TestEvalObjectLiteral(t *testing.T) {
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
		{60, `{{ user = {"father": {"name": {"first": "Serhii"}}}; user.father.name.first }}`, "Serhii"},
		{70, `{{ u = {"father": {name: {"first": "Serhii",},},}; u['father']['name'].first }}`, "Serhii"},

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

func TestTypeMismatchError(t *testing.T) {
	cases := []struct {
		id   uint
		inp  string
		objL object.ObjectType
		op   string
		objR object.ObjectType
	}{
		// Int + Float (both directions)
		{10, "{{ 3 + 2.0 }}", object.INT_OBJ, "+", object.FLOAT_OBJ},
		{20, "{{ 2.0 + 3 }}", object.FLOAT_OBJ, "+", object.INT_OBJ},
		{30, "{{ 5 - 1.5 }}", object.INT_OBJ, "-", object.FLOAT_OBJ},
		{40, "{{ 4.5 - 2 }}", object.FLOAT_OBJ, "-", object.INT_OBJ},
		{50, "{{ 3 * 2.5 }}", object.INT_OBJ, "*", object.FLOAT_OBJ},
		{60, "{{ 2.5 * 4 }}", object.FLOAT_OBJ, "*", object.INT_OBJ},
		{70, "{{ 10 / 2.5 }}", object.INT_OBJ, "/", object.FLOAT_OBJ},
		{80, "{{ 7.5 / 3 }}", object.FLOAT_OBJ, "/", object.INT_OBJ},
		// Arithmetic with strings
		{130, "{{ 5 * 'x' }}", object.INT_OBJ, "*", object.STR_OBJ},
		{140, "{{ 'x' * 3 }}", object.STR_OBJ, "*", object.INT_OBJ},
		{150, "{{ 'x' / 2 }}", object.STR_OBJ, "/", object.INT_OBJ},
		{160, "{{ 2 / 'x' }}", object.INT_OBJ, "/", object.STR_OBJ},
		{170, "{{ 'a' - 2 }}", object.STR_OBJ, "-", object.INT_OBJ},
		{180, "{{ 'a' / 3.0 }}", object.STR_OBJ, "/", object.FLOAT_OBJ},
		// Arithmetic with booleans
		{190, "{{ true + 5 }}", object.BOOL_OBJ, "+", object.INT_OBJ},
		{200, "{{ 5 - true }}", object.INT_OBJ, "-", object.BOOL_OBJ},
		{210, "{{ false - 2.0 }}", object.BOOL_OBJ, "-", object.FLOAT_OBJ},
		{220, "{{ 2.0 + true }}", object.FLOAT_OBJ, "+", object.BOOL_OBJ},
		{230, "{{ true * 1 }}", object.BOOL_OBJ, "*", object.INT_OBJ},
		{240, "{{ true / 2 }}", object.BOOL_OBJ, "/", object.INT_OBJ},
		// Boolean with strings
		{250, "{{ true + 'str' }}", object.BOOL_OBJ, "+", object.STR_OBJ},
		{260, "{{ 'str' - false }}", object.STR_OBJ, "-", object.BOOL_OBJ},
		// Modulo operator
		{390, "{{ 5 % 2.0 }}", object.INT_OBJ, "%", object.FLOAT_OBJ},
		{400, "{{ 5.0 % 2 }}", object.FLOAT_OBJ, "%", object.INT_OBJ},
		{410, "{{ 'a' % 2 }}", object.STR_OBJ, "%", object.INT_OBJ},
		{420, "{{ true % 2 }}", object.BOOL_OBJ, "%", object.INT_OBJ},
		// Mixed type comparisons with == and !=
		{480, "{{ 1 == '1' }}", object.INT_OBJ, "==", object.STR_OBJ},
		{490, "{{ '1' == 1 }}", object.STR_OBJ, "==", object.INT_OBJ},
		{500, "{{ true == 1 }}", object.BOOL_OBJ, "==", object.INT_OBJ},
		{510, "{{ 1 == true }}", object.INT_OBJ, "==", object.BOOL_OBJ},
		{520, "{{ false == 0 }}", object.BOOL_OBJ, "==", object.INT_OBJ},
		{530, "{{ 0 == false }}", object.INT_OBJ, "==", object.BOOL_OBJ},
		{540, "{{ 1.0 == 1 }}", object.FLOAT_OBJ, "==", object.INT_OBJ},
		{550, "{{ 1 == 1.0 }}", object.INT_OBJ, "==", object.FLOAT_OBJ},
	}

	for _, tc := range cases {
		evaluated, failure := testEval(tc.inp)
		if failure != nil {
			t.Fatalf("Case: %d. evaluation failed: %s", tc.id, failure)
		}

		err, ok := evaluated.(*object.Error)
		if !ok {
			t.Fatalf("Case: %d. Evaluation failed, got error %q", tc.id, err)
		}

		expect := fail.New(
			1,
			"/path/to/file",
			"evaluator",
			fail.ErrTypeMismatch,
			tc.objL,
			tc.op,
			tc.objR,
		)
		if err.String() != expect.String() {
			t.Fatalf("Case: %d. Error message is not %q, got %q", tc.id, expect, err)
		}
	}
}

func TestNotSupportedTypeError(t *testing.T) {
	cases := []struct {
		id  uint
		inp string
		op  string
		t   object.ObjectType
	}{
		// Array/Object operations
		{270, "{{ [] * {} }}", "*", object.ARR_OBJ},
		{280, "{{ {} / [] }}", "/", object.OBJ_OBJ},
		{290, "{{ [] + {} }}", "+", object.ARR_OBJ},
		{300, "{{ {} - [] }}", "-", object.OBJ_OBJ},
		{310, "{{ [] * 1 }}", "*", object.ARR_OBJ},
		{320, "{{ {} / 1 }}", "/", object.OBJ_OBJ},
		// Int/Object operations
		{330, "{{ 3 + {} }}", "+", object.OBJ_OBJ},
		{340, "{{ {} - 3 }}", "-", object.OBJ_OBJ},
		{350, "{{ 5 * {} }}", "*", object.OBJ_OBJ},
		// Array with arithmetic
		{360, "{{ [] + 5 }}", "+", object.ARR_OBJ},
		{370, "{{ 10 - [] }}", "-", object.ARR_OBJ},
		{380, "{{ [] * 3 }}", "*", object.ARR_OBJ},
		// Float with arrays/objects
		{460, "{{ 3.14 + [] }}", "+", object.ARR_OBJ},
		{470, "{{ {} / 2.5 }}", "/", object.FLOAT_OBJ},
		// String with arrays/objects
		{90, "{{ 'x' - [] }}", "-", object.ARR_OBJ},
		{100, "{{ {} - 'x' }}", "-", object.STR_OBJ},
		{110, "{{ 'str' + [] }}", "+", object.ARR_OBJ},
		{120, "{{ {} + 'str' }}", "+", object.STR_OBJ},
		// Mixed type comparisons with == and !=
		{560, "{{ [] == [] }}", "==", object.ARR_OBJ},
		{570, "{{ [1, 2] == [1, 2] }}", "==", object.ARR_OBJ},
		{580, "{{ [1] == [2] }}", "==", object.ARR_OBJ},
		{590, "{{ [] != [] }}", "!=", object.ARR_OBJ},
		{600, "{{ [1] != [1] }}", "!=", object.ARR_OBJ},
		{610, "{{ {} == {} }}", "==", object.OBJ_OBJ},
		{620, "{{ {a: 1} == {a: 1} }}", "==", object.OBJ_OBJ},
		{630, "{{ {a: 1} == {a: 2} }}", "==", object.OBJ_OBJ},
		{640, "{{ {} != {} }}", "!=", object.OBJ_OBJ},
		{650, "{{ {a: 1} != {a: 1} }}", "!=", object.OBJ_OBJ},
	}

	for _, tc := range cases {
		evaluated, failure := testEval(tc.inp)
		if failure != nil {
			t.Fatalf("Case: %d. evaluation failed: %s", tc.id, failure)
		}

		err, ok := evaluated.(*object.Error)
		if !ok {
			t.Fatalf("Case: %d. Evaluation failed, got error %q", tc.id, err)
		}

		expect := fail.New(
			1,
			"/path/to/file",
			"evaluator",
			fail.ErrUnknownTypeForOp,
			tc.t,
			tc.op,
		)
		if err.String() != expect.String() {
			t.Fatalf("Case: %d. Error message is not %q, got %q", tc.id, expect, err)
		}
	}
}
