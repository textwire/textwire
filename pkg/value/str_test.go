package value

import (
	"fmt"
	"testing"
)

func TestStrJSON(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{name: "plain string", input: "hello", want: `"hello"`},
		{name: "quote", input: `a"b`, want: `"a\"b"`},
		{name: "backslash", input: `a\b`, want: `"a\\b"`},
		{name: "html characters are not escaped", input: `<script>&`, want: `"<script>&"`},
		{name: "control character", input: "a\x01b", want: `"a\u0001b"`},
		{name: "invalid utf8", input: "a\xffb", want: `"a\ufffdb"`},
		{name: "newline and tab", input: "a\nb\tc", want: `"a\nb\tc"`},
		{name: "newline and tab", input: `<b>"`, want: `"<b>\""`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := (&Str{Val: tc.input}).JSON()
			if err != nil {
				t.Fatal(err)
			}

			if got != tc.want {
				t.Fatalf("expected %q, got %q", tc.want, got)
			}
		})
	}
}

func TestStrString(t *testing.T) {
	cases := []struct {
		name  string
		input *Str
		want  string
	}{
		{
			name:  "plain string",
			input: &Str{Val: "hello"},
			want:  `hello`,
		},
		{
			name:  "ignores raw mark",
			input: &Str{Val: `<b>"`, raw: true},
			want:  `<b>"`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.input.String(); got != tc.want {
				t.Fatalf("expected %q, got %q", tc.want, got)
			}
		})
	}
}

func TestStrNative(t *testing.T) {
	cases := []struct {
		name  string
		input *Str
		want  string
	}{
		{
			name:  "plain string",
			input: &Str{Val: "hello"},
			want:  `hello`,
		},
		{
			name:  "ignores raw mark",
			input: &Str{Val: `<b>"`, raw: true},
			want:  `<b>"`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.input.Native(); got != tc.want {
				t.Fatalf("expected %q, got %q", tc.want, got)
			}
		})
	}
}

func TestStrDump(t *testing.T) {
	cases := []struct {
		name  string
		input *Str
		want  string
	}{
		{
			name:  "plain string",
			input: &Str{Val: "hello"},
			want:  fmt.Sprintf(`<span style="%s">%q</span>`, DUMP_STR, "hello"),
		},
		{
			name:  "plain non-raw string escapes HTML",
			input: &Str{Val: "<br>"},
			want:  fmt.Sprintf(`<span style="%s">%q</span>`, DUMP_STR, "&lt;br&gt;"),
		},
		{
			name:  "string marked as raw still escapes HTML with marker",
			input: &Str{Val: "<br>", raw: true},
			want:  fmt.Sprintf(`<span style="%s">raw </span><span style="%s">%q</span>`, DUMP_META, DUMP_STR, "&lt;br&gt;"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.input.Dump(0); got != tc.want {
				t.Fatalf("expected %#q, got %#q", tc.want, got)
			}
		})
	}
}
