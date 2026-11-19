package apitest

import (
	"strings"
	"testing"

	"github.com/textwire/textwire/v6"
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
			name:   "embedded code alone on a line keeps whitespace",
			inp:    "    {{ item }}   \n",
			data:   map[string]any{"item": "apple"},
			expect: "    apple   \n",
		},
		{
			name:   "closer with text after it keeps whitespace",
			inp:    "@each(item in items)\n    {{ item }}\n    @end hello\n",
			data:   map[string]any{"items": []string{"apple", "pear"}},
			expect: "    apple\n        pear\n     hello\n",
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
			name:   "single-line @if with CRLF line endings is removed",
			inp:    "@if(true)\r\nHello\r\n@end\r\n",
			data:   nil,
			expect: "Hello\r\n",
		},
		{
			name:   "standalone directive with spaces indentation and trailing spaces is removed",
			inp:    "x\n    @if(true)   \ny\n    @end   \n",
			data:   nil,
			expect: "x\ny\n",
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
			name:   "two directives on one line keep whitespace (1)",
			inp:    "@if(true)@end\nz",
			data:   nil,
			expect: "\nz",
		},
		{
			name:   "two directives on one line keep whitespace (2)",
			inp:    "z\n@if(true)@end",
			data:   nil,
			expect: "z\n",
		},
		{
			name:   "two directives with @else keep whitespace",
			inp:    "@if(false)@else\n| Fallback\n@end\n",
			data:   nil,
			expect: "\n| Fallback\n",
		},
		{
			name:   "text before a directive keeps its line",
			inp:    "| Prefix @if(true)\nText\n@end\n",
			data:   nil,
			expect: "| Prefix \nText\n",
		},
		{
			name:   "text before a multiline directive keeps its line",
			inp:    "text @if(a &&\nb)\ncontent\n@end",
			data:   map[string]any{"a": true, "b": true},
			expect: "text content\n",
		},
		{
			name:   "spaces inside directive arguments do not matter",
			inp:    "@if( x == 1 )\nYes\n@end\n",
			data:   map[string]any{"x": 1},
			expect: "Yes\n",
		},
		{
			name:   "standalone directive at the start and end of a template",
			inp:    "@if(true)\nHello\n@end",
			data:   nil,
			expect: "Hello\n",
		},
		{
			name:   "embedded statement and multiline embedded code keep their lines",
			inp:    "{{ names = [\"Anna\", \"Serhii\", \"Vladimir\"] }}\n{{ test\n      }}\nHere",
			data:   map[string]any{"test": ""},
			expect: "\n\nHere",
		},
		{
			name:   "raw embedded code alone on a line keeps its line",
			inp:    "{!! test !!}\nHere",
			data:   map[string]any{"test": "<b>"},
			expect: "<b>\nHere",
		},
		{
			name:   "empty embedded code alone on a line keeps its line",
			inp:    "a\n{{ empty }}\nb",
			data:   map[string]any{"empty": ""},
			expect: "a\n\nb",
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
			name:   "multiline @each with embedded body keeps its output lines",
			inp:    "@each(x in\n    items)\n{{ x }}\n@end\n",
			data:   map[string]any{"items": []string{"a", "b"}},
			expect: "a\nb\n",
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
		{
			name:   "standalone @else is removed with its line",
			inp:    "@if(false)\na\n@else\nb\n@end\n",
			data:   nil,
			expect: "b\n",
		},
		{
			name:   "standalone @elseif is removed with its line",
			inp:    "@if(false)\na\n@elseif(true)\nb\n@else\nc\n@end\n",
			data:   nil,
			expect: "b\n",
		},
		{
			name:   "@else with text after it keeps whitespace",
			inp:    "@if(false)\na\n@else x\nb\n@end\n",
			data:   nil,
			expect: " x\nb\n",
		},
		{
			name:   "standalone @for is removed with its line",
			inp:    "@for(i = 0; i < 2; i++)\n- {{ i }}\n@end\n",
			data:   nil,
			expect: "- 0\n- 1\n",
		},
		{
			name:   "standalone @break is removed with its line",
			inp:    "@each(x in items)\n@if(x == \"b\")\n@break\n@end\n- {{ x }}\n@end\n",
			data:   map[string]any{"items": []string{"a", "b", "c"}},
			expect: "- a\n",
		},
		{
			name:   "standalone @continue is removed with its line",
			inp:    "@each(x in items)\n@if(x == \"b\")\n@continue\n@end\n- {{ x }}\n@end\n",
			data:   map[string]any{"items": []string{"a", "b", "c"}},
			expect: "- a\n- c\n",
		},
		{
			name:   "standalone @breakif is removed with its line",
			inp:    "@each(x in items)\n@breakif(x == \"b\")\n- {{ x }}\n@end\n",
			data:   map[string]any{"items": []string{"a", "b", "c"}},
			expect: "- a\n",
		},
		{
			name:   "standalone @continueif is removed with its line",
			inp:    "@each(x in items)\n@continueif(x == \"b\")\n- {{ x }}\n@end\n",
			data:   map[string]any{"items": []string{"a", "b", "c"}},
			expect: "- a\n- c\n",
		},
		{
			name:   "multiline raw embedded code keeps its line",
			inp:    "{!!\ntest\n!!}\nHere",
			data:   map[string]any{"test": "<b>"},
			expect: "<b>\nHere",
		},
		{
			name:   "blank lines around a standalone directive are kept",
			inp:    "a\n\n@if(true)\n\nb\n\n@end\n\nc\n",
			data:   nil,
			expect: "a\n\n\nb\n\n\nc\n",
		},
		{
			name:   "escaped directive alone on a line is not a directive",
			inp:    "\\@if\nx\n",
			data:   nil,
			expect: "@if\nx\n",
		},
		{
			name:   "multiline embedded code with text after the closer keeps whitespace",
			inp:    "    {{ x\n}} hello\n",
			data:   map[string]any{"x": "v"},
			expect: "    v hello\n",
		},
		{
			name:   "mixed directive and embedded keeps whitespace",
			inp:    " @if(true){{ \"hi\" }}@end\n",
			data:   nil,
			expect: " hi\n",
		},
		{
			name:   "mixed directive and embedded with spaces keeps whitespace",
			inp:    " @if(true) {{ \"hi\" }} @end\n",
			data:   nil,
			expect: "  hi \n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, failure := textwire.EvaluateString(tc.inp, tc.data)
			if failure != nil {
				t.Fatalf("evaluation failed: %s", failure)
			}

			if actual != tc.expect {
				t.Fatalf("expect %q, got %q", tc.expect, actual)
			}
		})
	}
}

