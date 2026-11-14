package apitest

import (
	"math"
	"strings"
	"testing"
	"time"

	"codeberg.org/textwire/textwire/v6"

	"codeberg.org/textwire/textwire/v6/pkg/fail"
	"codeberg.org/textwire/textwire/v6/pkg/file"
	"codeberg.org/textwire/textwire/v6/pkg/value"
)

type namedRole string

type namedUser struct {
	Name string
	Role namedRole
}

type namedUserWithHooks struct {
	Name string
	Role namedRole
	Fn   func()
	Ch   chan int
}

type namedUserWithChan struct {
	Name string
	Ch   chan int
}

type namedUserWithComplex struct {
	Name string
	Num  complex128
}

type namedUserWithUint struct {
	Name string
	Big  uint64
}

func TestEvalStr(t *testing.T) {
	datetime, _ := time.Parse("2006-01-02 15:04:05", "1990-12-23 11:22:33")
	date, _ := time.Parse("2006-01-02", "1990-12-23")
	var age *int

	cases := []struct {
		name   string
		inp    string
		expect string
		data   map[string]any
	}{
		{
			name:   "Negate boolean from data",
			inp:    `{{ !isAdmin }}`,
			expect: "1",
			data:   map[string]any{"isAdmin": false},
		},
		{
			name:   "Negate nil from data",
			inp:    `{{ !nothing }}`,
			expect: "1",
			data:   map[string]any{"nothing": nil},
		},
		{
			name:   "Iterating over mixed Go values",
			inp:    `@each(x in items){{ x }}|@end`,
			expect: "1|two|3.5||",
			data:   map[string]any{"items": []any{1, "two", 3.5, nil}},
		},
		{
			name:   "Accessing propery name on empty obj variable",
			inp:    `@if(item && item.items && item.items > 1)YES@elseNO@end`,
			expect: "NO",
			data:   map[string]any{"item": nil},
		},
		{
			name:   "Accessing propery name on empty obj variable",
			inp:    `@if(item && item.items && item.items > 1)YES@elseNO@end`,
			expect: "NO",
			data:   map[string]any{"item": nil},
		},
		{
			name:   "Accessing propery name on empty obj variable",
			inp:    `{{ long = user.name.len() > 0 }}{{ long }}`,
			expect: "1",
			data: map[string]any{
				"user": struct{ Name string }{Name: "Harry"},
			},
		},
		{
			name:   "Accessing propery name on empty obj variable",
			inp:    `{{ nameLen = user.name.len() }}{{ nameLen }}`,
			expect: "5",
			data: map[string]any{
				"user": struct{ Name string }{Name: "Harry"},
			},
		},
		{
			name:   "Integer pointer",
			inp:    `<p>{{ age }}</p>`,
			expect: `<p></p>`,
			data:   map[string]any{"age": age},
		},
		{
			name:   "Accessing propery name on empty obj variable",
			inp:    `<p>{{ obj = {}; obj.name }}</p>`,
			expect: "<p></p>",
			data:   nil,
		},
		{
			name:   "Accessing property test on empty object {}",
			inp:    `<h2>{{ {}.test }}</h2>`,
			expect: "<h2></h2>",
			data:   nil,
		},
		{
			name:   "Simple math operation with integers",
			inp:    "{{ 1 + 2 }}",
			expect: "3",
			data:   nil,
		},
		{
			name:   "Simple math operation with identifiers",
			inp:    "{{ n1 * n2 }}",
			expect: "2",
			data:   map[string]any{"n1": 1, "n2": 2},
		},
		{
			name:   "First letter of the object property is case insensitive",
			inp:    "{{ user.iD.str() + ' ' + user.ID.str() }}",
			expect: "1 1",
			data:   map[string]any{"user": struct{ ID uint }{1}},
		},
		{
			name:   "Increment lowercase struct field updates capitalized field",
			inp:    "{{ user.count++; user.Count }}",
			expect: "2",
			data:   map[string]any{"user": struct{ Count int }{1}},
		},
		{
			name:   "Index assign lowercase key updates capitalized struct field",
			inp:    `{{ user['count'] = 5; user.Count }}`,
			expect: "5",
			data:   map[string]any{"user": struct{ Count int }{1}},
		},
		{
			name:   "Index assign on issing uppercase field doesn't change original",
			inp:    `{{ user['Count'] = 5; user.count }}`,
			expect: "1",
			data:   map[string]any{"user": map[string]int{"count": 1}},
		},
		{
			name:   "Accessing user.name.firstName property",
			inp:    "{{ user.name.firstName }}",
			expect: "Ann",
			data: map[string]any{
				"user": struct {
					Name struct{ FirstName string }
				}{
					Name: struct{ FirstName string }{"Ann"},
				},
			},
		},
		{
			name:   "Empty global object is defined",
			inp:    "<span>{{ global }}</span>",
			expect: "<span>{}</span>",
			data:   nil,
		},
		{
			name:   "Properly parses date time with time.Time struct",
			inp:    "{{ date1 }} and {{ date2 }}",
			expect: "1990-12-23 11:22:33 and 1990-12-23 11:22:33",
			data:   map[string]any{"date1": datetime, "date2": &datetime},
		},
		{
			name:   "Properly parses date with time.Time struct",
			inp:    "{{ date1 }} and {{ date2 }}",
			expect: "1990-12-23 00:00:00 and 1990-12-23 00:00:00",
			data:   map[string]any{"date1": date, "date2": &date},
		},
		{
			name:   "Raw string print",
			inp:    "{!! str !!}",
			expect: "<span>Skirk</span>",
			data:   map[string]any{"str": "<span>Skirk</span>"},
		},
		{
			name:   "Named type in struct field",
			inp:    `{{ u.role }}`,
			expect: "editor",
			data:   map[string]any{"u": namedUser{Name: "Bob", Role: "editor"}},
		},
		{
			name:   "Named type in struct field via index",
			inp:    `{{ u["role"] }}`,
			expect: "editor",
			data:   map[string]any{"u": namedUser{Name: "Bob", Role: "editor"}},
		},
		{
			name:   "Top-level named string type",
			inp:    `{{ role }}`,
			expect: "admin",
			data:   map[string]any{"role": namedRole("admin")},
		},
		{
			name:   "Fixed-size array",
			inp:    `{{ nums }}`,
			expect: "1, 2, 3",
			data:   map[string]any{"nums": [3]int{1, 2, 3}},
		},
		{
			name:   "Largest int64 as uint64",
			inp:    `{{ n }}`,
			expect: "9223372036854775807",
			data:   map[string]any{"n": uint64(math.MaxInt64)},
		},
		{
			name:   "Fixed-size array with negative ints",
			inp:    `{{ nums }}`,
			expect: "-1, -2, 0",
			data:   map[string]any{"nums": [3]int{-1, -2, 0}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := textwire.EvaluateString(tc.inp, tc.data)
			if err != nil {
				t.Errorf("error evaluating template: %s", err)
			}

			if actual != tc.expect {
				t.Errorf("wrong evaluation result. Expect:\n'%s'\ngot:\n'%s'", tc.expect, actual)
			}
		})
	}
}

