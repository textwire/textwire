package value

import (
	"fmt"
	"testing"
)

func TestStrJSON(t *testing.T) {
	cases := []struct {
		name string
		inp  string
		want string
	}{
		{name: "plain string", inp: "hello", want: `"hello"`},
		{name: "quote", inp: `a"b`, want: `"a\"b"`},
		{name: "backslash", inp: `a\b`, want: `"a\\b"`},
		{name: "html characters are not escaped", inp: `<script>&`, want: `"<script>&"`},
		{name: "control character", inp: "a\x01b", want: `"a\u0001b"`},
		{name: "invalid utf8", inp: "a\xffb", want: `"a\ufffdb"`},
		{name: "newline and tab", inp: "a\nb\tc", want: `"a\nb\tc"`},
		{name: "newline and tab", inp: `<b>"`, want: `"<b>\""`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := (&Str{Val: tc.inp}).JSON()
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
		name string
		inp  *Str
		want string
	}{
		{
			name: "plain string",
			inp:  &Str{Val: "hello"},
			want: `hello`,
		},
		{
			name: "ignores raw mark",
			inp:  &Str{Val: `<b>"`, raw: true},
			want: `<b>"`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.inp.String(); got != tc.want {
				t.Fatalf("expected %q, got %q", tc.want, got)
			}
		})
	}
}

func TestStrNative(t *testing.T) {
	cases := []struct {
		name string
		inp  *Str
		want string
	}{
		{
			name: "plain string",
			inp:  &Str{Val: "hello"},
			want: `hello`,
		},
		{
			name: "ignores raw mark",
			inp:  &Str{Val: `<b>"`, raw: true},
			want: `<b>"`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.inp.Native(); got != tc.want {
				t.Fatalf("expected %q, got %q", tc.want, got)
			}
		})
	}
}

func TestStrDump(t *testing.T) {
	cases := []struct {
		name string
		inp  *Str
		want string
	}{
		{
			name: "plain string",
			inp:  &Str{Val: "hello"},
			want: fmt.Sprintf(`<span style="%s">%q</span>`, DUMP_STR, "hello"),
		},
		{
			name: "plain non-raw string escapes HTML",
			inp:  &Str{Val: "<br>"},
			want: fmt.Sprintf(`<span style="%s">%q</span>`, DUMP_STR, "&lt;br&gt;"),
		},
		{
			name: "string marked as raw still escapes HTML with marker",
			inp:  &Str{Val: "<br>", raw: true},
			want: fmt.Sprintf(
				`<span style="%s">raw </span><span style="%s">%q</span>`,
				DUMP_META,
				DUMP_STR,
				"&lt;br&gt;",
			),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.inp.Dump(0); got != tc.want {
				t.Fatalf("expected %#q, got %#q", tc.want, got)
			}
		})
	}
}
