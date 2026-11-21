package apitest

import (
	"path/filepath"
	"testing"

	"github.com/textwire/textwire/v6/config"
	"github.com/textwire/textwire/v6/pkg/fail"
	"github.com/textwire/textwire/v6/pkg/position"
)

func TestTemplateErrorHandling(t *testing.T) {
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
			errPos:  &position.Pos{StartCol: 35, EndCol: 40},
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
			name: "@use inside component",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: "A@component('card')@end B",
				},
				{
					relPath: "card.tw",
					content: "C@use('layout')",
				},
				{relPath: "layout.tw", content: "L"},
			},
			data:    nil,
			errFile: "card.tw",
			errPos:  &position.Pos{StartCol: 1, EndCol: 14},
			origin:  fail.OriginLink,
			errMsg:  fail.ErrUseDirIsNotAllowedInComp,
			errArgs: nil,
		},
		{
			name: "@reserve inside component",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: "@use('layout')\n@insert('t', 'X')\n@component('card')@end",
				},
				{
					relPath: "layout.tw",
					content: "@reserve('t')",
				},
				{relPath: "card.tw", content: "C@reserve('t', 'fallback')D"},
			},
			data:    nil,
			errFile: "card.tw",
			errPos:  &position.Pos{StartCol: 1, EndCol: 25},
			origin:  fail.OriginLink,
			errMsg:  fail.ErrReserveIsNotAllowedInComp,
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
			name: "unknown default @pass",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: `@component("book")<p>Harry Potter</p>@end`,
				},
				{
					relPath: "book.tw",
					content: `@slot("content")`,
				},
			},
			data:    nil,
			errFile: "index.tw",
			errPos:  &position.Pos{StartCol: 18, EndCol: 36},
			origin:  fail.OriginLink,
			errMsg:  fail.ErrDefaultSlotNotDefined,
			errArgs: []any{"book", "book"},
		},
		{
			name: "duplicate @pass",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: "@component('user')\n@pass('content')<p>1</p>@end\n@pass('content')<p>2</p>@end\n@end",
				},
				{
					relPath: "user.tw",
					content: `@slot("content")`,
				},
			},
			data:    nil,
			errFile: "index.tw",
			errPos:  &position.Pos{StartLine: 2, StartCol: 6, EndLine: 2, EndCol: 14},
			origin:  fail.OriginPars,
			errMsg:  fail.ErrDuplicatePass,
			errArgs: []any{"content", 2, "user"},
		},
		{
			name: "unknown @component",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: "@component('unknown-name')@end",
				},
			},
			data:    nil,
			errFile: "index.tw",
			errPos:  &position.Pos{StartCol: 11, EndCol: 24},
			origin:  fail.OriginLink,
			errMsg:  fail.ErrUndefinedComponent,
			errArgs: []any{"unknown-name"},
		},
		{
			name: "undefined @insert",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: "@use('layout')\n@insert('missing', 'Study')",
				},
				{
					relPath: "layout.tw",
					content: "@reserve('title')@reserve('main')",
				},
			},
			data:    nil,
			errFile: "index.tw",
			errPos:  &position.Pos{StartLine: 1, EndLine: 1, StartCol: 8, EndCol: 16},
			origin:  fail.OriginLink,
			errMsg:  fail.ErrUnusedInsertDetected,
			errArgs: []any{"missing", "missing"},
		},
		{
			name: "undefined variable inside component file",
			files: []tplFile{
				{
					relPath: "index.tw",
					content: "@component('hero')<h1>Title</h1>@end",
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
			name: "undefined variable inside layout",
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
			name: "undefined @use",
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
			name: "undefined variable inside nested compoment file",
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
			name: "undefined variable inside layout",
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
			name: "component file calls itself recursively",
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
			name: "component calls each other recursively",
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
			name: "inserts without @use",
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
