package evaluator

import (
	"strings"
	"testing"

	"github.com/textwire/textwire/v2/ctx"
	"github.com/textwire/textwire/v2/fail"
	"github.com/textwire/textwire/v2/lexer"
	"github.com/textwire/textwire/v2/object"
	"github.com/textwire/textwire/v2/parser"
)

func testEval(inp string) object.Object {
	l := lexer.New(inp)
	p := parser.New(l, "")
	prog := p.ParseProgram()
	env := object.NewEnv()

	eval := New(&ctx.EvalCtx{
		AbsPath: "/path/to/file",
	})

	return eval.Eval(prog, env)
}

func evaluationExpected(t *testing.T, inp, expect string) {
	evaluated := testEval(inp)
	errObj, ok := evaluated.(*object.Error)

	if ok {
		t.Fatalf("evaluation failed: %s", errObj.String())
	}

	result := evaluated.String()

	if result != expect {
		t.Fatalf("result is not '%s', got '%s'", expect, result)
	}
}

func TestEvalHTML(t *testing.T) {
	tests := []struct {
		inp      string
		expected string
	}{
		{"<h1>Hello World</h1>", "<h1>Hello World</h1>"},
		{"<ul><li><span>Email: anna@protonmail.com</span></li></ul>",
			"<ul><li><span>Email: anna@protonmail.com</span></li></ul>"},
		{"<b>Nice</b>@foo", "<b>Nice</b>@foo"},
		{`<h1>\@continue</h1>`, "<h1>@continue</h1>"},
		{`<h1>@\@break</h1>`, "<h1>@@break</h1>"},
		{`<h1>@@@\@break</h1>`, "<h1>@@@@break</h1>"},
		{`\@`, `\@`},
		{`\\@`, `\\@`},
		{`\@if(true)`, `@if(true)`},
		{`\\@if(true)`, `\@if(true)`},
		{`\{{ 5 }}`, `{{ 5 }}`},
		{`\\{{ "nice" }}`, `\{{ "nice" }}`},
		{`\\\{{ x }}`, `\\{{ x }}`},
	}

	for _, tc := range tests {
		evaluationExpected(t, tc.inp, tc.expected)
	}
}

func TestEvalNumericExp(t *testing.T) {
	tests := []struct {
		inp      string
		expected string
	}{
		{"{{ 5; 5 }}", "55"},
		{"{{ 5 }}", "5"},
		{"{{ 10 }}", "10"},
		{"{{ -123 }}", "-123"},
		{`{{ 5 + 5 }}`, "10"},
		{`{{ 5 - 5 }}`, "0"},
		{`{{ 20 / 2 }}`, "10"},
		{`{{ 23 * 2 }}`, "46"},
		{`{{ 11 + 13 - 1 }}`, "23"},
		{"{{ 2 * (5 + 10) }}", "30"},
		{`{{ (3 + 5) * 2 }}`, "16"},
		{`{{ 3 * 3 * 3 + 10 }}`, "37"},
		{`{{ (5 + 10 * 2 + 15 / 3) * 2 + -10 }}`, "50"},
		{`{{ ((5 + 10) * ((2 + 15) / 3) + 2) }}`, "77"},
		{`{{ 10++ }}`, "11"},
		{`{{ 10-- }}`, "9"},
		{`{{ 3++ + 2-- }}`, "5"},
		{`{{ 3-- + 2-- * 3++ + (4--) }}`, "9"},
		// Float
		{`{{ 4.4++ }}`, "5.4"},
		{`{{ 4.4-- }}`, "3.4"},
		{`{{ 4.0-- }}`, "3.0"},
		{"{{ 5.11 }}", "5.11"},
		{"{{ -12.3 }}", "-12.3"},
		{`{{ 2.123 + 1.111 }}`, "3.234"},
		{`{{ 2.0 + 1.2 }}`, "3.2"},
	}

	for _, tc := range tests {
		evaluationExpected(t, tc.inp, tc.expected)
	}
}

