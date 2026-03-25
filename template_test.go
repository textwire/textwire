package textwire

import (
	"fmt"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/textwire/textwire/v4/config"
	"github.com/textwire/textwire/v4/pkg/fail"
	"github.com/textwire/textwire/v4/pkg/file"
	"github.com/textwire/textwire/v4/pkg/position"
	"github.com/textwire/textwire/v4/pkg/value"
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
			dir: "undefined-default-slotif",
			err: fail.New(
				&position.Pos{StartCol: 21, EndCol: 41},
				absPath+"undefined-default-slotif/index.tw",
				fail.OriginLink,
				fail.ErrDefaultSlotNotDefined,
				"user",
			),
			data: nil,
		},
		{
			dir: "undefined-named-slotif",
			err: fail.New(
				&position.Pos{StartCol: 21, EndCol: 49},
				absPath+"undefined-named-slotif/index.tw",
				fail.OriginLink,
				fail.ErrSlotNotDefined,
				"user",
				"name",
			),
			data: nil,
		},
		{
			dir: "use-inside-tpl",
			err: fail.New(
				&position.Pos{EndCol: 13},
				absPath+"use-inside-tpl/layout.tw",
				fail.OriginEval,
				fail.ErrUseDirIsNotAllowed,
			),
			data: nil,
		},
		{
			dir: "unknown-named-slot",
			err: fail.New(
				&position.Pos{StartLine: 1, StartCol: 4, EndLine: 3, EndCol: 7},
				absPath+"unknown-named-slot/index.tw",
				fail.OriginLink,
				fail.ErrSlotNotDefined,
				"user",
				"unknown",
			),
			data: nil,
		},
		{
			dir: "unknown-default-slot",
			err: fail.New(
				&position.Pos{StartLine: 1, StartCol: 4, EndLine: 3, EndCol: 7},
				absPath+"unknown-default-slot/index.tw",
				fail.OriginLink,
				fail.ErrDefaultSlotNotDefined,
				"book",
			),
			data: nil,
		},
		{
			dir: "duplicate-slot",
			err: fail.New(
				&position.Pos{StartLine: 2, StartCol: 4, EndLine: 2, EndCol: 39},
				absPath+"duplicate-slot/index.tw",
				fail.OriginLink,
				fail.ErrDuplicateSlot,
				"content",
				2,
				"user",
			),
			data: nil,
		},
		{
			dir: "unknown-comp",
			err: fail.New(
				&position.Pos{StartLine: 8, StartCol: 4, EndLine: 8, EndCol: 29},
				absPath+"unknown-comp/index.tw",
				fail.OriginLink,
				fail.ErrUndefinedComponent,
				"unknown-name",
			),
			data: nil,
		},
		{
			dir: "undefined-insert",
			err: fail.New(
				&position.Pos{StartLine: 4, EndLine: 6, EndCol: 3},
				absPath+"undefined-insert/index.tw",
				fail.OriginLink,
				fail.ErrUnusedInsertDetected,
				"some-name",
				"some-name",
			),
			data: nil,
		},
		{
			dir: "undefined-var-in-comp",
			err: fail.New(
				&position.Pos{StartCol: 3, EndCol: 14},
				absPath+"undefined-var-in-comp/hero.tw",
				fail.OriginEval,
				fail.ErrVariableIsUndefined,
				"undefinedVar",
			),
			data: nil,
		},
		{
			dir: "undefined-var-in-use",
			err: fail.New(
				&position.Pos{StartLine: 7, StartCol: 9, EndLine: 7, EndCol: 20},
				absPath+"undefined-var-in-use/base.tw",
				fail.OriginEval,
				fail.ErrVariableIsUndefined,
				"undefinedVar",
			),
			data: nil,
		},
		{
			dir: "undefined-use",
			err: fail.New(
				&position.Pos{StartCol: 5, EndCol: 22},
				absPath+"undefined-use/index.tw",
				fail.OriginLink,
				fail.ErrUseDirMissingLayout,
				"undefined-layout",
			),
			data: nil,
		},
		{
			dir: "undefined-var-in-nested-comp",
			err: fail.New(
				&position.Pos{StartCol: 9, EndCol: 12},
				absPath+"undefined-var-in-nested-comp/second.tw",
				fail.OriginEval,
				fail.ErrVariableIsUndefined,
				"name",
			),
			data: map[string]any{"name": "Amy"},
		},
		{
			dir: "var-in-layout",
			err: fail.New(
				&position.Pos{StartCol: 9, EndCol: 16},
				absPath+"var-in-layout/layout.tw",
				fail.OriginEval,
				fail.ErrVariableIsUndefined,
				"fullName",
			),
			data: map[string]any{"fullName": "Amy Adams"},
		},
		{
			dir: "inserts-without-use",
			err: fail.New(
				&position.Pos{StartLine: 3, EndLine: 3, EndCol: 31},
				absPath+"inserts-without-use/index.tw",
				fail.OriginEval,
				fail.ErrInsertRequiresUse,
				"title",
			),
			data: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.dir, func(t *testing.T) {
			tpl, tplFail := NewTemplate(
				&config.Config{TemplateDir: "testdata/bad/" + tc.dir},
			)

			if tplFail != nil {
				if err := compareFailures(tplFail, tc.err); err != nil {
					t.Fatal(err)
				}
				return
			}

			_, err := tpl.String("index", tc.data)
			if err == nil {
				t.Fatalf("Expected error but got none")
			}

			if err := compareFailures(err, tc.err); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func compareFailures(got, expect *fail.Error) error {
	if got.String() != expect.String() {
		return fmt.Errorf("wrong error message! Expect:\n%q\ngot:\n%q", expect, got)
	}

	if got.Origin() != expect.Origin() {
		return fmt.Errorf(
			"wrong origin on error message, expect %s, got: %s in error message:\n%q",
			expect.Origin(),
			got.Origin(),
			got,
		)
	}

	if !reflect.DeepEqual(got.Pos(), expect.Pos()) {
		return fmt.Errorf(
			"wrong position on error message, expect %v, got: %v in error message:\n%q",
			expect.Pos(),
			got.Pos(),
			got,
		)
	}

	return nil
}

func TestNewTemplate(t *testing.T) {
	cases := []struct {
		conf *config.Config
		view string
		data map[string]any
		dir  string
	}{
		{conf: &config.Config{}, view: "index", data: nil, dir: "slotif"},
		{conf: &config.Config{}, view: "index", data: nil, dir: "slots-optional"},
		{conf: &config.Config{}, view: "index", data: nil, dir: "reserve-inside-slot"},
		{conf: &config.Config{}, view: "~index", data: nil, dir: "no-stmts"},
		{conf: &config.Config{}, view: "index", data: nil, dir: "inserts"},
		{conf: &config.Config{}, view: "index", data: nil, dir: "use-inside-if"},
		{conf: &config.Config{}, view: "index", data: nil, dir: "comp"},
		{conf: &config.Config{}, view: "index", data: nil, dir: "inserts-and-html"},
		{conf: &config.Config{}, view: "index", data: nil, dir: "comp-no-args"},
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
			view: "~index",
			data: map[string]any{"names": []string{"Anna", "Serhii", "Vladimir"}},
			dir:  "each-and-comp",
		},
		{
			conf: &config.Config{},
			view: "index",
			data: map[string]any{"name": "Анна ♥️", "age": 20},
			dir:  "comp-and-slots",
		},
		{
			conf: &config.Config{},
			view: "index",
			data: map[string]any{
				"person": struct {
					First_Name string
					Age        uint8
				}{First_Name: "Anna", Age: 25},
			},
			dir: "json",
		},
	}

	for _, tc := range cases {
		t.Run(tc.dir, func(t *testing.T) {
			tc.conf.TemplateDir = "testdata/good/before/" + tc.dir
			tpl, tplFail := NewTemplate(tc.conf)
			if tplFail != nil {
				t.Fatalf("Error creating template: %q", tplFail)
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
				&position.Pos{StartLine: 1, EndLine: 1},
				absPath+"prod-error-page/home.tw",
				fail.OriginPars,
				fail.ErrEachDirWithNonArrArg,
				value.STR_VAL,
			),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.conf.TemplateDir += "testdata/good/before/" + tc.dir
			tpl, failure := NewTemplate(tc.conf)
			if failure != nil {
				t.Errorf("Error creating template: %q", failure)
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

				if tc.err.String() != respErr.String() {
					t.Fatalf("Wrong error message! Expect:\n%q\ngot:\n%q", tc.err, respErr)
				}
			}

			if actual != expect {
				t.Fatalf("Wrong result. Expect:\n'%s'\ngot:\n'%s'", expect, actual)
			}

			// Make sure you don't see error in actual response without debug mode
			if !tc.conf.DebugMode && respErr != nil {
				if contains := strings.Contains(actual, respErr.String()); contains {
					t.Fatalf(
						"Actual string should not contain error message. Actual:\n'%s'\nError msg:\n'%s'",
						actual,
						respErr,
					)
				}
			}

			// Make sure you see error in actual response with debug mode enabled
			if tc.conf.DebugMode && respErr != nil {
				if contains := strings.Contains(actual, respErr.String()); !contains {
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
	tpl, tplErr := NewTemplate(&config.Config{
		TemplateDir: "testdata/good/before/globals",
		GlobalData:  map[string]any{"env": "dev", "name": "Serhii", "age": 36},
	})

	if tplErr != nil {
		t.Fatalf("Unexpected template error: %s", tplErr)
	}

	err := RegisterStrFunc("_secondLetterUpper", func(s string, args ...any) any {
		if len(s) < 2 {
			return s
		}
		return string(s[0]) + string(s[1]-32) + s[2:]
	})
	if err != nil {
		t.Fatalf("Unexpected error registering function: %s", tplErr)
	}

	expect, fileErr := readFile("testdata/good/expected/globals.html")
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
