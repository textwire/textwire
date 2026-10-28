package textwire

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"codeberg.org/textwire/textwire/v5/config"
	"codeberg.org/textwire/textwire/v5/pkg/ast"
	"codeberg.org/textwire/textwire/v5/pkg/fail"
	"codeberg.org/textwire/textwire/v5/pkg/file"
	"codeberg.org/textwire/textwire/v5/pkg/linker"
)

const (
	validTemplate  = "Hello, {{ name }}"
	brokenTemplate = "@if(true) no end" // the @if is never closed
)

func TestReparseIfModified(t *testing.T) {
	t.Run("replaces program on valid change", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "home.tw")
		writeTemplate(t, path, validTemplate)

		f := file.New("home", path, path, config.New())
		old := mustParseFile(t, f)

		fw := newFileWatcher(linker.New([]*ast.Program{old}))
		fw.files = []*file.SourceFile{f}

		writeTemplate(t, path, validTemplate+"!")
		setModTimeAfter(t, path, f.ModTime)

		fw.reparseIfModified(f)

		programs := fw.linker.Programs
		if len(programs) != 1 || programs[0] == old {
			t.Fatal("expected program 'home' to be replaced")
		}
	})

	t.Run("skips unchanged file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "home.tw")
		writeTemplate(t, path, validTemplate)

		f := file.New("home", path, path, config.New())
		old := mustParseFile(t, f)
		f.ModTime = statModTime(t, path)

		ln := linker.New([]*ast.Program{old})

		fw := newFileWatcher(ln)
		fw.files = []*file.SourceFile{f}

		fw.reparseIfModified(f)

		programs := fw.linker.Programs
		if len(programs) != 1 || programs[0] != old {
			t.Fatalf("expected unchanged program to be kept")
		}
	})

	t.Run("keeps old program on syntax error", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "home.tw")
		writeTemplate(t, path, validTemplate)

		f := file.New("home", path, path, config.New())
		prog := mustParseFile(t, f)

		fw := newFileWatcher(linker.New([]*ast.Program{prog}))
		fw.files = []*file.SourceFile{f}

		// The file changes on disk into something that no longer parses
		writeTemplate(t, path, brokenTemplate)
		setModTimeAfter(t, path, f.ModTime)

		fw.reparseIfModified(f)

		// The watcher survives and keeps serving the last valid program
		programs := fw.linker.Programs
		if len(programs) != 1 {
			t.Fatalf("expected to keep the old program, got %d programs", len(programs))
		}

		if programs[0].Name != "home" {
			t.Fatalf("expected to keep program %q, got %q", "home", programs[0].Name)
		}
	})
}

func TestUpsertProgram(t *testing.T) {
	t.Run("appends new program", func(t *testing.T) {
		fw := newFileWatcher(linker.New(nil))
		fw.upsertProgram(&ast.Program{Name: "home"})

		programs := fw.linker.Programs
		if len(programs) != 1 {
			t.Fatalf("expected 1 program, got %d", len(programs))
		}

		if programs[0].Name != "home" {
			t.Fatalf("expected program %q to be added, got %v", "home", programs)
		}
	})

	t.Run("replaces program with same name", func(t *testing.T) {
		old := &ast.Program{Name: "home"}
		ln := linker.New([]*ast.Program{old})
		fw := newFileWatcher(ln)

		updated := &ast.Program{Name: "home"}
		fw.upsertProgram(updated)

		programs := fw.linker.Programs
		if len(programs) != 1 {
			t.Fatalf("expected 1 program, got %d", len(programs))
		}

		if programs[0] != updated {
			t.Fatal("expected program 'home' to be replaced")
		}
	})
}

func TestRemoveProgramByName(t *testing.T) {
	t.Run("removes only matching program", func(t *testing.T) {
		ln := linker.New([]*ast.Program{
			{Name: "home"},
			{Name: "about"},
		})

		fw := newFileWatcher(ln)
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

		fw := newFileWatcher(ln)
		fw.removeProgramByName("contact")

		programs := fw.linker.Programs
		if len(programs) != 2 || programs[0].Name != "home" || programs[1].Name != "about" {
			t.Fatalf("expected all programs to remain, got %v", programs)
		}
	})
}

