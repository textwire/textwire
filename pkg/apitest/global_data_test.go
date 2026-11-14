package apitest

import (
	"codeberg.org/textwire/textwire/v6"
	"fmt"
	"sync"
	"testing"

	"codeberg.org/textwire/textwire/v6/config"
	"codeberg.org/textwire/textwire/v6/pkg/fail"
)

func TestGlobalData_UnsupportedValue(t *testing.T) {
	unsupportedFn := func() {}
	unsupportedChan := make(chan int)

	fnType := fmt.Sprintf("%T", unsupportedFn)
	chanType := fmt.Sprintf("%T", unsupportedChan)

	cases := []struct {
		name   string
		auth   any
		expect string
	}{
		{
			name:   "func in map",
			auth:   map[string]any{"bad": unsupportedFn},
			expect: fmt.Sprintf(fail.ErrUnsupportedType, fnType),
		},
		{
			name:   "chan in map",
			auth:   map[string]any{"bad": unsupportedChan},
			expect: fmt.Sprintf(fail.ErrUnsupportedType, chanType),
		},
		{
			name:   "func in slice",
			auth:   []any{unsupportedFn},
			expect: fmt.Sprintf(fail.ErrUnsupportedType, fnType),
		},
		{
			name:   "chan in slice",
			auth:   []any{unsupportedChan},
			expect: fmt.Sprintf(fail.ErrUnsupportedType, chanType),
		},
		{
			name:   "func struct field",
			auth:   struct{ Callback func() }{unsupportedFn},
			expect: fmt.Sprintf(fail.ErrUnsupportedFieldType, fnType, "Callback"),
		},
		{
			name:   "chan struct field",
			auth:   struct{ Events chan int }{unsupportedChan},
			expect: fmt.Sprintf(fail.ErrUnsupportedFieldType, chanType, "Events"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resetConfigAfterTest(t)

			_, failure := textwire.NewTemplate(&config.Config{
				GlobalData: map[string]any{"env": "dev", "auth": tc.auth},
			})

			if failure == nil {
				t.Fatal("expected error for unsupported global data, got none")
			}

			if failure.Message() != tc.expect {
				t.Fatalf("wrong error message. Expect %q but got %q", tc.expect, failure.Message())
			}
		})
	}
}

func TestGlobalData(t *testing.T) {
	t.Run("global is accessible inside layout", func(t *testing.T) {
		files := []tplFile{
			{relPath: "index.tw", content: `@use('layout')@insert('title', 'Home')`},
			{relPath: "layout.tw", content: `<h1>@reserve('title')</h1>{{ global.site }}`},
		}

		conf := &config.Config{GlobalData: map[string]any{"site": "textwire"}}
		tpl, err := setupTemplate(t, files, conf)
		if err != nil {
			t.Fatal(err)
		}

		actual, failure := tpl.String("index", nil)
		assertOutput(t, actual, failure, "<h1>Home</h1>textwire")
	})

	t.Run("writing to global inside layout triggers error", func(t *testing.T) {
		files := []tplFile{
			{relPath: "index.tw", content: `@use('layout')@insert('title', 'Home')`},
			{relPath: "layout.tw", content: `<h1>@reserve('title')</h1>{{ global.site = "x" }}`},
		}

		conf := &config.Config{GlobalData: map[string]any{"site": "textwire"}}
		tpl, err := setupTemplate(t, files, conf)
		if err != nil {
			t.Fatal(err)
		}

		_, failure := tpl.String("index", nil)
		assertReadOnlyValueError(t, failure)
	})

	t.Run("global is accessible inside component", func(t *testing.T) {
		files := []tplFile{
			{relPath: "index.tw", content: `@component('~header')@end`},
			{relPath: "components/header.tw", content: `<b>{{ global.site }}</b>`},
		}

		conf := &config.Config{GlobalData: map[string]any{"site": "textwire"}}
		tpl, err := setupTemplate(t, files, conf)
		if err != nil {
			t.Fatal(err)
		}

		actual, failure := tpl.String("index", nil)
		assertOutput(t, actual, failure, "<b>textwire</b>")
	})

	t.Run("writing to global inside component triggers error", func(t *testing.T) {
		files := []tplFile{
			{relPath: "index.tw", content: `@component('~header')@end`},
			{relPath: "components/header.tw", content: `<b>{{ global.site = "x" }}</b>`},
		}

		conf := &config.Config{GlobalData: map[string]any{"site": "textwire"}}
		tpl, err := setupTemplate(t, files, conf)
		if err != nil {
			t.Fatal(err)
		}

		_, failure := tpl.String("index", nil)
		assertReadOnlyValueError(t, failure)
	})

	t.Run("writing to global inside @each triggers error", func(t *testing.T) {
		conf := &config.Config{GlobalData: map[string]any{
			"users": []any{map[string]any{"name": "Alice"}},
		}}

		files := []tplFile{
			{relPath: "index.tw", content: `@each(u in global.users){{ u.name = "x" }}@end`},
		}

		tpl, err := setupTemplate(t, files, conf)
		if err != nil {
			t.Fatal(err)
		}

		_, failure := tpl.String("index", nil)
		assertReadOnlyValueError(t, failure)
	})
}

func TestGlobalData_ConcurrentRenders(t *testing.T) {
	conf := &config.Config{GlobalData: map[string]any{
		"site": "textwire",
		"user": map[string]any{"name": "Alice"},
	}}

	files := []tplFile{
		{relPath: "read.tw", content: `{{ global.site; global.user.name }}`},
		{relPath: "write.tw", content: `{{ global.site = "x" }}`},
	}

	tpl, err := setupTemplate(t, files, conf)
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup

	for range 20 {
		wg.Go(func() {
			actual, failure := tpl.String("read", nil)
			if failure != nil {
				t.Errorf("unexpected error: %s", failure)
				return
			}

			if actual != "textwireAlice" {
				t.Errorf("wrong result. Expect %q, got %q", "textwireAlice", actual)
			}
		})

		wg.Go(func() {
			_, failure := tpl.String("write", nil)
			if failure == nil {
				t.Error("expected read-only error, got none")
				return
			}

			if failure.Message() != fail.ErrReadOnlyValue {
				t.Errorf(
					"wrong error message. Expect %q, got %q",
					fail.ErrReadOnlyValue,
					failure.Message(),
				)
			}
		})
	}

	wg.Wait()
}
