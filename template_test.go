package textwire

import (
	"testing"

	"github.com/textwire/textwire/v3/config"
	"github.com/textwire/textwire/v3/fail"
)

func TestErrorHandlingEvaluatingTemplate(t *testing.T) {
	path, err := getFullPath("")
	path += "/textwire/testdata/bad/"

	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	cases := []struct {
		dir  string
		err  *fail.Error
		data map[string]any
	}{
		{
			dir: "use-inside-tpl",
			err: fail.New(
				1,
				path+"use-inside-tpl/index.tw",
				"evaluator",
				fail.ErrUseStmtNotAllowed,
			),
			data: nil,
		},
		{
			dir: "unknown-named-slot",
			err: fail.New(
				2,
				path+"unknown-named-slot/index.tw",
				"parser",
				fail.ErrSlotNotDefined,
				"unknown",
				"user",
			),
			data: nil,
		},
		{
			dir: "unknown-default-slot",
			err: fail.New(
				2,
				path+"unknown-default-slot/index.tw",
				"parser",
				fail.ErrDefaultSlotNotDefined,
				"book",
			),
			data: nil,
		},
		{
			dir: "duplicate-slot",
			err: fail.New(
				2,
				path+"duplicate-slot/index.tw",
				"parser",
				fail.ErrDuplicateSlotUsage,
				"content",
				2,
				"user",
			),
			data: nil,
		},
		{
			dir: "duplicate-default-slot",
			err: fail.New(
				2,
				path+"duplicate-default-slot/index.tw",
				"parser",
				fail.ErrDuplicateDefaultSlotUsage,
				2,
				"user",
			),
			data: nil,
		},
		{
			dir: "unknown-comp",
			err: fail.New(
				9,
				path+"unknown-comp/index.tw",
				"template",
				fail.ErrUndefinedComponent,
				"unknown-name",
			),
			data: nil,
		},
		{
			dir: "undefined-insert",
			err: fail.New(
				5,
				path+"undefined-insert/index.tw",
				"parser",
				fail.ErrUndefinedInsert,
				"some-name",
			),
			data: nil,
		},
		{
			dir: "duplicate-inserts",
			err: fail.New(
				4,
				path+"duplicate-inserts/index.tw",
				"parser",
				fail.ErrDuplicateInserts,
				"title",
			),
			data: nil,
		},
		{
			dir: "undefined-var-in-comp",
			err: fail.New(
				1,
				path+"undefined-var-in-comp/hero.tw",
				"parser",
				fail.ErrIdentifierIsUndefined,
				"undefinedVar",
			),
			data: nil,
		},
		{
			dir: "undefined-var-in-use",
			err: fail.New(
				8,
				path+"undefined-var-in-use/base.tw",
				"parser",
				fail.ErrIdentifierIsUndefined,
				"undefinedVar",
			),
			data: nil,
		},
		{
			dir: "undefined-use",
			err: fail.New(
				1,
				path+"undefined-use/index.tw",
				"parser",
				fail.ErrUseStmtMissingLayout,
				"undefined-layout",
			),
			data: nil,
		},
		{
			dir: "undefined-var-in-nested-comp",
			err: fail.New(
				1,
				path+"undefined-var-in-nested-comp/second.tw",
				"parser",
				fail.ErrIdentifierIsUndefined,
				"name",
			),
			data: map[string]any{"name": "Amy"},
		},
		{
			dir: "var-in-layout",
			err: fail.New(
				1,
				path+"var-in-layout/layout.tw",
				"parser",
				fail.ErrIdentifierIsUndefined,
				"fullName",
			),
			data: map[string]any{"fullName": "Amy Adams"},
		},
		{
			dir:  "duplicate-use",
			err:  fail.New(2, path+"duplicate-use/index.tw", "parser", fail.ErrOnlyOneUseDir),
			data: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.dir, func(t *testing.T) {
			tpl, tplErr := NewTemplate(&config.Config{
				TemplateDir: "textwire/testdata/bad/" + tc.dir,
			})

			if tplErr != nil {
				if tplErr.Error() != tc.err.String() {
					t.Fatalf("wrong error message. expect:\n\"%s\"\ngot:\n\"%s\"", tc.err, tplErr)
				}
				return
			}

			_, err := tpl.String("index", tc.data)
			if err == nil {
				t.Fatalf("expected error but got none")
				return
			}

			if err.String() != tc.err.String() {
				t.Fatalf("wrong error message. expect:\n\"%s\"\ngot:\n\"%s\"", tc.err, err)
			}
		})
	}
}

