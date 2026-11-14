package apitest

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"codeberg.org/textwire/textwire/v6"

	"codeberg.org/textwire/textwire/v6/config"
	"codeberg.org/textwire/textwire/v6/pkg/fail"
	"codeberg.org/textwire/textwire/v6/pkg/file"
	"codeberg.org/textwire/textwire/v6/pkg/position"
	"codeberg.org/textwire/textwire/v6/pkg/value"
)

func TestErrorHandlingEvaluatingTemplate(t *testing.T) {
	cases := []struct {
		name    string
		files   []tplFile
		data    map[string]any
		errFile string
		errPos  *position.Pos
		origin  fail.ErrOrigin
		errMsg  string
		errArgs []any
	}{
		{
			name: "undefined named @passif",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `<p>@component('user')@passif(true, 'name')TEST@end@end</p>`,
				},
				{relPath: "user.tw", content: `<section>@slot</section>`},
			},
			data:    nil,
			errFile: "index.tw",
			errPos:  &position.Pos{StartCol: 21, EndCol: 49},
			origin:  fail.OriginLink,
			errMsg:  fail.ErrSlotNotDefined,
			errArgs: []any{"user", "name"},
		},
		{
			name: "@use inside layout",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: "@use('layout')\n@insert('title', 'test')",
				},
				{
					relPath: "layout.tw",
					content: "@use('layout')\n@reserve('title')",
				},
			},
			data:    nil,
			errFile: "layout.tw",
			errPos:  &position.Pos{EndCol: 13},
			origin:  fail.OriginEval,
			errMsg:  fail.ErrUseDirIsNotAllowed,
			errArgs: nil,
		},
		{
			name: "unknown named @pass",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@component('user')@pass('unknown')<p>test</p>@end@end`,
				},
				{
					relPath: "user.tw",
					content: `@slot("content")`,
				},
			},
			data:    nil,
			errFile: "index.tw",
			errPos:  &position.Pos{StartCol: 24, EndCol: 32},
			origin:  fail.OriginLink,
			errMsg:  fail.ErrSlotNotDefined,
			errArgs: []any{"user", "unknown"},
		},
		{
			name: "unknown-default-pass",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@component("book", { title: "Harry Potter" })
    <p>Some content is here</p>
@end`,
				},
				{
					relPath: "book.tw",
					content: `<div class="book">
    <h2>{{ title }}</h2>
    @slot("content")
</div>
`,
				},
			},
			data:    nil,
			errFile: "index.tw",
			errPos:  &position.Pos{StartLine: 0, StartCol: 45, EndLine: 1, EndCol: 31},
			origin:  fail.OriginLink,
			errMsg:  fail.ErrDefaultSlotNotDefined,
			errArgs: []any{"book", "book"},
		},
		{
			name: "duplicate-pass",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@component("user", { name: "Anna" })
    @pass("content")<p>Content 1</p>@end
    @pass("content")<p>Content 2</p>@end
@end
`,
				},
				{
					relPath: "user.tw",
					content: `<div class="user">
    <h2>{{ name }}</h2>
    @slot("content")
</div>
`,
				},
			},
			data:    nil,
			errFile: "index.tw",
			errPos:  &position.Pos{StartLine: 2, StartCol: 4, EndLine: 2, EndCol: 39},
			origin:  fail.OriginPars,
			errMsg:  fail.ErrDuplicatePass,
			errArgs: []any{"content", 2, "user"},
		},
		{
			name: "unknown-comp",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Document</title>
</head>
<body>
    @component('unknown-name')@end
</body>
</html>
`,
				},
			},
			data:    nil,
			errFile: "index.tw",
			errPos:  &position.Pos{StartLine: 8, StartCol: 4, EndLine: 8, EndCol: 33},
			origin:  fail.OriginLink,
			errMsg:  fail.ErrUndefinedComponent,
			errArgs: []any{"unknown-name"},
		},
		{
			name: "undefined-insert",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@use('layout')

@insert('title', 'Some title')

@insert('some-name')
<h1>Some title is here</h1>
@end
`,
				},
				{
					relPath: "layout.tw",
					content: `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>@reserve('title')</title>
</head>
<body>
    @reserve('main')
</body>
</html>
`,
				},
			},
			data:    nil,
			errFile: "index.tw",
			errPos:  &position.Pos{StartLine: 4, EndLine: 6, EndCol: 3},
			origin:  fail.OriginLink,
			errMsg:  fail.ErrUnusedInsertDetected,
			errArgs: []any{"some-name", "some-name"},
		},
		{
			name: "undefined-var-in-comp",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@component("hero")
    <h1>Title</h1>
@end
`,
				},
				{relPath: "hero.tw", content: "{{ undefinedVar }}\n@slot"},
			},
			data:    nil,
			errFile: "hero.tw",
			errPos:  &position.Pos{StartCol: 3, EndCol: 14},
			origin:  fail.OriginEval,
			errMsg:  fail.ErrVariableIsUndefined,
			errArgs: []any{"undefinedVar"},
		},
		{
			name: "undefined var inside layout",
			files: []tplFile{
				{relPath: "index.tw", content: `@use("base")`},
				{
					relPath: "base.tw",
					content: "<!DOCTYPE html>\n<body>{{ undefinedVar }}</body>",
				},
			},
			data:    nil,
			errFile: "base.tw",
			errPos:  &position.Pos{StartLine: 1, StartCol: 9, EndLine: 1, EndCol: 20},
			origin:  fail.OriginEval,
			errMsg:  fail.ErrVariableIsUndefined,
			errArgs: []any{"undefinedVar"},
		},
		{
			name: "undefined-use",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@use('undefined-layout')`,
				},
			},
			data:    nil,
			errFile: "index.tw",
			errPos:  &position.Pos{StartCol: 5, EndCol: 22},
			origin:  fail.OriginLink,
			errMsg:  fail.ErrUseDirMissingLayout,
			errArgs: []any{"undefined-layout"},
		},
		{
			name: "undefined-var-in-nested-comp",
			files: []tplFile{
				{relPath: "index.tw", content: `@component('first', { name })@end`},
				{relPath: "first.tw", content: `@component('second')@end`},
				{relPath: "second.tw", content: `<span>{{ name }}</span>`},
			},
			data:    map[string]any{"name": "Amy"},
			errFile: "second.tw",
			errPos:  &position.Pos{StartCol: 9, EndCol: 12},
			origin:  fail.OriginEval,
			errMsg:  fail.ErrVariableIsUndefined,
			errArgs: []any{"name"},
		},
		{
			name: "var-in-layout",
			files: []tplFile{
				{relPath: "index.tw", content: `@use('layout')`},
				{relPath: "layout.tw", content: `<span>{{ fullName }}</span>`},
			},
			data:    map[string]any{"fullName": "Amy Adams"},
			errFile: "layout.tw",
			errPos:  &position.Pos{StartCol: 9, EndCol: 16},
			origin:  fail.OriginEval,
			errMsg:  fail.ErrVariableIsUndefined,
			errArgs: []any{"fullName"},
		},
		{
			name: "comp-calls-itself",
			files: []tplFile{
				{relPath: "index.tw", content: "@component('~tree')@end\n"},
				{
					relPath: "components/tree.tw",
					content: "<li>@component('~tree')@end</li>\n",
				},
			},
			data:    nil,
			errFile: "components/tree.tw",
			errPos:  &position.Pos{StartCol: 4, EndCol: 26},
			origin:  fail.OriginEval,
			errMsg:  fail.ErrMaxComponentDepth,
			errArgs: []any{"components/tree", config.DefaultMaxComponentDepth},
		},
		{
			name: "comp-calls-each-other",
			files: []tplFile{
				{relPath: "index.tw", content: "@component('~first')@end\n"},
				{
					relPath: "components/first.tw",
					content: "<b>@component('~second')@end</b>\n",
				},
				{
					relPath: "components/second.tw",
					content: "<i>@component('~first')@end</i>\n",
				},
			},
			data:    nil,
			errFile: "components/second.tw",
			errPos:  &position.Pos{StartCol: 3, EndCol: 26},
			origin:  fail.OriginEval,
			errMsg:  fail.ErrMaxComponentDepth,
			errArgs: []any{"components/first", config.DefaultMaxComponentDepth},
		},
		{
			name: "inserts-without-use",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `<h1>Amy Lou Adams</h1>@insert("title", "Causes error")`,
				},
			},
			data:    nil,
			errFile: "index.tw",
			errPos:  &position.Pos{StartCol: 22, EndCol: 53},
			origin:  fail.OriginEval,
			errMsg:  fail.ErrInsertRequiresUse,
			errArgs: []any{"title"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conf := &config.Config{}
			tpl, tplFail := setupTemplate(t, tc.files, conf)

			// setupTemplate points conf.TemplateDir to the temporary directory
			expect := fail.New(
				tc.errPos,
				filepath.Join(conf.TemplateDir, tc.errFile),
				tc.origin,
				tc.errMsg,
				tc.errArgs...,
			)

			if tplFail != nil {
				if err := compareFailures(tplFail, expect); err != nil {
					t.Fatal(err)
				}
				return
			}

			_, failure := tpl.String("index", tc.data)
			if failure == nil {
				t.Fatalf("expected error, got none")
			}

			if err := compareFailures(failure, expect); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestNewTemplate(t *testing.T) {
	assertNoFailure(t, textwire.RegisterStrFunc("_wrapStr", func(s string, args ...any) any {
		return "[" + s + "]"
	}))

	assertNoFailure(t, textwire.RegisterIntFunc("_twiceInt", func(num int, args ...any) any {
		return num * 2
	}))

	err := textwire.RegisterStrFunc("_secondLetterUpper", func(s string, args ...any) any {
		if len(s) < 2 {
			return s
		}
		return string(s[0]) + string(s[1]-32) + s[2:]
	})

	assertNoFailure(t, err)

	cases := []struct {
		conf *config.Config
		view string
		data map[string]any
		dir  string
	}{
		{conf: &config.Config{}, view: "index", data: nil, dir: "passif"},
		{conf: &config.Config{}, view: "index", data: nil, dir: "slots-optional"},
		{conf: &config.Config{}, view: "index", data: nil, dir: "reserve-inside-slot"},
		{conf: &config.Config{}, view: "~index", data: nil, dir: "no-stmts"},
		{conf: &config.Config{}, view: "index", data: nil, dir: "inserts"},
		{conf: &config.Config{}, view: "index", data: nil, dir: "use-inside-if"},
		{conf: &config.Config{}, view: "index", data: nil, dir: "comp"},
		{conf: &config.Config{}, view: "index", data: nil, dir: "inserts-and-html"},
		{conf: &config.Config{}, view: "index", data: nil, dir: "comp-no-args"},
		{conf: &config.Config{}, view: "index", data: nil, dir: "insert-is-optional"},
		{conf: &config.Config{}, view: "index", data: nil, dir: "nested-reserve-slot"},
		{conf: &config.Config{}, view: "index", data: nil, dir: "use-with-comp-inside"},
		{conf: &config.Config{}, view: "home", data: nil, dir: "comp-in-other-comp"},
		{conf: &config.Config{}, view: "index", data: nil, dir: "custom-funcs"},
		{
			conf: &config.Config{
				GlobalData: map[string]any{"env": "dev", "name": "Serhii", "age": 36},
			},
			view: "index",
			data: nil,
			dir:  "globals",
		},
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
			data: map[string]any{"names": []string{"Anna", "Serhii", "Vladimir"}},
			dir:  "comp-in-loop",
		},
		{
			conf: &config.Config{},
			view: "index",
			data: map[string]any{"name": "Анна ♥️", "age": 20},
			dir:  "comp-and-passes",
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
			tpl, tplFail := textwire.NewTemplate(tc.conf)
			if tplFail != nil {
				t.Fatalf("error creating template: %q", tplFail)
			}

			actual, failure := tpl.String(tc.view, tc.data)
			if failure != nil {
				t.Fatalf("error evaluating template: %q", failure)
			}

			expect, err := readFile("testdata/good/expected/" + tc.dir + ".html")
			if err != nil {
				t.Fatalf("error reading file: %s", err)
			}

			if actual != expect {
				t.Fatalf("wrong result. Expect:\n'%s'\ngot:\n'%s'", expect, actual)
			}
		})
	}
}

func TestNewTemplate_WithAbsoluteTemplateDir(t *testing.T) {
	// setupTemplate uses t.TempDir(), which is an absolute path
	tpl, err := setupTemplate(t, []tplFile{{relPath: "index.tw", content: "hello"}}, nil)
	if err != nil {
		t.Fatal(err)
	}

	actual, failure := tpl.String("index", nil)
	if failure != nil {
		t.Fatal(failure)
	}

	if actual != "hello" {
		t.Fatalf("expected 'hello', got %q", actual)
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
			tpl, failure := textwire.NewTemplate(tc.conf)
			if failure != nil {
				t.Errorf("error creating template: %q", failure)
				return
			}

			expect, err := readFile("testdata/good/expected/" + tc.dir + ".html")
			if err != nil {
				t.Fatalf("error reading file: %s", err)
			}

			rr := httptest.NewRecorder()
			respErr := tpl.Response(rr, tc.view, tc.data)
			actual := rr.Body.String()

			if tc.err != nil {
				if respErr == nil {
					t.Fatalf("expected error from response, got nil")
				}

				if tc.err.String() != respErr.String() {
					t.Fatalf("wrong error message. Expect:\n%q\ngot:\n%q", tc.err, respErr)
				}
			}

			if actual != expect {
				t.Fatalf("wrong result. Expect:\n'%s'\ngot:\n'%s'", expect, actual)
			}

			// Make sure you don't see error in actual response without debug mode
			if !tc.conf.DebugMode && respErr != nil {
				if contains := strings.Contains(actual, respErr.String()); contains {
					t.Fatalf(
						"actual string should not contain error message. Actual:\n'%s'\nError msg:\n'%s'",
						actual,
						respErr,
					)
				}
			}

			// Make sure you see error in actual response with debug mode enabled
			if tc.conf.DebugMode && respErr != nil {
				if contains := strings.Contains(actual, respErr.String()); !contains {
					t.Fatalf(
						"actual string should contain error message. Actual:\n'%s'\nError msg:\n'%s'",
						actual,
						respErr,
					)
				}
			}
		})
	}
}

func TestTemplateResponseStatusCodes(t *testing.T) {
	tpl, tplErr := textwire.NewTemplate(&config.Config{
		TemplateDir: "testdata/bad/undefined-var-in-comp",
	})
	if tplErr != nil {
		t.Fatalf("unexpected template error: %s", tplErr)
	}

	t.Run("render error returns 500", func(t *testing.T) {
		rr := httptest.NewRecorder()
		tpl.Response(rr, "index", nil)
		if rr.Code != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
		}
	})

	t.Run("missing template returns 404", func(t *testing.T) {
		rr := httptest.NewRecorder()
		tpl.Response(rr, "nonexistent", nil)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d", http.StatusNotFound, rr.Code)
		}
	})
}

func TestLocateFiles_IgnoresNonTemplateFiles(t *testing.T) {
	files := []tplFile{
		{relPath: "index.tw", content: "hello"},
		{relPath: "notes.twig", content: "junk @end"},
		{relPath: "some.tw.twig", content: "@end"},
		{relPath: "card.tw.html", content: "@end"},
	}

	tpl, err := setupTemplate(t, files, nil)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := tpl.String("index", nil)
	if err != nil {
		t.Fatalf("unexpected eval error: %s", err)
	}

	if actual != "hello" {
		t.Fatalf("expected 'hello', got %q", actual)
	}
}

func TestTwoTemplates(t *testing.T) {
	tpl, tplErr := textwire.NewTemplate(&config.Config{
		TemplateDir: "testdata/good/before/two-templates",
	})
	if tplErr != nil {
		t.Fatalf("unexpected template error: %s", tplErr)
	}

	expectHome, homeFileErr := readFile("testdata/good/expected/two-templates-home.html")
	if homeFileErr != nil {
		t.Errorf("error reading file: %s", homeFileErr)
		return
	}

	actualHome, evalHomeErr := tpl.String("home", map[string]any{"titleHome": "home"})
	if evalHomeErr != nil {
		t.Fatalf("error evaluating home.tw template: %s", evalHomeErr)
	}

	if actualHome != expectHome {
		t.Errorf("wrong result for home.tw. Expect:\n'%s'\ngot:\n'%s'", expectHome, actualHome)
	}

	expectAbout, aboutFileErr := readFile(
		"testdata/good/expected/two-templates-about.html",
	)
	if aboutFileErr != nil {
		t.Errorf("error reading file: %s", aboutFileErr)
		return
	}

	actualAbout, evalAboutErr := tpl.String("about", map[string]any{"titleAbout": "about"})
	if evalAboutErr != nil {
		t.Fatalf("error evaluating about.tw template: %s", evalAboutErr)
	}

	if actualAbout != expectAbout {
		t.Errorf("wrong result for about.tw. Expect:\n'%s'\ngot:\n'%s'", expectAbout, actualAbout)
	}
}

func TestInsertAndReserveEscaping(t *testing.T) {
	cases := []struct {
		name   string
		files  []tplFile
		expect string
	}{
		// Insert argument
		{
			name: "argument is escaped",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@use('layout')@insert('title', '<span>Title</span>')`,
				},
				{
					relPath: "layout.tw",
					content: `<div>@reserve('title')</div>`,
				},
			},
			expect: "<div>&lt;span&gt;Title&lt;/span&gt;</div>",
		},
		{
			name: "raw argument is not escaped",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@use('layout')@insert('title', raw('<span>Title</span>'))`,
				},
				{
					relPath: "layout.tw",
					content: `<div>@reserve('title')</div>`,
				},
			},
			expect: "<div><span>Title</span></div>",
		},
		{
			name: "string literal argument is escaped",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@use('layout')@insert('title', 'A & B')`,
				},
				{
					relPath: "layout.tw",
					content: `<div>@reserve('title')</div>`,
				},
			},
			expect: "<div>A &amp; B</div>",
		},
		{
			name: "raw escaped text argument stays as is",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@use('layout')@insert('title', raw('&amp;'))`,
				},
				{
					relPath: "layout.tw",
					content: `<div>@reserve('title')</div>`,
				},
			},
			expect: "<div>&amp;</div>",
		},
		{
			name: "integer argument is printed as is",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@use('layout')@insert('title', 5)`,
				},
				{
					relPath: "layout.tw",
					content: `<div>@reserve('title')</div>`,
				},
			},
			expect: "<div>5</div>",
		},
		{
			name: "array argument is escaped",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@use('layout')@insert('title', ['<a>', '<b>'])`,
				},
				{
					relPath: "layout.tw",
					content: `<div>@reserve('title')</div>`,
				},
			},
			expect: "<div>&lt;a&gt;, &lt;b&gt;</div>",
		},
		{
			name: "raw array argument is not escaped",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@use('layout')@insert('title', raw(['<a>', '<b>']))`,
				},
				{
					relPath: "layout.tw",
					content: `<div>@reserve('title')</div>`,
				},
			},
			expect: "<div><a>, <b></div>",
		},
		{
			name: "array argument escapes only plain strings",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@use('layout')@insert('title', [raw('<span>Title</span>'), '<span>Title</span>'])`,
				},
				{
					relPath: "layout.tw",
					content: `<div>@reserve('title')</div>`,
				},
			},
			expect: "<div><span>Title</span>, &lt;span&gt;Title&lt;/span&gt;</div>",
		},
		{
			name: "object argument escapes only plain strings",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@use('layout')@insert('title', {a: raw('<span>Title</span>'), b: '<span>Title</span>'})`,
				},
				{
					relPath: "layout.tw",
					content: `<div>@reserve('title')</div>`,
				},
			},
			expect: `<div>{a: "<span>Title</span>", b: &#34;&lt;span&gt;Title&lt;/span&gt;&#34;}</div>`,
		},
		{
			name: "argument inside component slot is escaped",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@use('layout')@insert('title', '<span>Title</span>')`,
				},
				{
					relPath: "layout.tw",
					content: `<div>@component('~box')@reserve('title')@end</div>`,
				},
				{
					relPath: "components/box.tw",
					content: "<i>@slot</i>",
				},
			},
			expect: "<div><i>&lt;span&gt;Title&lt;/span&gt;</i></div>",
		},
		{
			name: "raw argument inside component slot is not escaped",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@use('layout')@insert('title', raw('<span>Title</span>'))`,
				},
				{
					relPath: "layout.tw",
					content: `<div>@component('~box')@reserve('title')@end</div>`,
				},
				{
					relPath: "components/box.tw",
					content: "<i>@slot</i>",
				},
			},
			expect: "<div><i><span>Title</span></i></div>",
		},
		// Insert block
		{
			name: "block is escaped once",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@use('layout')@insert('title')<b>&amp; {{ '<span>Title</span>' }}</b>@end`,
				},
				{
					relPath: "layout.tw",
					content: `<div>@reserve('title')</div>`,
				},
			},
			expect: "<div><b>&amp; &lt;span&gt;Title&lt;/span&gt;</b></div>",
		},
		{
			name: "block with raw string is not escaped",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@use('layout')@insert('title'){{ raw('<span>Title</span>') }}@end`,
				},
				{
					relPath: "layout.tw",
					content: `<div>@reserve('title')</div>`,
				},
			},
			expect: "<div><span>Title</span></div>",
		},
		// Reserve fallback
		{
			name: "fallback is escaped",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@use('layout')`,
				},
				{
					relPath: "layout.tw",
					content: `<div>@reserve('title', '<i>')</div>`,
				},
			},
			expect: "<div>&lt;i&gt;</div>",
		},
		{
			name: "raw fallback is not escaped",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@use('layout')`,
				},
				{
					relPath: "layout.tw",
					content: `<div>@reserve('title', raw('<i>'))</div>`,
				},
			},
			expect: "<div><i></div>",
		},
		{
			name: "string literal fallback is escaped",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@use('layout')`,
				},
				{
					relPath: "layout.tw",
					content: `<div>@reserve('title', 'A & B')</div>`,
				},
			},
			expect: "<div>A &amp; B</div>",
		},
		{
			name: "nil fallback prints nothing",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@use('layout')`,
				},
				{
					relPath: "layout.tw",
					content: `<div>[@reserve('title', nil)]</div>`,
				},
			},
			expect: "<div>[]</div>",
		},
		{
			name: "fallback array escapes only plain strings",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@use('layout')`,
				},
				{
					relPath: "layout.tw",
					content: `<div>@reserve('title', [raw('<i>'), '<i>'])</div>`,
				},
			},
			expect: "<div><i>, &lt;i&gt;</div>",
		},
		{
			name: "raw fallback is not used when page inserts",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@use('layout')@insert('title', '<br>')`,
				},
				{
					relPath: "layout.tw",
					content: `<div>@reserve('title', raw('<i>'))</div>`,
				},
			},
			expect: "<div>&lt;br&gt;</div>",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tpl, err := setupTemplate(t, tc.files, nil)
			if err != nil {
				t.Fatal(err)
			}

			actual, failure := tpl.String("index", nil)
			assertOutput(t, actual, failure, tc.expect)
		})
	}
}

