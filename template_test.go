package textwire

import (
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/textwire/textwire/v6/config"
)

func TestLocateFiles_IgnoresNonTemplateFiles(t *testing.T) {
	dir := t.TempDir()
	useTemplateDir(t, dir)

	writeTemplate(t, filepath.Join(dir, "index.tw"), "hello")
	writeTemplate(t, filepath.Join(dir, "notes.twig"), "junk @end")
	writeTemplate(t, filepath.Join(dir, "some.tw.twig"), "@end")
	writeTemplate(t, filepath.Join(dir, "card.tw.html"), "@end")

	files, err := locateFiles()
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
