package textwire

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"codeberg.org/textwire/textwire/v6/config"
	"codeberg.org/textwire/textwire/v6/pkg/ast"
	"codeberg.org/textwire/textwire/v6/pkg/fail"
	"codeberg.org/textwire/textwire/v6/pkg/file"
	"codeberg.org/textwire/textwire/v6/pkg/linker"
)

const (
	validTpl         = "Hello {{ name }}"
	validModifiedTpl = validTpl + "!"
	brokenTpl        = "@if(true) no end" // the @if is never closed
)

// newTestFileWatcher returns a watcher that doesn't print to stdout, keeping
// test output clean when tests use intentionally broken templates.
func newTestFileWatcher(ln *linker.NodeLinker) *fileWatcher {
	fw := newFileWatcher(ln)
	fw.logger.out = io.Discard
	return fw
}

// newStartedWatcher mirrors Watch setup without starting the ticker, so
// tests can drive tick() directly.
func newStartedWatcher(t *testing.T, ln *linker.NodeLinker) *fileWatcher {
	t.Helper()

	fw := newTestFileWatcher(ln)

	files, err := locateFiles()
	if err != nil {
		t.Fatalf("failed to locate files: %v", err)
	}

	fw.files = files

	return fw
}

func TestRelinkAllPrograms(t *testing.T) {
	t.Run("replaces program with changed file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "home.tw")
		writeTemplate(t, path, validTpl)

		fw := newTestFileWatcher(linker.New(nil))
		fw.files = []*file.SourceFile{file.New("home", path, path, config.New())}
		fw.relinkAllPrograms()

		writeTemplate(t, path, validModifiedTpl)
		fw.relinkAllPrograms()

		programs := fw.linker.Programs
		if len(programs) != 1 {
			t.Fatalf("expected 1 program, got %d", len(programs))
		}

		expect := evalTpl(t, validModifiedTpl)
		got := evalTpl(t, programs[0])
		if got != expect {
			t.Fatalf("expected rendered %q, got %q", expect, got)
		}
	})

	t.Run("re-parses last valid source into a fresh program on syntax error", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "home.tw")
		writeTemplate(t, path, validTpl)

		fw := newTestFileWatcher(linker.New(nil))
		fw.files = []*file.SourceFile{file.New("home", path, path, config.New())}

		// First pass: file parses, so "home" is published and its source cached.
		fw.relinkAllPrograms()
		firstProg := fw.linker.Programs[0]

		// File becomes syntactically broken.
		writeTemplate(t, path, brokenTpl)
		fw.relinkAllPrograms()

		programs := fw.linker.Programs

		// The broken file is not dropped: the cached last-valid source is
		// re-parsed, so "home" keeps rendering its last good content.
		if len(programs) != 1 || programs[0].Name != "home" {
			t.Fatalf("expected to keep program %q, got %v", "home", programs)
		}

		// The re-parsed program must be a NEW object: the old published
		// program is never mutated/re-linked, so renders already running
		// that still hold it.
		if programs[0] == firstProg {
			t.Fatal("expected last valid source to be parsed into a fresh program, not reused")
		}

		expect := evalTpl(t, validTpl)
		got := evalTpl(t, programs[0])
		if got != expect {
			t.Fatalf("expected rendered %q, got %q", expect, got)
		}
	})

	t.Run("omits file with no valid source yet", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "home.tw")
		writeTemplate(t, path, brokenTpl)

		fw := newTestFileWatcher(linker.New(nil))
		fw.files = []*file.SourceFile{file.New("home", path, path, config.New())}
		fw.relinkAllPrograms()

		if programs := fw.linker.Programs; len(programs) != 0 {
			t.Fatalf("broken program should not be in linker, got %v", programs)
		}
	})

	t.Run("valid template still updates when another has a syntax error", func(t *testing.T) {
		dir := t.TempDir()
		homePath := filepath.Join(dir, "home.tw")
		aboutPath := filepath.Join(dir, "about.tw")
		writeTemplate(t, homePath, validTpl)
		writeTemplate(t, aboutPath, validTpl)

		fw := newTestFileWatcher(linker.New(nil))
		fw.files = []*file.SourceFile{
			file.New("home", homePath, homePath, config.New()),
			file.New("about", aboutPath, aboutPath, config.New()),
		}

		fw.relinkAllPrograms()

		writeTemplate(t, homePath, brokenTpl)
		writeTemplate(t, aboutPath, validModifiedTpl)
		fw.relinkAllPrograms()

		programs := fw.linker.Programs
		if len(programs) != 2 {
			t.Fatalf("expected 2 programs, got %d", len(programs))
		}

		expect := evalTpl(t, validModifiedTpl)
		got := evalTpl(t, programs[1])
		if got != expect {
			t.Fatalf("expected program 'about' to render %q, got %q", expect, got)
		}
	})
}

