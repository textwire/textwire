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
	lastError string
}

// newFileWatcher creates a new file watcher instance.
func newFileWatcher(oldLinker *linker.NodeLinker) *fileWatcher {
	return &fileWatcher{
		linker: oldLinker,
		logger: NewWatcherLogger(),
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
			fw.checkForChanges()
		}
	}()
}

// checkForChanges runs on every tick. It re-scans the template dir when files
// are added or deleted, reparses changed files on their own so one broken
// template doesn't block the rest, then relinks everything so link errors are
// tracked.
func (fw *fileWatcher) checkForChanges() {
	if fw.filesAddedOrDeleted() {
		fw.reScanTemplateDir()
	}

	for _, f := range fw.files {
		fw.reparseIfModified(f)
	}

	fw.relinkAllPrograms()
}

// filesAddedOrDeleted reports whether files were added or deleted on disk.
func (fw *fileWatcher) filesAddedOrDeleted() bool {
	return fw.countFiles() != len(fw.files) || slices.ContainsFunc(fw.files, fw.fileWasDeleted)
}

// reScanTemplateDir re-locates template files and updates tracking of new and
// deleted files.
func (fw *fileWatcher) reScanTemplateDir() {
	fw.logger.Info("file count changed, updating...")
	oldFiles := fw.files

	files, err := locateFiles()
	if err != nil {
		fw.logger.Error("error locating files " + err.Error())
		return
	}

	fw.files = files
	fw.markNewFilesForParsing(oldFiles)
	fw.removeDeletedPrograms(oldFiles)
}

func (fw *fileWatcher) fileWasDeleted(f *file.SourceFile) bool {
	_, err := os.Stat(f.Abs)
	return os.IsNotExist(err)
}

// reparseIfModified reparses a file if it has been modified since last check.
func (fw *fileWatcher) reparseIfModified(f *file.SourceFile) {
	modTime := fw.getFileModTime(f)
	if modTime.IsZero() || !modTime.After(f.ModTime) {
		return
	}

	fw.logger.Info("updated " + f.Rel)
	f.ModTime = modTime

	prog, failure, parseErr := parseFile(f)
	if parseErr != nil {
		fw.logger.Error(parseErr.Error())
		fw.removeProgramByName(f.Name)
		return
	}

	if failure != nil {
		fw.logger.Error(failure.Error().Error())
		return
	}

	fw.upsertProgram(prog)
}

// upsertProgram replaces the program with the same name, or appends it
// when none exists. The slice is copied so renders holding it stay consistent.
func (fw *fileWatcher) upsertProgram(prog *ast.Program) {
	fw.withLock(func() {
		programs := slices.Clone(fw.linker.Programs)

		i := slices.IndexFunc(programs, func(p *ast.Program) bool {
			return p.Name == prog.Name
		})

		if i == -1 {
			programs = append(programs, prog)
		} else {
			programs[i] = prog
		}

		fw.linker.Programs = programs
	})
}

// relinkAllPrograms parses the current files into fresh programs, links them, and
// swaps the whole slice under the write lock. Existing programs are never
// mutated, so renders that hold older programs stay consistent.
func (fw *fileWatcher) relinkAllPrograms() {
	if fw.files == nil {
		return
	}

	programs, parseFailure := parseFiles(fw.files)
	if parseFailure != nil {
		fw.logger.Error(parseFailure.String())
		return
	}

	ln := linker.New(programs)
	failure := ln.LinkNodes()

	fw.withLock(func() {
		fw.linker.Programs = ln.Programs
		fw.trackLinkingError(failure)
	})
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

// removeDeletedPrograms removes programs of files that existed in oldFiles but
// are no longer in fw.files.
func (fw *fileWatcher) removeDeletedPrograms(oldFiles []*file.SourceFile) {
	deletedNames := make(map[string]bool)
	for _, f := range fw.findDeletedFiles(oldFiles) {
		fw.logger.Info("removed " + f.Rel)
		deletedNames[f.Name] = true
	}

	fw.filterPrograms(func(p *ast.Program) bool {
		return !deletedNames[p.Name]
	})
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

// getFileModTime retrieves the last modification time of a file.
// Returns zero time if the file does not exist or cannot be accessed.
func (fw *fileWatcher) getFileModTime(f *file.SourceFile) time.Time {
	info, err := os.Stat(f.Abs)
	if err != nil {
		if !os.IsNotExist(err) {
			fw.logger.Error("failed to stat file " + f.Abs)
		}
		return time.Time{}
	}

	return info.ModTime()
}

// markNewFilesForParsing sets ModTime to zero for newly created files to force reparse.
func (fw *fileWatcher) markNewFilesForParsing(oldFiles []*file.SourceFile) {
	oldSet := makeFileNameSet(oldFiles)
	for _, f := range fw.files {
		if !oldSet[f.Name] {
			f.ModTime = time.Time{}
		}
	}
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
