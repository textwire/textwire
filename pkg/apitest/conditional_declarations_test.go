package apitest

import "testing"

func TestInsertInsideCondition(t *testing.T) {
	files := []tplFile{
		{
			relPath: "layout.tw",
			content: `[@reserve("t", "guest")]`,
		},
		{
			relPath: "index.tw",
			content: `@use("layout")@if(admin)@insert("t", "admin")@end`,
		},
	}

	cases := []struct {
		name   string
		admin  bool
		expect string
	}{
		{
			name:   "false condition keeps fallback",
			admin:  false,
			expect: "[guest]",
		},
		{
			name:   "true condition applies insert",
			admin:  true,
			expect: "[admin]",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tpl, tplFail := setupTemplate(t, files, nil)
			if tplFail != nil {
				t.Fatalf("error creating template: %q", tplFail)
			}

			actual, failure := tpl.String("index", map[string]any{"admin": tc.admin})
			assertOutput(t, actual, failure, tc.expect)
		})
	}
}

func TestUseInsideCondition(t *testing.T) {
	files := []tplFile{
		{
			relPath: "layout.tw",
			content: `[@reserve("t", "guest")]`,
		},
		{
			relPath: "index.tw",
			content: `@if(useLayout)@use("layout")@end plain`,
		},
	}

	cases := []struct {
		name      string
		useLayout bool
		expect    string
	}{
		{
			name:      "false condition skips layout",
			useLayout: false,
			expect:    " plain",
		},
		{
			name:      "true condition applies layout",
			useLayout: true,
			expect:    "[guest]",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tpl, tplFail := setupTemplate(t, files, nil)
			if tplFail != nil {
				t.Fatalf("error creating template: %q", tplFail)
			}

			actual, failure := tpl.String("index", map[string]any{"useLayout": tc.useLayout})
			assertOutput(t, actual, failure, tc.expect)
		})
	}
}
