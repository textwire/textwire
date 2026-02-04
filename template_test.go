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
		fileName string
		data     map[string]any
	}{
		{"1.no-stmts", nil},
		{"2.with-inserts", nil},
		{"3.without-layout", map[string]any{
			"pageTitle": "Test Page",
			"NAME_1":    "Anna Korotchaeva",
			"name_2":    "Serhii Cho",
		}},
		{"4.loops", map[string]any{
			"names": []string{"Anna", "Serhii", "Vladimir"},
		}},
		{"5.with-component", map[string]any{
			"names": []string{"Anna", "Serhii", "Vladimir"},
		}},
		{"6.use-inside-if", nil},
		{"7.insert-without-use", nil},
		{"8.with-component", nil},
		{"9.with-inserts-and-html", nil},
		{"10.with-component-and-slots", nil},
		{"11.with-component-no-args", nil},
		{"13.insert-is-optional", nil},
		{"15.use-layout-with-comp-inside", nil},
	}

	tpl, err := NewTemplate(&config.Config{
		TemplateDir: "textwire/testdata/good/before",
	})

	if err != nil {
		t.Errorf("error creating template: %s", err)
		return
	}

	for _, tc := range cases {
		t.Run(tc.fileName, func(t *testing.T) {
			actual, evalErr := tpl.String(tc.fileName, tc.data)
			if evalErr != nil {
				t.Fatalf("error evaluating template: %s", evalErr)
				return
			}

			expect, err := readFile("textwire/testdata/good/expected/" + tc.fileName + ".html")
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
		TemplateDir: "textwire/testdata/good/before/",
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

	expect, err := readFile("textwire/testdata/good/expected/12.with-custom-function.html")
	if err != nil {
		t.Errorf("error reading expected file: %s", err)
		return
	}

	actual, evalErr := tpl.String("12.with-custom-function", nil)
	if evalErr != nil {
		t.Fatalf("error evaluating template: %s", evalErr)
	}

	if actual != expect {
		t.Errorf("wrong result. expect: '%s' got: '%s'", expect, actual)
	}
}
