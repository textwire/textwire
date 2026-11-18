package apitest

import (
	"testing"

	"github.com/textwire/textwire/v6/config"
)

func TestReadonly_AssignmentFails(t *testing.T) {
	cases := []struct {
		name    string
		content string
		conf    *config.Config
	}{
		{
			name:    "property",
			content: `{{ global.limit = 99 }}`,
			conf:    &config.Config{GlobalData: map[string]any{"limit": 5}},
		},
		{
			name:    "nested property",
			content: `{{ global.user.name = "Bob" }}`,
			conf: &config.Config{
				GlobalData: map[string]any{
					"user": map[string]any{"name": "Alice"},
				},
			},
		},
		{
			name:    "array element",
			content: `{{ global.tags[0] = "z" }}`,
			conf: &config.Config{
				GlobalData: map[string]any{
					"tags": []string{"one", "two"},
				},
			},
		},
		{
			name:    "increment",
			content: `{{ global.limit++ }}`,
			conf:    &config.Config{GlobalData: map[string]any{"limit": 5}},
		},
		{
			name:    "aliased object",
			content: `{{ u = global.user; u.name = "Bob" }}`,
			conf: &config.Config{
				GlobalData: map[string]any{
					"user": map[string]any{"name": "Alice"},
				},
			},
		},
		{
			name:    "aliased array",
			content: `{{ tags = global.tags; tags[0] = "three" }}`,
			conf: &config.Config{
				GlobalData: map[string]any{
					"tags": []string{"one", "two"},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			files := []tplFile{{relPath: "index.tw", content: tc.content}}
			tpl, err := setupTemplate(t, files, tc.conf)
			if err != nil {
				t.Fatal(err)
			}

			_, failure := tpl.String("index", nil)
			assertReadOnlyValueError(t, failure)
		})
	}
}