func TestRelinkAllPrograms_LogsRefreshedFileOnModTimeChange(t *testing.T) {
	path := filepath.Join(t.TempDir(), "home.tw")
	writeTemplate(t, path, validTpl)

	var out bytes.Buffer
	fw := newFileWatcher(linker.New(nil))
	fw.logger.out = &out
	fw.files = []*file.SourceFile{file.New("home", path, path, config.New())}

	fw.relinkAllPrograms()
	if strings.Contains(out.String(), "refreshed") {
		t.Fatalf("expected no refresh log on first parse, got %q", out.String())
	}

	fw.relinkAllPrograms()
	if strings.Contains(out.String(), "refreshed") {
		t.Fatalf("expected no refresh log when file is unchanged, got %q", out.String())
	}

	writeTemplate(t, path, validModifiedTpl)
	modifiedAt := time.Now().Add(time.Second)
	if err := os.Chtimes(path, modifiedAt, modifiedAt); err != nil {
		t.Fatalf("failed to set modification time: %v", err)
	}

	fw.relinkAllPrograms()

	if !strings.Contains(out.String(), "refreshed "+fw.files[0].Rel) {
		t.Fatalf("expected a refresh log for modified file, got %q", out.String())
	}
}

func TestRemoveProgramByName(t *testing.T) {
	t.Run("removes only matching program", func(t *testing.T) {
		ln := linker.New([]*ast.Program{
			{Name: "home"},
			{Name: "about"},
		})

		fw := newTestFileWatcher(ln)
		fw.removeProgramByName("home")

		programs := fw.linker.Programs
		if len(programs) != 1 {
			t.Fatalf("expected 1 program, got %d", len(programs))
		}

		if programs[0].Name != "about" {
			t.Fatalf("expected only %q to remain, got %v", "about", programs)
		}
	})

	t.Run("keeps all programs when name not found", func(t *testing.T) {
		ln := linker.New([]*ast.Program{
			{Name: "home"},
			{Name: "about"},
		})

		fw := newTestFileWatcher(ln)
		fw.removeProgramByName("contact")

		programs := fw.linker.Programs
		if len(programs) != 2 || programs[0].Name != "home" || programs[1].Name != "about" {
			t.Fatalf("expected all programs to remain, got %v", programs)
		}
	})
}

func TestFindDeletedFiles_ReturnsFilesMissingFromCurrentList(t *testing.T) {
	fw := newTestFileWatcher(linker.New(nil))
	fw.files = []*file.SourceFile{{Name: "home"}}

	oldFiles := []*file.SourceFile{{Name: "home"}, {Name: "about"}}

	deleted := fw.findDeletedFiles(oldFiles)
	if len(deleted) != 1 {
		t.Fatalf("expected 1 program, got %d", len(deleted))
	}

	if deleted[0].Name != "about" {
		t.Fatalf("expected only %q to be deleted, got %v", "about", deleted)
	}
}

func TestTrackLinkingError_SetsAndClearsLinkError(t *testing.T) {
	ln := linker.New(nil)
	fw := newTestFileWatcher(ln)

	failure := fail.New(nil, "", fail.OriginLink, "undefined component")
	fw.trackLinkingError(failure)

	if ln.LinkError != failure {
		t.Fatalf("expected link error to be stored")
	}

	fw.trackLinkingError(nil)

	if ln.LinkError != nil {
		t.Fatalf("expected link error to be cleared")
	}
}

func TestCountFiles_CountsOnlyTemplateFiles(t *testing.T) {
	dir := t.TempDir()
	useTemplateDir(t, dir)

	writeTemplate(t, filepath.Join(dir, "home.tw"), validTpl)
	writeTemplate(t, filepath.Join(dir, "notes.txt"), "not a template")
	if err := os.Mkdir(filepath.Join(dir, "layouts"), 0o755); err != nil {
		t.Fatalf("failed to create layouts dir: %v", err)
	}
	writeTemplate(t, filepath.Join(dir, "layouts", "main.tw"), validTpl)

	fw := newTestFileWatcher(linker.New(nil))

	if count := fw.countFiles(); count != 2 {
		t.Fatalf("expected 2 template files, got %d", count)
	}
}

func TestTick(t *testing.T) {
	t.Run("adds program for new file", func(t *testing.T) {
		dir := t.TempDir()
		useTemplateDir(t, dir)

		fw := newStartedWatcher(t, linker.New(nil))

		writeTemplate(t, filepath.Join(dir, "home.tw"), validTpl)
		fw.tick()

		programs := fw.linker.Programs
		if len(programs) != 1 || programs[0].Name != "home" {
			t.Fatalf("expected program %q to be added, got %v", "home", programs)
		}
	})

	t.Run("removes program for deleted file", func(t *testing.T) {
		dir := t.TempDir()
		useTemplateDir(t, dir)

		path := filepath.Join(dir, "home.tw")
		writeTemplate(t, path, validTpl)

		fw := newStartedWatcher(t, linker.New(nil))
		fw.tick()

		if err := os.Remove(path); err != nil {
			t.Fatalf("failed to remove template: %v", err)
		}
		fw.tick()

		if programs := fw.linker.Programs; len(programs) != 0 {
			t.Fatalf("expected program to be removed, got %v", programs)
		}
	})

	t.Run("drops last valid source of deleted file", func(t *testing.T) {
		dir := t.TempDir()
		useTemplateDir(t, dir)

		path := filepath.Join(dir, "home.tw")
		writeTemplate(t, path, validTpl)

		fw := newStartedWatcher(t, linker.New(nil))
		fw.tick()

		if err := os.Remove(path); err != nil {
			t.Fatalf("failed to remove template: %v", err)
		}

		fw.tick()

		// A new file with the same name that doesn't parse must not fall back
		// to the source of the deleted file
		writeTemplate(t, path, brokenTpl)
		fw.tick()

		if programs := fw.linker.Programs; len(programs) != 0 {
			t.Fatalf("broken program should not be in linker, got %v", programs)
		}
	})
}
