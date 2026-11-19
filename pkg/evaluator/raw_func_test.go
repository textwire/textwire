package evaluator

import "testing"

func TestEvalRawGlobalFunc(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		// Strings
		{
			name:   "raw string is not escaped",
			inp:    `{{ raw('<b>') }}`,
			expect: "<b>",
		},
		{
			name:   "raw string in raw braces is not escaped",
			inp:    `{!! raw('<b>') !!}`,
			expect: "<b>",
		},
		{
			name:   "raw called twice keeps string raw",
			inp:    `{{ raw(raw('<b>')) }}`,
			expect: "<b>",
		},
		{
			name:   "already escaped text is kept the same",
			inp:    `{{ raw('&amp;') }}`,
			expect: "&amp;",
		},
		{
			name:   "raw string next to plain string escapes only plain one",
			inp:    `{{ raw('<b>'); '<b>' }}`,
			expect: "<b>&lt;b&gt;",
		},
		{
			name:   "one embedded escapes only plain string",
			inp:    `{{ html = '<b>'; raw(html); html }}`,
			expect: "<b>&lt;b&gt;",
		},
		{
			name:   "raw of concatenation is not escaped",
			inp:    `{{ raw('<i>' + '<b>') }}`,
			expect: "<i><b>",
		},
		{
			name:   "ternary returning raw string is not escaped",
			inp:    `{{ html = '<b>'; true ? raw(html) : html }}`,
			expect: "<b>",
		},
		{
			name:   "ternary returning plain string is escaped",
			inp:    `{{ html = '<b>'; false ? raw(html) : html }}`,
			expect: "&lt;b&gt;",
		},
		// Other literal types
		{
			name:   "integer is unchanged",
			inp:    `{{ raw(5) + 1 }}`,
			expect: "6",
		},
		{
			name:   "float is unchanged",
			inp:    `{{ raw(1.5) }}`,
			expect: "1.5",
		},
		{
			name:   "booleans are unchanged",
			inp:    `{{ raw(true) }}|{{ raw(false) }}`,
			expect: "1|0",
		},
		{
			name:   "nil prints nothing",
			inp:    `[{{ raw(nil) }}]`,
			expect: "[]",
		},
		// Arrays and objects
		{
			name:   "empty string and array print nothing, empty object prints braces",
			inp:    `{{ raw(''); raw([]); raw({}) }}`,
			expect: "{}",
		},
		{
			name:   "raw array strings are not escaped",
			inp:    `{{ raw(['<a>', 1, '<b>']) }}`,
			expect: "<a>, 1, <b>",
		},
		{
			name:   "array escapes only plain strings",
			inp:    `{{ html = '<b>'; [raw(html), html] }}`,
			expect: "<b>, &lt;b&gt;",
		},
		{
			name:   "raw object strings are not escaped",
			inp:    `{{ raw({ a: '<a>', n: 1 }) }}`,
			expect: `{a: "<a>", n: 1}`,
		},
		{
			name:   "object escapes only plain strings",
			inp:    `{{ html = '<b>'; { a: raw(html), b: html } }}`,
			expect: `{a: "<b>", b: &#34;&lt;b&gt;&#34;}`,
		},
		{
			name:   "object key stays escaped",
			inp:    `{{ raw({ '<k>': '<b>' }) }}`,
			expect: `{&lt;k&gt;: "<b>"}`,
		},
		{
			name:   "deeply nested strings are not escaped",
			inp:    `{{ html = '<b>'; raw({a: {b: [html, {c: html}]}}) }}`,
			expect: `{a: {b: <b>, {c: "<b>"}}}`,
		},
		{
			name:   "array escapes already escaped text only when plain",
			inp:    `{{ amp = '&amp;'; [raw(amp), amp] }}`,
			expect: "&amp;, &amp;amp;",
		},
		{
			name:   "object escapes already escaped text only when plain",
			inp:    `{{ {a: raw('&amp;'), b: '&amp;'} }}`,
			expect: `{a: "&amp;", b: &#34;&amp;amp;&#34;}`,
		},
		{
			name:   "raw braces print raw and plain strings the same",
			inp:    `{!! [raw('<b>'), '<b>']; ' | '; {a: raw('<b>'), b: '<b>'} !!}`,
			expect: `<b>, <b> | {a: "<b>", b: "<b>"}`,
		},
		// Nested access
		{
			name:   "nested string by index stays raw",
			inp:    `{{ raw([['<a>', '<b>']])[0][1] }}`,
			expect: "<b>",
		},
		{
			name:   "nested string by property stays raw",
			inp:    `{{ raw({user: {name: '<b>'}}).user.name }}`,
			expect: "<b>",
		},
		{
			name:   "nested string by get function stays raw",
			inp:    `{{ raw({user: {name: '<b>'}}).get('user.name') }}`,
			expect: "<b>",
		},
		{
			name:   "each over raw array gives raw items",
			inp:    `@each(item in raw(['<a>', '<b>']))[{{ item }}]@end`,
			expect: "[<a>][<b>]",
		},
		{
			name:   "raw item in each over plain array is not escaped",
			inp:    `@each(item in ['<a>', '<b>']){{ raw(item); item }};@end`,
			expect: "<a>&lt;a&gt;;<b>&lt;b&gt;;",
		},
		// JSON
		{
			name:   "json in raw braces ignores raw mark",
			inp:    `{!! raw({a: '<a>'}).json() !!}`,
			expect: `{"a":"<a>"}`,
		},
		{
			name:   "json of raw object in braces is escaped",
			inp:    `{{ raw({a: '<a>'}).json() }}`,
			expect: `{&#34;a&#34;:&#34;&lt;a&gt;&#34;}`,
		},
		{
			name:   "raw of json result is not escaped",
			inp:    `{{ raw({a: '<a>'}.json()) }}`,
			expect: `{"a":"<a>"}`,
		},
		{
			name:   "json of array ignores raw mark in raw braces",
			inp:    `{!! html = '<b>'; [raw(html), html].json() !!}`,
			expect: `["<b>","<b>"]`,
		},
		// Other global functions
		{
			name:   "defined is false for raw missing property",
			inp:    `{{ defined(raw({a: '<a>'}.missing)) }}`,
			expect: "0",
		},
		{
			name:   "defined is true for raw nil",
			inp:    `{{ defined(raw(nil)) }}`,
			expect: "1",
		},
		{
			name:   "hasValue checks raw values like plain ones",
			inp:    `{{ hasValue(raw('')); '|'; hasValue(raw(['<a>', '<b>'])) }}`,
			expect: "0|1",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			evaluationExpected(t, tc.inp, tc.expect, nil, "")
		})
	}
}