func TestStandaloneCommentWhitespace(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		{
			name:   "standalone comment line is removed",
			inp:    "a\n{{-- comment --}}\nb\n",
			expect: "a\nb\n",
		},
		{
			name:   "spaces around standalone comment are removed",
			inp:    "a\n    {{-- comment --}}   \nb\n",
			expect: "a\nb\n",
		},
		{
			name:   "tabs and CRLF around standalone comment are removed",
			inp:    "a\r\n\t{{-- comment --}}\t\r\nb\r\n",
			expect: "a\r\nb\r\n",
		},
		{
			name:   "standalone comment at the start of a template",
			inp:    "{{-- comment --}}\nb",
			expect: "b",
		},
		{
			name:   "final standalone comment line without newline",
			inp:    "a\n{{-- comment --}}",
			expect: "a\n",
		},
		{
			name:   "blank lines around a standalone comment are kept",
			inp:    "a\n\n{{-- comment --}}\n\nb\n",
			expect: "a\n\n\nb\n",
		},
		{
			name:   "text before a comment keeps whitespace",
			inp:    "a {{-- comment --}}\nb\n",
			expect: "a \nb\n",
		},
		{
			name:   "text after a comment keeps whitespace",
			inp:    "a\n{{-- comment --}} b\n",
			expect: "a\n b\n",
		},
		{
			name:   "two comments on one line keep whitespace",
			inp:    "a\n{{-- first --}}{{-- second --}}\nb\n",
			expect: "a\n\nb\n",
		},
		{
			name:   "directive and comment on one line keep whitespace",
			inp:    "@if(true){{-- comment --}}\nb\n@end\n",
			expect: "\nb\n",
		},
		{
			name:   "embedded code and comment on one line keep whitespace",
			inp:    "{{ \"x\" }}{{-- comment --}}\nb\n",
			expect: "x\nb\n",
		},
		{
			name:   "standalone multiline comment lines are removed",
			inp:    "a\n  {{-- first\n  second\n  --}}\nb\n",
			expect: "a\nb\n",
		},
		{
			name:   "text before a multiline comment keeps its line",
			inp:    "a {{-- first\nsecond --}}\nb\n",
			expect: "a b\n",
		},
		{
			name:   "multiline comment with text after the closer keeps whitespace",
			inp:    "a\n{{-- first\nsecond --}} b\n",
			expect: "a\n b\n",
		},
		{
			name:   "standalone nested comment line is removed",
			inp:    "a\n{{-- outer {{-- inner --}} --}}\nb\n",
			expect: "a\nb\n",
		},
		{
			name:   "escaped comment alone on a line is not a comment",
			inp:    "a\n\\{{-- comment --}}\nb\n",
			expect: "a\n{{-- comment --}}\nb\n",
		},
		{
			name:   "standalone comment inside a loop body is removed",
			inp:    "@each(x in [1, 2])\n    {{-- comment --}}\n- {{ x }}\n@end\n",
			expect: "- 1\n- 2\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, failure := textwire.EvaluateString(tc.inp, nil)
			if failure != nil {
				t.Fatalf("evaluation failed: %s", failure)
			}

			if actual != tc.expect {
				t.Fatalf("expect %q, got %q", tc.expect, actual)
			}
		})
	}
}

