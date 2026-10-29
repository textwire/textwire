package textwire

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"codeberg.org/textwire/textwire/v5/pkg/ast"
	"codeberg.org/textwire/textwire/v5/pkg/fail"
	"codeberg.org/textwire/textwire/v5/pkg/file"
	"codeberg.org/textwire/textwire/v5/pkg/linker"
)

// fileWatcher monitors template files for changes and refreshes parsed AST nodes.
// It is designed for development use only due to performance implications.
type fileWatcher struct {
	linker    *linker.NodeLinker
	logger    *WatcherLogger
	ticker    *time.Ticker
	files     []*file.SourceFile
	sources   map[string]string // last valid source of each template, by name
	lastError string
}

// newFileWatcher creates a new file watcher instance.
func newFileWatcher(oldLinker *linker.NodeLinker) *fileWatcher {
	return &fileWatcher{
		linker:  oldLinker,
		logger:  NewWatcherLogger(),
		sources: make(map[string]string),
	}
}

// Watch starts monitoring files in a background goroutine.
// It detects file creation, deletion, and modifications, then reparses and
// relinks accordingly.
func (fw *fileWatcher) Watch() {
	if userConf.UsesFS() {
		fw.logger.Fatal("cannot use config.FileWatcher when using config.TemplateFS")
	}

	fw.logger.Info("watching files for changes...")

	var err error
	fw.files, err = locateFiles()
	if err != nil {
		fw.logger.Fatal("error locating files " + err.Error())
	}

	fw.ticker = time.NewTicker(userConf.WatcherInterval)

	go func() {
		for range fw.ticker.C {
			fw.tick()
		}
	}()
}

// tick runs on every tick. It rescans the template dir when files
// are added or deleted, then relinks everything so changed templates and link
// errors are tracked.
func (fw *fileWatcher) tick() {
	if fw.filesAddedOrDeleted() {
		fw.rescanTemplateDir()
	}

	fw.relinkAllPrograms()
}

// filesAddedOrDeleted reports whether files were added or deleted on disk.
func (fw *fileWatcher) filesAddedOrDeleted() bool {
	return fw.countFiles() != len(fw.files) || slices.ContainsFunc(fw.files, fw.fileWasDeleted)
}

// rescanTemplateDir re-locates template files and removes programs of files
// that were deleted.
func (fw *fileWatcher) rescanTemplateDir() {
	fw.logger.Info("file count changed, updating...")
	oldFiles := fw.files

	files, err := locateFiles()
	if err != nil {
		fw.logger.Error("error locating files " + err.Error())
		return
	}

	fw.files = files
	fw.removeDeletedPrograms(oldFiles)
}

func (fw *fileWatcher) fileWasDeleted(f *file.SourceFile) bool {
	_, err := os.Stat(f.Abs)
	return os.IsNotExist(err)
}

// relinkAllPrograms parses the current files into fresh programs, links them, and
// swaps the whole slice under the write lock. Existing programs are never
// mutated, so renders that hold older programs stay consistent.
func (fw *fileWatcher) relinkAllPrograms() {
	if fw.files == nil {
		return
	}

	programs := make([]*ast.Program, 0, len(fw.files))
	for _, f := range fw.files {
		if prog := fw.parseFileOrLastValid(f); prog != nil {
			programs = append(programs, prog)
		}
	}

	ln := linker.New(programs)
	failure := ln.LinkNodes()

	fw.withLock(func() {
		fw.linker.Programs = ln.Programs
		fw.trackLinkingError(failure)
	})
}

// parseFileOrLastValid parses f into a fresh program. When f fails to parse, it
// parses the last valid source of f instead, so one broken template doesn't
// block the rest. Returns nil when f has never parsed.
func (fw *fileWatcher) parseFileOrLastValid(f *file.SourceFile) *ast.Program {
	content, err := f.Content()
	if err != nil {
		fw.logger.Error(err.Error())
		return fw.parseLastValid(f)
	}

	prog, failure := parseContent(f, content)
	if failure != nil {
		fw.logger.Error(failure.Error().Error())
		return fw.parseLastValid(f)
	}

	fw.sources[f.Name] = content
	return prog
}

// parseLastValid parses the last valid source of f into a fresh program.
// Returns nil when f has never parsed.
func (fw *fileWatcher) parseLastValid(f *file.SourceFile) *ast.Program {
	content, ok := fw.sources[f.Name]
	if !ok {
		return nil
	}

	prog, _ := parseContent(f, content)
	return prog
}

// trackLinkingError logs linking errors once and stores them for Template access.
func (fw *fileWatcher) trackLinkingError(failure *fail.Error) {
	fw.linker.LinkError = failure

	if failure == nil {
		if fw.lastError != "" {
			fw.logger.Info("all templates are valid!")
			fw.lastError = ""
		}
		return
	}

	errMsg := failure.Error().Error()
	if errMsg == fw.lastError {
		return
	}

	fw.logger.Error(errMsg)
	fw.lastError = errMsg
}

// removeDeletedPrograms removes programs and cached sources of files that
// existed in oldFiles but are no longer in fw.files.
func (fw *fileWatcher) removeDeletedPrograms(oldFiles []*file.SourceFile) {
	for _, f := range fw.findDeletedFiles(oldFiles) {
		fw.logger.Info("removed " + f.Rel)
		delete(fw.sources, f.Name)
		fw.removeProgramByName(f.Name)
	}
}

// findDeletedFiles returns files from oldFiles that are not in fw.files.
func (fw *fileWatcher) findDeletedFiles(oldFiles []*file.SourceFile) []*file.SourceFile {
	current := makeFileNameSet(fw.files)

	var deleted []*file.SourceFile
	for _, f := range oldFiles {
		if !current[f.Name] {
			deleted = append(deleted, f)
		}
	}

	return deleted
}

// removeProgramByName removes the program with the specified name from the linker.
func (fw *fileWatcher) removeProgramByName(name string) {
	fw.filterPrograms(func(p *ast.Program) bool {
		return p.Name != name
	})
}

// filterPrograms replaces linker programs with those for which keep returns true.
func (fw *fileWatcher) filterPrograms(keep func(*ast.Program) bool) {
	fw.withLock(func() {
		programs := make([]*ast.Program, 0, len(fw.linker.Programs))
		for _, p := range fw.linker.Programs {
			if keep(p) {
				programs = append(programs, p)
			}
		}

		fw.linker.Programs = programs
	})
}

// withLock executes the given function while holding the linker's write lock.
func (fw *fileWatcher) withLock(fn func()) {
	fw.linker.Lock()
	defer fw.linker.Unlock()
	fn()
}

// countFiles counts template files in the templates directory on disk.
func (fw *fileWatcher) countFiles() int {
	count := 0
	err := filepath.WalkDir(
		userConf.TemplateDir,
		func(path string, d os.DirEntry, err error) error {
			if err == nil && !d.IsDir() && strings.HasSuffix(path, userConf.TemplateExt) {
				count++
			}
			return nil
		},
	)

	if err != nil {
		fw.logger.Error("error counting files " + err.Error())
	}

	return count
}