func TestEvalRawAssign(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		// Variables
		{
			name:   "raw string in variable stays raw",
			inp:    `{{ unsafe = raw('<b>'); unsafe }}`,
			expect: "<b>",
		},
		{
			name:   "variable reassigned to its raw version is not escaped",
			inp:    `{{ x = '<b>'; x = raw(x); x }}`,
			expect: "<b>",
		},
		{
			name:   "variable reassigned to its plain version is escaped",
			inp:    `{{ x = raw('<b>'); x = '<b>'; x }}`,
			expect: "&lt;b&gt;",
		},
		{
			name:   "original array stays escaped",
			inp:    `{{ arr = ['<a>']; unsafe = raw(arr); arr }}`,
			expect: "&lt;a&gt;",
		},
		{
			name:   "variable sharing original array stays escaped",
			inp:    `{{ arr = ['<a>']; same = arr; unsafe = raw(arr); same }}`,
			expect: "&lt;a&gt;",
		},
		// Array elements
		{
			name:   "plain array keeps raw string assigned to it",
			inp:    `{{ x = '<b>'; list = [x, x]; list[1] = raw(list[1]); list }}`,
			expect: "&lt;b&gt;, <b>",
		},
		{
			name:   "raw array escapes plain string assigned to it",
			inp:    `{{ x = '<b>'; list = raw([x, x]); list[1] = x; list }}`,
			expect: "<b>, &lt;b&gt;",
		},
		{
			name:   "array keeps raw variable assigned to it",
			inp:    `{{ x = '<b>'; unsafe = raw(x); list = [x, x]; list[0] = unsafe; list }}`,
			expect: "<b>, &lt;b&gt;",
		},
		// Object properties
		{
			name:   "plain object keeps raw string assigned by property notation",
			inp:    `{{ x = '<b>'; obj = {a: x}; obj.b = raw(x); obj }}`,
			expect: `{a: &#34;&lt;b&gt;&#34;, b: "<b>"}`,
		},
		{
			name:   "plain object keeps raw string assigned by index notation",
			inp:    `{{ x = '<b>'; obj = {a: x}; obj['b'] = raw(x); obj }}`,
			expect: `{a: &#34;&lt;b&gt;&#34;, b: "<b>"}`,
		},
		{
			name:   "raw string replacing plain property is not escaped",
			inp:    `{{ x = '<b>'; obj = {a: x}; obj.a = raw(x); obj.a }}`,
			expect: "<b>",
		},
		{
			name:   "plain string replacing raw property is escaped",
			inp:    `{{ x = '<b>'; obj = raw({a: x}); obj.a = x; obj.a }}`,
			expect: "&lt;b&gt;",
		},
		{
			name:   "object literal keeps raw variable",
			inp:    `{{ x = '<b>'; unsafe = raw(x); obj = {a: unsafe, b: x}; obj.a; obj.b }}`,
			expect: "<b>&lt;b&gt;",
		},
		// Nested values
		{
			name:   "nested property keeps raw string assigned to it",
			inp:    `{{ obj = {inner: {}}; obj.inner.v = raw('<b>'); obj.inner.v }}`,
			expect: `<b>`,
		},
		{
			name:   "array in object keeps raw string assigned to it",
			inp:    `{{ x = '<b>'; obj = {list: [x, x]}; obj.list[0] = raw(x); obj }}`,
			expect: "{list: <b>, &lt;b&gt;}",
		},
		{
			name:   "object in array keeps raw string assigned to it",
			inp:    `{{ x = '<b>'; list = [{a: x}]; list[0].a = raw(x); list }}`,
			expect: `{a: "<b>"}`,
		},
		// Shared arrays and copies
		{
			name:   "raw element assigned through shared array shows in both",
			inp:    `{{ x = '<b>'; a = [x]; b = a; b[0] = raw(x); a }}`,
			expect: "<b>",
		},
		{
			name:   "plain element assigned through shared raw array shows in both",
			inp:    `{{ x = '<b>'; a = raw([x, x]); b = a; b[0] = x; a }}`,
			expect: "&lt;b&gt;, <b>",
		},
		{
			name:   "changing raw copy leaves original array unchanged",
			inp:    `{{ list = ['<b>']; unsafe = raw(list); unsafe[0] = 'z'; list; '|'; unsafe }}`,
			expect: "&lt;b&gt;|z",
		},
		{
			name:   "changing nested raw copy leaves original object unchanged",
			inp:    `{{ obj = {u: {n: '<b>'}}; unsafe = raw(obj); unsafe.u.n = 'z'; obj.u.n; '|'; unsafe.u.n }}`,
			expect: "&lt;b&gt;|z",
		},
		// Loops
		{
			name:   "raw elements assigned in each loop are not escaped",
			inp:    `{{ x = '<b>'; list = [x, x] }}@each(v in list){{ list[loop.index] = raw(v) }}@end{{ list }}`,
			expect: "<b>, <b>",
		},
		{
			name:   "raw element assigned in for loop is not escaped",
			inp:    `{{ x = '<b>'; list = [x, x] }}@for(i = 0; i < 1; i++){{ list[i] = raw(list[i]) }}@end{{ list }}`,
			expect: "<b>, &lt;b&gt;",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			evaluationExpected(t, tc.inp, tc.expect, nil, "")
		})
	}
}

