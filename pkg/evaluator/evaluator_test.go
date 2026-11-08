package evaluator

import (
	"strings"
	"testing"

	"codeberg.org/textwire/textwire/v5/config"
	"codeberg.org/textwire/textwire/v5/pkg/fail"
	"codeberg.org/textwire/textwire/v5/pkg/file"
	"codeberg.org/textwire/textwire/v5/pkg/lexer"
	"codeberg.org/textwire/textwire/v5/pkg/parser"
	"codeberg.org/textwire/textwire/v5/pkg/value"
)

func testEval(inp string, data map[string]any) (value.Value, *fail.Error) {
	l := lexer.New(inp)
	p := parser.New(l, file.New("file", "to/file", "/path/to/file", nil))
	prog := p.ParseProgram()

	scope, failure := value.NewScopeFromMap(data)
	if failure != nil {
		return nil, failure
	}

	if p.HasErrors() {
		return nil, p.Errors()[0]
	}

	e := New(&config.Func{}, nil)
	ctx := NewContext(scope, prog.AbsPath)

	return e.Eval(prog, ctx), nil
}

func evaluationExpected(
	t *testing.T,
	inp, expect string,
	data map[string]any,
	name string,
) {
	evaluated, failure := testEval(inp, data)
	if failure != nil {
		t.Fatalf("Case: %s. evaluation failed: %s", name, failure)
	}

	errObj, ok := evaluated.(*value.Error)
	if ok {
		t.Fatalf("Case: %s. evaluation failed: %s", name, errObj)
	}

	res := evaluated.String()
	if res != expect {
		t.Fatalf("Case: %s. Result is not '%s', got '%s'", name, expect, res)
	}
}

func TestEvalText(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		{"plain HTML heading", "<h1>Hello World</h1>", "<h1>Hello World</h1>"},
		{
			"HTML list with email",
			"<ul><li><span>Email: anna@protonmail.com</span></li></ul>",
			"<ul><li><span>Email: anna@protonmail.com</span></li></ul>",
		},
		{"non-directive at-sign in HTML", "<b>Nice</b>@foo", "<b>Nice</b>@foo"},
		{"escaped @continue", `<h1>\@continue</h1>`, "<h1>@continue</h1>"},
		{"escaped @break after at-sign", `<h1>@\@break</h1>`, "<h1>@@break</h1>"},
		{"multiple at-signs before break", `<h1>@@@\@break</h1>`, "<h1>@@@@break</h1>"},
		{"single escaped at-sign", `\@`, `\@`},
		{"double escaped at-sign", `\\@`, `\\@`},
		{"escaped @if", `\@if(true)`, `@if(true)`},
		{"double escaped @if", `\\@if(true)`, `\@if(true)`},
		{"escaped interpolation", `\{{ 5 }}`, `{{ 5 }}`},
		{"double escaped interpolation", `\\{{ "nice" }}`, `\{{ "nice" }}`},
		{"triple escaped interpolation", `\\\{{ x }}`, `\\{{ x }}`},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, nil, tc.name)
	}
}

func TestEvalNumericExp(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		// Ints
		{"two int statements", "{{ 5; 5 }}", "55"},
		{"int literal", "{{ 5 }}", "5"},
		{"int literal ten", "{{ 10 }}", "10"},
		{"negative int literal", "{{ -123 }}", "-123"},
		{"int addition", `{{ 5 + 5 }}`, "10"},
		{"int subtraction", `{{ 5 - 5 }}`, "0"},
		{"int division", `{{ 20 / 2 }}`, "10"},
		{"int multiplication", `{{ 23 * 2 }}`, "46"},
		{"int addition and subtraction", `{{ 11 + 13 - 1 }}`, "23"},
		{"multiplication after addition", "{{ 2 * (5 + 10) }}", "30"},
		{"addition after multiplication", `{{ (3 + 5) * 2 }}`, "16"},
		{"chained multiplication with addition", `{{ 3 * 3 * 3 + 10 }}`, "37"},
		{"mixed precedence expression", `{{ (5 + 10 * 2 + 15 / 3) * 2 + -10 }}`, "50"},
		{"nested parentheses expression", `{{ ((5 + 10) * ((2 + 15) / 3) + 2) }}`, "77"},
		{"int min literal", "{{ -9223372036854775808 }}", "-9223372036854775808"},
		{"int min plus one", "{{ -9223372036854775808 + 1 }}", "-9223372036854775807"},
		{"int min times one", "{{ -9223372036854775808 * 1 }}", "-9223372036854775808"},
		{"int max overflow", "{{ 9223372036854775807 + 1 }}", "-9223372036854775808"},
		// Floats
		{"float literal", "{{ 5.11 }}", "5.11"},
		{"negative float literal", "{{ -12.3 }}", "-12.3"},
		{"float addition", `{{ 2.123 + 1.111 }}`, "3.234"},
		{"float addition with trailing zero", `{{ 2.0 + 1.2 }}`, "3.2"},
		{"float zero addition", `{{ 0.0 + 0.0 }}`, "0.0"},
		{"float zero literal", `{{ 0.0 }}`, "0.0"},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, nil, tc.name)
	}
}

func TestEvalFloatOverflow(t *testing.T) {
	// 308 nines is just below math.MaxFloat64, so multiplying by 10 overflows.
	huge := strings.Repeat("9", 308) + ".0"

	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		{"positive float overflow", "{{ " + huge + " * 10.0 }}", "+Inf"},
		{"negative float overflow", "{{ -" + huge + " * 10.0 }}", "-Inf"},
		{
			"subtracting infinities yields NaN",
			"{{ (" + huge + " * 10.0) - (" + huge + " * 10.0) }}",
			"NaN",
		},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, nil, tc.name)
	}
}

