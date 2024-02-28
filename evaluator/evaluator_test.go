package evaluator

import (
	"strings"
	"testing"

	"github.com/textwire/textwire/lexer"
	"github.com/textwire/textwire/object"
	"github.com/textwire/textwire/parser"
)

func testEval(inp string) object.Object {
	l := lexer.New(inp)
	p := parser.New(l, "")
	prog := p.ParseProgram()
	env := object.NewEnv()

	eval := New(&EvalContext{
		absPath: "/path/to/file",
	})

	return eval.Eval(prog, env)
}

func evaluationExpected(t *testing.T, inp, expect string) {
	evaluated := testEval(inp)
	errObj, ok := evaluated.(*object.Error)

	if ok {
		t.Errorf("evaluation failed: %s", errObj.String())
	}

	result := evaluated.String()

	if result != expect {
		t.Errorf("result is not '%s', got '%s'", expect, result)
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

	for _, tt := range tests {
		evaluationExpected(t, tt.inp, tt.expected)
	}
}

func TestEvalNumericExpression(t *testing.T) {
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

	for _, tt := range tests {
		evaluationExpected(t, tt.inp, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
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

	for _, tt := range tests {
		evaluated := testEval(tt.inp)

		errObj, ok := evaluated.(*object.Error)

		if ok {
			t.Errorf("evaluation failed: %s", errObj.String())
			return
		}

		result := evaluated.String()

		if result != tt.expected {
			t.Errorf("result is not %s, got %s", tt.expected, result)
		}
	}
}

func TestEvalNilExpression(t *testing.T) {
	inp := "<h1>{{ nil }}</h1>"
	evaluationExpected(t, inp, "<h1></h1>")
}

func TestEvalStringExpression(t *testing.T) {
	tests := []struct {
		inp      string
		expected string
	}{
		{`{{ "Hello World" }}`, "Hello World"},
		{`{{ "She \"is\" pretty" }}`, `She &#34;is&#34; pretty`},
		{`{{ "Korotchaeva" + " " + "Anna" }}`, "Korotchaeva Anna"},
		{`{{ "She" + " " + "is" + " " + "nice" }}`, "She is nice"},
		{"{{ '' }}", ""},
		{`{{ "<h1>Test</h1>" }}`, "&lt;h1&gt;Test&lt;/h1&gt;"},
	}

	for _, tt := range tests {
		evaluationExpected(t, tt.inp, tt.expected)
	}
}

func TestEvalTernaryExpression(t *testing.T) {
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

	for _, tt := range tests {
		evaluationExpected(t, tt.inp, tt.expected)
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

	for _, tt := range tests {
		evaluated := testEval(tt.inp)
		errObj, ok := evaluated.(*object.Error)

		if ok {
			t.Errorf("evaluation failed: %s", errObj.String())
		}

		result := strings.TrimSpace(evaluated.String())

		if result != tt.expect {
			t.Errorf("result is not %q, got %q", tt.expect, result)
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

	for _, tt := range tests {
		evaluationExpected(t, tt.inp, tt.expected)
	}
}

func TestEvalIndexExpression(t *testing.T) {
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

	for _, tt := range tests {
		evaluationExpected(t, tt.inp, tt.expected)
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

	for _, tt := range tests {
		evaluationExpected(t, tt.inp, tt.expected)
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
	}

	for _, tt := range tests {
		evaluationExpected(t, tt.inp, tt.expected)
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
	}

	for _, tt := range tests {
		evaluationExpected(t, tt.inp, tt.expected)
	}
}

func TestEvalObjectLiteral(t *testing.T) {
	tests := []struct {
		inp      string
		expected string
	}{
		{`{{ {"name": "John"}['name'] }}`, "John"},
		{`{{ {"name": "John"}.name }}`, "John"},
		{`{{ obj = {"name": "John"}; obj.name }}`, "John"},
		{`{{ o = {"name": "John", "age": 22}; o.age }}`, "22"},
		{`{{ user = {"father": {"name": "John"}}; user.father.name }}`, "John"},
		{`{{ user = {"father": {"name": {"first": "Sam"}}}; user.father.name.first }}`, "Sam"},
		{`{{ u = {"father": {"name": {"first": "Sam"}}}; u['father']['name'].first }}`, "Sam"},
	}

	for _, tt := range tests {
		evaluationExpected(t, tt.inp, tt.expected)
	}
}
