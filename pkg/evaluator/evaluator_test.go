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

func testEval(inp string) object.Object {
	l := lexer.New(inp)
	p := parser.New(l, file.New("file", "to/file", "/path/to/file", nil))
	prog := p.ParseProgram()
	scope := object.NewScope()

	e := New(&config.Func{}, nil)
	ctx := NewContext(scope, prog.AbsPath)

	return e.Eval(prog, ctx)
}

func evaluationExpected(t *testing.T, inp, expect string, idx int) {
	evaluated := testEval(inp)

	errObj, ok := evaluated.(*object.Error)
	if ok {
		t.Fatalf("Case: %d. evaluation failed: %s", idx, errObj.String())
	}

	res := evaluated.String()
	if res != expect {
		t.Fatalf("Case: %d. Result is not %s, got %s", idx, expect, res)
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
		evaluated := testEval(tc.inp)
		err, ok := evaluated.(*object.Error)
		if ok {
			t.Errorf("Case: %d. Evaluation failed: %s", tc.id, err.String())
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
		{10, `@if(true)Hello@end`, "Hello"},
		{20, `@if(true.binary())Hello@end`, "Hello"},
		{30, `@if(false.binary())Hello@end`, ""},
		{40, `@if(false)Hello@end`, ""},
		{50, `@if(true)Anna@elseif(true)Lili@end`, "Anna"},
		{60, `@if(false)Alan@elseif(true)Serhii@end`, "Serhii"},
		{70, `@if(false)Ana De Armaz@elseif(false)David@elseVladimir@end`, "Vladimir"},
		{80, `@if(false)Will@elseif(false)Daria@elseif(true)Poll@end`, "Poll"},
		{90, `@if(false)Lara@elseif(true)Susan@elseif(true)Smith@end`, "Susan"},
		{100, `<h2>@if(true)Hello@end</h2>`, "<h2>Hello</h2>"},
		{110, `<h2>@if(false)Hello@end</h2>`, "<h2></h2>"},
		{120, `@if(true)Hello@end`, "Hello"},
		{
			130,
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
		evaluated := testEval(tc.inp)
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
		{10, `{{ [] }}`, ""},
		{20, `{{ [[[[[]]]]] }}`, ""},
		{30, `{{ [1, 2, 3] }}`, "1, 2, 3"},
		{40, `{{ ["Anna", "Serhii" ] }}`, "Anna, Serhii"},
		{50, `{{ [true, false] }}`, "1, 0"},
		{60, `{{ [[1, [2]], 3] }}`, "1, 2, 3"},
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
		{10, `{{ [1, 2, 3][0] }}`, "1"},
		{20, `{{ [1, 2, 3][1] }}`, "2"},
		{30, `{{ [1, 2, 3][2] }}`, "3"},
		{40, `{{ ["Some string"][0] }}`, "Some string"},
		{50, `{{ [[[11]]][0][0][0] }}`, "11"},
		{60, `{{ [][2] }}`, ""},
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
		{10, `{{ age = 18 }}`, ""},
		{20, `{{ age = 18; age }}`, "18"},
		{30, `{{ myAge = 33; herAge = 25; myAge + herAge }}`, "58"},
		{40, `{{ age = 18; age + age }}`, "36"},
		{50, `{{ herName = "Anna"; herName }}`, "Anna"},
		{60, `{{ age = 18; age }}`, "18"},
		{70, `{{ age = 18; age + 2 }}`, "20"},
		{80, `{{ age = 18; age + age }}`, "36"},
		{90, `{{ herName = "Anna"; herName }}`, "Anna"},
		{100, `{{ she = "Anna"; me = "Serhii"; she + " " + me }}`, "Anna Serhii"},
		{110, `{{ names = ["Anna", "Serhii"] }}`, ""},
		{120, `{{ names = ["Anna", "Serhii"]; names }}`, "Anna, Serhii"},
		{130, `{{ age = 18; age = 2; age }}`, "2"},
		{140, `{{ city = "Kiev"; city = "Moscow"; city }}`, "Moscow"},
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
		{10, `@for(i = 0; i < 2; i++){{ i }}@end`, "01"},
		{20, `@for(i = 1; i <= 3; i++){{ i }}@end`, "123"},
		{30, `@for(i = 0; i < 2; i++){{ i }}@end`, "01"},
		{40, `@for(; false;)Here@end`, ""},
		{50, `@for(c = 1; false; c++){{ c }}@end`, ""},
		{60, `@for(c = 1; c == 1; c++){{ c }}@end`, "1"},
		// test @else directive
		{70, `@for(c = 1; false; c++){{ c }}@else<b>Empty</b>@end`, "<b>Empty</b>"},
		{80, `@for(c = 0; c < 0; c++){{ c }}@elseEmpty@end`, "Empty"},
		// test @break directive
		{90, `@for(i = 1; i <= 3; i++){{ i }}@break@end`, "1"},
		{100, `@for(i = 1; i <= 3; i++)@break{{ i }}@end`, ""},
		{110, `@for(i = 1; i <= 3; i++)@if(i == 3)@break@end{{ i }}@end`, "12"},
		// test @continue directive
		{120, `@for(i = 1; i <= 3; i++)@continue{{ i }}@end`, ""},
		{130, `@for(i = 1; i <= 3; i++){{ i }}@continue@end`, "123"},
		{140, `@for(i = 1; i <= 3; i++)@if(i == 2)@continue@end{{ i }}@end`, "13"},
		// test @breakif directive
		{150, `@for(i = 1; i <= 3; i++)@breakif(i == 3){{ i }}@end`, "12"},
		{160, `@for(i = 1; i <= 3; i++)@breakif(i == 2){{ i }}@end`, "1"},
		// test @continueif directive
		{170, `@for(i = 1; i <= 3; i++)@continueif(i == 3){{ i }}@end`, "12"},
		{180, `@for(i = 1; i <= 3; i++)@continueif(i == 2){{ i }}@end`, "13"},
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
		{30, "{{ 'nice' - [] }}", object.STR_OBJ, "-", object.ARR_OBJ},
		{40, "{{ {} - 'bad' }}", object.OBJ_OBJ, "-", object.STR_OBJ},
		{50, "{{ 5 * 'bad' }}", object.INT_OBJ, "*", object.STR_OBJ},
		{60, "{{ 'nice' / 2 }}", object.STR_OBJ, "/", object.INT_OBJ},
		{70, "{{ true + 5 }}", object.BOOL_OBJ, "+", object.INT_OBJ},
		{80, "{{ false - 2.0 }}", object.BOOL_OBJ, "-", object.FLOAT_OBJ},
		{90, "{{ [] * {} }}", object.ARR_OBJ, "*", object.OBJ_OBJ},
		{100, "{{ {} / [] }}", object.OBJ_OBJ, "/", object.ARR_OBJ},
	}

	for _, tc := range cases {
		evaluated := testEval(tc.inp)
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
