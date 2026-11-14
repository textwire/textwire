package apitest

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"codeberg.org/textwire/textwire/v6/config"
)

func TestComponentDepth(t *testing.T) {
	t.Run("component hits default limit with error", func(t *testing.T) {
		content := "@component('~tree')@end"
		files := []tplFile{{relPath: "index.tw", content: content},
			{relPath: "components/tree.tw", content: content},
		}

		tpl := setupTemplate(t, files, &config.Config{})

		_, failure := tpl.String("index", nil)
		assertDepthError(t, failure, "components/tree", config.DefaultMaxComponentDepth)
	})

	t.Run("components call each other with limit 10", func(t *testing.T) {
		files := []tplFile{
			{relPath: "index.tw", content: "@component('~ping')@end"},
			{relPath: "components/ping.tw", content: "ping @component('~pong')@end"},
			{relPath: "components/pong.tw", content: "pong @component('~ping')@end"},
		}

		tpl := setupTemplate(t, files, &config.Config{MaxComponentDepth: 10})

		_, failure := tpl.String("index", nil)
		assertDepthError(t, failure, "components/ping", 10)
	})

	t.Run("nesting 3 with limit 3 succeeds", func(t *testing.T) {
		files := []tplFile{
			{relPath: "index.tw", content: "@component('~card')@end"},
			{relPath: "components/card.tw", content: "<div>@component('~avatar')@end</div>"},
			{relPath: "components/avatar.tw", content: "<span>@component('~icon')@end</span>"},
			{relPath: "components/icon.tw", content: "<i></i>"},
		}

		tpl := setupTemplate(t, files, &config.Config{MaxComponentDepth: 3})

		actual, failure := tpl.String("index", nil)
		assertOutput(t, actual, failure, "<div><span><i></i></span></div>")
	})

	t.Run("nesting over limit of 2", func(t *testing.T) {
		files := []tplFile{
			{relPath: "index.tw", content: "@component('~card')@end"},
			{relPath: "components/card.tw", content: "<div>@component('~avatar')@end</div>"},
			{relPath: "components/avatar.tw", content: "<span>@component('~icon')@end</span>"},
			{relPath: "components/icon.tw", content: "<i></i>"},
		}

		tpl := setupTemplate(t, files, &config.Config{MaxComponentDepth: 2})

		_, failure := tpl.String("index", nil)
		assertDepthError(t, failure, "components/icon", 2)
	})

	t.Run("limit of 1 allows no nesting", func(t *testing.T) {
		files := []tplFile{
			{relPath: "index.tw", content: "@component('~hero')@end"},
			{relPath: "components/hero.tw", content: "Skirk @component('~avatar')@end"},
			{relPath: "components/avatar.tw", content: "<img>"},
		}

		tpl := setupTemplate(t, files, &config.Config{MaxComponentDepth: 1})

		_, failure := tpl.String("index", nil)
		assertDepthError(t, failure, "components/avatar", 1)
	})

	t.Run("recursive menu within limit", func(t *testing.T) {
		menu := `<ul>@each(item in items)<li>{{ item.name }}@if(item.children.len() > 0)@component('~menu', { items: item.children })@end@end</li>@end</ul>`

		files := []tplFile{
			{relPath: "index.tw", content: "@component('~menu', { items })@end"},
			{relPath: "components/menu.tw", content: menu},
		}

		items := []map[string]any{
			{ // 1st depth
				"name": "Docs",
				"children": []map[string]any{
					{ // 2nd depth
						"name": "Guides",
						"children": []map[string]any{
							{ // 3rd depth
								"name":     "Components",
								"children": []map[string]any{},
							},
						},
					},
				},
			},
			{
				"name":     "Blog",
				"children": []map[string]any{},
			},
		}

		tpl := setupTemplate(t, files, &config.Config{MaxComponentDepth: 3})

		actual, failure := tpl.String("index", map[string]any{"items": items})
		expect := "<ul><li>Docs<ul><li>Guides<ul><li>Components</li></ul></li></ul></li><li>Blog</li></ul>"
		assertOutput(t, actual, failure, expect)
	})

	t.Run("siblings calls do not accumulate depth", func(t *testing.T) {
		files := []tplFile{
			{relPath: "index.tw", content: "@for(i = 0; i < 50; i++)@component('~badge')@end@end"},
			{relPath: "components/badge.tw", content: "x"},
		}

		tpl := setupTemplate(t, files, &config.Config{MaxComponentDepth: 1})

		actual, failure := tpl.String("index", nil)
		assertOutput(t, actual, failure, strings.Repeat("x", 50))
	})

	t.Run("nested siblings do not accumulate depth", func(t *testing.T) {
		files := []tplFile{
			{
				relPath: "index.tw",
				content: "@component('~card')@end@component('~card')@end@component('~card')@end",
			},
			{relPath: "components/card.tw", content: "[@component('~avatar')@end]"},
			{relPath: "components/avatar.tw", content: "x"},
		}

		tpl := setupTemplate(t, files, &config.Config{MaxComponentDepth: 2})

		actual, failure := tpl.String("index", nil)
		assertOutput(t, actual, failure, "[x][x][x]")
	})

	t.Run("limit above default", func(t *testing.T) {
		node := "x@if(n < limit)@component('~node', { n: n + 1, limit })@end@end"

		files := []tplFile{
			{relPath: "index.tw", content: "@component('~node', { n: 1, limit })@end"},
			{relPath: "components/node.tw", content: node},
		}

		tpl := setupTemplate(t, files, &config.Config{MaxComponentDepth: 1500})

		actual, failure := tpl.String("index", map[string]any{"limit": 1500})
		assertOutput(t, actual, failure, strings.Repeat("x", 1500))

		_, failure = tpl.String("index", map[string]any{"limit": 1501})
		assertDepthError(t, failure, "components/node", 1500)
	})

	t.Run("failed render does not affect next render", func(t *testing.T) {
		node := "x@if(n < limit)@component('~node', { n: n + 1, limit })@end@end"

		files := []tplFile{
			{relPath: "index.tw", content: "@component('~node', { n: 1, limit })@end"},
			{relPath: "components/node.tw", content: node},
		}

		tpl := setupTemplate(t, files, &config.Config{MaxComponentDepth: 3})

		_, failure := tpl.String("index", map[string]any{"limit": 4})
		assertDepthError(t, failure, "components/node", 3)

		actual, failure := tpl.String("index", map[string]any{"limit": 3})
		assertOutput(t, actual, failure, "xxx")
	})

}

