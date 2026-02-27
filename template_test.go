package textwire

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/textwire/textwire/v3/config"
	"github.com/textwire/textwire/v3/pkg/fail"
	"github.com/textwire/textwire/v3/pkg/file"
	"github.com/textwire/textwire/v3/pkg/object"
)

func TestErrorHandlingEvaluatingTemplate(t *testing.T) {
	absPath, err := file.ToFullPath("")
	absPath += "/testdata/bad/"
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
			dir: "duplicate-reserves",
			err: fail.New(
				3,
				absPath+"duplicate-reserves/base.tw",
				"parser",
				fail.ErrDuplicateReserves,
				"title",
				absPath+"duplicate-reserves/base.tw",
			),
			data: nil,
		},
		{
			dir: "use-inside-tpl",
			err: fail.New(
				1,
				absPath+"use-inside-tpl/index.tw",
				"evaluator",
				fail.ErrUseStmtNotAllowed,
			),
			data: nil,
		},
		{
			dir: "unknown-named-slot",
			err: fail.New(
				2,
				absPath+"unknown-named-slot/index.tw",
				"parser",
				fail.ErrSlotNotDefined,
				"user",
				"unknown",
			),
			data: nil,
		},
		{
			dir: "unknown-default-slot",
			err: fail.New(
				2,
				absPath+"unknown-default-slot/index.tw",
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
				absPath+"duplicate-slot/index.tw",
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
				absPath+"duplicate-default-slot/index.tw",
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
				absPath+"unknown-comp/index.tw",
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
				absPath+"undefined-insert/index.tw",
				"parser",
				fail.ErrUnusedInsertDetected,
				"some-name",
				"some-name",
			),
			data: nil,
		},
		{
			dir: "duplicate-inserts",
			err: fail.New(
				4,
				absPath+"duplicate-inserts/index.tw",
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
				absPath+"undefined-var-in-comp/hero.tw",
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
				absPath+"undefined-var-in-use/base.tw",
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
				absPath+"undefined-use/index.tw",
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
				absPath+"undefined-var-in-nested-comp/second.tw",
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
				absPath+"var-in-layout/layout.tw",
				"evaluator",
				fail.ErrVariableIsUndefined,
				"fullName",
			),
			data: map[string]any{"fullName": "Amy Adams"},
		},
		{
			dir:  "duplicate-use",
			err:  fail.New(2, absPath+"duplicate-use/index.tw", "parser", fail.ErrOnlyOneUseDir),
			data: nil,
		},
		{
			dir: "inserts-without-use",
			err: fail.New(
				4,
				absPath+"inserts-without-use/index.tw",
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
				&config.Config{TemplateDir: "testdata/bad/" + tc.dir},
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
	cases := []struct {
		conf *config.Config
		view string
		data map[string]any
		dir  string
	}{
		{conf: &config.Config{}, view: "index", data: nil, dir: "slots-optional"},
		{conf: &config.Config{}, view: "index", data: nil, dir: "reserve-inside-slot"},
		{conf: &config.Config{}, view: "index", data: nil, dir: "no-stmts"},
		{conf: &config.Config{}, view: "index", data: nil, dir: "inserts"},
		{conf: &config.Config{}, view: "index", data: nil, dir: "use-inside-if"},
		{conf: &config.Config{}, view: "index", data: nil, dir: "with-comp"},
		{conf: &config.Config{}, view: "index", data: nil, dir: "with-inserts-and-html"},
		{conf: &config.Config{}, view: "index", data: nil, dir: "with-comp-no-args"},
		{conf: &config.Config{}, view: "index", data: nil, dir: "insert-is-optional"},
		{conf: &config.Config{}, view: "index", data: nil, dir: "use-with-comp-inside"},
		{conf: &config.Config{}, view: "home", data: nil, dir: "comp-in-other-comp"},
		{
			conf: &config.Config{},
			view: "index",
			data: map[string]any{
				"pageTitle": "Test Page",
				"NAME_1":    "Anna Korotchaeva",
				"name_2":    "Serhii Cho",
			},
			dir: "without-use",
		},
		{
			conf: &config.Config{},
			view: "index",
			data: map[string]any{"names": []string{"Anna", "Serhii", "Vladimir"}},
			dir:  "loops",
		},
		{
			conf: &config.Config{},
			view: "views/index",
			data: map[string]any{"names": []string{"Anna", "Serhii", "Vladimir"}},
			dir:  "with-each-and-comp",
		},
		{
			conf: &config.Config{},
			view: "index",
			data: map[string]any{"name": "Анна ♥️", "age": 20},
			dir:  "comp-and-slots",
		},
	}

	for _, tc := range cases {
		t.Run(tc.dir, func(t *testing.T) {
			tc.conf.TemplateDir = "testdata/good/before/" + tc.dir
			tpl, err := NewTemplate(tc.conf)
			if err != nil {
				t.Fatalf("Error creating template: %q", err)
			}

			actual, failure := tpl.String(tc.view, tc.data)
			if failure != nil {
				t.Fatalf("Error evaluating template: %q", failure)
			}

			expect, err := readFile("testdata/good/expected/" + tc.dir + ".html")
			if err != nil {
				t.Fatalf("Error reading file. Error: %s", err)
			}

			if actual != expect {
				t.Fatalf("Wrong result. Expect:\n'%s'\ngot:\n'%s'", expect, actual)
			}
		})
	}
}

func TestTemplateResponse(t *testing.T) {
	absPath, err := file.ToFullPath("")
	absPath += "/testdata/good/before/"
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	cases := []struct {
		name string
		conf *config.Config
		view string
		data map[string]any
		dir  string
		err  *fail.Error
	}{
		{
			name: "Should show custom error page",
			conf: &config.Config{
				ErrorPagePath: "custom-error-page",
				GlobalData:    map[string]any{"year": "2020"},
			},
			view: "home",
			dir:  "prod-error-page",
			data: map[string]any{"arr": "some string"},
			err: fail.New(
				2,
				absPath+"prod-error-page/home.tw",
				"parser",
				fail.ErrEachDirWithNonArrArg,
				object.STR_OBJ,
			),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.conf.TemplateDir += "testdata/good/before/" + tc.dir
			tpl, err := NewTemplate(tc.conf)
			if err != nil {
				t.Errorf("Error creating template: %q", err)
				return
			}

			expect, err := readFile("testdata/good/expected/" + tc.dir + ".html")
			if err != nil {
				t.Fatalf("Error reading file. Error: %s", err)
			}

			rr := httptest.NewRecorder()
			respErr := tpl.Response(rr, tc.view, tc.data)
			actual := rr.Body.String()

			if tc.err != nil {
				if respErr == nil {
					t.Fatalf("We expect error from response, got nil")
				}

				if tc.err.String() != respErr.Error() {
					t.Fatalf("Wrong error message! Expect:\n%q\ngot:\n%q", tc.err, respErr)
				}
			}

			if actual != expect {
				t.Fatalf("Wrong result. Expect:\n'%s'\ngot:\n'%s'", expect, actual)
			}

			// Make sure you don't see error in actual response without debug mode
			if !tc.conf.DebugMode && respErr != nil {
				if contains := strings.Contains(actual, respErr.Error()); contains {
					t.Fatalf(
						"Actual string should not contain error message. Actual:\n'%s'\nError msg:\n'%s'",
						actual,
						respErr,
					)
				}
			}

			// Make sure you see error in actual response with debug mode enabled
			if tc.conf.DebugMode && respErr != nil {
				if contains := strings.Contains(actual, respErr.Error()); !contains {
					t.Fatalf(
						"Actual string should contain error message. Actual:\n'%s'\nError msg:\n'%s'",
						actual,
						respErr,
					)
				}
			}
		})
	}
}

func TestRegisteringCustomFunction(t *testing.T) {
	tpl, fileErr := NewTemplate(&config.Config{
		TemplateDir: "testdata/good/before/with-customs",
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

	expect, fileErr := readFile("testdata/good/expected/with-customs.html")
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
		TemplateDir: "testdata/good/before/two-templates",
	})
	if tplErr != nil {
		t.Fatalf("Unexpected template error: %s", tplErr)
	}

	expectHome, homeFileErr := readFile("testdata/good/expected/two-templates-home.html")
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
		"testdata/good/expected/two-templates-about.html",
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