func TestEvalBooleanExp(t *testing.T) {
	tests := []struct {
		inp      string
		expected string
	}{
		// Booleans
		{"{{ true }}", "1"},
		{"{{ false }}", "0"},
		{"{{ !true }}", "0"},
		{"{{ !false }}", "1"},
		{"{{ !nil }}", "1"},
		{"{{ !!true }}", "1"},
		{"{{ !!false }}", "0"},
		// Ints
		{`{{ 1 == 1 }}`, "1"},
		{`{{ 1 == 2 }}`, "0"},
		{`{{ 1 != 1 }}`, "0"},
		{`{{ 1 != 2 }}`, "1"},
		{`{{ 1 < 2 }}`, "1"},
		{`{{ 1 > 2 }}`, "0"},
		{`{{ 1 <= 2 }}`, "1"},
		{`{{ 1 >= 2 }}`, "0"},
		// Floats
		{`{{ 1.1 == 1.1 }}`, "1"},
		{`{{ 1.1 == 2.1 }}`, "0"},
		{`{{ 1.1 != 1.1 }}`, "0"},
		{`{{ 1.1 != 2.1 }}`, "1"},
		{`{{ 1.1 < 2.1 }}`, "1"},
		{`{{ 1.1 > 2.1 }}`, "0"},
		{`{{ 1.1 <= 2.1 }}`, "1"},
		{`{{ 1.1 >= 2.1 }}`, "0"},
	}

	for _, tc := range tests {
		evaluated := testEval(tc.inp)

		errObj, ok := evaluated.(*object.Error)

		if ok {
			t.Errorf("evaluation failed: %s", errObj.String())
			return
		}

		result := evaluated.String()

		if result != tc.expected {
			t.Errorf("result is not %s, got %s", tc.expected, result)
		}
	}
}

func TestEvalNilExp(t *testing.T) {
	inp := "<h1>{{ nil }}</h1>"
	evaluationExpected(t, inp, "<h1></h1>")
}

func TestEvalStringExp(t *testing.T) {
	tests := []struct {
		inp      string
		expected string
	}{
		{`{{ "Hello World" }}`, "Hello World"},
		{`<div {{ 'data-attr="Test"' }}></div>`, `<div data-attr="Test"></div>`},
		{`<div {{ "data-attr='Test'" }}></div>`, `<div data-attr='Test'></div>`},
		{`{{ "She \"is\" pretty" }}`, `She "is" pretty`},
		{`{{ "Korotchaeva" + " " + "Anna" }}`, "Korotchaeva Anna"},
		{`{{ "She" + " " + "is" + " " + "nice" }}`, "She is nice"},
		{"{{ '' }}", ""},
		{`{{ "<h1>Test</h1>" }}`, "&lt;h1&gt;Test&lt;/h1&gt;"},
	}

	for _, tc := range tests {
		evaluationExpected(t, tc.inp, tc.expected)
	}
}

func TestEvalTernaryExp(t *testing.T) {
	tests := []struct {
		inp      string
		expected string
	}{
		{`{{ true ? "Yes" : "No" }}`, "Yes"},
		{`{{ false ? "Yes" : "No" }}`, "No"},
		{`{{ nil ? "Yes" : "No" }}`, "No"},
		{`{{ 1 ? "Yes" : "No" }}`, "Yes"},
		{`{{ 0 ? "Yes" : "No" }}`, "No"},
		{`{{ "" ? "Yes" : "No" }}`, "No"},
		{`{{ !true ? "Yes" : "No" }}`, "No"},
		{`{{ !false ? "Yes" : "No" }}`, "Yes"},
		{`{{ !!true ? 1 : 0 }}`, "1"},
		{`{{ !!false ? 1 : 0 }}`, "0"},
	}

	for _, tc := range tests {
		evaluationExpected(t, tc.inp, tc.expected)
	}
}