func TestFindDeletedFiles_ReturnsFilesMissingFromCurrentList(t *testing.T) {
	fw := newFileWatcher(linker.New(nil))
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

func TestMarkNewFilesForParsing_ZeroesModTimeOnlyForNewFiles(t *testing.T) {
	lastCheck := time.Now()
	oldFile := &file.SourceFile{Name: "home", ModTime: lastCheck}
	newFile := &file.SourceFile{Name: "about", ModTime: lastCheck}

	fw := newFileWatcher(linker.New(nil))
	fw.files = []*file.SourceFile{oldFile, newFile}

	fw.markNewFilesForParsing([]*file.SourceFile{oldFile})

	if !oldFile.ModTime.Equal(lastCheck) {
		t.Fatalf("expected existing file ModTime to be kept")
	}

	if !newFile.ModTime.IsZero() {
		t.Fatalf("expected new file ModTime to be zeroed, got %v", newFile.ModTime)
	}
}

func TestGetFileModTime(t *testing.T) {
	t.Run("returns zero for missing file", func(t *testing.T) {
		fw := newFileWatcher(linker.New(nil))
		f := &file.SourceFile{Abs: filepath.Join(t.TempDir(), "missing.tw")}

		if modTime := fw.getFileModTime(f); !modTime.IsZero() {
			t.Fatalf("expected zero time for missing file, got %v", modTime)
		}
	})

	t.Run("returns file mod time", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "home.tw")
		writeTemplate(t, path, validTemplate)

		fw := newFileWatcher(linker.New(nil))
		f := &file.SourceFile{Abs: path}

		want := statModTime(t, path)
		if modTime := fw.getFileModTime(f); !modTime.Equal(want) {
			t.Fatalf("expected mod time %v, got %v", want, modTime)
		}
	})
}

func TestTrackLinkingError_SetsAndClearsLinkError(t *testing.T) {
	ln := linker.New(nil)
	fw := newFileWatcher(ln)

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

	writeTemplate(t, filepath.Join(dir, "home.tw"), validTemplate)
	writeTemplate(t, filepath.Join(dir, "notes.txt"), "not a template")
	if err := os.Mkdir(filepath.Join(dir, "layouts"), 0o755); err != nil {
		t.Fatalf("failed to create layouts dir: %v", err)
	}
	writeTemplate(t, filepath.Join(dir, "layouts", "main.tw"), validTemplate)

	fw := newFileWatcher(linker.New(nil))

	if count := fw.countFiles(); count != 2 {
		t.Fatalf("expected 2 template files, got %d", count)
	}
}

func TestRelinkAllPrograms_KeepsOldProgramsOnParseFailure(t *testing.T) {
	dir := t.TempDir()
	useTemplateDir(t, dir)

	path := filepath.Join(dir, "home.tw")
	writeTemplate(t, path, brokenTemplate)

	old := &ast.Program{Name: "home"}
	fw := newFileWatcher(linker.New([]*ast.Program{old}))
	fw.files = []*file.SourceFile{file.New("home", path, path, userConf)}

	fw.relinkAllPrograms()

	programs := fw.linker.Programs
	if len(programs) != 1 || programs[0] != old {
		t.Fatalf("expected old program to be kept on parse failure")
	}
}

func TestCheckForChanges(t *testing.T) {
	t.Run("adds program for new file", func(t *testing.T) {
		dir := t.TempDir()
		useTemplateDir(t, dir)

		fw := newStartedWatcher(t, linker.New(nil))

		writeTemplate(t, filepath.Join(dir, "home.tw"), validTemplate)
		fw.checkForChanges()

		programs := fw.linker.Programs
		if len(programs) != 1 || programs[0].Name != "home" {
			t.Fatalf("expected program %q to be added, got %v", "home", programs)
		}
	})

	t.Run("removes program for deleted file", func(t *testing.T) {
		dir := t.TempDir()
		useTemplateDir(t, dir)

		path := filepath.Join(dir, "home.tw")
		writeTemplate(t, path, validTemplate)

		fw := newStartedWatcher(t, linker.New(nil))
		fw.checkForChanges()

		if err := os.Remove(path); err != nil {
			t.Fatalf("failed to remove template: %v", err)
		}
		fw.checkForChanges()

		if programs := fw.linker.Programs; len(programs) != 0 {
			t.Fatalf("expected program to be removed, got %v", programs)
		}
	})
}

// newStartedWatcher mirrors Watch setup without starting the ticker, so
// tests can drive checkForChanges directly.
func newStartedWatcher(t *testing.T, ln *linker.NodeLinker) *fileWatcher {
	t.Helper()

	fw := newFileWatcher(ln)

	files, err := locateFiles()
	if err != nil {
		t.Fatalf("failed to locate files: %v", err)
	}

	fw.files = files

	return fw
}

func useTemplateDir(t *testing.T, dir string) {
	t.Helper()

	prev := *userConf
	t.Cleanup(func() { *userConf = prev })

	userConf.Configure(&config.Config{TemplateDir: dir})
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

func statModTime(t *testing.T, path string) time.Time {
	t.Helper()

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("failed to stat template: %v", err)
	}

	return info.ModTime()
}

func setModTimeAfter(t *testing.T, path string, lastCheck time.Time) {
	t.Helper()

	modTime := lastCheck.Add(time.Second)
	if err := os.Chtimes(path, modTime, modTime); err != nil {
		t.Fatalf("failed to change template mod time: %v", err)
	}
}
