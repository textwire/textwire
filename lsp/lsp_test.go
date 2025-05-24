package lsp

import "testing"

func TestIsInLoop(t *testing.T) {
	tests := []struct {
		doc     string
		linePos uint
		colPos  uint
		expect  bool
	}{
		{doc: `s`, linePos: 0, colPos: 0, expect: false},
		{
			doc:     `@each(x in users){{x}}@end`,
			linePos: 0,
			colPos:  19,
			expect:  true,
		},
	}

	for _, tc := range tests {
		actual := IsInLoop(tc.doc, "", tc.linePos, tc.colPos)

		if actual != tc.expect {
			t.Errorf("Expect IsInLoop to return %v, but got %v in %q",
				tc.expect, actual, tc.doc)
		}
	}
}
