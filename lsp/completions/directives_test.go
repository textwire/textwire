package completions

import "testing"

func TestGetDirectives(t *testing.T) {
	directives, err := GetDirectives("en")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(directives) == 0 {
		t.Fatal("expected non-empty directives")
	}
}
