package textwire

import (
	"testing"

	"github.com/textwire/textwire/v3/config"
	"github.com/textwire/textwire/v3/fail"
	"github.com/textwire/textwire/v3/file"
)

func TestErrorHandlingEvaluatingTemplate(t *testing.T) {
	path, err := file.ToFullPath("")
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
				fail.ErrDuplicateSlot,
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
				fail.ErrDuplicateDefaultSlot,
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
				fail.ErrAddMatchingReserve,
				"some-name",
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
				"evaluator",
				fail.ErrVariableIsUndefined,
				"undefinedVar",
			),
			data: nil,
		},
		{
			dir: "undefined-var-in-use",
			err: fail.New(
				8,
				path+"undefined-var-in-use/base.tw",
				"evaluator",
				fail.ErrVariableIsUndefined,
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
				"evaluator",
				fail.ErrVariableIsUndefined,
				"name",
			),
			data: map[string]any{"name": "Amy"},
		},
		{
			dir: "var-in-layout",
			err: fail.New(
				1,
				path+"var-in-layout/layout.tw",
				"evaluator",
				fail.ErrVariableIsUndefined,
				"fullName",
			),
			data: map[string]any{"fullName": "Amy Adams"},
		},
		{
			dir:  "duplicate-use",
			err:  fail.New(2, path+"duplicate-use/index.tw", "parser", fail.ErrOnlyOneUseDir),
			data: nil,
		},
		{
			dir: "inserts-without-use",
			err: fail.New(
				4,
				path+"inserts-without-use/index.tw",
				"evaluator",
				fail.ErrInsertRequiresUse,
				"title",
			),
			data: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.dir, func(t *testing.T) {
			tpl, tplErr := NewTemplate(
				&config.Config{TemplateDir: "textwire/testdata/bad/" + tc.dir},
			)
			if tplErr != nil {
				if tplErr.Error() != tc.err.String() {
					t.Fatalf("Wrong error message! Expect:\n%q\ngot:\n%q", tc.err, tplErr)
				}
				return
			}

			_, err := tpl.String("index", tc.data)
			if err == nil {
				t.Fatalf("Expected error but got none")
			}

			if err.String() != tc.err.String() {
				t.Fatalf("Wrong error message! Expect:\n%s\ngot:\n%q", tc.err, err)
			}

			if err.Origin() != tc.err.Origin() {
				t.Fatalf(
					"Wrong origin on error message, expect %s, got: %s in error message:\n%q",
					tc.err.Origin(),
					err.Origin(),
					err,
				)
			}
		})
	}
}

func TestNewTemplate(t *testing.T) {
	path := "textwire/testdata/good/before/"
	cases := []struct {
		config   *config.Config
		viewName string
		data     map[string]any
		dirName  string
	}{
		{
			config:   &config.Config{TemplateDir: path + "no-stmts"},
			viewName: "index",
			data:     nil,
			dirName:  "no-stmts",
		},
		{
			config:   &config.Config{TemplateDir: path + "with-inserts"},
			viewName: "index",
			data:     nil,
			dirName:  "with-inserts",
		},
		{
			config:   &config.Config{TemplateDir: path + "without-use"},
			viewName: "index",
			data: map[string]any{
				"pageTitle": "Test Page",
				"NAME_1":    "Anna Korotchaeva",
				"name_2":    "Serhii Cho",
			},
			dirName: "without-use",
		},
		{
			config:   &config.Config{TemplateDir: path + "loops"},
			viewName: "index",
			data:     map[string]any{"names": []string{"Anna", "Serhii", "Vladimir"}},
			dirName:  "loops",
		},
		{
			config:   &config.Config{TemplateDir: path + "with-each-and-comp"},
			viewName: "views/index",
			data:     map[string]any{"names": []string{"Anna", "Serhii", "Vladimir"}},
			dirName:  "with-each-and-comp",
		},
		{
			config:   &config.Config{TemplateDir: path + "use-inside-if"},
			viewName: "index",
			data:     nil,
			dirName:  "use-inside-if",
		},
		{
			config:   &config.Config{TemplateDir: path + "with-comp"},
			viewName: "index",
			data:     nil,
			dirName:  "with-comp",
		},
		{
			config:   &config.Config{TemplateDir: path + "with-inserts-and-html"},
			viewName: "index",
			data:     nil,
			dirName:  "with-inserts-and-html",
		},
		{
			config:   &config.Config{TemplateDir: path + "with-comp-and-slots"},
			viewName: "index",
			data:     map[string]any{"name": "Anna", "age": 20},
			dirName:  "with-comp-and-slots",
		},
		{
			config:   &config.Config{TemplateDir: path + "with-comp-no-args"},
			viewName: "index",
			data:     nil,
			dirName:  "with-comp-no-args",
		},
		{
			config:   &config.Config{TemplateDir: path + "insert-is-optional"},
			viewName: "index",
			data:     nil,
			dirName:  "insert-is-optional",
		},
		{
			config:   &config.Config{TemplateDir: path + "use-with-comp-inside"},
			viewName: "index",
			data:     nil,
			dirName:  "use-with-comp-inside",
		},
		{
			config:   &config.Config{TemplateDir: path + "comp-in-other-comp"},
			viewName: "home",
			data:     nil,
			dirName:  "comp-in-other-comp",
		},
	}

	for _, tc := range cases {
		t.Run(tc.dirName, func(t *testing.T) {
			tpl, err := NewTemplate(tc.config)
			if err != nil {
				t.Errorf("Error creating template: %q", err)
				return
			}

			actual, failure := tpl.String(tc.viewName, tc.data)
			if failure != nil {
				t.Fatalf("Error evaluating template: %q", failure)
				return
			}

			expect, err := readFile("textwire/testdata/good/expected/" + tc.dirName + ".html")
			if err != nil {
				t.Fatalf("Error reading file. Error: %s", err)
				return
			}

			if actual != expect {
				t.Fatalf("Wrong result. Expect:\n'%s'\ngot:\n'%s'", expect, actual)
			}
		})
	}
}