func TestDefinedGlobalFunc(t *testing.T) {
	cases := []struct {
		id     int
		inp    string
		expect string
		data   map[string]any
	}{
		{10, `{{ defined('') }}`, "1", nil},
		{20, `{{ defined("") }}`, "1", nil},
		{30, `{{ defined(-0) }}`, "1", nil},
		{40, `{{ defined(0) }}`, "1", nil},
		{50, `{{ defined(1) }}`, "1", nil},
		{60, `{{ defined(0.0) }}`, "1", nil},
		{70, `{{ defined(-0.0) }}`, "1", nil},
		{80, `{{ defined(1.0) }}`, "1", nil},
		{90, `{{ defined({}) }}`, "1", nil},
		{100, `{{ defined([]) }}`, "1", nil},
		{110, `{{ defined(true) }}`, "1", nil},
		{120, `{{ defined(false) }}`, "1", nil},
		{130, `{{ defined(nil) }}`, "1", nil},
		{140, `{{ defined(undefinedVar) }}`, "0", nil},
		{150, `@if(!defined(definedVar))YES@end`, "YES", nil},
		{160, `{{ defined(definedVar) }}`, "1", map[string]any{"definedVar": "nice"}},
		{170, `@if(defined(definedVar))YES@end`, "YES", map[string]any{"definedVar": "nice"}},
		{
			180,
			`{{ defined(definedVar).then("Yes", "No") }}`,
			"Yes",
			map[string]any{"definedVar": "nice"},
		},
		{190, `{{ defined(nilVar) }}`, "1", map[string]any{"nilVar": nil}},
		{200, `@if(defined(nilVar))YES@end`, "YES", map[string]any{"nilVar": nil}},
		{210, `{{ defined(emptyStr) }}`, "1", map[string]any{"emptyStr": ""}},
		{220, `@if(defined(emptyStr))YES@end`, "YES", map[string]any{"emptyStr": ""}},
		{230, `{{ defined(falseVar) }}`, "1", map[string]any{"falseVar": false}},
		{240, `@if(defined(falseVar))YES@end`, "YES", map[string]any{"falseVar": false}},
		{250, `{{ defined(zeroInt) }}`, "1", map[string]any{"zeroInt": 0}},
		{260, `@if(defined(zeroInt))YES@end`, "YES", map[string]any{"zeroInt": 0}},
		{270, `{{ defined(zeroFloat) }}`, "1", map[string]any{"zeroFloat": 0.0}},
		{
			280,
			`{{ defined(obj.prop) }}`,
			"1",
			map[string]any{"obj": map[string]any{"prop": "value"}},
		},
		{290, `{{ defined(obj.prop) }}`, "0", map[string]any{"obj": map[string]any{}}},
		{
			300,
			`{{ defined(obj.prop.test.nice.cool) }}`,
			"0",
			map[string]any{"obj": map[string]any{}},
		},
		{
			310,
			`{{ defined(obj.nested.prop) }}`,
			"1",
			map[string]any{"obj": map[string]any{"nested": map[string]any{"prop": "value"}}},
		},
		{320, `{{ defined(arr[0]) }}`, "1", map[string]any{"arr": []any{"first", "second"}}},
		{
			330,
			`@if(defined(obj.prop))YES@end`,
			"YES",
			map[string]any{"obj": map[string]any{"prop": "value"}},
		},
		{
			340,
			`@if(defined(definedVar) && defined(nilVar))YES@end`,
			"YES",
			map[string]any{"definedVar": "nice", "nilVar": nil},
		},
		{
			350,
			`@if(defined(definedVar) || defined(undefinedVar))YES@end`,
			"YES",
			map[string]any{"definedVar": "nice"},
		},
		{
			360,
			`@if(defined(obj.prop, obj.nested.prop))YES@end`,
			"YES",
			map[string]any{
				"obj": map[string]any{"prop": "value", "nested": map[string]any{"prop": "value"}},
			},
		},
		{370, `{{ defined(-5) }}`, "1", nil},
		{380, `{{ defined(-3.14) }}`, "1", nil},
		{390, `{{ defined("   ") }}`, "1", nil},
		{400, `{{ defined("\t\n") }}`, "1", nil},
		{410, `{{ defined(boolVar) }}`, "1", map[string]any{"boolVar": true}},
		{420, `{{ defined(boolVar) }}`, "1", map[string]any{"boolVar": false}},
		{430, `{{ defined(floatVar) }}`, "1", map[string]any{"floatVar": 3.14}},
		{440, `{{ defined(floatVar) }}`, "1", map[string]any{"floatVar": -2.5}},
		{450, `{{ defined(a, b, c) }}`, "0", map[string]any{"a": "hello", "c": "world"}},
		{460, `{{ defined(arr[0]) }}`, "1", map[string]any{"arr": []any{nil, ""}}},
		{470, `{{ defined(arr[0][0]) }}`, "1", map[string]any{"arr": []any{[]any{"nested"}}}},
		{480, `{{ defined(arr[99]) }}`, "0", map[string]any{"arr": []any{"first", "second"}}},
		{490, `{{ defined(999999999) }}`, "1", nil},
		{500, `{{ defined("hello") }}`, "1", nil},
		{510, `{{ defined("str", 123, true) }}`, "1", nil},
		{520, `{{ defined("str", undefinedVar, true) }}`, "0", nil},
		{530, `{{ defined(arr[-1]) }}`, "0", map[string]any{"arr": []any{"first"}}},
		{540, `{{ defined(0.000001) }}`, "1", nil},
		{550, `{{ defined(-0.000001) }}`, "1", nil},
		{
			580,
			`{{ defined(obj.arr[0].prop) }}`,
			"1",
			map[string]any{"obj": map[string]any{"arr": []any{map[string]any{"prop": "val"}}}},
		},
		{
			590,
			`{{ defined(obj.arr[0].prop) }}`,
			"0",
			map[string]any{"obj": map[string]any{"arr": []any{map[string]any{}}}},
		},
		{
			600,
			`{{ defined(obj.arr[5].prop) }}`,
			"0",
			map[string]any{"obj": map[string]any{"arr": []any{map[string]any{"prop": "val"}}}},
		},
		{
			610,
			`@if(defined(a, b) && defined(c))YES@end`,
			"YES",
			map[string]any{"a": "x", "b": "y", "c": "z"},
		},
		{
			620,
			`@if(defined(a, undefinedVar) && defined(c))YES@elseNO@end`,
			"NO",
			map[string]any{"a": "x", "c": "z"},
		},
		{630, `{{ defined(user.name) }}`, "0", map[string]any{"user": map[string]any{}}},
		{
			640,
			`{{ defined(user.name) }}`,
			"1",
			map[string]any{"user": map[string]any{"name": nil}},
		},
		{650, `{{ defined(user["name"]) }}`, "0", map[string]any{"user": map[string]any{}}},
		{
			660,
			`{{ defined(user["name"]) }}`,
			"1",
			map[string]any{"user": map[string]any{"name": nil}},
		},
		{670, `{{ defined(user.name.age) }}`, "0", map[string]any{"user": map[string]any{}}},
	}

	for _, tc := range cases {
		res, err := textwire.EvaluateString(tc.inp, tc.data)
		if err != nil {
			t.Fatalf("case %d: unexpected error: %s", tc.id, err)
		}

		if tc.expect != res {
			t.Errorf("case %d: wrong result. Expect %q, got %q", tc.id, tc.expect, res)
		}
	}
}

