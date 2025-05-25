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
		{doc: `@each(x in users){{x}}@end`, linePos: 0, colPos: 16, expect: false},
		{doc: `@each(x in users){{x}}@end`, linePos: 0, colPos: 17, expect: true},
		{doc: `@each(x in users){{x}}@end`, linePos: 0, colPos: 18, expect: true},
		{doc: `@each(x in users){{x}}@end`, linePos: 0, colPos: 19, expect: true},
		{doc: `@each(x in users){{x}}@end`, linePos: 0, colPos: 20, expect: true},
		{doc: `@each(x in users){{x}}@end`, linePos: 0, colPos: 21, expect: true},
		{doc: `@each(x in users){{x}}@end`, linePos: 0, colPos: 22, expect: false},
		{doc: `@for(i = 0; i < 10; i++){{x}}@end`, linePos: 0, colPos: 23, expect: false},
		{doc: `@for(i = 0; i < 10; i++){{x}}@end`, linePos: 0, colPos: 24, expect: true},
		{doc: `@for(i = 0; i < 10; i++){{x}}@end`, linePos: 0, colPos: 25, expect: true},
		{doc: `@for(i = 0; i < 10; i++){{x}}@end`, linePos: 0, colPos: 26, expect: true},
		{doc: `@for(i = 0; i < 10; i++){{x}}@end`, linePos: 0, colPos: 27, expect: true},
		{doc: `@for(i = 0; i < 10; i++){{x}}@end`, linePos: 0, colPos: 28, expect: true},
		{doc: `@for(i = 0; i < 10; i++){{x}}@end`, linePos: 0, colPos: 29, expect: false},
		{doc: `@for(;;)x@end`, linePos: 0, colPos: 9, expect: false},
		{doc: `@for(;;)x@end`, linePos: 0, colPos: 8, expect: true},
		{doc: `@for(;;)x@end`, linePos: 0, colPos: 7, expect: false},
	}

	for _, tc := range tests {
		actual := IsInLoop(tc.doc, "", tc.linePos, tc.colPos)

		if actual != tc.expect {
			t.Errorf("Expect IsInLoop to return %v, but got %v in %q, line %d, col %d",
				tc.expect, actual, tc.doc, tc.linePos, tc.colPos)
		}
	}
}
