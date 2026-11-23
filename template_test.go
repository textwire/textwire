package textwire

import (
	"fmt"
	"path/filepath"
	"sync"
	"testing"
	"testing/fstest"

	"github.com/textwire/textwire/v6/config"
	"github.com/textwire/textwire/v6/pkg/fail"
)

func TestLocateFiles_IgnoresNonTemplateFiles(t *testing.T) {
	dir := t.TempDir()
	conf := useTemplateDir(t, dir)

	writeTemplate(t, filepath.Join(dir, "index.tw"), "hello")
	writeTemplate(t, filepath.Join(dir, "notes.twig"), "junk @end")
	writeTemplate(t, filepath.Join(dir, "some.tw.twig"), "@end")
	writeTemplate(t, filepath.Join(dir, "card.tw.html"), "@end")

	files, err := locateFiles(conf)
	if err != nil {
		t.Fatalf("failed to locate files: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("expected 1 template file, got %d", len(files))
	}

	if files[0].Name != "index" {
		t.Fatalf("expected template %q, got %q", "index", files[0].Name)
	}
}

func TestNewTemplate_WithRootTemplateDirAndFS(t *testing.T) {
	tpl, failure := NewTemplate(&config.Config{
		TemplateDir: ".",
		TemplateFS: fstest.MapFS{
			"index.tw": &fstest.MapFile{Data: []byte("hello")},
		},
	})

	if failure != nil {
		t.Fatalf("failed to create template: %v", failure)
	}

	actual, failure := tpl.String("index", nil)
	if failure != nil {
		t.Fatalf("failed to render template: %v", failure)
	}

	if actual != "hello" {
		t.Fatalf("expected %q, got %q", "hello", actual)
	}
}

func TestNewTemplate_WithFSIgnoresFilesOutsideTemplateDir(t *testing.T) {
	tpl, failure := NewTemplate(&config.Config{
		TemplateDir: "templates",
		TemplateFS: fstest.MapFS{
			"templates/index.tw": &fstest.MapFile{Data: []byte("hello")},
			"docs/readme.tw":     &fstest.MapFile{Data: []byte("hello")},
		},
	})
	if failure != nil {
		t.Fatalf("failed to create template: %v", failure)
	}

	actual, failure := tpl.String("index", nil)
	if failure != nil {
		t.Fatalf("failed to render template: %v", failure)
	}

	if actual != "hello" {
		t.Fatalf("expected %q, got %q", "hello", actual)
	}

	if _, failure := tpl.String("docs/readme", nil); failure == nil {
		t.Fatal("expected file outside of template dir to be ignored")
	}
}

func TestNewTemplate_KeepsOwnConfig(t *testing.T) {
	t.Run("component depth is not changed by another template", func(t *testing.T) {
		recFS := fstest.MapFS{
			"c/rec.tw": &fstest.MapFile{Data: []byte(`x@component("rec")@end`)},
		}

		first, failure := NewTemplate(&config.Config{
			TemplateDir:       "c",
			TemplateFS:        recFS,
			MaxComponentDepth: 3,
		})

		if failure != nil {
			t.Fatalf("failed to create template: %v", failure)
		}

		assertCompDepthLimit(t, first, "rec", 3)

		_, failure = NewTemplate(&config.Config{
			TemplateDir:       "b",
			TemplateFS:        fstest.MapFS{"b/page.tw": &fstest.MapFile{Data: []byte("page")}},
			MaxComponentDepth: 10,
		})

		if failure != nil {
			t.Fatalf("failed to create second template: %v", failure)
		}

		assertCompDepthLimit(t, first, "rec", 3)
	})

	t.Run("global data does not leak from another template", func(t *testing.T) {
		pageFS := fstest.MapFS{
			"b/page.tw": &fstest.MapFile{Data: []byte("{{ global.env }}")},
		}

		_, failure := NewTemplate(&config.Config{
			TemplateDir: "b",
			TemplateFS:  pageFS,
			GlobalData:  map[string]any{"env": "prod"},
		})

		if failure != nil {
			t.Fatalf("failed to create first template: %v", failure)
		}

		next, failure := NewTemplate(&config.Config{
			TemplateDir: "b",
			TemplateFS:  pageFS,
		})

		if failure != nil {
			t.Fatalf("failed to create second template: %v", failure)
		}

		actual, failure := next.String("page", nil)
		if failure != nil {
			t.Fatalf("failed to render template: %v", failure)
		}

		if actual != "" {
			t.Fatalf("expected no global data, got %q", actual)
		}
	})

	t.Run("templates keep their own directory", func(t *testing.T) {
		first, failure := NewTemplate(&config.Config{
			TemplateDir: "a",
			TemplateFS:  fstest.MapFS{"a/home.tw": &fstest.MapFile{Data: []byte("A home")}},
		})
		if failure != nil {
			t.Fatalf("failed to create first template: %v", failure)
		}

		second, failure := NewTemplate(&config.Config{
			TemplateDir: "b",
			TemplateFS:  fstest.MapFS{"b/page.tw": &fstest.MapFile{Data: []byte("B page")}},
		})
		if failure != nil {
			t.Fatalf("failed to create second template: %v", failure)
		}

		actual, failure := first.String("home", nil)
		if failure != nil {
			t.Fatalf("failed to render first template: %v", failure)
		}

		if actual != "A home" {
			t.Fatalf("expected %q, got %q", "A home", actual)
		}

		actual, failure = second.String("page", nil)
		if failure != nil {
			t.Fatalf("failed to render second template: %v", failure)
		}

		if actual != "B page" {
			t.Fatalf("expected %q, got %q", "B page", actual)
		}

		// page exists only in the second template's directory, so the first
		// template must report it under its own directory.
		assertTemplateNotFoundError(t, first, "page", "a/page.tw")
	})
}

func TestNewTemplate_OptionsDoNotBecomeDefaults(t *testing.T) {
	resetConfig(t)

	chainFS := fstest.MapFS{
		"x/index.tw":         &fstest.MapFile{Data: []byte("@component('~c1')@end")},
		"x/components/c1.tw": &fstest.MapFile{Data: []byte("@component('~c2')@end")},
		"x/components/c2.tw": &fstest.MapFile{Data: []byte("@component('~c3')@end")},
		"x/components/c3.tw": &fstest.MapFile{Data: []byte("@component('~c4')@end")},
		"x/components/c4.tw": &fstest.MapFile{Data: []byte("ok")},
	}

	_, failure := NewTemplate(&config.Config{
		TemplateDir:       "x",
		TemplateFS:        chainFS,
		MaxComponentDepth: 3,
	})

	if failure != nil {
		t.Fatalf("failed to create first template: %v", failure)
	}

	next, failure := NewTemplate(&config.Config{
		TemplateDir: "x",
		TemplateFS:  chainFS,
	})

	if failure != nil {
		t.Fatalf("failed to create second template: %v", failure)
	}

	actual, failure := next.String("index", nil)
	if failure != nil {
		t.Fatalf("failed to render template: %v", failure)
	}

	if actual != "ok" {
		t.Fatalf("expected %q, got %q", "ok", actual)
	}
}

func TestConfigure_DoesNotAffectExistingTemplates(t *testing.T) {
	resetConfig(t)

	dirA := t.TempDir()
	writeTemplate(t, filepath.Join(dirA, "home.tw"), "A home")

	existing, failure := NewTemplate(&config.Config{TemplateDir: dirA})
	if failure != nil {
		t.Fatalf("failed to create existing template: %v", failure)
	}

	dirB := t.TempDir()
	writeTemplate(t, filepath.Join(dirB, "home.tw"), "B home")
	writeTemplate(t, filepath.Join(dirB, "page.tw"), "{{ global.env }}")

	Configure(&config.Config{
		TemplateDir: dirB,
		GlobalData:  map[string]any{"env": "prod"},
	})

	future, failure := NewTemplate(nil)
	if failure != nil {
		t.Fatalf("failed to create future template: %v", failure)
	}

	actual, failure := future.String("home", nil)
	if failure != nil {
		t.Fatalf("failed to render future template: %v", failure)
	}

	if actual != "B home" {
		t.Fatalf("expected %q, got %q", "B home", actual)
	}

	actual, failure = future.String("page", nil)
	if failure != nil {
		t.Fatalf("failed to render future template: %v", failure)
	}

	if actual != "prod" {
		t.Fatalf("expected %q, got %q", "prod", actual)
	}

	actual, failure = existing.String("home", nil)
	if failure != nil {
		t.Fatalf("failed to render existing template: %v", failure)
	}

	if actual != "A home" {
		t.Fatalf("expected %q, got %q", "A home", actual)
	}
}

func TestConfigure_NilKeepsTemplateFS(t *testing.T) {
	resetConfig(t)

	Configure(&config.Config{
		TemplateDir: "embedded",
		TemplateFS: fstest.MapFS{
			"embedded/index.tw": &fstest.MapFile{Data: []byte("embedded")},
		},
	})

	Configure(nil)

	tpl, failure := NewTemplate(nil)
	if failure != nil {
		t.Fatalf("failed to create template: %v", failure)
	}

	actual, failure := tpl.String("index", nil)
	if failure != nil {
		t.Fatalf("failed to render template: %v", failure)
	}

	if actual != "embedded" {
		t.Fatalf("expected %q, got %q", "embedded", actual)
	}
}

func TestNewTemplate_ConcurrentConfigKeepsPerTemplateDepth(t *testing.T) {
	resetConfig(t)

	// Finite chain: a -> b -> c -> d -> e
	chainFS := fstest.MapFS{
		"c/a.tw": &fstest.MapFile{Data: []byte(`a@component("b")@end`)},
		"c/b.tw": &fstest.MapFile{Data: []byte(`b@component("c")@end`)},
		"c/c.tw": &fstest.MapFile{Data: []byte(`c@component("d")@end`)},
		"c/d.tw": &fstest.MapFile{Data: []byte(`d@component("e")@end`)},
		"c/e.tw": &fstest.MapFile{Data: []byte(`e`)},
	}

	newConfig := func(depth int) *config.Config {
		return &config.Config{
			TemplateDir:       "c",
			TemplateFS:        chainFS,
			MaxComponentDepth: depth,
		}
	}

	shallow, failure := NewTemplate(newConfig(3))
	if failure != nil {
		t.Fatalf("failed to create shallow template: %v", failure)
	}

	deep, failure := NewTemplate(newConfig(10))
	if failure != nil {
		t.Fatalf("failed to create deep template: %v", failure)
	}

	// Sequential baselines, so we don't have to guess exact messages.
	_, shallowFailure := shallow.String("a", nil)
	if shallowFailure == nil {
		t.Fatal("baseline: expected depth error with MaxComponentDepth=3")
	}

	expectErr := shallowFailure.String()
	expect, deepFailure := deep.String("a", nil)
	if deepFailure != nil {
		t.Fatalf("baseline: expected success with MaxComponentDepth=10, got: %v", deepFailure)
	}

	var wg sync.WaitGroup

	for range 10 {
		// Existing template with a low limit must keep failing the same way,
		// no matter what the global config becomes.
		wg.Go(func() {
			_, failure := shallow.String("a", nil)
			if failure == nil {
				t.Error("shallow: expected a depth error, got none")
				return
			}
			if got := failure.String(); got != expectErr {
				t.Errorf("shallow: error changed\n got: %s\nwant: %s", got, expectErr)
			}
		})

		// Existing template with a high limit must keep succeeding
		// with identical output.
		wg.Go(func() {
			actual, failure := deep.String("a", nil)
			if failure != nil {
				t.Errorf("deep: unexpected error: %v", failure)
				return
			}
			if actual != expect {
				t.Errorf("deep: output changed: got %q, want %q", actual, expect)
			}
		})

		// Creating new templates concurrently must work and honor its own config.
		wg.Go(func() {
			tpl, failure := NewTemplate(newConfig(3))
			if failure != nil {
				t.Errorf("new template: %v", failure)
				return
			}
			if _, failure := tpl.String("a", nil); failure == nil {
				t.Error("new template (depth 3): expected a depth error, got none")
			}
		})

		// Global reconfiguration racing with everything above.
		wg.Go(func() {
			Configure(newConfig(20))
		})
	}

	wg.Wait()
}

func assertCompDepthLimit(t *testing.T, tpl *Template, compName string, maxDepth int) {
	t.Helper()

	_, failure := tpl.String(compName, nil)
	if failure == nil {
		t.Fatal("expected a component depth error, got none")
	}

	expect := fmt.Sprintf(fail.ErrMaxComponentDepth, compName, maxDepth)
	if failure.Message() != expect {
		t.Fatalf("expected %q, got %q", expect, failure.Message())
	}
}

// assertTemplateNotFoundError checks that rendering a missing template returns
// the template-not-found error and reports the given path.
func assertTemplateNotFoundError(t *testing.T, tpl *Template, name, expectPath string) {
	t.Helper()

	_, failure := tpl.String(name, nil)
	if failure == nil {
		t.Fatal("expected a template not found error, got none")
	}

	expectMsg := fmt.Sprintf(fail.ErrTemplateNotFound, name)
	if failure.Message() != expectMsg {
		t.Fatalf("expected message %q, got %q", expectMsg, failure.Message())
	}

	if failure.Filepath() != expectPath {
		t.Fatalf("expected file path %q, got %q", expectPath, failure.Filepath())
	}
}