func TestStandaloneTemplateWhitespace(t *testing.T) {
	cases := []struct {
		name   string
		files  []tplFile
		file   string
		expect string
	}{
		{
			name: "standalone @use and @insert are removed with their lines",
			files: []tplFile{
				{relPath: "index.tw", content: "@use('layout')\n@insert('title', 'T')\nX\n"},
				{relPath: "layout.tw", content: "<head>\n    @reserve('title')\n</head>\n"},
			},
			file:   "index",
			expect: "<head>\nT</head>\n",
		},
		{
			name: "standalone @reserve and @slot are removed with their lines",
			files: []tplFile{
				{relPath: "index.tw", content: "@use('layout')@insert('title', 'T')\n"},
				{
					relPath: "layout.tw",
					content: "<head>\n    @component('comp')\n        @reserve('title')\n    @end\n</head>\n",
				},
				{relPath: "comp.tw", content: "<title>@slot</title>"},
			},
			file:   "index",
			expect: "<head>\n<title>T</title></head>\n",
		},
		{
			name: "standalone @slot is removed with its line",
			files: []tplFile{
				{relPath: "index.tw", content: "@component('comp')\n\nX\n@end\n"},
				{relPath: "comp.tw", content: "<div>\n    @slot\n</div>\n"},
			},
			file:   "index",
			expect: "<div>\n\nX\n</div>\n",
		},
		{
			name: "standalone @pass is removed with its line",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: "@component('comp')\n@pass('title')\n    T\n@end\nX\n@end\n",
				},
				{
					relPath: "comp.tw",
					content: "<title>@slot('title')</title>\n<main>@slot</main>\n",
				},
			},
			file:   "index",
			expect: "<title>    T\n</title>\n<main>X\n</main>\n",
		},
		{
			name: "standalone @passif is removed with its line",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: "@component('comp')\n@passif(true, 'title')\n    T\n@end\nX\n@end\n",
				},
				{
					relPath: "comp.tw",
					content: "<title>@slot('title')</title>\n<main>@slot</main>\n",
				},
			},
			file:   "index",
			expect: "<title>    T\n</title>\n<main>X\n</main>\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tpl, err := setupTemplate(t, tc.files, nil)
			if err != nil {
				t.Fatal(err)
			}

			actual, failure := tpl.String(tc.file, nil)
			if failure != nil {
				t.Fatalf("evaluation failed: %s", failure)
			}

			if actual != tc.expect {
				t.Fatalf("expect %q, got %q", tc.expect, actual)
			}
		})
	}
}

func TestStandaloneDumpDirectiveWhitespace(t *testing.T) {
	actual, failure := textwire.EvaluateString("    @dump(1)   \nX\n", nil)
	if failure != nil {
		t.Fatalf("evaluation failed: %s", failure)
	}

	if !strings.HasPrefix(actual, "\n<div") {
		t.Fatalf("expect output to start with the dump markup, got %q", actual)
	}

	if !strings.HasSuffix(actual, "X\n") {
		t.Fatalf("expect output to end with %q, got %q", "X\n", actual)
	}
}
