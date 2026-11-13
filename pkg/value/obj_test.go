package value

import "testing"

func TestObjJSON(t *testing.T) {
	t.Run("escapes keys", func(t *testing.T) {
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
	})
}
