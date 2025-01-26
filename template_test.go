package textwire

import (
	"testing"

	"github.com/textwire/textwire/v2/config"
	"github.com/textwire/textwire/v2/fail"
)

func TestErrorHandlingEvaluatingTemplate(t *testing.T) {
	path, err := getFullPath("", false)
	path += "/testdata/bad/"

	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	tests := []struct {
		dirName string
		err     *fail.Error
		data    map[string]interface{}
	}{
		{
			"use-inside-tpl",
			fail.New(1, path+"use-inside-tpl/index.tw", "evaluator", fail.ErrUseStmtNotAllowed),
			nil,
		},
		{
			"unknown-slot",
			fail.New(
				2,
				path+"unknown-slot/index.tw",
				"parser",
				fail.ErrSlotNotDefined, "unknown", "user",
			),
			nil,
		},
		{
			"unknown-default-slot",
			fail.New(
				2,
				path+"unknown-default-slot/index.tw",
				"parser",
				fail.ErrDefaultSlotNotDefined, "book",
			),
			nil,
		},
		{
			"duplicate-slot",
			fail.New(
				2,
				path+"duplicate-slot/index.tw",
				"parser",
				fail.ErrDuplicateSlotUsage, "content", 2, "user",
			),
			nil,
		},
		{
			"duplicate-default-slot",
			fail.New(
				2,
				path+"duplicate-default-slot/index.tw",
				"parser",
				fail.ErrDuplicateSlotUsage, "", 2, "user",
			),
			nil,
		},
		{
			"unknown-component",
			fail.New(
				9,
				path+"unknown-component/index.tw",
				"template",
				fail.ErrUndefinedComponent, "unknown-name",
			),
			nil,
		},
		{
			"undefined-insert",
			fail.New(5, path+"undefined-insert/index.tw", "parser", fail.ErrUndefinedInsert, "some-name"),
			nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.dirName, func(t *testing.T) {
			tpl, tplErr := NewTemplate(&config.Config{
				TemplateDir: "testdata/bad/" + tc.dirName,
				TemplateExt: ".tw",
			})

			if tplErr != nil {
				if tplErr.Error() != tc.err.String() {
					t.Errorf("wrong error message. EXPECTED:\n\"%s\"\nGOT:\n\"%s\"", tc.err, tplErr)
				}
				return
			}

			_, err := tpl.String("index", tc.data)

			if err == nil {
				t.Errorf("expected error but got none")
				return
			}

			if err.String() != tc.err.String() {
				t.Errorf("wrong error message. EXPECTED:\n\"%s\"\nGOT:\n\"%s\"", tc.err, err)
			}
		})
	}
}

func TestFiles(t *testing.T) {
	tests := []struct {
		fileName string
		data     map[string]interface{}
	}{
		{"1.no-stmts", nil},
		{"2.with-inserts", nil},
		{"3.without-layout", map[string]interface{}{
			"pageTitle": "Test Page",
			"NAME_1":    "Anna Korotchaeva",
			"name_2":    "Serhii Cho",
		}},
		{"4.loops", map[string]interface{}{
			"names": []string{"Anna", "Serhii", "Vladimir"},
		}},
		{"5.with-component", map[string]interface{}{
			"names": []string{"Anna", "Serhii", "Vladimir"},
		}},
		{"6.use-inside-if", nil},
		{"7.insert-without-use", nil},
		{"8.with-component", nil},
		{"9.with-inserts-and-html", nil},
		{"10.with-component-and-slots", nil},
		{"11.with-component-no-args", nil},
		{"13.insert-is-optional", nil},
	}

	tpl, err := NewTemplate(&config.Config{
		TemplateDir: "testdata/good/before",
		TemplateExt: ".tw",
	})

	if err != nil {
		t.Errorf("error creating template: %s", err)
		return
	}

	for _, tc := range tests {
		actual, evalErr := tpl.String(tc.fileName, tc.data)

		if evalErr != nil {
			t.Errorf("error evaluating template: %s", evalErr)
			return
		}

		expected, err := readFile("testdata/good/expected/" + tc.fileName + ".html")

		if err != nil {
			t.Errorf("error reading expected file: %s", err)
			return
		}

		if actual != expected {
			t.Errorf("wrong result. EXPECTED:\n\"%s\"\nGOT:\n\"%s\"",
				expected, actual)
		}
	}
}

func TestRegisteringCustomFunction(t *testing.T) {
	fileName, err := getFullPath("", false)
	fileName += "/testdata/good/12.with-custom-function"

	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	tpl, tplErr := NewTemplate(&config.Config{
		TemplateDir: "testdata/good/before/",
	})

	RegisterStrFunc("secondLetterUppercase", func(s string, args ...interface{}) string {
		if len(s) < 2 {
			return s
		}

		return string(s[0]) + string(s[1]-32) + s[2:]
	})

	if tplErr != nil {
		t.Fatalf("unexpected error: %s", tplErr)
	}

	expected, err := readFile("testdata/good/expected/12.with-custom-function.html")

	if err != nil {
		t.Errorf("error reading expected file: %s", err)
		return
	}

	actual, evalErr := tpl.String("12.with-custom-function", nil)

	if evalErr != nil {
		t.Fatalf("error evaluating template: %s", evalErr)
	}

	if actual != expected {
		t.Errorf("wrong result. EXPECTED: '%s' GOT: '%s'", expected, actual)
	}
}