func TestHasValueCallExpression(t *testing.T) {
	cases := []struct {
		id     int
		inp    string
		expect string
		data   map[string]any
	}{
		{10, `{{ hasValue('') }}`, "0", nil},
		{20, `{{ hasValue("") }}`, "0", nil},
		{21, `{{ hasValue("", true) }}`, "0", nil},
		{30, `{{ hasValue(0) }}`, "0", nil},
		{40, `{{ hasValue(-0) }}`, "0", nil},
		{50, `{{ hasValue(1) }}`, "1", nil},
		{60, `{{ hasValue(0.0) }}`, "0", nil},
		{70, `{{ hasValue(-0.0) }}`, "0", nil},
		{80, `{{ hasValue(1.0) }}`, "1", nil},
		{90, `{{ hasValue({}) }}`, "0", nil},
		{100, `{{ hasValue([]) }}`, "0", nil},
		{110, `{{ hasValue(true) }}`, "1", nil},
		{110, `{{ hasValue(true, "nice") }}`, "1", nil},
		{120, `{{ hasValue(false) }}`, "0", nil},
		{130, `{{ hasValue(nil) }}`, "0", nil},
		{140, `{{ hasValue(undefinedVar) }}`, "0", nil},
		{150, `@if(!hasValue(definedVar))YES@end`, "YES", nil},
		{160, `{{ hasValue(emptyStr) }}`, "0", map[string]any{"emptyStr": ""}},
		{170, `{{ hasValue(zeroInt) }}`, "0", map[string]any{"zeroInt": 0}},
		{180, `{{ hasValue(zeroFloat) }}`, "0", map[string]any{"zeroFloat": 0.0}},
		{190, `{{ hasValue(falseVar) }}`, "0", map[string]any{"falseVar": false}},
		{200, `{{ hasValue(nilVar) }}`, "0", map[string]any{"nilVar": nil}},
		{210, `{{ hasValue(emptyObj) }}`, "0", map[string]any{"emptyObj": map[string]any{}}},
		{220, `{{ hasValue(emptyArr) }}`, "0", map[string]any{"emptyArr": []any{}}},
		{230, `{{ hasValue(definedVar) }}`, "1", map[string]any{"definedVar": "nice"}},
		{240, `{{ hasValue(nonEmptyStr) }}`, "1", map[string]any{"nonEmptyStr": "hello"}},
		{250, `{{ hasValue(positiveInt) }}`, "1", map[string]any{"positiveInt": 42}},
		{260, `{{ hasValue(positiveFloat) }}`, "1", map[string]any{"positiveFloat": 3.14}},
		{270, `{{ hasValue(trueVar) }}`, "1", map[string]any{"trueVar": true}},
		{
			280,
			`{{ hasValue(nonEmptyObj) }}`,
			"1",
			map[string]any{"nonEmptyObj": map[string]any{"key": "val"}},
		},
		{290, `{{ hasValue(nonEmptyArr) }}`, "1", map[string]any{"nonEmptyArr": []any{1, 2, 3}}},
		{300, `@if(hasValue(definedVar))YES@end`, "YES", map[string]any{"definedVar": "nice"}},
		{310, `@if(hasValue(emptyStr))YES@elseNO@end`, "NO", map[string]any{"emptyStr": ""}},
		{320, `@if(hasValue(nilVar))YES@elseNO@end`, "NO", map[string]any{"nilVar": nil}},
		{330, `@if(hasValue(zeroInt))YES@elseNO@end`, "NO", map[string]any{"zeroInt": 0}},
		{340, `@if(hasValue(falseVar))YES@elseNO@end`, "NO", map[string]any{"falseVar": false}},
		{
			350,
			`{{ hasValue(definedVar).then("Yes", "No") }}`,
			"Yes",
			map[string]any{"definedVar": "nice"},
		},
		{360, `{{ hasValue(emptyStr).then("Yes", "No") }}`, "No", map[string]any{"emptyStr": ""}},
		{
			370,
			`{{ hasValue(obj.prop) }}`,
			"1",
			map[string]any{"obj": map[string]any{"prop": "value"}},
		},
		{380, `{{ hasValue(obj.prop) }}`, "0", map[string]any{"obj": map[string]any{"prop": ""}}},
		{390, `{{ hasValue(obj.prop) }}`, "0", map[string]any{"obj": map[string]any{"prop": 0}}},
		{400, `{{ hasValue(obj.prop) }}`, "0", map[string]any{"obj": map[string]any{}}},
		{410, `{{ hasValue(obj.missing) }}`, "0", map[string]any{"obj": map[string]any{}}},
		{
			420,
			`{{ hasValue(obj.nested.prop) }}`,
			"1",
			map[string]any{"obj": map[string]any{"nested": map[string]any{"prop": "value"}}},
		},
		{
			430,
			`{{ hasValue(obj.nested.prop) }}`,
			"0",
			map[string]any{"obj": map[string]any{"nested": map[string]any{"prop": nil}}},
		},
		{440, `{{ hasValue(arr[0]) }}`, "1", map[string]any{"arr": []any{"first", "second"}}},
		{450, `{{ hasValue(arr[0]) }}`, "0", map[string]any{"arr": []any{""}}},
		{460, `{{ hasValue(arr[0]) }}`, "0", map[string]any{"arr": []any{}}},
		{
			470,
			`@if(defined(definedVar) && hasValue(definedVar))YES@end`,
			"YES",
			map[string]any{"definedVar": "nice"},
		},
		{
			480,
			`@if(defined(definedVar) && hasValue(definedVar))YES@elseNO@end`,
			"NO",
			map[string]any{"definedVar": ""},
		},
		{
			490,
			`@if(hasValue(obj.prop) && hasValue(obj.prop2))YES@end`,
			"YES",
			map[string]any{"obj": map[string]any{"prop": "a", "prop2": "b"}},
		},
		{
			500,
			`@if(hasValue(obj.prop) || hasValue(obj.prop2))YES@end`,
			"YES",
			map[string]any{"obj": map[string]any{"prop": "a"}},
		},
		{510, `{{ hasValue(a, b) }}`, "1", map[string]any{"a": "hello", "b": "world"}},
		{520, `{{ hasValue(a, b) }}`, "0", map[string]any{"a": "hello", "b": ""}},
		{530, `{{ hasValue(a, b) }}`, "0", map[string]any{"a": "", "b": "hello"}},
		{540, `{{ hasValue(a, b) }}`, "0", map[string]any{"a": "", "b": ""}},
		{550, `{{ hasValue(a, b) }}`, "0", map[string]any{"a": 0, "b": 1}},
		{560, `{{ hasValue(a, b) }}`, "0", map[string]any{"a": 1, "b": 0}},
		{570, `{{ hasValue(a, b) }}`, "0", map[string]any{"a": nil, "b": "hello"}},
		{580, `{{ hasValue(a, b, c) }}`, "1", map[string]any{"a": "a", "b": "b", "c": "c"}},
		{590, `{{ hasValue(a, b, c) }}`, "0", map[string]any{"a": "a", "b": "", "c": "c"}},
		{600, `@if(hasValue(a, b))YES@end`, "YES", map[string]any{"a": "x", "b": "y"}},
		{610, `@if(hasValue(a, b))YES@elseNO@end`, "NO", map[string]any{"a": "x", "b": ""}},
		{
			620,
			`@if(hasValue(obj.prop, obj.prop2))YES@end`,
			"YES",
			map[string]any{"obj": map[string]any{"prop": "a", "prop2": "b"}},
		},
		{
			630,
			`@if(hasValue(obj.prop, obj.prop2))YES@elseNO@end`,
			"NO",
			map[string]any{"obj": map[string]any{"prop": "a"}},
		},
	}

	for _, tc := range cases {
		res, err := textwire.EvaluateString(tc.inp, tc.data)
		if err != nil {
			t.Fatalf("case %d: unexpected error: %s", tc.id, err)
		}

		if tc.expect != res {
			t.Errorf("case %d: wrong result. Expect %q, got %q", tc.id, tc.expect, res)
		}
	}
}

