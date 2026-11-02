package fail

import (
	"errors"
	"testing"

	"codeberg.org/textwire/textwire/v5/pkg/position"
)

func TestFromError(t *testing.T) {
	t.Run("keeps percent in message", func(t *testing.T) {
		err := FromError(errors.New("open /tmp/100%tw"), &position.Pos{}, "", OriginTpl)

		want := "open /tmp/100%tw"
		if err.Message() != want {
			t.Fatalf("want %q, got %q", want, err.Message())
		}
	})

	t.Run("keeps percent verbs literal", func(t *testing.T) {
		msgs := []string{"100%%", "%d items", "%s", "%v"}

		for _, msg := range msgs {
			err := FromError(errors.New(msg), &position.Pos{}, "", OriginTpl)

			if err.Message() != msg {
				t.Fatalf("want %q, got %q", msg, err.Message())
			}
		}
	})

	t.Run("passes filepath and origin", func(t *testing.T) {
		err := FromError(errors.New("boom"), &position.Pos{}, "/tpl/index.tw", OriginEval)

		if err.Filepath() != "/tpl/index.tw" {
			t.Fatalf("want filepath %q, got %q", "/tpl/index.tw", err.Filepath())
		}
		if err.Origin() != OriginEval {
			t.Fatalf("want origin %q, got %q", OriginEval, err.Origin())
		}
	})

	t.Run("nil pos gets empty pos", func(t *testing.T) {
		err := FromError(errors.New("boom"), nil, "", OriginTpl)

		if err.Pos() == nil {
			t.Fatal("want non-nil pos, got nil")
		}
	})

	t.Run("panics on nil error", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("want panic on nil error, got none")
			}
		}()

		FromError(nil, &position.Pos{}, "", OriginTpl)
	})
}
