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

func evaluationExpected(t *testing.T, inp, expect string, idx int) {
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
		id     int
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
		id     int
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
		id     int
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
		// Ints
		{330, "{{ 1 == 1 }}", "1"},
		{340, "{{ 1 == 2 }}", "0"},
		{350, "{{ 1 != 1 }}", "0"},
		{360, "{{ 1 != 2 }}", "1"},
		{370, "{{ 1 < 2 }}", "1"},
		{380, "{{ 1 > 2 }}", "0"},
		{390, "{{ 1 <= 2 }}", "1"},
		{400, "{{ 1 >= 2 }}", "0"},
		// Floats
		{410, "{{ 1.1 == 1.1 }}", "1"},
		{420, "{{ 1.1 == 2.1 }}", "0"},
		{430, "{{ 1.1 != 1.1 }}", "0"},
		{440, "{{ 1.1 != 2.1 }}", "1"},
		{450, "{{ 1.1 < 2.1 }}", "1"},
		{460, "{{ 1.1 > 2.1 }}", "0"},
		{470, "{{ 1.1 <= 2.1 }}", "1"},
		{480, "{{ 1.1 >= 2.1 }}", "0"},
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
		id     int
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
		id     int
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
		id     int
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
		id     int
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
		id     int
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

func TestEvalAssignVariable(t *testing.T) {
	cases := []struct {
		id     int
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
		{71, `{{ f = 0.0; f = f + 0.5; f }}`, "0.5"},
		// String assignment
		{80, `{{ name = "Anna"; name }}`, "Anna"},
		{90, `{{ empty = ""; empty }}`, ""},
		{100, `{{ quote = "He said \"Hello\""; quote }}`, `He said "Hello"`},
		// Boolean assignment
		{110, `{{ flag = true; flag }}`, "1"},
		{120, `{{ flag = false; flag }}`, "0"},
		// Nil assignment
		{130, `{{ nothing = nil; nothing }}`, ""},
		// Array assignment
		{140, `{{ names = ["Anna", "Serhii"]; names }}`, "Anna, Serhii"},
		{150, `{{ empty = []; empty }}`, ""},
		{160, `{{ nested = [[1, 2], [3, 4]]; nested }}`, "1, 2, 3, 4"},
		// Object assignment
		{170, `{{ user = {name: "John"}; user.name }}`, "John"},
		{180, `{{ data = {"age": 25}; data.age }}`, "25"},
		// Multiple assignments
		{190, `{{ a = 1; b = 2; c = 3; a + b + c }}`, "6"},
		{200, `{{ x = "Hello"; y = "World"; x + " " + y }}`, "Hello World"},
		// Reassignment
		{210, `{{ age = 18; age = 25; age }}`, "25"},
		{220, `{{ name = "Anna"; name = "Maria"; name }}`, "Maria"},
		{230, `{{ x = 1; x = x + 1; x = x + 1; x }}`, "3"},
		// Assignment with expression
		{240, `{{ sum = 5 + 3; sum }}`, "8"},
		{250, `{{ calc = 10 * 2 - 5; calc }}`, "15"},
		{260, `{{ result = (2 + 3) * 4; result }}`, "20"},
		// Assignment with string concatenation
		{270, `{{ full = "John" + " " + "Doe"; full }}`, "John Doe"},
		{280, `{{ msg = "Count: " + 5.str(); msg }}`, "Count: 5"},
		// Assignment with ternary
		{290, `{{ val = true ? 1 : 0; val }}`, "1"},
		{300, `{{ val = false ? "yes" : "no"; val }}`, "no"},
		// Assignment with method call
		{310, `{{ upper = "hello".upper(); upper }}`, "HELLO"},
		{320, `{{ len = [1, 2, 3].len(); len }}`, "3"},
		// Assignment with index
		{330, `{{ arr = [10, 20, 30]; val = arr[1]; val }}`, "20"},
		{340, `{{ str = "abc"; val = str.at(0); val }}`, "a"},
		// Assignment in conditional
		{350, `@if(true){{ x = 5 }}{{ x }}@end`, "5"},
		{360, `@if(false){{ x = 5 }}{{ x }}@end`, ""},
		// Assignment with loop variable
		{370, `@each(n in [1, 2, 3]){{ x = n }}{{ x }}@end`, "123"},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, tc.id)
	}
}

