package completions

import (
	"testing"

	"github.com/textwire/textwire/v3/pkg/token"
)

func TestGetDirectives(t *testing.T) {
	directives, err := GetDirectives("en")
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if len(directives) == 0 {
		t.Fatal("expect non-empty directives")
	}

	hasInsertDir := false
	for _, dir := range directives {
		if dir.Label == "@insert" {
			hasInsertDir = true
		}
	}

	if !hasInsertDir {
		t.Fatal("GetDirectives() should return slice that contain @insert directive")
	}

	directivesCount := len(token.GetDirectives())
	if directivesCount != len(directives) {
		t.Fatalf("GetDirectives() should return %d directives, got %d",
			directivesCount, len(directives))
	}
}