func TestNewTemplate(t *testing.T) {
	cases := []struct {
		dirName  string
		viewName string
		data     map[string]any
	}{
		{"no-stmts", "index", nil},
		{"with-inserts", "index", nil},
		{
			"without-use",
			"index",
			map[string]any{
				"pageTitle": "Test Page",
				"NAME_1":    "Anna Korotchaeva",
				"name_2":    "Serhii Cho",
			},
		},
		{"loops", "index", map[string]any{"names": []string{"Anna", "Serhii", "Vladimir"}}},
		{
			"with-each-and-comp",
			"views/index",
			map[string]any{"names": []string{"Anna", "Serhii", "Vladimir"}},
		},
		{"use-inside-if", "index", nil},
		{"insert-without-use", "index", nil},
		{"with-comp", "index", nil},
		{"with-inserts-and-html", "index", nil},
		{
			"with-comp-and-slots",
			"index",
			map[string]any{"head": "Header", "foot": "Footer"},
		},
		{"with-comp-no-args", "index", nil},
		{"insert-is-optional", "index", nil},
		{"use-with-comp-inside", "index", nil},
		{"comp-in-other-comp", "home", nil},
	}

	for _, tc := range cases {
		t.Run(tc.dirName, func(t *testing.T) {
			tpl, err := NewTemplate(&config.Config{
				TemplateDir: "textwire/testdata/good/before/" + tc.dirName,
			})

			if err != nil {
				t.Errorf("error creating template: %s", err)
				return
			}

			actual, evalErr := tpl.String(tc.viewName, tc.data)
			if evalErr != nil {
				t.Fatalf("error evaluating template: %s", evalErr)
				return
			}

			expect, err := readFile("textwire/testdata/good/expected/" + tc.dirName + ".html")
			if err != nil {
				t.Fatalf("error reading expected file: %s", err)
				return
			}

			if actual != expect {
				t.Fatalf("wrong result. expect:\n\"%s\"\ngot:\n\"%s\"", expect, actual)
			}
		})
	}
}

func TestRegisteringCustomFunction(t *testing.T) {
	tpl, err := NewTemplate(&config.Config{
		TemplateDir: "textwire/testdata/good/before/with-customs",
		GlobalData:  map[string]any{"env": "dev", "name": "Serhii", "age": 36},
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	err = RegisterStrFunc("_secondLetterUpper", func(s string, args ...any) any {
		if len(s) < 2 {
			return s
		}

		return string(s[0]) + string(s[1]-32) + s[2:]
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	expect, err := readFile("textwire/testdata/good/expected/with-customs.html")
	if err != nil {
		t.Errorf("error reading expected file: %s", err)
		return
	}

	actual, evalErr := tpl.String("index", nil)
	if evalErr != nil {
		t.Fatalf("error evaluating template: %s", evalErr)
	}

	if actual != expect {
		t.Errorf("wrong result. expect: '%s' got: '%s'", expect, actual)
	}
}

func TestTwoTemplates(t *testing.T) {
	tpl, err := NewTemplate(&config.Config{
		TemplateDir: "textwire/testdata/good/before/two-templates",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	expectHome, err := readFile("textwire/testdata/good/expected/two-templates-home.html")
	if err != nil {
		t.Errorf("error reading expected file: %s", err)
		return
	}

	actualHome, evalHomeErr := tpl.String("home", map[string]any{"titleHome": "home"})
	if evalHomeErr != nil {
		t.Fatalf("error evaluating home.tw template: %s", evalHomeErr)
	}

	if actualHome != expectHome {
		t.Errorf("wrong result for home.tw. expect: '%s' got: '%s'", expectHome, actualHome)
	}

	expectAbout, err := readFile("textwire/testdata/good/expected/two-templates-about.html")
	if err != nil {
		t.Errorf("error reading expected file: %s", err)
		return
	}

	actualAbout, evalAboutErr := tpl.String("about", map[string]any{"titleAbout": "about"})
	if evalAboutErr != nil {
		t.Fatalf("error evaluating home.tw template: %s", evalAboutErr)
	}

	if actualAbout != expectAbout {
		t.Errorf("wrong result for about.tw. expect: '%s' got: '%s'", expectAbout, actualAbout)
	}
}
