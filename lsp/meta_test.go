package lsp

import (
	"strings"
	"testing"

	"github.com/textwire/textwire/v2/token"
)

func TestGetTokenMeta(t *testing.T) {
	t.Run("Invlid locale", func(t *testing.T) {
		_, err := GetTokenMeta(token.IF, "invalid")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("Token missing docs", func(t *testing.T) {
		_, err := GetTokenMeta(token.EOF, "en")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("Valid @if token meta", func(t *testing.T) {
		meta, err := GetTokenMeta(token.IF, "en")
		if err != nil {
			t.Errorf("expected err to be nil, got error %v", err)
		}

		expect := "@if(condition)"
		if !strings.Contains(meta, expect) {
			t.Errorf("expected %q in meta, got %q", expect, meta)
		}

		expect = "(directive)"
		if !strings.Contains(meta, expect) {
			t.Errorf("expected %q in meta, got %q", expect, meta)
		}
	})
}