func TestComponentDepth_LimitBelowOneUsesDefault(t *testing.T) {
	// Recursive component: renders "x" once per level and stops
	// recursing once the depth counter `n` reaches `limit`.
	files := []tplFile{
		{
			relPath: "index.tw",
			content: "@component('~node', { n: 1, limit })@end",
		},
		{
			relPath: "components/node.tw",
			content: "x@if(n < limit)@component('~node', { n: n + 1, limit })@end@end",
		},
	}

	defaultDepth := config.DefaultMaxComponentDepth

	for _, configuredDepth := range []int{0, -1} {
		t.Run(
			fmt.Sprintf("limit below 1 uses default, limit %d", configuredDepth),
			func(t *testing.T) {
				tpl := setupTemplate(t, files, &config.Config{MaxComponentDepth: configuredDepth})

				// Recursing exactly to the default depth succeeds.
				actual, failure := tpl.String("index", map[string]any{"limit": defaultDepth})
				assertOutput(t, actual, failure, strings.Repeat("x", defaultDepth))

				// Recursing one level past the default depth fails.
				_, failure = tpl.String("index", map[string]any{"limit": defaultDepth + 1})
				assertDepthError(t, failure, "components/node", defaultDepth)
			},
		)
	}
}

func TestComponentDepth_ConcurrentRendersDoNotShareDepth(t *testing.T) {
	files := []tplFile{
		{
			relPath: "index.tw",
			content: "@component('~node', { n: 1, limit })@end",
		},
		{
			relPath: "components/node.tw",
			content: "x@if(n < limit)@component('~node', { n: n + 1, limit })@end@end",
		},
	}

	tpl := setupTemplate(t, files, &config.Config{MaxComponentDepth: 5})

	var wg sync.WaitGroup

	for range 20 {
		wg.Go(func() {
			if _, failure := tpl.String("index", map[string]any{"limit": 6}); failure == nil {
				t.Error("expected depth error, got none")
			}
		})

		wg.Go(func() {
			actual, failure := tpl.String("index", map[string]any{"limit": 5})
			if failure != nil {
				t.Errorf("unexpected error: %s", failure)
				return
			}

			if actual != "xxxxx" {
				t.Errorf("wrong result. Expect %q, got %q", "xxxxx", actual)
			}
		})
	}

	wg.Wait()
}