func TestErrorHandling(t *testing.T) {
	cases := []struct {
		inp  string
		err  *fail.Error
		data map[string]any
	}{
		{
			inp:  "@for(i 0; i < 2; i++)x@end",
			err:  fail.New(nil, "", fail.OriginPars, fail.ErrForLoopExpectInitStmt, "i"),
			data: nil,
		},
		{
			inp: `{{ defined(name.undefinedFunc()) }}`,
			err: fail.New(
				nil,
				"",
				fail.OriginEval,
				fail.ErrFuncNotDefined,
				value.STR_VAL,
				"undefinedFunc",
			),
			data: map[string]any{"name": "Anna"},
		},
		{
			inp:  `@use("someTemplate")`,
			err:  fail.New(nil, "", fail.OriginEval, fail.ErrTemplateDirectives),
			data: nil,
		},
		{
			inp:  `@insert("title", "hi")`,
			err:  fail.New(nil, "", fail.OriginEval, fail.ErrTemplateDirectives),
			data: nil,
		},
		{
			inp:  `@reserve("content")`,
			err:  fail.New(nil, "", fail.OriginEval, fail.ErrTemplateDirectives),
			data: nil,
		},
		{
			inp:  `@component("~header")@end`,
			err:  fail.New(nil, "", fail.OriginEval, fail.ErrTemplateDirectives),
			data: nil,
		},
		{
			inp: `{{ 1 + "a" }}`,
			err: fail.New(
				nil,
				"",
				fail.OriginEval,
				fail.ErrCannotUseOperator,
				"+",
				value.INT_VAL,
				"+",
				value.STR_VAL,
			),
			data: nil,
		},
		{
			inp:  `{{ loop = "test" }}`,
			err:  fail.New(nil, "", fail.OriginEval, fail.ErrReservedIdentifiers),
			data: nil,
		},
		{
			inp:  `{{ global = "test" }}`,
			err:  fail.New(nil, "", fail.OriginEval, fail.ErrReservedIdentifiers),
			data: nil,
		},
		{
			inp:  `{{ global.age = 20 }}`,
			err:  fail.New(nil, "", fail.OriginEval, fail.ErrReadOnlyValue),
			data: nil,
		},
		{
			inp:  `{{ raw() }}`,
			err:  fail.New(nil, "", fail.OriginPars, fail.ErrGlobalFuncFewArgs, "raw", 1, 0),
			data: nil,
		},
		{
			inp:  `{{ raw("a", "b") }}`,
			err:  fail.New(nil, "", fail.OriginPars, fail.ErrGlobalFuncLotsOfArgs, "raw", 1, 2),
			data: nil,
		},
		{
			inp:  `{{ raw(undefinedVar) }}`,
			err:  fail.New(nil, "", fail.OriginEval, fail.ErrVariableIsUndefined, "undefinedVar"),
			data: nil,
		},
		{
			inp: `{{ loop }}`,
			err: fail.New(
				nil,
				"",
				fail.OriginEval,
				fail.ErrReservedIdentifiers,
			), data: map[string]any{"loop": "test"},
		},
		{
			inp: `{{ n = 1; n = "test" }}`,
			err: fail.New(
				nil,
				"",
				fail.OriginEval,
				fail.ErrIdentTypeMismatch,
				"n",
				value.INT_VAL,
				value.STR_VAL,
			),
			data: nil,
		},
		{
			inp:  `{{ user = {}; user.address.zip }}`,
			err:  fail.New(nil, "", fail.OriginEval, fail.ErrKeyOnNonObj, value.NIL_VAL, "zip"),
			data: nil,
		},
		{
			inp:  `@each(i in [1]){{ loop.parent.parent }}@end`,
			err:  fail.New(nil, "", fail.OriginEval, fail.ErrKeyOnNonObj, value.NIL_VAL, "parent"),
			data: nil,
		},
		{
			inp:  `@for(i = 0; i < 1; i++){{ loop.parent.parent }}@end`,
			err:  fail.New(nil, "", fail.OriginEval, fail.ErrKeyOnNonObj, value.NIL_VAL, "parent"),
			data: nil,
		},
		{
			inp: `{{ 5.someFunction() }}`,
			err: fail.New(
				nil,
				"",
				fail.OriginEval,
				fail.ErrFuncNotDefined,
				value.INT_VAL,
				"someFunction",
			),
			data: nil,
		},
		{
			inp:  `{{ 3 / 0 }}`,
			err:  fail.New(nil, "", fail.OriginEval, fail.ErrDivisionByZero),
			data: nil,
		},
		{
			inp:  `{{ 3 % 0 }}`,
			err:  fail.New(nil, "", fail.OriginEval, fail.ErrDivisionByZero),
			data: nil,
		},
		{
			inp:  `{{ 1.0 / 0.0 }}`,
			err:  fail.New(nil, "", fail.OriginEval, fail.ErrDivisionByZero),
			data: nil,
		},
		{
			inp:  `{{ 0.0 / 0.0 }}`,
			err:  fail.New(nil, "", fail.OriginEval, fail.ErrDivisionByZero),
			data: nil,
		},
		{
			inp:  `{{ undefinedVar }}`,
			err:  fail.New(nil, "", fail.OriginEval, fail.ErrVariableIsUndefined, "undefinedVar"),
			data: nil,
		},
		{
			inp:  `{{ obj = {name: "Amy"}; obj.name.id }}`,
			err:  fail.New(nil, "", fail.OriginEval, fail.ErrKeyOnNonObj, value.STR_VAL, "id"),
			data: nil,
		},
		{
			inp:  `@each(v in {}){{ v }}@end`,
			err:  fail.New(nil, "", fail.OriginEval, fail.ErrEachDirWithNonArrArg, value.OBJ_VAL),
			data: nil,
		},
		{
			inp:  `{{ 1 = 10 }}`,
			err:  fail.New(nil, "", fail.OriginEval, fail.ErrNotSupportedAssign, value.INT_VAL),
			data: nil,
		},
		{
			inp:  `{{ {} = 10 }}`,
			err:  fail.New(nil, "", fail.OriginEval, fail.ErrNotSupportedAssign, value.OBJ_VAL),
			data: nil,
		},
		{
			inp:  `{{ [] = 10 }}`,
			err:  fail.New(nil, "", fail.OriginEval, fail.ErrNotSupportedAssign, value.ARR_VAL),
			data: nil,
		},
		{
			inp:  `{{ 1.1 = 10 }}`,
			err:  fail.New(nil, "", fail.OriginEval, fail.ErrNotSupportedAssign, value.FLOAT_VAL),
			data: nil,
		},
		{
			inp:  `{{ true = 10 }}`,
			err:  fail.New(nil, "", fail.OriginEval, fail.ErrNotSupportedAssign, value.BOOL_VAL),
			data: nil,
		},
		{
			inp:  `{{ 'anna' = 10 }}`,
			err:  fail.New(nil, "", fail.OriginEval, fail.ErrNotSupportedAssign, value.STR_VAL),
			data: nil,
		},
		{
			inp:  `{{ "serhii" = 10 }}`,
			err:  fail.New(nil, "", fail.OriginEval, fail.ErrNotSupportedAssign, value.STR_VAL),
			data: nil,
		},
		{
			inp:  `{{ false = 10 }}`,
			err:  fail.New(nil, "", fail.OriginEval, fail.ErrNotSupportedAssign, value.BOOL_VAL),
			data: nil,
		},
		{
			inp:  `{{ nil = 10 }}`,
			err:  fail.New(nil, "", fail.OriginEval, fail.ErrNotSupportedAssign, value.NIL_VAL),
			data: nil,
		},
		{
			inp: `{{ -x = 2 }}`,
			err: fail.New(
				nil,
				"",
				fail.OriginEval,
				fail.ErrNotSupportedAssign,
				"prefix expression",
			),
			data: nil,
		},
		{
			inp: `{{ a = [1]; a[0] + 1 = 3 }}`,
			err: fail.New(
				nil,
				"",
				fail.OriginEval,
				fail.ErrNotSupportedAssign,
				"infix expression",
			),
			data: nil,
		},
		{
			inp: `{{ x = 1; x + 1 = 2 }}`,
			err: fail.New(
				nil,
				"",
				fail.OriginEval,
				fail.ErrNotSupportedAssign,
				"infix expression",
			),
			data: nil,
		},
		{
			inp: `{{ obj = {}; obj[0] = "Anna" }}`,
			err: fail.New(
				nil,
				"",
				fail.OriginEval,
				fail.ErrObjIndexStr,
				value.INT_VAL,
			),
			data: nil,
		},
		{
			inp: `{{ arr = [1]; arr["name"] = 2 }}`,
			err: fail.New(
				nil,
				"",
				fail.OriginEval,
				fail.ErrArrIndexInt,
				value.STR_VAL,
			),
			data: nil,
		},
		{
			inp: `{{ obj = {}; obj[nil] = "Skirk" }}`,
			err: fail.New(
				nil,
				"",
				fail.OriginEval,
				fail.ErrObjIndexStr,
				value.NIL_VAL,
			),
			data: nil,
		},
		{
			inp: `{{ obj = {}; obj[0] }}`,
			err: fail.New(
				nil,
				"",
				fail.OriginEval,
				fail.ErrObjIndexStr,
				value.INT_VAL,
			),
			data: nil,
		},
		// Global functions
		{
			inp: `{{ formatDate(date, '2006-01-02') }}`,
			err: fail.New(
				nil,
				"",
				fail.OriginEval,
				fail.ErrFormatDateWrongDate,
				"2026-03-24 13:03",
			),
			data: map[string]any{"date": "2026-03-24 13:03"},
		},
		{
			inp: `{{ formatDate(date, '2006-01-02') }}`,
			err: fail.New(
				nil,
				"",
				fail.OriginEval,
				fail.ErrFormatDateWrongDate,
				"2026-03-4 13:03:00",
			),
			data: map[string]any{"date": "2026-03-4 13:03:00"},
		},
		{
			inp: `{{ formatDate(date, '2006-01-02') }}`,
			err: fail.New(
				nil,
				"",
				fail.OriginEval,
				fail.ErrFormatDateWrongDate,
				" 2026-03-34 13:03:00",
			),
			data: map[string]any{"date": " 2026-03-34 13:03:00"},
		},
		{
			inp: "x",
			err: fail.New(
				nil,
				"",
				fail.OriginTpl,
				fail.ErrUnsupportedFieldType,
				"func()",
				"Fn",
			),
			data: map[string]any{
				"u": namedUserWithHooks{
					Name: "Bob",
					Role: "editor",
					Fn:   func() {},
				},
			},
		},
		{
			inp:  "x",
			err:  fail.New(nil, "", fail.OriginTpl, fail.ErrUnsupportedType, "func()"),
			data: map[string]any{"callback": func() {}},
		},
		{
			inp: "x",
			err: fail.New(
				nil,
				"",
				fail.OriginTpl,
				fail.ErrUnsupportedFieldType,
				"chan int",
				"Ch",
			),
			data: map[string]any{"u": namedUserWithChan{Name: "Bob", Ch: make(chan int)}},
		},
		{
			inp: "x",
			err: fail.New(
				nil,
				"",
				fail.OriginTpl,
				fail.ErrUnsupportedFieldType,
				"complex128",
				"Num",
			),
			data: map[string]any{"u": namedUserWithComplex{Name: "Bob", Num: complex(1, 2)}},
		},
		{
			inp:  "x",
			err:  fail.New(nil, "", fail.OriginTpl, fail.ErrUintOverflow, uint64(math.MaxUint64)),
			data: map[string]any{"big": uint64(math.MaxUint64)},
		},
		{
			inp: "x",
			err: fail.New(
				nil,
				"",
				fail.OriginTpl,
				fail.ErrUintOverflowField,
				uint64(math.MaxUint64),
				"Big",
			),
			data: map[string]any{"u": namedUserWithUint{Name: "Bob", Big: math.MaxUint64}},
		},
	}

	for _, tc := range cases {
		_, err := textwire.EvaluateString(tc.inp, tc.data)
		if err == nil {
			t.Fatalf("expected error, got none")
		}

		if err.String() != tc.err.String() {
			t.Errorf("wrong error message. Expected:\n%q\ngot:\n%q", tc.err, err)
		}
	}
}

