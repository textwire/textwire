package textwire

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"codeberg.org/textwire/textwire/v5/config"
	"codeberg.org/textwire/textwire/v5/pkg/ast"
	"codeberg.org/textwire/textwire/v5/pkg/file"
	"codeberg.org/textwire/textwire/v5/pkg/linker"
)

func TestUpdateFileIfModified_KeepsOldProgramOnSyntaxError(t *testing.T) {
	const (
		validTemplate  = "Hello, {{ name }}"
		brokenTemplate = "@if(true) no end" // the @if is never closed
	)

	path := filepath.Join(t.TempDir(), "home.tw")
	writeTemplate(t, path, validTemplate)

	f := file.New("home", path, path, config.New())
	prog := mustParseFile(t, f)

	fw := newFileWatcher(linker.New([]*ast.Program{prog}))
	fw.files = []*file.SourceFile{f}

	// The file changes on disk into something that no longer parses
	writeTemplate(t, path, brokenTemplate)
	setModTimeAfterLastCheck(t, path, f.ModTime)

	fw.updateFileIfModified(f)

	// The watcher survives and keeps serving the last valid program
	programs := fw.linker.Programs
	if len(programs) != 1 {
		t.Fatalf("expected to keep the old program, got %d programs", len(programs))
	}

	if programs[0].Name != "home" {
		t.Fatalf("expected to keep program %q, got %q", "home", programs[0].Name)
	}
}

func mustParseFile(t *testing.T, f *file.SourceFile) *ast.Program {
	t.Helper()

	prog, failure, err := parseFile(f)
	if err != nil {
		t.Fatalf("failed to parse valid template: %v", err)
	}
	if failure != nil {
		t.Fatalf("failed to parse valid template: %v", failure.Error())
	}

	return prog
}

func writeTemplate(t *testing.T, path, content string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write template: %v", err)
	}
}

func setModTimeAfterLastCheck(t *testing.T, path string, lastCheck time.Time) {
	t.Helper()

	modTime := lastCheck.Add(time.Second)
	if err := os.Chtimes(path, modTime, modTime); err != nil {
		t.Fatalf("failed to change template mod time: %v", err)
	}
}
