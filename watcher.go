package textwire

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/textwire/textwire/v3/pkg/ast"
	"github.com/textwire/textwire/v3/pkg/fail"
	"github.com/textwire/textwire/v3/pkg/file"
	"github.com/textwire/textwire/v3/pkg/linker"
)

// fileWatcher monitors template files for changes and refreshes parsed AST nodes.
// It is designed for development use only due to performance implications.
type fileWatcher struct {
	oldLinker *linker.NodeLinker
	ticker    *time.Ticker
	files     []*file.SourceFile
	fileCount int
	lastError string
}

// newFileWatcher creates a new file watcher instance.
func newFileWatcher(oldLinker *linker.NodeLinker) *fileWatcher {
	return &fileWatcher{
		oldLinker: oldLinker,
		files:     nil,
		fileCount: 0,
	}
}

// Watch starts monitoring files in a background goroutine.
// It detects file creation, deletion, and modifications, then reparses and relinks accordingly.
func (fw *fileWatcher) Watch() {
	if userConf.UsesFS() {
		fw.fatal("Cannot use config.FileWatcher when using config.TemplateFS")
	}

	fw.info("Watching files for changes...")

	var err error
	fw.files, err = locateFiles()
	if err != nil {
		fw.fatal("Error locating files: " + err.Error())
	}

	fw.fileCount = fw.countTemplateFiles()
	fw.ticker = time.NewTicker(userConf.WatcherInterval)

	go func() {
		for range fw.ticker.C {
			fw.tick()
		}
	}()
}

func (fw *fileWatcher) tick() {
	currentCount := fw.countTemplateFiles()
	if currentCount != fw.fileCount {
		fw.handleNewOrDeletedFiles()
		fw.fileCount = currentCount
	}

	for i := range fw.files {
		fw.updateFileIfModified(fw.files[i])
	}

	fw.relinkPrograms()
}

// handleNewOrDeletedFiles re-locates files and updates tracking when file count changes.
func (fw *fileWatcher) handleNewOrDeletedFiles() {
	fw.info("File count changed, re-locating files...")
	oldFiles := fw.files

	files, err := locateFiles()
	if err != nil {
		fw.fatal("Error locating files: " + err.Error())
	}
	fw.files = files

	fw.markNewFilesForParsing(oldFiles)
	fw.cleanupDeletedPrograms(oldFiles)
}

// updateFileIfModified reparses a file if it has been modified since last check.
func (fw *fileWatcher) updateFileIfModified(f *file.SourceFile) {
	modTime := fw.getFileModTime(f)
	if modTime.IsZero() || !modTime.After(f.ModTime) {
		return
	}

	fw.info("Refreshing file: " + f.Rel)
	f.ModTime = modTime

	prog, failure, parseErr := parseFile(f)
	if parseErr != nil {
		fw.info(parseErr.Error())
		fw.removeProgramsByName(f.Name)
		return
	}

	if failure != nil {
		fw.info(failure.Error().Error())
	}

	fw.updateOrAddProgram(prog)
}

// updateOrAddProgram updates an existing program or adds a new one to the linker.
func (fw *fileWatcher) updateOrAddProgram(prog *ast.Program) {
	fw.withLock(func() {
		for i := range fw.oldLinker.Programs {
			if fw.oldLinker.Programs[i].Name == prog.Name {
				fw.oldLinker.Programs[i] = prog
				return
			}
		}
		fw.oldLinker.Programs = append(fw.oldLinker.Programs, prog)
	})
}

// relinkPrograms links all programs together and tracks any linking errors.
func (fw *fileWatcher) relinkPrograms() {
	fw.withLock(func() {
		ln := linker.New(fw.oldLinker.Programs)
		failure := ln.LinkNodes()
		fw.oldLinker.Programs = ln.Programs

		fw.trackLinkingError(failure)
	})
}

// trackLinkingError logs linking errors once and stores them for Template access.
func (fw *fileWatcher) trackLinkingError(failure *fail.Error) {
	if failure == nil {
		if fw.lastError != "" {
			fw.info("Errors resolved")
			fw.lastError = ""
		}
		fw.oldLinker.LinkError = nil
		return
	}

	errMsg := failure.Error().Error()
	if errMsg != fw.lastError {
		fw.info(errMsg)
		fw.lastError = errMsg
	}

	fw.oldLinker.LinkError = failure
}

func (fw *fileWatcher) cleanupDeletedPrograms(oldFiles []*file.SourceFile) {
	deletedFiles := fw.findDeletedFiles(oldFiles)

	fw.withLock(func() {
		newProgs := fw.oldLinker.Programs[:0]
		for _, prog := range fw.oldLinker.Programs {
			if !deletedFiles[prog.Name] {
				newProgs = append(newProgs, prog)
			}
		}
		fw.oldLinker.Programs = newProgs
	})
}

func (fw *fileWatcher) findDeletedFiles(oldFiles []*file.SourceFile) map[string]bool {
	deleted := make(map[string]bool)
	for _, f := range oldFiles {
		if !fw.fileExists(f.Name) {
			deleted[f.Name] = true
		}
	}
	return deleted
}

// removeProgramsByName removes programs with the specified names from the linker.
func (fw *fileWatcher) removeProgramsByName(names ...string) {
	removeSet := make(map[string]bool, len(names))
	for _, name := range names {
		removeSet[name] = true
	}

	fw.withLock(func() {
		newProgs := fw.oldLinker.Programs[:0]
		for _, prog := range fw.oldLinker.Programs {
			if !removeSet[prog.Name] {
				newProgs = append(newProgs, prog)
			}
		}
		fw.oldLinker.Programs = newProgs
	})
}

// fileExists checks if a file with the given name is in the current file list.
func (fw *fileWatcher) fileExists(name string) bool {
	for _, f := range fw.files {
		if f.Name == name {
			return true
		}
	}
	return false
}

// withLock executes the given function while holding the linker's write lock.
func (fw *fileWatcher) withLock(fn func()) {
	fw.oldLinker.Lock()
	defer fw.oldLinker.Unlock()
	fn()
}

// getFileModTime retrieves the last modification time of a file.
// Returns zero time if the file cannot be accessed.
func (fw *fileWatcher) getFileModTime(f *file.SourceFile) time.Time {
	info, err := os.Stat(f.Abs)
	if err != nil {
		fw.info("Failed to stat file: " + f.Abs)
		return time.Time{}
	}

	return info.ModTime()
}

// markNewFilesForParsing sets ModTime to zero for newly created files to force reparse.
func (fw *fileWatcher) markNewFilesForParsing(oldFiles []*file.SourceFile) {
	oldMap := make(map[string]bool, len(oldFiles))
	for _, f := range oldFiles {
		oldMap[f.Name] = true
	}

	for _, f := range fw.files {
		if !oldMap[f.Name] {
			f.ModTime = time.Time{}
		}
	}
}

func (fw *fileWatcher) countTemplateFiles() int {
	count := 0
	filepath.WalkDir(userConf.TemplateDir, func(path string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() && strings.HasSuffix(path, userConf.TemplateExt) {
			count++
		}
		return nil
	})
	return count
}

func (fw *fileWatcher) info(text string) {
	fmt.Printf("[Textwire File Watch]: %s\n", text)
}

func (fw *fileWatcher) fatal(text string) {
	fw.info(text)
	os.Exit(1)
}