func TestRegisteringCustomFunction(t *testing.T) {
	tpl, fileErr := NewTemplate(&config.Config{
		TemplateDir: "textwire/testdata/good/before/with-customs",
		GlobalData:  map[string]any{"env": "dev", "name": "Serhii", "age": 36},
	})
	if fileErr != nil {
		t.Fatalf("Unexpected template error: %s", fileErr)
	}

	err := RegisterStrFunc("_secondLetterUpper", func(s string, args ...any) any {
		if len(s) < 2 {
			return s
		}
		return string(s[0]) + string(s[1]-32) + s[2:]
	})
	if err != nil {
		t.Fatalf("Unexpected error registering function: %s", fileErr)
	}

	expect, fileErr := readFile("textwire/testdata/good/expected/with-customs.html")
	if fileErr != nil {
		t.Errorf("Error reading file: %s", fileErr)
		return
	}

	actual, evalErr := tpl.String("index", nil)
	if evalErr != nil {
		t.Fatalf("Error evaluating template: %s", evalErr)
	}

	if actual != expect {
		t.Errorf("Wrong result. Expect:\n'%s'\ngot:\n'%s'", expect, actual)
	}
}

func TestTwoTemplates(t *testing.T) {
	tpl, tplErr := NewTemplate(&config.Config{
		TemplateDir: "textwire/testdata/good/before/two-templates",
	})
	if tplErr != nil {
		t.Fatalf("Unexpected template error: %s", tplErr)
	}

	expectHome, homeFileErr := readFile("textwire/testdata/good/expected/two-templates-home.html")
	if homeFileErr != nil {
		t.Errorf("Error reading file: %s", homeFileErr)
		return
	}

	actualHome, evalHomeErr := tpl.String("home", map[string]any{"titleHome": "home"})
	if evalHomeErr != nil {
		t.Fatalf("Error evaluating home.tw template: %s", evalHomeErr)
	}

	if actualHome != expectHome {
		t.Errorf("Wrong result for home.tw. Expect\n'%s'\ngot:\n'%s'", expectHome, actualHome)
	}

	expectAbout, aboutFileErr := readFile(
		"textwire/testdata/good/expected/two-templates-about.html",
	)
	if aboutFileErr != nil {
		t.Errorf("Error reading file: %s", aboutFileErr)
		return
	}

	actualAbout, evalAboutErr := tpl.String("about", map[string]any{"titleAbout": "about"})
	if evalAboutErr != nil {
		t.Fatalf("Error evaluating home.tw file template: %s", evalAboutErr)
	}

	if actualAbout != expectAbout {
		t.Errorf("Wrong result for about.tw. Expect:\n'%s'\ngot:\n'%s'", expectAbout, actualAbout)
	}
}
