package lsp

import "testing"

func TestIsInLoop(t *testing.T) {
	cases := []struct {
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
		{
			doc: `
            @each(name in names)
                {{ loop }}
            @end`,
			linePos: 2,
			colPos:  23,
			expect:  true,
		},
		{
			doc: `
            @each(name in names)
                {{ loop }}
            @end`,
			linePos: 2,
			colPos:  23,
			expect:  true,
		},
		{
			doc: `
            @each(name in names)
                {{ loop }}
            @end`,
			linePos: 1,
			colPos:  30,
			expect:  false,
		},
		{
			doc: `
            @each(name in names)
                {{ loop }}
            @end`,
			linePos: 3,
			colPos:  10,
			expect:  true,
		},
		{
			doc: `
            @each(name in names)
                @each(name in names)
                    {{ loop }}
                @end
            @end`,
			linePos: 3,
			colPos:  27,
			expect:  true,
		},
		{
			doc: `
            @each(name in names)
                {{ loop }}
            @end`,
			linePos: 3,
			colPos:  15,
			expect:  false,
		},
		{
			doc: `
            @each(name in names)
                @if(loop.
                    {{ loop.first }}
                @end
            @end`,
			linePos: 2,
			colPos:  25,
			expect:  true,
		},
		{
			doc:     `@each(name in names)x@end`,
			linePos: 0,
			colPos:  20,
			expect:  true,
		},
		{
			doc:     `@each(name in names)x@end`,
			linePos: 0,
			colPos:  19,
			expect:  false,
		},
		{
			doc:     `@each(name in names)x@end`,
			linePos: 0,
			colPos:  21,
			expect:  false,
		},
		{
			doc: `@use('~main')

		@insert('title', "About Us")

		@insert('content')
		    @component('~header', {
		        title: "About Us",
		        description: "We have %n potatoes",
		    })

		    @each(name in names)
		        {{ loop }}
		    @end

		    <h2>{{ "Our Team" }}</h2>
		@end`,
			linePos: 11,
			colPos:  15,
			expect:  true,
		},
	}

	for _, tc := range cases {
		actual, errors := IsInLoop(tc.doc, "", tc.linePos, tc.colPos)

		if len(errors) > 0 {
			for _, msg := range errors {
				t.Logf("parser error: %q", msg.Error())
			}
		}

		if actual != tc.expect {
			t.Errorf("Expect IsInLoop to return %v, but got %v in %q, line %d, col %d",
				tc.expect, actual, tc.doc, tc.linePos, tc.colPos)
		}
	}
}