func TestGlobalFuncRaw_InComponents(t *testing.T) {
	cases := []struct {
		name   string
		files  []tplFile
		expect string
	}{
		{
			name: "raw prop is not escaped",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@component('~header', {title: raw('<span>Title</span>')})@end`,
				},
				{
					relPath: "components/header.tw",
					content: `[{{ title }}]`,
				},
			},
			expect: "[<span>Title</span>]",
		},
		{
			name: "plain prop is escaped",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@component('~header', {title: '<span>Title</span>'})@end`,
				},
				{
					relPath: "components/header.tw",
					content: `[{{ title }}]`,
				},
			},
			expect: "[&lt;span&gt;Title&lt;/span&gt;]",
		},
		{
			name: "plain prop marked raw in component is not escaped",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@component('~header', {title: '<span>Title</span>'})@end`,
				},
				{
					relPath: "components/header.tw",
					content: `[{{ raw(title) }}]`,
				},
			},
			expect: "[<span>Title</span>]",
		},
		{
			name: "array prop escapes only plain string",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@component('~header', {list: [raw('<span>Title</span>'), '<span>Title</span>']})@end`,
				},
				{
					relPath: "components/header.tw",
					content: `@each(item in list)[{{ item }}]@end`,
				},
			},
			expect: "[<span>Title</span>][&lt;span&gt;Title&lt;/span&gt;]",
		},
		{
			name: "raw object prop strings are not escaped",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@component('~header', {user: raw({name: '<n>'})})@end`,
				},
				{
					relPath: "components/header.tw",
					content: `{{ user.name }}|{{ user }}`,
				},
			},
			expect: `<n>|{name: "<n>"}`,
		},
		{
			name: "default slot escapes only plain string",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@component('~header'){{ raw('<span>Title</span>') }}|{{ '<span>Title</span>' }}@end`,
				},
				{
					relPath: "components/header.tw",
					content: `<i>@slot</i>`,
				},
			},
			expect: "<i><span>Title</span>|&lt;span&gt;Title&lt;/span&gt;</i>",
		},
		{
			name: "named slot escapes only plain string",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@component('~header')@pass('s'){{ raw('<span>Title</span>') }}|{{ '<span>Title</span>' }}@end@end`,
				},
				{
					relPath: "components/header.tw",
					content: `<i>@slot('s')</i>`,
				},
			},
			expect: "<i><span>Title</span>|&lt;span&gt;Title&lt;/span&gt;</i>",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tpl, err := setupTemplate(t, tc.files, nil)
			if err != nil {
				t.Fatal(err)
			}

			actual, failure := tpl.String("index", nil)
			assertOutput(t, actual, failure, tc.expect)
		})
	}
}
