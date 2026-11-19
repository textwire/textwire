package value

import "testing"

func TestEscapedString(t *testing.T) {
	cases := []struct {
		name   string
		val    Literal
		expect string
	}{
		{
			name:   "string is escaped",
			val:    &Str{Val: `<b>"Tom" & 'Jerry'</b>`},
			expect: `&lt;b&gt;&#34;Tom&#34; &amp; &#39;Jerry&#39;&lt;/b&gt;`,
		},
		{
			name:   "raw string is not escaped",
			val:    &Str{Val: `<b>"Tom" & 'Jerry'</b>`, raw: true},
			expect: `<b>"Tom" & 'Jerry'</b>`,
		},
		{
			name:   "integer is unchanged",
			val:    &Int{Val: 12},
			expect: "12",
		},
		{
			name:   "nil is empty",
			val:    new(Nil),
			expect: "",
		},
		{
			name:   "array strings are escaped",
			val:    &Arr{Elements: []Literal{&Str{Val: "<a>"}, &Str{Val: "b&c"}}},
			expect: "&lt;a&gt;, b&amp;c",
		},
		{
			name:   "array escapes only plain strings",
			val:    &Arr{Elements: []Literal{&Str{Val: "<a>", raw: true}, &Str{Val: "<b>"}}},
			expect: "<a>, &lt;b&gt;",
		},
		{
			name:   "object keys and strings are escaped",
			val:    NewObj(map[string]Literal{"<k>": &Str{Val: "<v>"}, "n": &Int{Val: 1}}),
			expect: "{&lt;k&gt;: &#34;&lt;v&gt;&#34;, n: 1}",
		},
		{
			name:   "object escapes key but not raw string",
			val:    NewObj(map[string]Literal{"<k>": &Str{Val: "<v>", raw: true}}),
			expect: `{&lt;k&gt;: "<v>"}`,
		},
		{
			name: "object in array escapes only plain strings",
			val: &Arr{Elements: []Literal{NewObj(map[string]Literal{
				"a": &Str{Val: "<a>", raw: true},
				"b": &Str{Val: "<b>"},
			})}},
			expect: `{a: "<a>", b: &#34;&lt;b&gt;&#34;}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := escapedString(tc.val); got != tc.expect {
				t.Fatalf("expect %q, got %q", tc.expect, got)
			}
		})
	}
}

func TestEscapedString_WithArrayThatContainsItself(t *testing.T) {
	arr := &Arr{Elements: []Literal{&Str{Val: "<a>", raw: true}}}
	arr.Elements = append(arr.Elements, arr)

	if got := escapedString(arr); got != circularReference {
		t.Fatalf("expect %q, got %q", circularReference, got)
	}
}
