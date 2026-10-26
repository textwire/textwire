package value

import "testing"

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

func TestObjJSON_EscapesKeys(t *testing.T) {
	obj := NewObj(map[string]Literal{
		`a"b<c>`: &Str{Val: `<script>`},
	})

	got, err := obj.JSON()
	if err != nil {
		t.Fatal(err)
	}

	want := `{"a\"b<c>":"<script>"}`
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
