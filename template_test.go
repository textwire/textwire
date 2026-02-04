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
		dirName string
		err     *fail.Error
		data    map[string]any
	}{
		{
			dirName: "use-inside-tpl",
			err: fail.New(
				1,
				path+"use-inside-tpl/index.tw",
				"evaluator",
				fail.ErrUseStmtNotAllowed,
			),
			data: nil,
		},
		{
			dirName: "unknown-slot",
			err: fail.New(
				2,
				path+"unknown-slot/index.tw",
				"parser",
				fail.ErrSlotNotDefined,
				"unknown",
				"user",
			),
			data: nil,
		},
		{
			dirName: "unknown-default-slot",
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
			dirName: "duplicate-slot",
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
			dirName: "duplicate-default-slot",
			err: fail.New(
				2,
				path+"duplicate-default-slot/index.tw",
				"parser",
				fail.ErrDuplicateSlotUsage,
				"",
				2,
				"user",
			),
			data: nil,
		},
		{
			dirName: "unknown-component",
			err: fail.New(
				9,
				path+"unknown-component/index.tw",
				"template",
				fail.ErrUndefinedComponent,
				"unknown-name",
			),
			data: nil,
		},
		{
			dirName: "undefined-insert",
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
			dirName: "duplicate-inserts",
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
			dirName: "component-error",
			err: fail.New(
				1,
				path+"component-error/hero.tw",
				"parser",
				fail.ErrIdentifierIsUndefined,
				"undefinedVar",
			),
			data: nil,
		},
		{
			dirName: "use-stmt-error",
			err: fail.New(
				8,
				path+"use-stmt-error/base.tw",
				"parser",
				fail.ErrIdentifierIsUndefined,
				"undefinedVar",
			),
			data: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.dirName, func(t *testing.T) {
			tpl, tplErr := NewTemplate(&config.Config{
				TemplateDir: "textwire/testdata/bad/" + tc.dirName,
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
		{"001.no-stmts", "index", nil},
		{"002.with-inserts", "index", nil},
		{"003.without-layout", "index", map[string]any{
			"pageTitle": "Test Page",
			"NAME_1":    "Anna Korotchaeva",
			"name_2":    "Serhii Cho",
		}},
		{"004.loops", "index", map[string]any{
			"names": []string{"Anna", "Serhii", "Vladimir"},
		}},
		{"005.with-component", "views/index", map[string]any{
			"names": []string{"Anna", "Serhii", "Vladimir"},
		}},
		{"006.use-inside-if", "index", nil},
		{"007.insert-without-use", "index", nil},
		{"008.with-component", "index", nil},
		{"009.with-inserts-and-html", "index", nil},
		{"010.with-component-and-slots", "index", nil},
		{"011.with-component-no-args", "index", nil},
		{"013.insert-is-optional", "index", nil},
		{"015.use-layout-with-comp-inside", "index", nil},
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
		TemplateDir: "textwire/testdata/good/before/012.with-custom-function",
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

	expect, err := readFile("textwire/testdata/good/expected/012.with-custom-function.html")
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