func TestEvalForStmt(t *testing.T) {
	cases := []struct {
		id     int
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
		id     int
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
		id     int
		inp    string
		expect string
	}{
		{10, `{{ {"name": "John"}['name'] }}`, "John"},
		{20, `{{ {"name": "John"}.name }}`, "John"},
		{30, `{{ obj = {name: "John"}; obj.name }}`, "John"},
		{40, `{{ o = {"name": "John", "age": 22}; o.age }}`, "22"},
		{50, `{{ user = {"father": {"name": "John"}}; user.father.name }}`, "John"},
		{60, `{{ user = {"father": {"name": {"first": "Sam"}}}; user.father.name.first }}`, "Sam"},
		{70, `{{ u = {"father": {name: {"first": "Sam",},},}; u['father']['name'].first }}`, "Sam"},
		{80, `{{ name = "Sam"; age = 12; obj = { name, age }; obj.name }}`, "Sam"},
		{90, `{{ name = "Sam"; age = 12; obj = { name, age }; obj.age }}`, "12"},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, tc.id)
	}
}

func TestEvalComments(t *testing.T) {
	cases := []struct {
		id     int
		inp    string
		expect string
	}{
		{10, "{{-- This is a comment --}}", ""},
		{20, "<section>{{-- This is a comment --}}</section>", "<section></section>"},
		{30, "Some {{-- --}}text", "Some text"},
		{40, "{{-- @each(u in users){{ u }}@end --}}", ""},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, tc.id)
	}
}

func TestTypeMismatchErrors(t *testing.T) {
	cases := []struct {
		id   int
		inp  string
		objL object.ObjectType
		op   string
		objR object.ObjectType
	}{
		{10, "{{ 3 + 2.0 }}", object.INT_OBJ, "+", object.FLOAT_OBJ},
		{20, "{{ 2.0 + 3 }}", object.FLOAT_OBJ, "+", object.INT_OBJ},
		{30, "{{ 'x' - [] }}", object.STR_OBJ, "-", object.ARR_OBJ},
		{40, "{{ {} - 'x' }}", object.OBJ_OBJ, "-", object.STR_OBJ},
		{50, "{{ 5 * 'x' }}", object.INT_OBJ, "*", object.STR_OBJ},
		{60, "{{ 'x' * 3 }}", object.STR_OBJ, "*", object.INT_OBJ},
		{70, "{{ 'x' / 2 }}", object.STR_OBJ, "/", object.INT_OBJ},
		{80, "{{ 2 / 'x' }}", object.INT_OBJ, "/", object.STR_OBJ},
		{90, "{{ true + 5 }}", object.BOOL_OBJ, "+", object.INT_OBJ},
		{100, "{{ 5 - true }}", object.INT_OBJ, "-", object.BOOL_OBJ},
		{110, "{{ false - 2.0 }}", object.BOOL_OBJ, "-", object.FLOAT_OBJ},
		{120, "{{ 2.0 + true }}", object.FLOAT_OBJ, "+", object.BOOL_OBJ},
		{130, "{{ [] * {} }}", object.ARR_OBJ, "*", object.OBJ_OBJ},
		{140, "{{ {} / [] }}", object.OBJ_OBJ, "/", object.ARR_OBJ},
		{150, "{{ [] + {} }}", object.ARR_OBJ, "+", object.OBJ_OBJ},
		{160, "{{ 3 + {} }}", object.INT_OBJ, "+", object.OBJ_OBJ},
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