func TestEvalRawWithFunctions(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		// Functions that return a new string
		{
			name:   "upper makes new escaped string",
			inp:    `{{ raw('<b>').upper() }}`,
			expect: "&lt;B&gt;",
		},
		{
			name:   "raw of upper result is not escaped",
			inp:    `{{ raw('<b>'.upper()) }}`,
			expect: "<B>",
		},
		{
			name:   "concatenation makes new escaped string",
			inp:    `{{ x = '<b>'; raw(x) + x }}`,
			expect: "&lt;b&gt;&lt;b&gt;",
		},
		{
			name:   "at makes new escaped string",
			inp:    `{{ raw('<b>').at(0) }}`,
			expect: "&lt;",
		},
		{
			name:   "join makes new escaped string",
			inp:    `{{ raw(['<a>', '<b>']).join(' ') }}`,
			expect: "&lt;a&gt; &lt;b&gt;",
		},
		{
			name:   "len of raw string counts its text",
			inp:    `{{ raw('<b>').len() }}`,
			expect: "3",
		},
		// Function contains
		{
			"contains matches raw array element to plain array",
			`{{ html = '<br>'; [[raw(html)]].contains([html]) }}`,
			"1",
		},
		{
			"contains matches raw object element to plain object",
			`{{ obj = { html: '<br>' }; [raw(obj)].contains(obj) }}`,
			"1",
		},
		{
			"contains in raw array finds plain string",
			`{{ raw(['<a>', '<b>', '<c>']).contains('<b>') }}`,
			"1",
		},
		// Functions that keep the elements
		{
			name:   "reverse keeps raw strings",
			inp:    `{{ raw(['<a>', '<b>']).reverse() }}`,
			expect: "<b>, <a>",
		},
		{
			name:   "slice keeps raw strings",
			inp:    `{{ raw(['<a>', '<b>']).slice(1) }}`,
			expect: "<b>",
		},
		{
			name:   "rand returns raw string",
			inp:    `{{ raw(['<a>']).rand() }}`,
			expect: "<a>",
		},
		{
			name:   "append to raw array escapes only new plain string",
			inp:    `{{ raw(['<a>', '<b>']).append('<b>') }}`,
			expect: "<a>, <b>, &lt;b&gt;",
		},
		{
			name:   "prepend to plain array keeps new raw string",
			inp:    `{{ ['<a>', '<b>'].prepend(raw('<b>')) }}`,
			expect: "<b>, &lt;a&gt;, &lt;b&gt;",
		},
		{
			name:   "camel keeps raw strings",
			inp:    `{{ x = '<b>'; raw({first_name: x, inner: {last_name: x}}).camel() }}`,
			expect: `{firstName: "<b>", inner: {lastName: "<b>"}}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			evaluationExpected(t, tc.inp, tc.expect, nil, "")
		})
	}
}

func TestEvalRawComparison(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		{
			name:   "raw string == plain string",
			inp:    `{{ html = '<b>'; raw(html) == html }}`,
			expect: "1",
		},
		{
			name:   "raw string != plain string is false",
			inp:    `{{ html = '<b>'; raw(html) != html }}`,
			expect: "0",
		},
		{
			name:   "raw array == plain array",
			inp:    `{{ arr = ['<a>', '<b>']; raw(arr) == arr }}`,
			expect: "1",
		},
		{
			name:   "raw array != plain array is false",
			inp:    `{{ arr = ['<a>', '<b>']; raw(arr) != arr }}`,
			expect: "0",
		},
		{
			name:   "raw object == plain object",
			inp:    `{{ obj = {a: '<a>'}; raw(obj) == obj }}`,
			expect: "1",
		},
		{
			name:   "array with raw string == array with plain string",
			inp:    `{{ html = '<b>'; [raw(html)] == [html] }}`,
			expect: "1",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			evaluationExpected(t, tc.inp, tc.expect, nil, "")
		})
	}
}