func TestJSONEscaping(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		data   map[string]any
		expect string
	}{
		{
			name:   "html characters are not escaped",
			inp:    `{!! {k: x}.json() !!}`,
			data:   map[string]any{"x": `</script><script>alert(1)</script>`},
			expect: `{"k":"</script><script>alert(1)</script>"}`,
		},
		{
			name:   "control character is escaped",
			inp:    `{!! [x].json() !!}`,
			data:   map[string]any{"x": "a\x01b"},
			expect: `["a\u0001b"]`,
		},
		{
			name:   "invalid utf8 is replaced",
			inp:    `{!! [x].json() !!}`,
			data:   map[string]any{"x": "a\xffb"},
			expect: `["a\ufffdb"]`,
		},
		{
			name:   "object key is escaped",
			inp:    `{!! {"a\"b": 1}.json() !!}`,
			data:   nil,
			expect: `{"a\"b":1}`,
		},
		{
			name:   "escaped braces html-escape the json",
			inp:    `{{ {k: x}.json() }}`,
			data:   map[string]any{"x": `<&>`},
			expect: `{&#34;k&#34;:&#34;&lt;&amp;&gt;&#34;}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assertEvalOutput(t, tc.inp, tc.data, tc.expect)
		})
	}
}

func TestDumpEscapesHTML(t *testing.T) {
	t.Run("string value", func(t *testing.T) {
		got, err := textwire.EvaluateString(`@dump(x)`, map[string]any{
			"x": `</span><script>alert(1)</script>`,
		})
		if err != nil {
			t.Fatal(err)
		}

		if strings.Contains(got, "<script>alert(1)</script>") {
			t.Fatalf("dump output contains a raw script tag:\n%s", got)
		}

		if !strings.Contains(got, "&lt;script&gt;alert(1)&lt;/script&gt;") {
			t.Fatalf("dump output doesn't contain the escaped script tag:\n%s", got)
		}
	})

	t.Run("object key and value", func(t *testing.T) {
		got, err := textwire.EvaluateString(`@dump(x)`, map[string]any{
			"x": map[string]any{"<b>key</b>": "<i>value</i>"},
		})
		if err != nil {
			t.Fatal(err)
		}

		if strings.Contains(got, "<b>key</b>") || strings.Contains(got, "<i>value</i>") {
			t.Fatalf("dump output contains raw HTML from the object:\n%s", got)
		}

		if !strings.Contains(got, "&lt;b&gt;key&lt;/b&gt;") {
			t.Fatalf("dump output doesn't contain the escaped key:\n%s", got)
		}

		if !strings.Contains(got, "&lt;i&gt;value&lt;/i&gt;") {
			t.Fatalf("dump output doesn't contain the escaped value:\n%s", got)
		}
	})
}

func TestEvaluateFile(t *testing.T) {
	absPath, fileErr := file.ToFullPath("testdata/good/before/two-vars-no-use/index.tw")
	if fileErr != nil {
		t.Errorf("error getting full path: %s", fileErr)
		return
	}

	actual, err := textwire.EvaluateFile(absPath, map[string]any{
		"title":   "Textwire is Awesome",
		"visible": true,
	})

	if err != nil {
		t.Errorf("error evaluating file: %q", err)
	}

	expect, err2 := readFile("testdata/good/expected/two-vars-no-use.html")
	if err2 != nil {
		t.Errorf("error reading file: %q", err2)
		return
	}

	if actual != expect {
		t.Errorf("wrong output. Expect:\n%s\ngot:\n%s", expect, actual)
	}
}