func TestEvalBooleanExp(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		// --------------------------- Booleans ---------------------------
		{"true literal", "{{ true }}", "1"},
		{"false literal", "{{ false }}", "0"},
		{"not true", "{{ !true }}", "0"},
		{"not false", "{{ !false }}", "1"},
		{"not nil", "{{ !nil }}", "1"},
		{"not grouped conjunction", "{{ !(1 == 1 && 2 == 2) }}", "0"},
		{"not contains call", `{{ !"Hello".contains("World") }}`, "1"},
		{"double not true", "{{ !!true }}", "1"},
		{"double not false", "{{ !!false }}", "0"},
		// --------------------------- Logical ---------------------------
		// Logical OR
		{"true or false", "{{ true || false }}", "1"},
		{"false or true", "{{ false || true }}", "1"},
		{"false or false", "{{ false || false }}", "0"},
		{"false or false or true", "{{ false || false || true }}", "1"},
		{"empty string or non-empty", "{{ '' || '3' }}", "1"},
		{"int or zero", "{{ 3 || 0 }}", "1"},
		{"empty array or empty array", "{{ [] || [] }}", "0"},
		{"empty object or empty object", "{{ {} || {} }}", "0"},
		{"float or float", "{{ 0.2 || 2.3 }}", "1"},
		{"string or string", "{{ 'a' || 'b' }}", "1"},
		{"nil or nil", "{{ nil || nil }}", "0"},
		// Logical AND
		{"false and false or true", "{{ false && false || true }}", "1"},
		{"true and false or false", "{{ true && false || false }}", "0"},
		{"false and grouped false or true", "{{ false && (false || true) }}", "0"},
		{"int and zero", "{{ 3 && 0 }}", "0"},
		{"empty array and empty array", "{{ [] && [] }}", "0"},
		{"empty object and empty object", "{{ {} && {} }}", "0"},
		{"float and float", "{{ 0.2 && 2.3 }}", "1"},
		{"string and string", "{{ 'a' && 'b' }}", "1"},
		{"empty string and empty string", "{{ '' && '' }}", "0"},
		{"nil and nil", "{{ nil && nil }}", "0"},
		{"true and true", "{{ true && true }}", "1"},
		{"not false and not false", "{{ !false && !false }}", "1"},
		{"false and false", "{{ false && false }}", "0"},
		{"false and not false", "{{ false && !false }}", "0"},
		// --------------------------- Integers ---------------------------
		// Integers - equality operators
		{"int equal", "{{ 1 == 1 }}", "1"},
		{"zero equal", "{{ 0 == 0 }}", "1"},
		{"int not equal", "{{ 1 == 2 }}", "0"},
		{"negative int equal", "{{ -1 == -1 }}", "1"},
		{"negative and positive int not equal", "{{ -1 == 1 }}", "0"},
		{"hundred equal", "{{ 100 == 100 }}", "1"},
		{"hundred not equal", "{{ 100 == 99 }}", "0"},
		// Integers - inequality operators
		{"int inequality false", "{{ 1 != 1 }}", "0"},
		{"zero inequality false", "{{ 0 != 0 }}", "0"},
		{"int inequality true", "{{ 1 != 2 }}", "1"},
		{"negative int inequality", "{{ -1 != 1 }}", "1"},
		{"hundred inequality false", "{{ 100 != 100 }}", "0"},
		{"hundred inequality true", "{{ 100 != 99 }}", "1"},
		// Integers - less than
		{"int less than", "{{ 1 < 2 }}", "1"},
		{"zero less than one", "{{ 0 < 1 }}", "1"},
		{"negative less than zero", "{{ -1 < 0 }}", "1"},
		{"equal ints less than false", "{{ 1 < 1 }}", "0"},
		{"greater int less than false", "{{ 2 < 1 }}", "0"},
		{"more negative less than false", "{{ -2 < -3 }}", "0"},
		{"hundred less than two hundred", "{{ 100 < 200 }}", "1"},
		// Integers - greater than
		{"int greater than false", "{{ 1 > 2 }}", "0"},
		{"zero greater than one false", "{{ 0 > 1 }}", "0"},
		{"negative greater than zero false", "{{ -1 > 0 }}", "0"},
		{"int greater than true", "{{ 2 > 1 }}", "1"},
		{"zero greater than negative", "{{ 0 > -1 }}", "1"},
		{"equal ints greater than false", "{{ 1 > 1 }}", "0"},
		{"two hundred greater than hundred", "{{ 200 > 100 }}", "1"},
		// Integers - less than or equal
		{"int less or equal", "{{ 1 <= 2 }}", "1"},
		{"equal ints less or equal", "{{ 1 <= 1 }}", "1"},
		{"zero less or equal zero", "{{ 0 <= 0 }}", "1"},
		{"greater int less or equal false", "{{ 2 <= 1 }}", "0"},
		{"negative less or equal zero", "{{ -1 <= 0 }}", "1"},
		{"equal negative less or equal", "{{ -1 <= -1 }}", "1"},
		{"hundred less or equal fifty false", "{{ 100 <= 50 }}", "0"},
		// Integers - greater than or equal
		{"int greater or equal false", "{{ 1 >= 2 }}", "0"},
		{"equal ints greater or equal", "{{ 1 >= 1 }}", "1"},
		{"zero greater or equal zero", "{{ 0 >= 0 }}", "1"},
		{"int greater or equal true", "{{ 2 >= 1 }}", "1"},
		{"zero greater or equal negative", "{{ 0 >= -1 }}", "1"},
		{"equal negative greater or equal", "{{ -1 >= -1 }}", "1"},
		{"fifty greater or equal hundred false", "{{ 50 >= 100 }}", "0"},
		// Integers with negative numbers
		{"negative equality", "{{ -5 == -5 }}", "1"},
		{"negative inequality", "{{ -5 != -3 }}", "1"},
		{"negative less than", "{{ -10 < -5 }}", "1"},
		{"negative greater than", "{{ -10 > -5 }}", "0"},
		{"negative less or equal", "{{ -5 <= -5 }}", "1"},
		{"negative greater or equal", "{{ -5 >= -10 }}", "1"},
		// --------------------------- Floats ---------------------------
		// Floats - equality operators
		{"float equal", "{{ 1.1 == 1.1 }}", "1"},
		{"zero float equal", "{{ 0.0 == 0.0 }}", "1"},
		{"float not equal", "{{ 1.1 == 2.1 }}", "0"},
		{"negative float equal", "{{ -1.5 == -1.5 }}", "1"},
		{"negative and positive float not equal", "{{ -1.5 == 1.5 }}", "0"},
		{"pi equal", "{{ 3.14159 == 3.14159 }}", "1"},
		{"pi not equal", "{{ 3.14 == 3.15 }}", "0"},
		// Floats - inequality operators
		{"float inequality false", "{{ 1.1 != 1.1 }}", "0"},
		{"zero float inequality false", "{{ 0.0 != 0.0 }}", "0"},
		{"float inequality true", "{{ 1.1 != 2.1 }}", "1"},
		{"negative and positive float inequality", "{{ -1.5 != 1.5 }}", "1"},
		{"float inequality different precision", "{{ 3.14 != 3.141 }}", "1"},
		// Floats - less than
		{"float less than", "{{ 1.1 < 2.1 }}", "1"},
		{"zero less than small float", "{{ 0.0 < 0.1 }}", "1"},
		{"negative float less than zero", "{{ -1.5 < 0.0 }}", "1"},
		{"equal floats less than false", "{{ 1.1 < 1.1 }}", "0"},
		{"greater float less than false", "{{ 2.5 < 1.5 }}", "0"},
		{"more negative float less than false", "{{ -2.5 < -3.5 }}", "0"},
		// Floats - greater than
		{"float greater than false", "{{ 1.1 > 2.1 }}", "0"},
		{"zero float greater than one float false", "{{ 0.0 > 1.0 }}", "0"},
		{"negative float greater than zero float false", "{{ -1.0 > 0.0 }}", "0"},
		{"float greater than true", "{{ 2.1 > 1.1 }}", "1"},
		{"positive greater than negative", "{{ 0.5 > -0.5 }}", "1"},
		{"equal floats greater than false", "{{ 1.1 > 1.1 }}", "0"},
		// Floats - less than or equal
		{"float less or equal", "{{ 1.1 <= 2.1 }}", "1"},
		{"equal floats less or equal", "{{ 1.1 <= 1.1 }}", "1"},
		{"zero float less or equal", "{{ 0.0 <= 0.0 }}", "1"},
		{"greater float less or equal false", "{{ 2.1 <= 1.1 }}", "0"},
		{"negative float less or equal zero", "{{ -1.0 <= 0.0 }}", "1"},
		{"equal negative floats less or equal", "{{ -1.0 <= -1.0 }}", "1"},
		// Floats - greater than or equal
		{"float greater or equal false", "{{ 1.1 >= 2.1 }}", "0"},
		{"equal floats greater or equal", "{{ 1.1 >= 1.1 }}", "1"},
		{"zero float greater or equal", "{{ 0.0 >= 0.0 }}", "1"},
		{"float greater or equal true", "{{ 2.1 >= 1.1 }}", "1"},
		{"zero float greater or equal negative float", "{{ 0.0 >= -1.0 }}", "1"},
		{"equal negative floats greater or equal", "{{ -1.0 >= -1.0 }}", "1"},
		// Floats with negative numbers
		{"negative float equality", "{{ -5.5 == -5.5 }}", "1"},
		{"negative float inequality", "{{ -5.5 != -3.5 }}", "1"},
		{"negative float less than", "{{ -10.5 < -5.5 }}", "1"},
		{"negative float greater than", "{{ -10.5 > -5.5 }}", "0"},
		{"negative float less or equal", "{{ -5.5 <= -5.5 }}", "1"},
		{"negative float greater or equal", "{{ -5.5 >= -10.5 }}", "1"},
		// --------------------------- Strings ---------------------------
		// Strings - equality operators
		{"string equal", "{{ 'hello' == 'hello' }}", "1"},
		{"empty strings equal", "{{ '' == '' }}", "1"},
		{"string not equal", "{{ 'hello' == 'world' }}", "0"},
		{"string case sensitive", "{{ 'abc' == 'ABC' }}", "0"},
		{"whitespace string equal", "{{ ' test ' == ' test ' }}", "1"},
		{"prefix string not equal", "{{ 'a' == 'ab' }}", "0"},
		// Strings - inequality operators
		{"string inequality false", "{{ 'hello' != 'hello' }}", "0"},
		{"empty strings inequality false", "{{ '' != '' }}", "0"},
		{"string inequality true", "{{ 'hello' != 'world' }}", "1"},
		{"string case inequality", "{{ 'abc' != 'ABC' }}", "1"},
		{"string length inequality", "{{ 'a' != 'ab' }}", "1"},
		// Strings - empty vs non-empty
		{"empty string not equal to non-empty", "{{ '' == 'a' }}", "0"},
		{"empty string inequality non-empty", "{{ '' != 'a' }}", "1"},
		// Strings - numbers as strings
		{"numeric strings equal", "{{ '10' == '10' }}", "1"},
		{"numeric strings different", "{{ '10' == '2' }}", "0"},
		// Strings - special characters
		{"string with space equal", "{{ 'hello world' == 'hello world' }}", "1"},
		{"string with newline equal", "{{ 'test\nline' == 'test\nline' }}", "1"},
		// Booleans - equality operators
		{"true equal true", "{{ true == true }}", "1"},
		{"false equal false", "{{ false == false }}", "1"},
		{"true not equal false", "{{ true == false }}", "0"},
		{"false not equal true", "{{ false == true }}", "0"},
		// Booleans - inequality operators
		{"true inequality false", "{{ true != true }}", "0"},
		{"false inequality false", "{{ false != false }}", "0"},
		{"true inequality different", "{{ true != false }}", "1"},
		{"false inequality different", "{{ false != true }}", "1"},
		// --------------------------- Nils ---------------------------
		// Nil with boolean
		{"true not equal nil", "{{ true == nil }}", "0"},
		{"nil not equal true", "{{ nil == true }}", "0"},
		{"false not equal nil", "{{ false == nil }}", "0"},
		{"nil not equal false", "{{ nil == false }}", "0"},
		// Nil with integers
		{"nil not equal zero", "{{ nil == 0 }}", "0"},
		{"zero not equal nil", "{{ 0 == nil }}", "0"},
		{"nil not equal one", "{{ nil == 1 }}", "0"},
		{"one not equal nil", "{{ 1 == nil }}", "0"},
		{"nil inequality zero", "{{ nil != 0 }}", "1"},
		{"zero inequality nil", "{{ 0 != nil }}", "1"},
		{"nil inequality five", "{{ nil != 5 }}", "1"},
		{"five inequality nil", "{{ 5 != nil }}", "1"},
		// Nil with floats
		{"nil not equal zero float", "{{ nil == 0.0 }}", "0"},
		{"zero float not equal nil", "{{ 0.0 == nil }}", "0"},
		{"nil not equal float", "{{ nil == 1.5 }}", "0"},
		{"float not equal nil", "{{ 1.5 == nil }}", "0"},
		{"nil inequality zero float", "{{ nil != 0.0 }}", "1"},
		{"zero float inequality nil", "{{ 0.0 != nil }}", "1"},
		{"nil inequality pi", "{{ nil != 3.14 }}", "1"},
		{"pi inequality nil", "{{ 3.14 != nil }}", "1"},
		// Nil with strings
		{"nil not equal empty string", "{{ nil == '' }}", "0"},
		{"empty string not equal nil", "{{ '' == nil }}", "0"},
		{"nil not equal string", "{{ nil == 'test' }}", "0"},
		{"string not equal nil", "{{ 'test' == nil }}", "0"},
		{"nil inequality empty string", "{{ nil != '' }}", "1"},
		{"empty string inequality nil", "{{ '' != nil }}", "1"},
		{"nil inequality hello", "{{ nil != 'hello' }}", "1"},
		{"hello inequality nil", "{{ 'hello' != nil }}", "1"},
		// Nil with arrays
		{"nil not equal empty array", "{{ nil == [] }}", "0"},
		{"empty array not equal nil", "{{ [] == nil }}", "0"},
		{"nil not equal int array", "{{ nil == [1, 2] }}", "0"},
		{"int array not equal nil", "{{ [1, 2] == nil }}", "0"},
		{"nil inequality empty array", "{{ nil != [] }}", "1"},
		{"empty array inequality nil", "{{ [] != nil }}", "1"},
		{"nil inequality single element array", "{{ nil != [1] }}", "1"},
		{"single element array inequality nil", "{{ [1] != nil }}", "1"},
		// Nil with objects
		{"nil not equal empty object", "{{ nil == {} }}", "0"},
		{"empty object not equal nil", "{{ {} == nil }}", "0"},
		{"nil not equal named object", "{{ nil == {name: 'test'} }}", "0"},
		{"named object not equal nil", "{{ {name: 'test'} == nil }}", "0"},
		{"nil inequality empty object", "{{ nil != {} }}", "1"},
		{"empty object inequality nil", "{{ {} != nil }}", "1"},
		{"nil inequality single key object", "{{ nil != {x: 1} }}", "1"},
		{"single key object inequality nil", "{{ {x: 1} != nil }}", "1"},
		// Nil with nil
		{"nil equal nil", "{{ nil == nil }}", "1"},
		{"nil inequality nil", "{{ nil != nil }}", "0"},
		// --------------------------- Arrays ---------------------------
		// Arrays with arrays
		{"empty arrays equal", "{{ [] == [] }}", "1"},
		{"empty arrays inequality false", "{{ [] != [] }}", "0"},
		{"arrays of different length not equal", "{{ [1, 2] == [1] }}", "0"},
		{"arrays equal", "{{ [1, 2] == [1, 2] }}", "1"},
		{"equal arrays inequality false", "{{ [1, 2] != [1, 2] }}", "0"},
		{"different single element arrays not equal", "{{ [1] == [2] }}", "0"},
		{"different single element arrays inequality true", "{{ [1] != [2] }}", "1"},
		{"equal single element arrays", "{{ [1] == [1] }}", "1"},
		{"equal single element arrays inequality false", "{{ [1] != [1] }}", "0"},
		{"deeply nested arrays equal", "{{ [[[1]]] == [[[1]]] }}", "1"},
		{"deeply nested arrays inequality false", "{{ [[[1]]] != [[[1]]] }}", "0"},
		{"deeply nested arrays with different values not equal", "{{ [[[1]]] == [[[2]]] }}", "0"},
		{
			"deeply nested arrays with different values inequality true",
			"{{ [[[1]]] != [[[2]]] }}",
			"1",
		},
		// --------------------------- Objects ---------------------------
		// Objects with objects
		{"empty objects equal", "{{ {} == {} }}", "1"},
		{"empty objects inequality false", "{{ {} != {} }}", "0"},
		{"objects with equal name pair equal", "{{ {name: 'test'} == {name: 'test'} }}", "1"},
		{
			"objects with equal name pair inequality false",
			"{{ {name: 'test'} != {name: 'test'} }}",
			"0",
		},
		{"objects with equal single key equal", "{{ {x: 1} == {x: 1} }}", "1"},
		{"objects with equal single key inequality false", "{{ {x: 1} != {x: 1} }}", "0"},
		{"nested objects equal", "{{ {x: {y: 2}} == {x: {y: 2}} }}", "1"},
		{"nested objects inequality false", "{{ {x: {y: 2}} != {x: {y: 2}} }}", "0"},
		{
			"deeply nested objects equal",
			"{{ {x: {y: {z: 'Anna'}}} == {x: {y: {z: 'Anna'}}} }}",
			"1",
		},
		{
			"deeply nested objects inequality false",
			"{{ {x: {y: {z: 'Anna'}}} != {x: {y: {z: 'Anna'}}} }}",
			"0",
		},
		// Mixed type comparisons with == and !=
		{"int and numeric string not equal", "{{ 1 == '1' }}", "0"},
		{"numeric string and int not equal", "{{ '1' == 1 }}", "0"},
		{"true and one not equal", "{{ true == 1 }}", "0"},
		{"one and true not equal", "{{ 1 == true }}", "0"},
		{"false and zero not equal", "{{ false == 0 }}", "0"},
		{"zero and false not equal", "{{ 0 == false }}", "0"},
		{"float and int not equal", "{{ 1.0 == 1 }}", "0"},
		{"int and float not equal", "{{ 1 == 1.0 }}", "0"},
		{"object and int not equal", "{{ {} == 1 }}", "0"},
		{"array and int not equal", "{{ [] == 1 }}", "0"},
	}

	for _, tc := range cases {
		evaluated, failure := testEval(tc.inp, nil)
		if failure != nil {
			t.Fatalf("Case: %s. evaluation failed: %s", tc.name, failure)
		}

		err, ok := evaluated.(*value.Error)
		if ok {
			t.Errorf("Case: %s. Evaluation failed: %s", tc.name, err)
			return
		}

		if res := evaluated.String(); res != tc.expect {
			t.Errorf("Case: %s. Result is not %q, got %q", tc.name, tc.expect, res)
		}
	}
}

