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
	p := parser.New(l)
	prog := p.ParseProgram()
	env := object.NewEnv()

	return Eval(prog, env)
}

func evaluationExpected(t *testing.T, inp, expect string) {
	evaluated := testEval(inp)
	errObj, ok := evaluated.(*object.Error)

	if ok {
		t.Errorf("evaluation failed: %s", errObj.Message)
	}

	result := evaluated.String()

	if result != expect {
		t.Errorf("result is not %q, got %q", expect, result)
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
		{`{{ 4.-- }}`, "3"},
		{"{{ 5.11 }}", "5.11"},
		{"{{ -12.3 }}", "-12.3"},
		{`{{ 2.123 + 1.111 }}`, "3.234"},
		{`{{ 2. + 1.2 }}`, "3.2"},
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
			t.Errorf("evaluation failed: %s", errObj.Message)
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
		{"{{ `` }}", ""},
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

func TestEvalIfStatement(t *testing.T) {
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
			t.Errorf("evaluation failed: %s", errObj.Message)
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

func TestEvalVariableDeclaration(t *testing.T) {
	tests := []struct {
		inp      string
		expected string
	}{
		{`{{ var age = 18 }}`, ""},
		{`{{ age := 18 }}`, ""},
		{`{{ var age = 18; age }}`, "18"},
		{`{{ var myAge = 33; var herAge = 25; myAge + herAge }}`, "58"},
		{`{{ var age = 18; age + age }}`, "36"},
		{`{{ var herName = "Anna"; herName }}`, "Anna"},
		{`{{ age := 18; age }}`, "18"},
		{`{{ age := 18; age + 2 }}`, "20"},
		{`{{ age := 18; age + age }}`, "36"},
		{`{{ herName := "Anna"; herName }}`, "Anna"},
		{`{{ she := "Anna"; var me = "Serhii"; she + " " + me }}`, "Anna Serhii"},
		{`{{ var names = ["Anna", "Serhii"] }}`, ""},
		{`{{ var names = ["Anna", "Serhii"]; names }}`, "Anna, Serhii"},
	}

	for _, tt := range tests {
		evaluationExpected(t, tt.inp, tt.expected)
	}
}
