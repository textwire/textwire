package textwire

import (
	"path/filepath"
	"testing"
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