func TestEvalNilExp(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
		data   map[string]any
	}{
		{"nil literal", `<h1>{{ nil }}</h1>`, "<h1></h1>", nil},
		{"repeated nil literals", `<h1>{{ nil; nil }}</h1>`, "<h1></h1>", nil},
		{"nil from data", `<h1>{{ data }}</h1>`, "<h1></h1>", map[string]any{"data": nil}},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, tc.data, tc.name)
	}
}

func TestEvalStringExp(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		// Basic string output
		{"double quoted string", `{{ "Hello World" }}`, "Hello World"},
		{"single quoted string", `{{ 'Hello World 2' }}`, "Hello World 2"},
		// String in HTML attributes
		{
			"string in attribute with double quotes",
			`<div {{ 'data-attr="Test"' }}></div>`,
			"<div data-attr=&#34;Test&#34;></div>",
		},
		{
			"string in attribute with single quotes",
			`<div {{ "data-attr='Test'" }}></div>`,
			"<div data-attr=&#39;Test&#39;></div>",
		},
		// String with escaped characters
		{"escaped double quotes", `{{ "She \"is\" pretty" }}`, "She &#34;is&#34; pretty"},
		{"escaped backslash", `{{ "a\\" }}`, `a\`},
		{"windows path", `{{ "C:\\new" }}`, `C:\new`},
		// String concatenation
		{
			"concatenation with spaces",
			`{{ "Korotchaeva" + " " + "Anna" }}`,
			"Korotchaeva Anna",
		},
		{
			"concatenation of multiple strings",
			`{{ "She" + " " + "is" + " " + "nice" }}`,
			"She is nice",
		},
		// Empty string
		{"single quoted empty string", "{{ '' }}", ""},
		{"double quoted empty string", `{{ "" }}`, ""},
		// String with HTML escaping
		{"escaped h1 tag", `{{ "<h1>Test</h1>" }}`, "&lt;h1&gt;Test&lt;/h1&gt;"},
		{"escaped div tag", `{{ "<div>Hello</div>" }}`, "&lt;div&gt;Hello&lt;/div&gt;"},
		{"escaped script tag", `{{ "<script>alert('xss')</script>" }}`,
			"&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;",
		},
		{
			"escaped img tag",
			`{{ "<img src='test.jpg'>" }}`,
			"&lt;img src=&#39;test.jpg&#39;&gt;",
		},
		{"escaped ampersand", `{{ "&amp;" }}`, "&amp;amp;"},
		{"escaped br tag", `{{ "<br/>" }}`, "&lt;br/&gt;"},
		{
			"escaped anchor tag",
			`{{ "<a href='#'>Link</a>" }}`,
			"&lt;a href=&#39;#&#39;&gt;Link&lt;/a&gt;",
		},
		{"escaped paragraph tag", `{{ "<p class='text'>Content</p>" }}`,
			"&lt;p class=&#39;text&#39;&gt;Content&lt;/p&gt;",
		},
		{
			"escaped list tag",
			`{{ "<ul><li>Item</li></ul>" }}`,
			"&lt;ul&gt;&lt;li&gt;Item&lt;/li&gt;&lt;/ul&gt;",
		},
		{"escaped input tag", `{{ "<input type='text' value='test'>" }}`,
			"&lt;input type=&#39;text&#39; value=&#39;test&#39;&gt;",
		},
		// String concatenation with variables
		{
			"concatenation with variable",
			`{{ name = "Anna"; "Hello " + name }}`,
			"Hello Anna",
		},
		{
			"concatenation of two variables",
			`{{ a = "Hello"; b = "World"; a + " " + b }}`,
			"Hello World",
		},
		// String with numbers
		{"string with int conversion", `{{ "Count: " + 5.str() }}`, "Count: 5"},
		{"string with float conversion", `{{ "Pi: " + 3.14.str() }}`, "Pi: 3.14"},
		// Empty string concatenation
		{"empty prefix concatenation", `{{ "" + "test" }}`, "test"},
		{"empty suffix concatenation", `{{ "test" + "" }}`, "test"},
		{"concatenation of empty strings", `{{ "" + "" }}`, ""},
		{"single quoted empty concatenation", `{{ '' + '' }}`, ""},
		{"mixed quote empty concatenation", `{{ '' + "" }}`, ""},
		{"mixed empty concatenation", `{{ "" + '' }}`, ""},
		// String with special characters
		{"newline in string", `{{ "test\nline" }}`, "test\nline"},
		{"tab in string", `{{ "tab\there" }}`, "tab\there"},
		// Long string concatenation
		{"long concatenation", `{{ "a" + "b" + "c" + "d" + "e" }}`, "abcde"},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, nil, tc.name)
	}
}

func TestEvalTernaryExp(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		// Basic boolean conditions
		{"true condition", `{{ true ? "Yes" : "No" }}`, "Yes"},
		{"false condition", `{{ false ? "Yes" : "No" }}`, "No"},
		// Falsy values
		{"nil condition", `{{ nil ? "Yes" : "No" }}`, "No"},
		{"truthy int condition", `{{ 1 ? "Yes" : "No" }}`, "Yes"},
		{"zero int condition", `{{ 0 ? "Yes" : "No" }}`, "No"},
		{"empty string condition", `{{ "" ? "Yes" : "No" }}`, "No"},
		// Negation
		{"negated true", `{{ !true ? "Yes" : "No" }}`, "No"},
		{"negated false", `{{ !false ? "Yes" : "No" }}`, "Yes"},
		// Double negation
		{"double negated true", `{{ !!true ? 1 : 0 }}`, "1"},
		{"double negated false", `{{ !!false ? 1 : 0 }}`, "0"},
		// Comparison operators
		{"int equality true", `{{ 1 == 1 ? "Yes" : "No" }}`, "Yes"},
		{"int equality false", `{{ 1 == 2 ? "Yes" : "No" }}`, "No"},
		{"greater than true", `{{ 5 > 3 ? "Yes" : "No" }}`, "Yes"},
		{"greater than false", `{{ 3 > 5 ? "Yes" : "No" }}`, "No"},
		// String comparison
		{"string equality true", `{{ "a" == "a" ? "Yes" : "No" }}`, "Yes"},
		{"string equality false", `{{ "a" == "b" ? "Yes" : "No" }}`, "No"},
		// Array truthiness
		{"empty array is falsy", `{{ [] ? "Yes" : "No" }}`, "No"},
		{"non-empty array is truthy", `{{ [1] ? "Yes" : "No" }}`, "Yes"},
		// Object truthiness
		{"empty object is falsy", `{{ {} ? "Yes" : "No" }}`, "No"},
		{"non-empty object is truthy", `{{ {x: 1} ? "Yes" : "No" }}`, "Yes"},
		// Nested ternary
		{"nested ternary true true", `{{ true ? (true ? "A" : "B") : "C" }}`, "A"},
		{"nested ternary true false", `{{ true ? (false ? "A" : "B") : "C" }}`, "B"},
		{"nested ternary false", `{{ false ? (true ? "A" : "B") : "C" }}`, "C"},
		// Unparenthesised nested ternary in the true branch
		{"unparenthesized nested ternary true", `{{ true ? true ? "A" : "B" : "C" }}`, "A"},
		{"unparenthesized nested ternary false", `{{ true ? false ? "A" : "B" : "C" }}`, "B"},
		{"unparenthesized nested ternary comparison", `{{ 1 > 0 ? 1 > 2 ? "x" : "y" : "z" }}`, "y"},
		// Arithmetic expressions
		{"arithmetic comparison true", `{{ 1 + 1 == 2 ? "Yes" : "No" }}`, "Yes"},
		{"subtraction comparison true", `{{ 10 - 5 > 3 ? "Yes" : "No" }}`, "Yes"},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, nil, tc.name)
	}
}

func TestEvalIfDir(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		// Basic if statements
		{"basic true if", `@if(true)Hello@end`, "Hello"},
		{"basic false if", `@if(false)Hello@end`, ""},
		{"true if with else", `@if(true)X@else@end`, "X"},
		{"false if with else", `@if(false)@elseX@end`, "X"},
		{"true if with empty else", `@if(true)@else@end`, ""},
		{"false if with else body", `@if(false)1@else@end`, ""},
		// HTML wrapper
		{"if inside heading true", `<h2>@if(true)Hello@end</h2>`, "<h2>Hello</h2>"},
		{"if inside heading false", `<h2>@if(false)Hello@end</h2>`, "<h2></h2>"},
		// @elseif statements
		{"elseif first true", `@if(true)Anna@elseif(true)Lili@end`, "Anna"},
		{
			"elseif second true",
			`@if(false)Alan@elseif(true)Serhii@end`,
			"Serhii",
		},
		{
			"else branch chosen",
			`@if(false)Ana De Armaz@elseif(false)David@elseVladimir@end`,
			"Vladimir",
		},
		{
			"third elseif true",
			`@if(false)Will@elseif(false)Daria@elseif(true)Poll@end`,
			"Poll",
		},
		{
			"first matching elseif wins",
			`@if(false)Lara@elseif(true)Susan@elseif(true)Smith@end`,
			"Susan",
		},
		{
			"else with email text",
			`@if(false)@elseif(false)@elsemy@mail.com@end`,
			"my@mail.com",
		},
		{"no branch matches", `@if(false)A@elseif(false)B@end`, ""},
		// Functions
		{"if with true binary function", `@if(true.binary())Hello@end`, "Hello"},
		{"if with false binary function", `@if(false.binary())Hello@end`, ""},
		{
			"if with empty string length",
			`@if("".len() > 0)Non empty@elseEmpty@end`,
			"Empty",
		},
		{
			"if with non-empty string length",
			`@if("x".len() > 0)Non empty@elseEmpty@end`,
			"Non empty",
		},
		// Truthy/falsy values
		{"truthy int", `@if(1)Yes@end`, "Yes"},
		{"falsy zero", `@if(0)Yes@end`, ""},
		{"falsy empty string", `@if("")Yes@end`, ""},
		{"falsy empty array", `@if([])Yes@end`, ""},
		{"falsy empty object", `@if({})Yes@end`, ""},
		{"falsy nil", `@if(nil)Yes@end`, ""},
		// Logical operators
		{"logical and true", `@if(true && true)Yes@end`, "Yes"},
		{"logical and false", `@if(true && false)Yes@end`, ""},
		{"logical or true", `@if(false || true)Yes@end`, "Yes"},
		{"logical or false", `@if(false || false)Yes@end`, ""},
		// Boolean negation
		{"not true", `@if(!true)Yes@end`, ""},
		{"not false", `@if(!false)Yes@end`, "Yes"},
		{"double not true", `@if(!!true)Yes@end`, "Yes"},
		{"double not false", `@if(!!false)Yes@end`, ""},
		// Comparison operators
		{"int equality true", `@if(1 == 1)Yes@end`, "Yes"},
		{"int equality false", `@if(1 == 2)Yes@end`, ""},
		{"int inequality true", `@if(1 != 2)Yes@end`, "Yes"},
		{"int less than", `@if(1 < 2)Yes@end`, "Yes"},
		{"int greater than", `@if(2 > 1)Yes@end`, "Yes"},
		{"int greater or equal", `@if(2 >= 2)Yes@end`, "Yes"},
		{"int less or equal", `@if(2 <= 2)Yes@end`, "Yes"},
		// Expression results as conditions
		{"arithmetic equality true", `@if(1 + 1 == 2)Yes@end`, "Yes"},
		{"arithmetic equality false", `@if(5 - 3 == 1)Yes@end`, ""},
		{"falsy arithmetic condition uses else", `@if(2 - 2)No@elseYes@end`, "Yes"},
		{"truthy negative arithmetic condition", `@if(-1 + 2)Yes@elseNo@end`, "Yes"},
		// Nested if statements
		{
			"nested if chooses innermost else",
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
		// Variable modification from outside scope
		{
			"increment modifies outer variable",
			`{{ sum = 0 }}@if(true){{ sum++ }}@end{{ sum }}`,
			"1",
		},
		// Nested @if
		{
			"nested if multiplies outer variable",
			`{{ val = 1 }}@if(true)@if(true){{ val = val * 5 }}@end@end{{ val }}`,
			"5",
		},
		// Modification inside @else
		{
			"modification inside else",
			`{{ status = "pending" }}@if(false)@else{{ status = "done" }}@end{{ status }}`,
			"done",
		},
		// Sequential modifications
		{
			"sequential if modifications",
			`{{ n = 10 }}@if(n > 5){{ n = n - 2 }}@end@if(n == 8){{ n = 0 }}@end{{ n }}`,
			"0",
		},
		// Deep nesting
		{
			"deep nested if modification",
			`{{ deep = 0 }}@if(true)@if(true)@if(true){{ deep = 42 }}@end@end@end{{ deep }}`,
			"42",
		},
		// Modification inside @elseif
		{
			"modification inside elseif",
			`{{ status = "pending" }}@if(false)@elseif(true){{ status = "done" }}@end{{ status }}`,
			"done",
		},
	}

	for _, tc := range cases {
		evaluated, failure := testEval(tc.inp, nil)
		if failure != nil {
			t.Fatalf("Case: %s. evaluation failed: %s", tc.name, failure)
		}

		err, ok := evaluated.(*value.Error)
		if ok {
			t.Errorf("Case: %s. Evaluation failed: %s", tc.name, err)
		}

		if res := strings.TrimSpace(evaluated.String()); res != tc.expect {
			t.Errorf("Case: %s. Result is not %q, got %q", tc.name, tc.expect, res)
		}
	}
}

func TestEvalArr(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		// Empty arrays
		{"empty array", `{{ [] }}`, ""},
		{"array with empty array", `{{ [[]] }}`, ""},
		{"deeply nested empty arrays", `{{ [[[[[]]]]] }}`, ""},
		// Integer arrays
		{"int array", `{{ [1, 2, 3] }}`, "1, 2, 3"},
		{"single zero", `{{ [0] }}`, "0"},
		{"negative int array", `{{ [-1, -2, -3] }}`, "-1, -2, -3"},
		// Float arrays
		{"float array", `{{ [1.5, 2.5, 3.5] }}`, "1.5, 2.5, 3.5"},
		{"single float zero", `{{ [0.0] }}`, "0.0"},
		// String arrays
		{"string array", `{{ ["Anna", "Serhii" ] }}`, "Anna, Serhii"},
		{"single string", `{{ ["hello"] }}`, "hello"},
		{"strings with spaces", `{{ ["with space", "another"] }}`, "with space, another"},
		// Boolean arrays
		{"boolean array", `{{ [true, false] }}`, "1, 0"},
		{"single true", `{{ [true] }}`, "1"},
		{"single false", `{{ [false] }}`, "0"},
		// Mixed type arrays
		{"mixed type array", `{{ [1, "hello", true] }}`, "1, hello, 1"},
		{"mixed falsy values", `{{ [0, "", false] }}`, "0, , 0"},
		// Nested arrays
		{"nested arrays", `{{ [[1, 2], [3, 4]] }}`, "1, 2, 3, 4"},
		{"nested array with mixed depth", `{{ [[1, [2]], 3] }}`, "1, 2, 3"},
		{"deeply nested single int", `{{ [[[11]]] }}`, "11"},
		// Arrays with expressions
		{"array with arithmetic expressions", `{{ [1 + 2, 3 * 2] }}`, "3, 6"},
		{"array with subtraction and division", `{{ [5 - 3, 10 / 2] }}`, "2, 5"},
		// Arrays with variables
		{"array with variables", `{{ a = 1; b = 2; [a, b] }}`, "1, 2"},
		{"array with repeated variable", `{{ x = "test"; [x, x] }}`, "test, test"},
		// Arrays with ternary
		{
			"array with ternary",
			`{{ [true ? 1 : 0, false ? "yes" : "no"] }}`,
			"1, no",
		},
		// Trailing comma
		{"trailing comma ints", `{{ [1, 2,] }}`, "1, 2"},
		{"trailing comma strings", `{{ ["a", "b",] }}`, "a, b"},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, nil, tc.name)
	}
}

func TestEvalIndexExp(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		// Basic indexing
		{"index zero", `{{ [1, 2, 3][0] }}`, "1"},
		{"index one", `{{ [1, 2, 3][1] }}`, "2"},
		{"index two", `{{ [1, 2, 3][2] }}`, "3"},
		// Negative indexing (not allowed)
		{"negative index", `{{ [1, 2, 3][-1] }}`, ""},
		{"negative index two", `{{ [1, 2, 3][-2] }}`, ""},
		{"negative index on strings", `{{ ["a", "b", "c"][-3] }}`, ""},
		// Out of bounds indexing
		{"index into empty array", `{{ [][2] }}`, ""},
		{"index out of bounds", `{{ [1, 2][5] }}`, ""},
		{"large negative index", `{{ [1, 2][-5] }}`, ""},
		// String indexing
		{"index string array", `{{ ["Some string"][0] }}`, "Some string"},
		{"index second string", `{{ ["hello", "world"][1] }}`, "world"},
		// Nested array indexing
		{"nested index zero one", `{{ [[1, 2], [3, 4]][0][1] }}`, "2"},
		{"nested index one zero", `{{ [[1, 2], [3, 4]][1][0] }}`, "3"},
		{"triple nested index", `{{ [[[11]]][0][0][0] }}`, "11"},
		// Index with expression
		{"index with addition expression", `{{ arr = [1, 2, 3]; arr[1 + 1] }}`, "3"},
		{"index with subtraction expression", `{{ arr = [1, 2, 3]; arr[3 - 1] }}`, "3"},
		// Variable as index
		{"index with variable", `{{ arr = [1, 2, 3]; i = 2; arr[i] }}`, "3"},
		{
			"index with variable on string array",
			`{{ arr = ["x", "y", "z"]; i = 0; arr[i] }}`,
			"x",
		},
		// Ternary as index
		{"index with true ternary", `{{ ["a", "b", "c"][true ? 0 : 1] }}`, "a"},
		{"index with false ternary", `{{ ["a", "b", "c"][false ? 0 : 1] }}`, "b"},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, nil, tc.name)
	}
}

func TestEvalAssign(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		// Integer assignment
		{"regular int assignment", `{{ age = 18 }}`, ""},
		{"int assignment and read", `{{ age = 18; age }}`, "18"},
		{"negative int assignment", `{{ age = -5; age }}`, "-5"},
		{"zero int assignment", `{{ age = 0; age }}`, "0"},
		// Float assignment
		{"float assignment", `{{ pi = 3.14; pi }}`, "3.14"},
		{"zero float assignment", `{{ price = 0.0; price }}`, "0.0"},
		{"negative float assignment", `{{ negative = -2.5; negative }}`, "-2.5"},
		{"float reassignment with addition", `{{ f = 0.0; f = f + 0.5; f }}`, "0.5"},
		// String assignment
		{"string assignment", `{{ name = "Anna"; name }}`, "Anna"},
		{"empty string assignment", `{{ empty = ""; empty }}`, ""},
		{
			"string with escaped quotes",
			`{{ quote = "He said \"Hello\""; quote }}`,
			"He said &#34;Hello&#34;",
		},
		// Boolean assignment
		{"true boolean assignment", `{{ flag = true; flag }}`, "1"},
		{"false boolean assignment", `{{ flag = false; flag }}`, "0"},
		// Nil assignment
		{"nil assignment", `{{ nothing = nil; nothing }}`, ""},
		// Multiple assignments
		{"multiple int assignments", `{{ a = 1; b = 2; c = 3; a + b + c }}`, "6"},
		{
			"multiple string assignments",
			`{{ x = "Hello"; y = "World"; x + " " + y }}`,
			"Hello World",
		},
		// Reassignment
		{"int reassignment", `{{ age = 18; age = 25; age }}`, "25"},
		{
			"string reassignment",
			`{{ name = "Anna"; name = "Maria"; name }}`,
			"Maria",
		},
		{"repeated increment expression assignment", `{{ x = 1; x = x + 1; x = x + 1; x }}`, "3"},
		// Assignment with expression
		{"assignment from addition", `{{ sum = 5 + 3; sum }}`, "8"},
		{"assignment from mixed arithmetic", `{{ calc = 10 * 2 - 5; calc }}`, "15"},
		{"assignment from parenthesized arithmetic", `{{ result = (2 + 3) * 4; result }}`, "20"},
		// Assignment with string concatenation
		{
			"assignment from string concatenation",
			`{{ full = "John" + " " + "Doe"; full }}`,
			"John Doe",
		},
		{
			"assignment from concatenated string with method call",
			`{{ msg = "Count: " + 5.str(); msg }}`,
			"Count: 5",
		},
		// Assignment with ternary
		{"assignment from true ternary", `{{ val = true ? 1 : 0; val }}`, "1"},
		{"assignment from false ternary", `{{ val = false ? "yes" : "no"; val }}`, "no"},
		// Assignment with method call
		{"assignment from method call", `{{ upper = "hello".upper(); upper }}`, "HELLO"},
		{"assignment from len method", `{{ len = [1, 2, 3].len(); len }}`, "3"},
		// Assignment with index
		{
			"assignment from array index",
			`{{ arr = [10, 20, 30]; val = arr[1]; val }}`,
			"20",
		},
		{"assignment from string method", `{{ str = "abc"; val = str.at(0); val }}`, "a"},
		// Assignment in conditional
		{"assignment inside true conditional", `@if(true){{ x = 5 }}{{ x }}@end`, "5"},
		{"assignment inside false conditional", `@if(false){{ x = 5 }}{{ x }}@end`, ""},
		// Assignment with loop variable
		{
			"assignment from each loop variable",
			`@each(n in [1, 2, 3]){{ x = n }}{{ x }}@end`,
			"123",
		},
		// Object assignment
		{"object assignment with dot access", `{{ user = {name: "John"}; user.name }}`, "John"},
		{"object assignment with quoted key", `{{ data = {"age": 25}; data.age }}`, "25"},
		{
			"object reassignment",
			`{{ user = {"name": "Ann"}; user = {"name": "Anna"}; user.name }}`,
			"Anna",
		},
		{
			"object property assignment",
			`{{ user = {}; user.name = "Anna"; user.name }}`,
			"Anna",
		},
		{
			"object property assignment duplicate",
			`{{ user = {}; user.name = "Anna"; user.name }}`,
			"Anna",
		},
		{
			name:   "nested object property assignment",
			inp:    `{{ user = {}; user.address = { street: "x" }; user.address.street = 'y'; user.address.street }}`,
			expect: "y",
		},
		{
			name:   "nested object property assignment after empty object",
			inp:    `{{ user = {}; user.address = {}; user.address.street = "Pushkina"; user.address.street }}`,
			expect: "Pushkina",
		},
		// Array assignment
		{
			"array assignment",
			`{{ names = ["Anna", "Serhii"]; names }}`,
			"Anna, Serhii",
		},
		{"empty array assignment", `{{ empty = []; empty }}`, ""},
		{
			"nested array assignment",
			`{{ nested = [[1, 2], [3, 4]]; nested }}`,
			"1, 2, 3, 4",
		},
		{
			"array element assignment",
			`{{ names = ['Serhii', 'Nastya']; names[1] = 'Anna'; names }}`,
			"Serhii, Anna",
		},
		{
			"multiple array element assignments",
			`{{ nums = [10, 20, 30]; nums[0] = 1; nums[1] = 2; nums }}`,
			"1, 2, 30",
		},
		{"nested array element assignment", `{{ x = [[[20]]]; x[0][0][0] = 30; x }}`, "30"},
		{
			"deeply nested array element assignment",
			`{{ x = ['1', ['2', ['3', ['4']]]]; x[1][1][1][0] = '5'; x }}`,
			"1, 2, 3, 5",
		},
		{
			"array element assigned variable",
			`{{ x = [0]; newVal = 10; x[0] = newVal; x[0] }}`,
			"10",
		},
		// Mixed assignment
		{
			"object property inside array",
			`{{ x = [{ name: 'Chiori' }]; x[0].name = 'Mavuika'; x[0].name }}`,
			"Mavuika",
		},
		{
			"object property inside array assigned variable",
			`{{ name = 'Mavuika'; x = [{ name: 'Chiori' }]; x[0].name = name; x[0].name }}`,
			"Mavuika",
		},
		// Index assignment edge cases
		{
			"array element assigned object",
			`{{ arr = [1, 2]; arr[0] = {name: 'x'}; arr[0].name }}`,
			"x",
		},
		{
			"array element assigned array",
			`{{ arr = [1, 2]; arr[0] = [3, 4]; arr[0][0] }}`,
			"3",
		},
		{
			"array element assigned expression index",
			`{{ arr = [1, 2, 3]; arr[1 + 1] = 5; arr[2] }}`,
			"5",
		},
		{
			"nested array element assigned int",
			`{{ arr = [[1], [2]]; arr[0][0] = 10; arr[0][0] }}`,
			"10",
		},
		{
			"array element assigned variable copy",
			`{{ x = 5; arr = [1, 2]; arr[0] = x; arr[0] }}`,
			"5",
		},
		{
			"array element assigned another element",
			`{{ arr = [1, 2]; arr[0] = arr[1]; arr[0] }}`,
			"2",
		},
		{
			"last array element assignment",
			`{{ arr = [1, 2, 3]; arr[2] = 10; arr[2] }}`,
			"10",
		},
		{
			"array element concatenation",
			`{{ arr = ['a', 'b', 'c']; arr[0] = arr[0] + 'x'; arr[0] }}`,
			"ax",
		},
		{
			"boolean array element assignment",
			`{{ arr = [true, false]; arr[0] = false; arr[0] }}`,
			"0",
		},
		// Object property assignment via index notation
		{
			"object property via bracket notation",
			`{{ user = {}; user["name"] = "Barbara"; user.name }}`,
			"Barbara",
		},
		{
			"object property via bracket reassignment",
			`{{ user = {name: "Skir"}; user["name"] = "Skirk"; user["name"] }}`,
			"Skirk",
		},
		{
			"object property via bracket assigned nil",
			`{{ user = {name: "Ann"}; user["name"] = nil; user["name"] }}`,
			"",
		},
		{
			"object property via computed key",
			`{{ user = {}; key = "name"; user[key] = "Ifa"; user.name }}`,
			"Ifa",
		},
		{
			"nested object via bracket notation",
			`{{ u = {a: {s: "x"}}; u["a"]["s"] = "y"; u.a.s }}`,
			"y",
		},
		{
			"array of objects property via bracket notation",
			`{{ x = [{ name: 'Chiori' }]; x[0]['name'] = 'Mavuika'; x[0].name }}`,
			"Mavuika",
		},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, nil, tc.name)
	}
}

func TestIncDecDir(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		// Integer
		{"int increment", `{{ i = 10; i++; i }}`, "11"},
		{"zero increment", `{{ x = 0; x++; x }}`, "1"},
		{"decrement copy does not affect original", `{{ a = 2; b = a; b--; a }}`, "2"},
		// Float
		{"float decrement", `{{ x = 4.4; x--; x }}`, "3.4000000000000004"},
		{"float zero decrement", `{{ x = 0.0; x--; x }}`, "-1.0"},
		{"float half decrement", `{{ x = 0.5; x--; x }}`, "-0.5"},
		{"negative float decrement", `{{ x = -1.5; x--; x }}`, "-2.5"},
		{"float decrement copy does not affect original", `{{ a = 2.5; b = a; b--; a }}`, "2.5"},
		// Object with property access
		{
			"object property increment",
			`{{ user = {age: 30}; user.age++; user.age }}`,
			"31",
		},
		{
			"object property decrement",
			`{{ user = {age: 30}; user.age--; user.age }}`,
			"29",
		},
		{
			"object property increment then decrement",
			`{{ user = {age: 30}; user.age++; user.age--; user.age }}`,
			"30",
		},
		// Array with index access
		{"array element increment", `{{ arr = [1, 2, 3]; arr[0]++; arr[0] }}`, "2"},
		{"array element decrement", `{{ arr = [1, 2, 3]; arr[1]--; arr[1] }}`, "1"},
		{
			"array element increment then decrement",
			`{{ arr = [1, 2, 3]; arr[2]++; arr[2]--; arr[2] }}`,
			"3",
		},
		// Nested object and array access
		{
			"nested object array property increment",
			`{{ data = {users: [{age: 20}]}; data.users[0].age++; data.users[0].age }}`,
			"21",
		},
		{
			"nested object array property decrement",
			`{{ data = {users: [{age: 20}]}; data.users[0].age--; data.users[0].age }}`,
			"19",
		},
		{
			"nested object array property increment then decrement",
			`{{ data = {users: [{age: 20}]}; data.users[0].age++; data.users[0].age--; data.users[0].age }}`,
			"20",
		},
		// Multiple increments and decrements
		{"double increment", `{{ x = 5; x++; x++; x }}`, "7"},
		{"double decrement", `{{ y = 10; y--; y--; y }}`, "8"},
		{"increment and decrement sequence", `{{ z = 3; z++; z--; z++; z }}`, "4"},
		// Edge cases
		{"zero decrement", `{{ x = 0; x--; x }}`, "-1"},
		{"negative increment", `{{ x = -1; x++; x }}`, "0"},
		{"float zero increment", `{{ x = 0.0; x++; x }}`, "1.0"},
		{"negative float decrement from small value", `{{ x = -0.5; x--; x }}`, "-1.5"},
		{"float increment", `{{ x = 1.5; x++; x }}`, "2.5"},
		{"float decrement above zero", `{{ x = 1.5; x--; x }}`, "0.5"},
		{"negative float increment", `{{ x = -1.5; x++; x }}`, "-0.5"},
		{"negative float decrement again", `{{ x = -1.5; x--; x }}`, "-2.5"},
		// Increment and decrement with expressions
		{"increment after addition assignment", `{{ x = 5; x = x + 2; x++; x }}`, "8"},
		{"decrement after subtraction assignment", `{{ y = 10; y = y - 3; y--; y }}`, "6"},
		{
			"increment and decrement after multiplication assignment",
			`{{ z = 3; z = z * 2; z++; z--; z }}`,
			"6",
		},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, nil, tc.name)
	}
}

func TestEvalForDir(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		// Basic for loops
		{"basic for ascending", `@for(i = 0; i < 2; i++){{ i }}@end`, "01"},
		{"for inclusive upper bound", `@for(i = 1; i <= 3; i++){{ i }}@end`, "123"},
		{"for descending", `@for(i = 5; i > 0; i--){{ i }}@end`, "54321"},
		// Empty loop body
		{"empty loop body", `@for(i = 0; i < 2; i++)@end`, ""},
		{"empty loop body with else", `@for(i = 0; i < 2; i++)@else@end`, ""},
		// Condition evaluates to nil
		{
			"nil loop condition",
			`{{ nilVar = nil }}@for(i = 0; nilVar; i++)x@end`,
			"",
		},
		{"nil literal loop condition", `@for(i = 0; nil; i++)Some content@end`, ""},
		{
			"missing object property loop condition",
			`{{ user = {} }}@for(i = 0; user.missing; i++)x@breakif(i == 3)@end`,
			"",
		},
		// Empty loop header parts
		{"false only loop header", `@for(; false;)Here@end`, ""},
		{"infinite loop with break", `@for(; true;)x@break@end`, "x"},
		{"false condition with init", `@for(c = 1; false; c++){{ c }}@end`, ""},
		{"zero iterations", `@for(i = 0; i < 0; i++){{ i }}@end`, ""},
		{"infinite loop with break and empty header", `@for(;;){{ 1 }}@break@end`, "1"},
		// Single iteration
		{"single iteration with equality condition", `@for(c = 1; c == 1; c++){{ c }}@end`, "1"},
		{"single iteration", `@for(i = 0; i < 1; i++){{ i }}@end`, "0"},
		// @else directive
		{"false condition with empty else", `@for(c = 1; false; c++){{ c }}@else@end`, ""},
		{
			"false condition with else content",
			`@for(c = 1; false; c++){{ c }}@else<b>Empty</b>@end`,
			"<b>Empty</b>",
		},
		{
			"zero iterations with else content",
			`@for(c = 0; c < 0; c++){{ c }}@elseEmpty@end`,
			"Empty",
		},
		// @break directive
		{
			"break after first iteration",
			`@for(i = 1; i <= 3; i++){{ i }}@break@end`,
			"1",
		},
		{
			"break before body",
			`@for(i = 1; i <= 3; i++)@break{{ i }}@end`,
			"",
		},
		{
			"break inside if",
			`@for(i = 1; i <= 3; i++)@if(i == 3)@break@end{{ i }}@end`,
			"12",
		},
		{"break immediately", `@for(i = 0; i < 10; i++)@break@end`, ""},
		// @continue directive
		{
			"continue before body",
			`@for(i = 1; i <= 3; i++)@continue{{ i }}@end`,
			"",
		},
		{
			"continue after body",
			`@for(i = 1; i <= 3; i++){{ i }}@continue@end`,
			"123",
		},
		{
			"continue inside if",
			`@for(i = 1; i <= 3; i++)@if(i == 2)@continue@end{{ i }}@end`,
			"13",
		},
		{
			"continue on even numbers",
			`@for(i = 1; i <= 5; i++)@if(i % 2 == 0)@continue@end{{ i }}@end`,
			"135",
		},
		// @breakif directive
		{
			"breakif at third iteration",
			`@for(i = 1; i <= 3; i++)@breakif(i == 3){{ i }}@end`,
			"12",
		},
		{
			"breakif at second iteration",
			`@for(i = 1; i <= 3; i++)@breakif(i == 2){{ i }}@end`,
			"1",
		},
		{
			"breakif greater than five",
			`@for(i = 1; i <= 10; i++)@breakif(i > 5){{ i }}@end`,
			"12345",
		},
		// @continueif directive
		{
			"continueif at third iteration",
			`@for(i = 1; i <= 3; i++)@continueif(i == 3){{ i }}@end`,
			"12",
		},
		{
			"continueif at second iteration",
			`@for(i = 1; i <= 3; i++)@continueif(i == 2){{ i }}@end`,
			"13",
		},
		{
			"continueif on even numbers",
			`@for(i = 1; i <= 5; i++)@continueif(i % 2 == 0){{ i }}@end`,
			"135",
		},
		// Nested for loops
		{
			"nested for loops output coordinates",
			`@for(i = 0; i < 2; i++)@for(j = 0; j < 2; j++){{ i }}{{ j }}@end@end`,
			"00011011",
		},
		{
			"nested for loops multiply",
			`@for(i = 1; i <= 2; i++)@for(j = 1; j <= 2; j++){{ i * j }}@end@end`,
			"1224",
		},
		// For loop with HTML
		{
			"for loop with html list",
			`<ul>@for(i = 1; i <= 3; i++)<li>{{ i }}</li>@end</ul>`,
			"<ul><li>1</li><li>2</li><li>3</li></ul>",
		},
		// Variable modification from outside scope
		{
			"accumulate sum in for loop",
			`{{ sum = 0 }}@for(i = 1; i <= 5; i++){{ sum = sum + i }}@end{{ sum }}`,
			"15",
		},
		{
			"variable declared before for",
			`{{ x = 0 }}@for(; x < 2; x++){{ x }}@end`,
			"01",
		},
		{
			"count iterations in for loop",
			`{{ count = 0 }}@for(i = 0; i < 3; i++){{ count = count + 1 }}@end{{ count }}`,
			"3",
		},
		{
			"incrementing variable before for",
			`{{ n = 0 }}@for(; true; n++){{ n }}@breakif(n == 2)@end`,
			"012",
		},
		{
			"for loop uses existing counter",
			`{{ i = 0 }}@for(; i < 3; i++){{ i }}@end`,
			"012",
		},
		// Multiple statements in loop body
		{
			"multiple statements in loop body",
			`@for(i = 0; i < 3; i++){{ i }};{{ i * 2 }}@end`,
			"0;01;22;4",
		},
		{
			name:   "for loop prints square",
			inp:    `@for(i = 1; i <= 3; i++){{ i }} -> {{ i * i }}|@end`,
			expect: "1 -> 1|2 -> 4|3 -> 9|",
		},
		// Control statements only affect the current loop
		{
			"break only affects inner loop",
			`@for(i = 0; i < 3; i++)@for(j = 0; j < 3; j++)@breakif(j == 1){{ i }}{{ j }}@end;@end`,
			"00;10;20;",
		},
		{
			"continue only affects inner loop",
			`@for(i = 0; i < 3; i++)@for(j = 0; j < 2; j++)@continueif(j == 0){{ i }}{{ j }}@end;@end`,
			"01;11;21;",
		},
		// @break must not run the post statement
		{
			"breakif does not run post statement",
			`{{ i = 0 }}@for(; i < 10; i++)@breakif(i == 3)@end{{ i }}`,
			"3",
		},
		{
			"loop variable shadowing string",
			`{{ i = "hi" }}@for(i = 0; i < 2; i++){{ i }}@end{{ i }}`,
			"01hi",
		},
		{
			"break inside if does not run post statement",
			`{{ i = 0 }}@for(; i < 10; i++)@if(i == 3)@break@end@end{{ i }}`,
			"3",
		},
		// test loop variable
		{"loop.index inside @for", `@for(i = 0; i < 3; i++){{ loop.index }}@end`, "012"},
		{"loop.iter inside @for", `@for(i = 0; i < 3; i++){{ loop.iter }}@end`, "123"},
		{"loop.first inside @for", `@for(i = 0; i < 3; i++){{ loop.first }}@end`, "100"},
		// last is unavailable because a for loop may not have a known length
		{"loop.last is nil inside @for", `@for(i = 0; i < 2; i++){{ loop.last }}|@end`, "||"},
		// the outermost loop has no parent
		{
			"loop.parent is nil at outermost @for",
			`@for(i = 0; i < 2; i++){{ loop.parent }}|@end`,
			"||",
		},
		// loop refers to the nearest enclosing loop
		{
			"nested @for uses nearest loop",
			`@for(i = 0; i < 2; i++)@for(j = 0; j < 2; j++){{ loop.index }}{{ loop.parent.index }};@end@end`,
			"00;10;01;11;",
		},
		// the outer loop binding is restored after the inner loop ends
		{
			"outer @for restored after inner loop",
			`@for(i = 0; i < 2; i++)@for(j = 0; j < 1; j++)@end{{ loop.index }}@end`,
			"01",
		},
		// loop.parent.parent reaches the outermost loop
		{
			"loop.parent.parent reaches outermost @for",
			`@for(i = 0; i < 1; i++)@for(j = 0; j < 1; j++)@for(k = 0; k < 1; k++){{ loop.parent.parent.index }}@end@end@end`,
			"0",
		},
		// loop metadata is available in body conditionals and is not affected by continue
		{
			"loop.index inside @for continueif",
			`@for(i = 0; i < 3; i++)@continueif(loop.index == 1){{ loop.index }}@end`,
			"02",
		},
		{
			"loop.iter inside @for breakif",
			`@for(i = 0; i < 5; i++)@breakif(loop.iter == 3){{ loop.index }}@end`,
			"01",
		},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, nil, tc.name)
	}
}

func TestEvalEachDir(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		{
			"each over strings",
			`@each(name in ["anna", "serhii"]) {{ name }} @end`,
			" anna  serhii ",
		},
		{
			"each over strings with space",
			`@each(name in ["anna", "serhii"]){{ name + " " }}@end`,
			"anna serhii ",
		},
		{
			"each over strings inside HTML",
			`<div>@each(name in ["anna", "serhii"]){{ name }} @end</div>`,
			"<div>anna serhii </div>",
		},
		{"each over ints", `@each(num in [1, 2, 3]){{ num }}@end`, "123"},
		{"each over empty array", `@each(num in []){{ num }}@end`, ""},
		{"each over nil elements", `@each(item in [nil, nil]){{ item }}@end`, ""},
		{"each with nil among ints", `@each(x in [1, nil, 2]){{ x }},@end`, "1,,2,"},
		{"each over mixed types", `@each(x in [1, "two"]){{ x }}@end`, "1two"},
		// test loop variable
		{"loop.index inside @each", `@each(num in [43, 12, 53]){{ loop.index }}@end`, "012"},
		{"loop.index for single item @each", `@each(num in [100]){{ loop.index }}@end`, "0"},
		{"loop.first inside @each", `@each(num in [1, 2, 3, 4]){{ loop.first }}@end`, "1000"},
		{"loop.last inside @each", `@each(num in [1, 2, 3, 4]){{ loop.last }}@end`, "0001"},
		{"loop.iter inside @each", `@each(num in [4, 2, 8]){{ loop.iter }}@end`, "123"},
		{
			"loop.iter across many @each items",
			`@each(num in [9, 3, 44, 24, 1, 3]){{ loop.iter }}@end`,
			"123456",
		},
		// test @else directive
		{
			"each with else on empty array",
			`@each(v in []){{ v }}@else<b>Empty array</b>@end`,
			"<b>Empty array</b>",
		},
		{"each with empty else", `@each(n in []){{ n }}@else@end`, ""},
		{"each with else text", `@each(n in []){{ n }}@elsetest@end`, "test"},
		{
			"each over five ints",
			`@each(n in [1, 2, 3, 4, 5]){{ n }}@end`,
			"12345",
		},
		// test @break directive
		{
			"break before body",
			`@each(n in [1, 2, 3, 4, 5])@break{{ n }}@end`,
			"",
		},
		{
			"break after body",
			`@each(n in [1, 2, 3, 4, 5]){{ n }}@break@end`,
			"1",
		},
		{
			"break inside if",
			`@each(n in [1, 2, 3, 4, 5])@if(n == 3)@break@end{{ n }}@end`,
			"12",
		},
		// test @continue directive
		{
			"continue before body",
			`@each(n in [1, 2, 3, 4, 5])@continue{{ n }}@end`,
			"",
		},
		{
			"continue after body",
			`@each(n in [1, 2, 3, 4, 5]){{ n }}@continue@end`,
			"12345",
		},
		{
			"continue inside if",
			`@each(n in [1, 2, 3, 4, 5])@if(n == 3)@continue@end{{ n }}@end`,
			"1245",
		},
		// test @breakif directive
		{
			"breakif at third element",
			`@each(n in [1, 2, 3, 4, 5])@breakif(n == 3){{ n }}@end`,
			"12",
		},
		{
			"breakif on string element",
			`@each(n in ["ann", "serhii", "sam"])@breakif(n == 'sam'){{ n }}{{ ' ' }}@end`,
			"ann serhii ",
		},
		// test @continueif directive
		{
			"continueif at third element",
			`@each(n in [1, 2, 3, 4, 5])@continueif(n == 3){{ n }}@end`,
			"1245",
		},
		{
			"continueif on string element",
			`@each(n in ["ann", "serhii", "sam"])@continueif(n == 'sam'){{ n }}{{' '}}@end`,
			"ann serhii ",
		},
		// Variable modification from outside scope
		{
			"each modifies outer variable",
			`{{ sum = 0 }}@each(num in [1]){{ sum = sum + num }}@end{{ sum }}`,
			"1",
		},
		{
			"each assignment does not modify outer sum",
			`{{ sum = 0 }}@each(num in [1]){{ a = sum + num }}@end{{ sum }}`,
			"0",
		},
		// Accumulation across multiple iterations
		{
			"accumulation across multiple iterations",
			`{{ sum = 0 }}@each(num in [1, 2, 3]){{ sum = sum + num }}@end{{ sum }}`,
			"6",
		},
		// String concatenation across iterations
		{
			"string concatenation across iterations",
			`{{ str = "" }}@each(char in ["x", "y", "z"]){{ str = str + char }}@end{{ str }}`,
			"xyz",
		},
		// Nested @each loops modifying outside scope
		{
			"nested each increments counter",
			`{{ count = 0 }}@each(i in [1, 2])@each(j in [1, 2]){{ count++ }}@end@end{{ count }}`,
			"4",
		},
		// Shadowing an outside variable with the loop iterator
		{
			"loop iterator shadows outer variable",
			`{{ fruit = "apple" }}@each(fruit in ["banana"])@end{{ fruit }}`,
			"apple",
		},
		// Loop variable does not clash with outer variable of another type
		{
			"loop variable does not clash with outer variable type",
			`{{ x = "hi" }}@each(x in [1, 2]){{ x }}@end`,
			"12",
		},
		// Loop variable does not leak outside the loop
		{
			"loop variable does not leak outside loop",
			`{{ x = "hi" }}@each(x in [1, 2]){{ x }}@end{{ x }}`,
			"12hi",
		},
		// Loop variable does not overwrite outer variable of the same type
		{
			"loop variable does not overwrite outer variable of same type",
			`{{ x = 9 }}@each(x in [1, 2]){{ x }}@end{{ x }}`,
			"129",
		},
		// Control statements only affect the current loop
		{
			"break only affects inner each",
			`@each(i in [0, 1, 2])@each(j in [0, 1, 2])@breakif(j == 1){{ i }}{{ j }}@end;@end`,
			"00;10;20;",
		},
		{
			"continue only affects inner each",
			`@each(i in [0, 1, 2])@each(j in [0, 1])@continueif(j == 0){{ i }}{{ j }}@end;@end`,
			"01;11;21;",
		},
		{
			"break only affects inner loop in mixed for",
			`@each(i in [0, 1, 2])@for(j = 0; j < 3; j++)@breakif(j == 1){{ i }}{{ j }}@end;@end`,
			"00;10;20;",
		},
		// loop refers to the nearest enclosing loop
		{
			"nested @each uses nearest loop",
			`@each(i in [1, 2])@each(j in [3, 4]){{ loop.parent.index }}:{{ loop.index }}|@end@end`,
			"0:0|0:1|1:0|1:1|",
		},
		// the outer loop binding is restored after the inner loop ends
		{
			"outer @each restored after inner loop",
			`@each(i in [1, 2])@each(j in [3])@end{{ loop.index }}@end`,
			"01",
		},
		// loop.parent.parent reaches the outermost loop
		{
			"loop.parent.parent reaches outermost @each",
			`@each(i in [1, 2])@each(j in [3])@each(k in [4]){{ loop.parent.parent.index }}@end@end@end`,
			"01",
		},
		// the outermost loop has no parent
		{"loop.parent is nil at outermost @each", `@each(i in [1]){{ loop.parent }}|@end`, "|"},
		// @for creates its own loop object
		{
			"@each inside @for uses @for as parent",
			`@each(i in [1, 2])@for(j = 0; j < 1; j++){{ loop.index }}{{ loop.parent.index }}@end@end`,
			"0001",
		},
		// an @each inside @for points to the @for loop as its parent
		{
			"nested @each inside @for uses parent chain",
			`@each(i in [1])@for(j = 0; j < 2; j++)@each(k in [1]){{ loop.parent.index }}{{ loop.parent.parent.index }};@end@end@end`,
			"00;10;",
		},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, nil, tc.name)
	}
}

func TestEvalObjLit(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		// Bracket and dot access
		{"bracket access", `{{ {"name": "Ann"}['name'] }}`, "Ann"},
		{"dot access", `{{ {"name": "Ann"}.name }}`, "Ann"},
		// Basic property access
		{"object assignment dot access", `{{ obj = {name: "Ann"}; obj.name }}`, "Ann"},
		{
			"object with quoted keys",
			`{{ o = {"name": "Ann", "age": 22}; o.age }}`,
			"22",
		},
		// Nested objects
		{
			"nested object access",
			`{{ user = {"father": {"name": "Ann"}}; user.father.name }}`,
			"Ann",
		},
		{
			"deeply nested object access",
			`{{ user = {"father": {"name": {"first": "Serhii"}}}; user.father.name.first }}`,
			"Serhii",
		},
		{
			"nested object with trailing commas",
			`{{ u = {"father": {name: {"first": "Serhii",},},}; u['father']['name'].first }}`,
			"Serhii",
		},
		// Shorthand properties
		{
			"shorthand property name",
			`{{ name = "Serhii"; age = 12; obj = { name, age }; obj.name }}`,
			"Serhii",
		},
		{
			"shorthand property age",
			`{{ name = "Serhii"; age = 12; obj = { name, age }; obj.age }}`,
			"12",
		},
		// Case-insensitive first character access
		{"capitalized key lower access", `{{ {"Name": "Ann"}.name }}`, "Ann"},
		{"lowercase key capitalized access", `{{ {"name": "Ann"}.Name }}`, ""},
		{"capitalized key bracket lower access", `{{ {"Name": "Ann"}["name"] }}`, "Ann"},
		{"lowercase key bracket capitalized access", `{{ {"name": "Ann"}["Name"] }}`, ""},
		// Non-existent keys
		{"missing dot key", `{{ obj = {"name": "Ann"}; obj.age }}`, ""},
		{
			"missing bracket key",
			`{{ obj = {"name": "Ann"}; obj['missing'] }}`,
			"",
		},
		// Empty object
		{"missing key on literal", `{{ {}.name }}`, ""},
		{"missing key on variable", `{{ obj = {}; obj.name }}`, ""},
		// Empty string key
		{"empty string key on literal", `{{ {a: 1}[""] }}`, ""},
		{"empty string key on variable", `{{ obj = {a: 1}; key = ""; obj[key] }}`, ""},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, nil, tc.name)
	}
}

func TestResolveObjKey(t *testing.T) {
	cases := []struct {
		name   string
		pairs  map[string]value.Literal
		key    string
		expect string
	}{
		{
			name: "exact key wins over capitalized key",
			pairs: map[string]value.Literal{
				"name": &value.Int{Val: 1},
				"Name": &value.Int{Val: 2},
			},
			key:    "name",
			expect: "name",
		},
		{
			name:   "lowercase key resolves to capitalized key",
			pairs:  map[string]value.Literal{"Name": &value.Int{Val: 2}},
			key:    "name",
			expect: "Name",
		},
		{
			name:   "missing key returns key as given",
			pairs:  map[string]value.Literal{"age": &value.Int{Val: 1}},
			key:    "name",
			expect: "name",
		},
		{
			name:   "empty object returns key as given",
			pairs:  map[string]value.Literal{},
			key:    "name",
			expect: "name",
		},
		{
			name:   "capitalized key that exists is returned",
			pairs:  map[string]value.Literal{"Name": &value.Int{Val: 1}},
			key:    "Name",
			expect: "Name",
		},
		{
			name:   "missing capitalized key returns key as given",
			pairs:  map[string]value.Literal{"name": &value.Int{Val: 1}},
			key:    "Name",
			expect: "Name",
		},
		{
			name:   "empty key returns empty key",
			pairs:  map[string]value.Literal{},
			key:    "",
			expect: "",
		},
		{
			name:   "key with special characters resolves correctly",
			pairs:  map[string]value.Literal{"user-name": &value.Int{Val: 1}},
			key:    "user-name",
			expect: "user-name",
		},
		{
			name:   "key with special characters that doesn't exist returns key as given",
			pairs:  map[string]value.Literal{"user-name": &value.Int{Val: 1}},
			key:    "user_name",
			expect: "user_name",
		},
		{
			name:   "numeric string key resolves correctly",
			pairs:  map[string]value.Literal{"123": &value.Int{Val: 1}},
			key:    "123",
			expect: "123",
		},
		{
			name:   "numeric string key that doesn't exist returns key as given",
			pairs:  map[string]value.Literal{"123": &value.Int{Val: 1}},
			key:    "456",
			expect: "456",
		},
		{
			name:   "key with spaces resolves correctly",
			pairs:  map[string]value.Literal{"first name": &value.Int{Val: 1}},
			key:    "first name",
			expect: "first name",
		},
		{
			name:   "key with spaces that doesn't exist returns key as given",
			pairs:  map[string]value.Literal{"first name": &value.Int{Val: 1}},
			key:    "last name",
			expect: "last name",
		},
		{
			name:   "key with mixed case resolves correctly",
			pairs:  map[string]value.Literal{"UserName": &value.Int{Val: 1}},
			key:    "UserName",
			expect: "UserName",
		},
		{
			name:   "key with mixed case that doesn't exist returns key as given",
			pairs:  map[string]value.Literal{"UserName": &value.Int{Val: 1}},
			key:    "username",
			expect: "username",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			obj := value.NewObj(tc.pairs)

			got := resolveObjKey(obj, tc.key)

			if got != tc.expect {
				t.Errorf("expected %q, got %q", tc.expect, got)
			}
		})
	}
}

func TestEvalComments(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		// Basic comments
		{"basic comment", "{{-- This is a comment --}}", ""},
		{"empty comment", "{{-- --}}", ""},
		{"compact comment", "{{--a--}}", ""},
		// Comments in HTML
		{
			"comment inside section",
			"<section>{{-- This is a comment --}}</section>",
			"<section></section>",
		},
		{"comment inside div", "<div>{{-- comment --}}</div>", "<div></div>"},
		// Comments with text around
		{"comment between text", "Some {{-- --}}text", "Some text"},
		{"comment inside word", "Hello{{-- comment --}}World", "HelloWorld"},
		{"comment between words", "Start{{-- middle --}}End", "StartEnd"},
		// Empty or minimal comments
		{"minimal comment", "{{----}}", ""},
		{"comment with spaces", "{{--     --}}", ""},
		{"comment with newline", "{{--\n--}}", ""},
		// Comments with directives
		{"commented each directive", "{{-- @each(u in users){{ u }}@end --}}", ""},
		{"commented if directive", "{{-- @if(true)Hello@end --}}", ""},
		{"commented for directive", "{{-- @for(i=0;i<10;i++){{i}}@end --}}", ""},
		// Comments with special characters
		{"comment with html and symbols", "{{-- <html> &amp; @#$%^&*() --}}", ""},
		{"comment with japanese", "{{-- Japanese: こんにちは --}}", ""},
		{"comment with emoji", "{{-- Emoji: 😀🎉🚀 --}}", ""},
		// Comments inside directives
		{"comment at start of if", "@if(true){{-- comment --}}Yes@end", "Yes"},
		{"comment at end of if", "@if(true)Yes{{-- comment --}}@end", "Yes"},
		{"comment inside each", "@each(n in [1,2]){{-- c --}}{{n}}@end", "12"},
		// Multiple comments
		{"comment between text fragments", "{{-- a --}}Text{{-- b --}}", "Text"},
		{"two comments before text", "{{-- x --}}{{-- y --}}Result", "Result"},
		{"comments between letters", "A{{-- 1 --}}B{{-- 2 --}}C", "ABC"},
		// Comments at boundaries
		{"comment at start", "{{-- start --}}End", "End"},
		{"comment at end", "Start{{-- end --}}", "Start"},
		// Multi-line comments
		{"multiline comment", "{{-- Line 1\nLine 2\nLine 3 --}}", ""},
		{
			"multiline comment between text",
			"Before{{--\nmulti\nline\n--}}After",
			"BeforeAfter",
		},
		// Nested comment-like text (should be treated as part of comment)
		{"nested comment-like text", "{{-- Contains {{-- and --}} inside --}}", ""},
		{"comment containing interpolation", "{{-- {{ 1 + 1 }} inside --}}", ""},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, nil, tc.name)
	}
}

func TestCannotUseOperatorError(t *testing.T) {
	cases := []struct {
		name  string
		inp   string
		left  value.ValueType
		op    string
		right value.ValueType
	}{
		// Int + Float (both directions)
		{"int plus float", "{{ 3 + 2.0 }}", value.INT_VAL, "+", value.FLOAT_VAL},
		{"float plus int", "{{ 2.0 + 3 }}", value.FLOAT_VAL, "+", value.INT_VAL},
		{"int minus float", "{{ 5 - 1.5 }}", value.INT_VAL, "-", value.FLOAT_VAL},
		{"float minus int", "{{ 4.5 - 2 }}", value.FLOAT_VAL, "-", value.INT_VAL},
		{"int times float", "{{ 3 * 2.5 }}", value.INT_VAL, "*", value.FLOAT_VAL},
		{"float times int", "{{ 2.5 * 4 }}", value.FLOAT_VAL, "*", value.INT_VAL},
		{"int divide float", "{{ 10 / 2.5 }}", value.INT_VAL, "/", value.FLOAT_VAL},
		{"float divide int", "{{ 7.5 / 3 }}", value.FLOAT_VAL, "/", value.INT_VAL},
		// Arithmetic with strings
		{"int times string", "{{ 5 * 'x' }}", value.INT_VAL, "*", value.STR_VAL},
		{"string times int", "{{ 'x' * 3 }}", value.STR_VAL, "*", value.INT_VAL},
		{"string divide int", "{{ 'x' / 2 }}", value.STR_VAL, "/", value.INT_VAL},
		{"int divide string", "{{ 2 / 'x' }}", value.INT_VAL, "/", value.STR_VAL},
		{"string minus int", "{{ 'a' - 2 }}", value.STR_VAL, "-", value.INT_VAL},
		{"string divide float", "{{ 'a' / 3.0 }}", value.STR_VAL, "/", value.FLOAT_VAL},
		// Arithmetic with booleans
		{"true plus int", "{{ true + 5 }}", value.BOOL_VAL, "+", value.INT_VAL},
		{"int minus true", "{{ 5 - true }}", value.INT_VAL, "-", value.BOOL_VAL},
		{"false minus float", "{{ false - 2.0 }}", value.BOOL_VAL, "-", value.FLOAT_VAL},
		{"float plus true", "{{ 2.0 + true }}", value.FLOAT_VAL, "+", value.BOOL_VAL},
		{"true times int", "{{ true * 1 }}", value.BOOL_VAL, "*", value.INT_VAL},
		{"true divide int", "{{ true / 2 }}", value.BOOL_VAL, "/", value.INT_VAL},
		// Boolean with strings
		{"true plus string", "{{ true + 'str' }}", value.BOOL_VAL, "+", value.STR_VAL},
		{"string minus false", "{{ 'str' - false }}", value.STR_VAL, "-", value.BOOL_VAL},
		// Modulo operator
		{"int modulo float", "{{ 5 % 2.0 }}", value.INT_VAL, "%", value.FLOAT_VAL},
		{"float modulo int", "{{ 5.0 % 2 }}", value.FLOAT_VAL, "%", value.INT_VAL},
		{"string modulo int", "{{ 'a' % 2 }}", value.STR_VAL, "%", value.INT_VAL},
		{"bool modulo int", "{{ true % 2 }}", value.BOOL_VAL, "%", value.INT_VAL},
		// Array/Object operations
		{"array times object", "{{ [] * {} }}", value.ARR_VAL, "*", value.OBJ_VAL},
		{"object divide array", "{{ {} / [] }}", value.OBJ_VAL, "/", value.ARR_VAL},
		{"array plus object", "{{ [] + {} }}", value.ARR_VAL, "+", value.OBJ_VAL},
		{"object minus array", "{{ {} - [] }}", value.OBJ_VAL, "-", value.ARR_VAL},
		{"array times one", "{{ [] * 1 }}", value.ARR_VAL, "*", value.INT_VAL},
		{"object divide one", "{{ {} / 1 }}", value.OBJ_VAL, "/", value.INT_VAL},
		// Int/Object operations
		{"int plus object", "{{ 3 + {} }}", value.INT_VAL, "+", value.OBJ_VAL},
		{"object minus int", "{{ {} - 3 }}", value.OBJ_VAL, "-", value.INT_VAL},
		{"int times object", "{{ 5 * {} }}", value.INT_VAL, "*", value.OBJ_VAL},
		// Array with arithmetic
		{"array plus int", "{{ [] + 5 }}", value.ARR_VAL, "+", value.INT_VAL},
		{"int minus array", "{{ 10 - [] }}", value.INT_VAL, "-", value.ARR_VAL},
		{"array times three", "{{ [] * 3 }}", value.ARR_VAL, "*", value.INT_VAL},
		// Float with arrays/objects
		{"float plus array", "{{ 3.14 + [] }}", value.FLOAT_VAL, "+", value.ARR_VAL},
		{"object divide float", "{{ {} / 2.5 }}", value.OBJ_VAL, "/", value.FLOAT_VAL},
		// String with arrays/objects
		{"string minus array", "{{ 'x' - [] }}", value.STR_VAL, "-", value.ARR_VAL},
		{"object minus string", "{{ {} - 'x' }}", value.OBJ_VAL, "-", value.STR_VAL},
		{"string plus array", "{{ 'str' + [] }}", value.STR_VAL, "+", value.ARR_VAL},
		{"object plus string", "{{ {} + 'str' }}", value.OBJ_VAL, "+", value.STR_VAL},
		// Strings
		{"string minus string", "{{ 'a' - 'b' }}", value.STR_VAL, "-", value.STR_VAL},
		{"string times string", "{{ 'a' * 'b' }}", value.STR_VAL, "*", value.STR_VAL},
		{"string divide string", "{{ 'a' / 'b' }}", value.STR_VAL, "/", value.STR_VAL},
		{"string less than string", "{{ 'a' < 'b' }}", value.STR_VAL, "<", value.STR_VAL},
		{"string greater than string", "{{ 'a' > 'b' }}", value.STR_VAL, ">", value.STR_VAL},
		{"string less or equal string", "{{ 'a' <= 'b' }}", value.STR_VAL, "<=", value.STR_VAL},
		{"string greater or equal string", "{{ 'a' >= 'b' }}", value.STR_VAL, ">=", value.STR_VAL},
		{"string modulo string", "{{ 'a' % 'b' }}", value.STR_VAL, "%", value.STR_VAL},
		// String + number (addition not supported for strings)
		{"string plus int", "{{ 'test' + 5 }}", value.STR_VAL, "+", value.INT_VAL},
		{"int plus string", "{{ 5 + 'test' }}", value.INT_VAL, "+", value.STR_VAL},
		{"string plus float", "{{ 'test' + 5.5 }}", value.STR_VAL, "+", value.FLOAT_VAL},
		{"float plus string", "{{ 5.5 + 'test' }}", value.FLOAT_VAL, "+", value.STR_VAL},
		// Boolean - string
		{"true minus string", "{{ true - 'str' }}", value.BOOL_VAL, "-", value.STR_VAL},
		{"string minus true", "{{ 'str' - true }}", value.STR_VAL, "-", value.BOOL_VAL},
		// Float modulo operations
		{"float modulo float", "{{ 3.14 % 2.0 }}", value.FLOAT_VAL, "%", value.FLOAT_VAL},
		{"fractional float modulo int", "{{ 3.14 % 2 }}", value.FLOAT_VAL, "%", value.INT_VAL},
		// Array / array
		{"array divide array", "{{ [] / [] }}", value.ARR_VAL, "/", value.ARR_VAL},
		// Object / object
		{"object divide object", "{{ {} / {} }}", value.OBJ_VAL, "/", value.OBJ_VAL},
		// Nil with arithmetic operators
		{"nil plus int", "{{ nil + 1 }}", value.NIL_VAL, "+", value.INT_VAL},
		{"nil minus int", "{{ nil - 1 }}", value.NIL_VAL, "-", value.INT_VAL},
		{"nil times int", "{{ nil * 1 }}", value.NIL_VAL, "*", value.INT_VAL},
		{"nil divide int", "{{ nil / 1 }}", value.NIL_VAL, "/", value.INT_VAL},
	}

	for _, tc := range cases {
		evaluated, failure := testEval(tc.inp, nil)
		if failure != nil {
			t.Fatalf("Case: %s. evaluation failed: %s", tc.name, failure)
		}

		err, ok := evaluated.(*value.Error)
		if !ok {
			t.Fatalf("Case: %s. Evaluation failed, got error %q", tc.name, err)
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
			t.Fatalf("Case: %s. Error message must be:\n%q\ngot:\n%q", tc.name, expect, err)
		}
	}
}

func TestEvalFormatDateGlobalFunc(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		// Test date time
		{
			"datetime time only",
			"{{ formatDate('2026-03-24 13:03:25', '15:04:05') }}",
			"13:03:25",
		},
		{
			"datetime full",
			"{{ formatDate('2027-12-13 10:00:09', '2006-01-02 15:04:05') }}",
			"2027-12-13 10:00:09",
		},
		{
			"datetime date only",
			"{{ formatDate('2023-10-12 01:10:25', '2006-01-02') }}",
			"2023-10-12",
		},
		// Test date
		{
			"date with date format",
			"{{ formatDate('2026-03-24', '2006-01-02') }}",
			"2026-03-24",
		},
		{
			"date with time format",
			"{{ formatDate('2026-03-24', '15:04:05') }}",
			"00:00:00",
		},
		{
			"date with full format",
			"{{ formatDate('2026-03-24', '2006-01-02 15:04:05') }}",
			"2026-03-24 00:00:00",
		},
		// Test time
		{
			"time with date format",
			"{{ formatDate('10:00:09', '2006-01-02') }}",
			"0000-01-01",
		},
		{
			"time with time format",
			"{{ formatDate('10:00:09', '15:04:05') }}",
			"10:00:09",
		},
		{
			"time with full format",
			"{{ formatDate('10:00:09', '2006-01-02 15:04:05') }}",
			"0000-01-01 10:00:09",
		},
		// Edge cases empty string
		{
			"empty input",
			"{{ formatDate('', '2006-01-02 15:04:05') }}",
			"",
		},
		// Edge cases midnight
		{
			"midnight time only",
			"{{ formatDate('2026-01-01 00:00:00', '15:04:05') }}",
			"00:00:00",
		},
		{
			"midnight full",
			"{{ formatDate('2026-01-01 00:00:00', '2006-01-02 15:04:05') }}",
			"2026-01-01 00:00:00",
		},
		{
			"midnight date only",
			"{{ formatDate('2026-01-01 00:00:00', '2006-01-02') }}",
			"2026-01-01",
		},
		// Edge cases end of day
		{
			"end of day time only",
			"{{ formatDate('2026-12-31 23:59:59', '15:04:05') }}",
			"23:59:59",
		},
		{
			"end of day full",
			"{{ formatDate('2026-12-31 23:59:59', '2006-01-02 15:04:05') }}",
			"2026-12-31 23:59:59",
		},
		{
			"end of day date only",
			"{{ formatDate('2026-12-31 23:59:59', '2006-01-02') }}",
			"2026-12-31",
		},
		// Leap year date
		{
			"leap day date only",
			"{{ formatDate('2024-02-29 12:30:45', '2006-01-02') }}",
			"2024-02-29",
		},
		{
			"leap day time only",
			"{{ formatDate('2024-02-29 12:30:45', '15:04:05') }}",
			"12:30:45",
		},
		{
			"leap day full",
			"{{ formatDate('2024-02-29 12:30:45', '2006-01-02 15:04:05') }}",
			"2024-02-29 12:30:45",
		},
		// Single-digit month/day padding
		{
			"padded date",
			"{{ formatDate('2026-03-05 01:02:03', '2006-01-02') }}",
			"2026-03-05",
		},
		{
			"padded time",
			"{{ formatDate('2026-03-05 01:02:03', '15:04:05') }}",
			"01:02:03",
		},
		{
			"padded full datetime",
			"{{ formatDate('2026-03-05 01:02:03', '2006-01-02 15:04:05') }}",
			"2026-03-05 01:02:03",
		},
		// Early morning times
		{
			"early morning time",
			"{{ formatDate('2026-06-15 01:23:45', '15:04:05') }}",
			"01:23:45",
		},
		{
			"early morning full datetime",
			"{{ formatDate('2026-06-15 02:00:00', '2006-01-02 15:04:05') }}",
			"2026-06-15 02:00:00",
		},
		// Afternoon/PM times
		{
			"afternoon time",
			"{{ formatDate('2026-08-20 13:45:30', '15:04:05') }}",
			"13:45:30",
		},
		{
			"afternoon full datetime",
			"{{ formatDate('2026-08-20 15:30:00', '2006-01-02 15:04:05') }}",
			"2026-08-20 15:30:00",
		},
		{
			"late evening time",
			"{{ formatDate('2026-08-20 23:00:00', '15:04:05') }}",
			"23:00:00",
		},
		// Date-only inputs with various formats
		{
			"year 2000 date",
			"{{ formatDate('2000-01-01', '2006-01-02') }}",
			"2000-01-01",
		},
		{
			"year 1999 date with full format",
			"{{ formatDate('1999-12-31', '2006-01-02 15:04:05') }}",
			"1999-12-31 00:00:00",
		},
		// Time-only inputs with various formats
		{
			"midnight time input",
			"{{ formatDate('00:00:00', '15:04:05') }}",
			"00:00:00",
		},
		{
			"noon time input with full format",
			"{{ formatDate('12:34:56', '2006-01-02 15:04:05') }}",
			"0000-01-01 12:34:56",
		},
		{
			"late night time input",
			"{{ formatDate('23:59:59', '15:04:05') }}",
			"23:59:59",
		},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, nil, tc.name)
	}
}
