package textwire

import (
	"fmt"
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
			tpl, failure := NewTemplate(&config.Config{
				TemplateDir: "testdata/bad/unsupported-global-data",
				GlobalData:  map[string]any{"env": "dev", "auth": tc.auth},
			})

			if failure != nil {
				t.Fatalf("Unexpected template error: %s", failure)
			}

			_, failure = tpl.String("index", nil)
			if failure == nil {
				t.Fatal("Expected error for unsupported global data, got none")
			}

			if failure.Message() != tc.expect {
				t.Fatalf("Wrong error message. Expect %q but got %q", tc.expect, failure.Message())
			}
		})
	}
}
