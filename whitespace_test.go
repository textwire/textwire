package textwire

import (
	"testing"
)

func TestStandaloneControlDirectiveWhitespace(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		data   map[string]any
		expect string
	}{
		{
			name:   "standalone loop lines are removed",
			inp:    "<ul>\n    @each(item in items)\n    <li> {{ item }} </li>\n    @end\n</ul>\n",
			data:   map[string]any{"items": []string{"apple", "pear"}},
			expect: "<ul>\n    <li> apple </li>\n    <li> pear </li>\n</ul>\n",
		},
		{
			name:   "inline loop keeps whitespace",
			inp:    "- @each(item in items){{ item }}, @end.\n",
			data:   map[string]any{"items": []string{"apple", "pear"}},
			expect: "- apple, pear, .\n",
		},
		{
			name:   "interpolation alone on a line is removed with its line",
			inp:    "    {{ item }}   \n",
			data:   map[string]any{"item": "apple"},
			expect: "apple",
		},
		{
			name:   "closer with text after it keeps whitespace",
			inp:    "@each(item in items)\n    {{ item }}\n    @end hello\n",
			data:   map[string]any{"items": []string{"apple", "pear"}},
			expect: "apple    pear     hello\n",
		},
		{
			name:   "YAML list keeps item indentation",
			inp:    "items:\n    @each(item in items)\n  - {{ item }}\n    @end\n",
			data:   map[string]any{"items": []string{"apple", "pear"}},
			expect: "items:\n  - apple\n  - pear\n",
		},
		{
			name:   "CRLF line endings are removed with the directive line",
			inp:    "@each(item in items)\r\n- {{ item }}\r\n@end\r\n",
			data:   map[string]any{"items": []string{"apple", "pear"}},
			expect: "- apple\r\n- pear\r\n",
		},
		{
			name:   "tabs around standalone directive are removed",
			inp:    "\t@if(true)\t\n\tkept\n\t@end\t\n",
			data:   nil,
			expect: "\tkept\n",
		},
		{
			name:   "final standalone line without newline",
			inp:    "x\n@if(true)\ny\n    @end",
			data:   nil,
			expect: "x\ny\n",
		},
		{
			name:   "two directives on one line are not standalone",
			inp:    "@if(true)@end\nz",
			data:   nil,
			expect: "\nz",
		},
		{
			name:   "statement and multiline embedded code alone on lines leave only text",
			inp:    "{{ names = [\"Anna\", \"Serhii\", \"Vladimir\"] }}\n{{ test\n      }}\nHere",
			data:   map[string]any{"test": ""},
			expect: "Here",
		},
		{
			name:   "raw embedded code alone on a line is removed with its line",
			inp:    "{!! test !!}\nHere",
			data:   map[string]any{"test": "<b>"},
			expect: "<b>Here",
		},
		{
			name:   "empty interpolation alone on a line is removed with its line",
			inp:    "a\n{{ empty }}\nb",
			data:   map[string]any{"empty": ""},
			expect: "a\nb",
		},
		{
			name:   "loop with zero items",
			inp:    "@each(x in items)\n- {{ x }}\n@end\nend",
			data:   map[string]any{"items": []string{}},
			expect: "end",
		},
		{
			name:   "loop with one item",
			inp:    "@each(x in items)\n- {{ x }}\n@end\nend",
			data:   map[string]any{"items": []string{"a"}},
			expect: "- a\nend",
		},
		{
			name:   "nested standalone blocks",
			inp:    "@if(true)\n@each(x in items)\n- {{ x }}\n@end\n@end\n",
			data:   map[string]any{"items": []string{"a", "b"}},
			expect: "- a\n- b\n",
		},
		{
			name:   "multiline @each standalone lines are removed",
			inp:    "@each(x in\n    items)\n{{ x }}\n@end\n",
			data:   map[string]any{"items": []string{"a", "b"}},
			expect: "ab",
		},
		{
			name:   "multiline @if standalone lines are removed",
			inp:    "@if(\n    true\n)\nsome text\n@end\n",
			data:   nil,
			expect: "some text\n",
		},
		{
			name:   "multiline @if with indentation is removed",
			inp:    "x\n    @if(\n      true\n    )\ny\n    @end\n",
			data:   nil,
			expect: "x\ny\n",
		},
		{
			name:   "multiline @if with text after closing parenthesis is kept",
			inp:    "@if(\n    true\n) hello\nsome text\n@end\n",
			data:   nil,
			expect: " hello\nsome text\n",
		},
		{
			name:   "multiline @if with CRLF line endings",
			inp:    "@if(\r\n    true\r\n)\r\nsome text\r\n@end\r\n",
			data:   nil,
			expect: "some text\r\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, failure := EvaluateString(tc.inp, tc.data)
			if failure != nil {
				t.Fatalf("evaluation failed: %s", failure)
			}

			if actual != tc.expect {
				t.Fatalf("expect %q, got %q", tc.expect, actual)
			}
		})
	}
}
