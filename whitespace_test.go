package textwire

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"codeberg.org/textwire/textwire/v5/config"
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
			name:   "two directives on one line are removed with the line (1)",
			inp:    "@if(true)@end\nz",
			data:   nil,
			expect: "z",
		},
		{
			name:   "two directives on one line are removed with the line (2)",
			inp:    "z\n@if(true)@end",
			data:   nil,
			expect: "z\n",
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
			name:   "multiline raw embedded code alone on lines leaves only its output",
			inp:    "{!!\ntest\n!!}\nHere",
			data:   map[string]any{"test": "<b>"},
			expect: "<b>Here",
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
			name:   "multiline embedded code with text after the closer is kept",
			inp:    "    {{ x\n}} hello\n",
			data:   map[string]any{"x": "v"},
			expect: "v hello\n",
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

func TestStandaloneTemplateDirectiveWhitespace(t *testing.T) {
	cases := []struct {
		name   string
		files  map[string]string
		file   string
		expect string
	}{
		{
			name: "standalone @use and @insert are removed with their lines",
			files: map[string]string{
				"index.tw":  "@use('layout')\n@insert('title', 'T')\nX\n",
				"layout.tw": "<head>\n    @reserve('title')\n</head>\n",
			},
			file:   "index",
			expect: "<head>\nT</head>\n",
		},
		{
			name: "standalone @reserve and @slot are removed with their lines",
			files: map[string]string{
				"index.tw":  "@use('layout')@insert('title', 'T')\n",
				"layout.tw": "<head>\n    @component('comp')\n        @reserve('title')\n    @end\n</head>\n",
				"comp.tw":   "<title>@slot</title>",
			},
			file:   "index",
			expect: "<head>\n<title>T</title></head>\n",
		},
		{
			name: "standalone @slot is removed with its line",
			files: map[string]string{
				"index.tw": "@component('comp')\n\nX\n@end\n",
				"comp.tw":  "<div>\n    @slot\n</div>\n",
			},
			file:   "index",
			expect: "<div>\n\nX\n</div>\n",
		},
		{
			name: "standalone @pass is removed with its line",
			files: map[string]string{
				"index.tw": "@component('comp')\n@pass('title')\n    T\n@end\nX\n@end\n",
				"comp.tw":  "<title>@slot('title')</title>\n<main>@slot</main>\n",
			},
			file:   "index",
			expect: "<title>    T\n</title>\n<main>X\n</main>\n",
		},
		{
			name: "standalone @passif is removed with its line",
			files: map[string]string{
				"index.tw": "@component('comp')\n@passif(true, 'title')\n    T\n@end\nX\n@end\n",
				"comp.tw":  "<title>@slot('title')</title>\n<main>@slot</main>\n",
			},
			file:   "index",
			expect: "<title>    T\n</title>\n<main>X\n</main>\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			writeTemplateFiles(t, dir, tc.files)

			tpl, failure := NewTemplate(&config.Config{TemplateDir: dir})
			if failure != nil {
				t.Fatalf("template creation failed: %s", failure)
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
	actual, failure := EvaluateString("    @dump(1)   \nX\n", nil)
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

func writeTemplateFiles(t *testing.T, dir string, files map[string]string) {
	t.Helper()

	for name, body := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o600); err != nil {
			t.Fatalf("cannot write template %s: %s", name, err)
		}
	}
}
