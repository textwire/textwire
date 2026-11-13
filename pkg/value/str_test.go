package value

import (
	"fmt"
	"testing"
)

func TestStrJSON(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		{name: "plain string", inp: "hello", expect: `"hello"`},
		{name: "quote", inp: `a"b`, expect: `"a\"b"`},
		{name: "backslash", inp: `a\b`, expect: `"a\\b"`},
		{name: "html characters are not escaped", inp: `<script>&`, expect: `"<script>&"`},
		{name: "control character", inp: "a\x01b", expect: `"a\u0001b"`},
		{name: "invalid utf8", inp: "a\xffb", expect: `"a\ufffdb"`},
		{name: "newline and tab", inp: "a\nb\tc", expect: `"a\nb\tc"`},
		{name: "newline and tab", inp: `<b>"`, expect: `"<b>\""`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := (&Str{Val: tc.inp}).JSON()
			if err != nil {
				t.Fatal(err)
			}

			if got != tc.expect {
				t.Fatalf("expected %q, got %q", tc.expect, got)
			}
		})
	}
}

func TestStrString(t *testing.T) {
	cases := []struct {
		name   string
		inp    *Str
		expect string
	}{
		{
			name:   "plain string",
			inp:    &Str{Val: "hello"},
			expect: `hello`,
		},
		{
			name:   "ignores raw mark",
			inp:    &Str{Val: `<b>"`, raw: true},
			expect: `<b>"`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.inp.String(); got != tc.expect {
				t.Fatalf("expected %q, got %q", tc.expect, got)
			}
		})
	}
}

func TestStrNative(t *testing.T) {
	cases := []struct {
		name   string
		inp    *Str
		expect string
	}{
		{
			name:   "plain string",
			inp:    &Str{Val: "hello"},
			expect: `hello`,
		},
		{
			name:   "ignores raw mark",
			inp:    &Str{Val: `<b>"`, raw: true},
			expect: `<b>"`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.inp.Native(); got != tc.expect {
				t.Fatalf("expected %q, got %q", tc.expect, got)
			}
		})
	}
}

func TestStrDump(t *testing.T) {
	cases := []struct {
		name   string
		inp    *Str
		expect string
	}{
		{
			name:   "plain string",
			inp:    &Str{Val: "hello"},
			expect: fmt.Sprintf(`<span style="%s">%q</span>`, DUMP_STR, "hello"),
		},
		{
			name:   "plain non-raw string escapes HTML",
			inp:    &Str{Val: "<br>"},
			expect: fmt.Sprintf(`<span style="%s">%q</span>`, DUMP_STR, "&lt;br&gt;"),
		},
		{
			name: "string marked as raw still escapes HTML with marker",
			inp:  &Str{Val: "<br>", raw: true},
			expect: fmt.Sprintf(
				`<span style="%s">raw </span><span style="%s">%q</span>`,
				DUMP_META,
				DUMP_STR,
				"&lt;br&gt;",
			),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.inp.Dump(0); got != tc.expect {
				t.Fatalf("expected %#q, got %#q", tc.expect, got)
			}
		})
	}
}