func TestEvalIfStmt(t *testing.T) {
	tests := []struct {
		inp    string
		expect string
	}{
		{`@if(true)Hello@end`, "Hello"},
		{`@if(false)Hello@end`, ""},
		{`@if(true)Anna@elseif(true)Lili@end`, "Anna"},
		{`@if(false)Alan@elseif(true)Serhii@end`, "Serhii"},
		{`@if(false)Ana De Armaz@elseif(false)David@elseVladimir@end`, "Vladimir"},
		{`@if(false)Will@elseif(false)Daria@elseif(true)Poll@end`, "Poll"},
		{`@if(false)Lara@elseif(true)Susan@elseif(true)Smith@end`, "Susan"},
		{`<h2>@if(true)Hello@end</h2>`, "<h2>Hello</h2>"},
		{`<h2>@if(false)Hello@end</h2>`, "<h2></h2>"},
		{`@if(true)Hello@end`, "Hello"},
		{
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

	for _, tc := range tests {
		evaluated := testEval(tc.inp)
		errObj, ok := evaluated.(*object.Error)

		if ok {
			t.Errorf("evaluation failed: %s", errObj.String())
		}

		result := strings.TrimSpace(evaluated.String())

		if result != tc.expect {
			t.Errorf("result is not %q, got %q", tc.expect, result)
		}
	}
}

func TestEvalArray(t *testing.T) {
	tests := []struct {
		inp      string
		expected string
	}{
		{`{{ [] }}`, ""},
		{`{{ [[[[[]]]]] }}`, ""},
		{`{{ [1, 2, 3] }}`, "1, 2, 3"},
		{`{{ ["Anna", "Serhii" ] }}`, "Anna, Serhii"},
		{`{{ [true, false] }}`, "1, 0"},
		{`{{ [[1, [2]], 3] }}`, "1, 2, 3"},
	}

	for _, tc := range tests {
		evaluationExpected(t, tc.inp, tc.expected)
	}
}

func TestEvalIndexExp(t *testing.T) {
	tests := []struct {
		inp      string
		expected string
	}{
		{`{{ [1, 2, 3][0] }}`, "1"},
		{`{{ [1, 2, 3][1] }}`, "2"},
		{`{{ [1, 2, 3][2] }}`, "3"},
		{`{{ ["Some string"][0] }}`, "Some string"},
		{`{{ [[[11]]][0][0][0] }}`, "11"},
	}

	for _, tc := range tests {
		evaluationExpected(t, tc.inp, tc.expected)
	}
}

func TestEvalAssignVariable(t *testing.T) {
	tests := []struct {
		inp      string
		expected string
	}{
		{`{{ age = 18 }}`, ""},
		{`{{ age = 18; age }}`, "18"},
		{`{{ myAge = 33; herAge = 25; myAge + herAge }}`, "58"},
		{`{{ age = 18; age + age }}`, "36"},
		{`{{ herName = "Anna"; herName }}`, "Anna"},
		{`{{ age = 18; age }}`, "18"},
		{`{{ age = 18; age + 2 }}`, "20"},
		{`{{ age = 18; age + age }}`, "36"},
		{`{{ herName = "Anna"; herName }}`, "Anna"},
		{`{{ she = "Anna"; me = "Serhii"; she + " " + me }}`, "Anna Serhii"},
		{`{{ names = ["Anna", "Serhii"] }}`, ""},
		{`{{ names = ["Anna", "Serhii"]; names }}`, "Anna, Serhii"},
		{`{{ age = 18; age = 2; age }}`, "2"},
		{`{{ city = "Kiev"; city = "Moscow"; city }}`, "Moscow"},
	}

	for _, tc := range tests {
		evaluationExpected(t, tc.inp, tc.expected)
	}
}

func TestEvalForStmt(t *testing.T) {
	tests := []struct {
		inp      string
		expected string
	}{
		{`@for(i = 0; i < 2; i++){{ i }}@end`, "01"},
		{`@for(i = 1; i <= 3; i++){{ i }}@end`, "123"},
		{`@for(i = 0; i < 2; i++){{ i }}@end`, "01"},
		{`@for(; false;)Here@end`, ""},
		{`@for(c = 1; false; c++){{ c }}@end`, ""},
		{`@for(c = 1; c == 1; c++){{ c }}@end`, "1"},
		// test @else directive
		{`@for(c = 1; false; c++){{ c }}@else<b>Empty</b>@end`, "<b>Empty</b>"},
		{`@for(c = 0; c < 0; c++){{ c }}@elseEmpty@end`, "Empty"},
		// test @break directive
		{`@for(i = 1; i <= 3; i++){{ i }}@break@end`, "1"},
		{`@for(i = 1; i <= 3; i++)@break{{ i }}@end`, ""},
		{`@for(i = 1; i <= 3; i++)@if(i == 3)@break@end{{ i }}@end`, "12"},
		// test @continue directive
		{`@for(i = 1; i <= 3; i++)@continue{{ i }}@end`, ""},
		{`@for(i = 1; i <= 3; i++){{ i }}@continue@end`, "123"},
		{`@for(i = 1; i <= 3; i++)@if(i == 2)@continue@end{{ i }}@end`, "13"},
		// test @breakIf directive
		{`@for(i = 1; i <= 3; i++)@breakIf(i == 3){{ i }}@end`, "12"},
		{`@for(i = 1; i <= 3; i++)@breakIf(i == 2){{ i }}@end`, "1"},
		// test @continueIf directive
		{`@for(i = 1; i <= 3; i++)@continueIf(i == 3){{ i }}@end`, "12"},
		{`@for(i = 1; i <= 3; i++)@continueIf(i == 2){{ i }}@end`, "13"},
	}

	for _, tc := range tests {
		evaluationExpected(t, tc.inp, tc.expected)
	}
}

func TestEvalEachStmt(t *testing.T) {
	tests := []struct {
		inp      string
		expected string
	}{
		{`@each(name in ["anna", "serhii"]){{ name }} @end`, "anna serhii "},
		{`@each(num in [1, 2, 3]){{ num }}@end`, "123"},
		{`@each(num in []){{ num }}@end`, ""},
		// test loop variable
		{`@each(num in [43, 12, 53]){{ loop.index }}@end`, "012"},
		{`@each(num in [100]){{ loop.index }}@end`, "0"},
		{`@each(num in [1, 2, 3, 4]){{ loop.first }}@end`, "1000"},
		{`@each(num in [1, 2, 3, 4]){{ loop.last }}@end`, "0001"},
		{`@each(num in [4, 2, 8]){{ loop.iter }}@end`, "123"},
		{`@each(num in [9, 3, 44, 24, 1, 3]){{ loop.iter }}@end`, "123456"},
		// test @else directive
		{`@each(v in []){{ v }}@else<b>Empty array</b>@end`, "<b>Empty array</b>"},
		{`@each(n in []){{ n }}@else@end`, ""},
		{`@each(n in []){{ n }}@elsetest@end`, "test"},
		{`@each(n in [1, 2, 3, 4, 5]){{ n }}@end`, "12345"},
		// test @break directive
		{`@each(n in [1, 2, 3, 4, 5])@break{{ n }}@end`, ""},
		{`@each(n in [1, 2, 3, 4, 5]){{ n }}@break@end`, "1"},
		{`@each(n in [1, 2, 3, 4, 5])@if(n == 3)@break@end{{ n }}@end`, "12"},
		// test @continue directive
		{`@each(n in [1, 2, 3, 4, 5])@continue{{ n }}@end`, ""},
		{`@each(n in [1, 2, 3, 4, 5]){{ n }}@continue@end`, "12345"},
		{`@each(n in [1, 2, 3, 4, 5])@if(n == 3)@continue@end{{ n }}@end`, "1245"},
		// test @breakIf directive
		{`@each(n in [1, 2, 3, 4, 5])@breakIf(n == 3){{ n }}@end`, "12"},
		{`@each(n in ["ann", "serhii", "sam"])@breakIf(n == 'sam'){{ n }} @end`, "ann serhii "},
		// test @continueIf directive
		{`@each(n in [1, 2, 3, 4, 5])@continueIf(n == 3){{ n }}@end`, "1245"},
		{`@each(n in ["ann", "serhii", "sam"])@continueIf(n == 'sam'){{ n }} @end`, "ann serhii "},
	}

	for _, tc := range tests {
		evaluationExpected(t, tc.inp, tc.expected)
	}
}

func TestEvalObjectLiteral(t *testing.T) {
	tests := []struct {
		inp      string
		expected string
	}{
		{`{{ {"name": "John"}['name'] }}`, "John"},
		{`{{ {"name": "John"}.name }}`, "John"},
		{`{{ obj = {name: "John"}; obj.name }}`, "John"},
		{`{{ o = {"name": "John", "age": 22}; o.age }}`, "22"},
		{`{{ user = {"father": {"name": "John"}}; user.father.name }}`, "John"},
		{`{{ user = {"father": {"name": {"first": "Sam"}}}; user.father.name.first }}`, "Sam"},
		{`{{ u = {"father": {name: {"first": "Sam",},},}; u['father']['name'].first }}`, "Sam"},
		{`{{ name = "Sam"; age = 12; obj = { name, age }; obj.name }}`, "Sam"},
		{`{{ name = "Sam"; age = 12; obj = { name, age }; obj.age }}`, "12"},
	}

	for _, tc := range tests {
		evaluationExpected(t, tc.inp, tc.expected)
	}
}

func TestEvalComments(t *testing.T) {
	tests := []struct {
		inp      string
		expected string
	}{
		{"{{-- This is a comment --}}", ""},
		{"<section>{{-- This is a comment --}}</section>", "<section></section>"},
		{"Some {{-- --}}text", "Some text"},
		{"{{-- @each(u in users){{ u }}@end --}}", ""},
	}

	for _, tc := range tests {
		evaluationExpected(t, tc.inp, tc.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	cases := []struct {
		inp string
		err *fail.Error
	}{
		{
			inp: "{{ 5.somefunction() }}",
			err: fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrNoFuncForThisType,
				"somefunction",
				object.INT_OBJ,
			),
		},
		{
			inp: "{{ 3 / 0 }}",
			err: fail.New(
				1,
				"/path/to/file",
				"evaluator",
				fail.ErrDivisionByZero,
			),
		},
	}

	for _, tc := range cases {
		evaluated := testEval(tc.inp)

		errObj, ok := evaluated.(*object.Error)

		if !ok {
			t.Fatalf("evaluation failed: %s", errObj.String())
		}

		expect := tc.err.String()
		if errObj.String() != expect {
			t.Fatalf("error message is not '%s', got '%s'", expect, errObj.String())
		}
	}
}
