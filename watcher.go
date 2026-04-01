package textwire

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/textwire/textwire/v4/pkg/ast"
	"github.com/textwire/textwire/v4/pkg/fail"
	"github.com/textwire/textwire/v4/pkg/file"
	"github.com/textwire/textwire/v4/pkg/linker"
)

// fileWatcher monitors template files for changes and refreshes parsed AST nodes.
// It is designed for development use only due to performance implications.
type fileWatcher struct {
	linker    *linker.NodeLinker
	logger    *WatcherLogger
	ticker    *time.Ticker
	files     []*file.SourceFile
	fileCount int
	lastError string
}

// newFileWatcher creates a new file watcher instance.
func newFileWatcher(oldLinker *linker.NodeLinker) *fileWatcher {
	return &fileWatcher{
		linker:    oldLinker,
		logger:    NewWatcherLogger(),
		files:     nil,
		fileCount: 0,
	}
}

// Watch starts monitoring files in a background goroutine.
// It detects file creation, deletion, and modifications, then reparses and relinks accordingly.
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

	fw.fileCount = fw.countFiles()
	fw.ticker = time.NewTicker(userConf.WatcherInterval)

	go func() {
		for range fw.ticker.C {
			fw.tick()
		}
	}()
}

func (fw *fileWatcher) tick() {
	fw.detectAndHandleFileChanges()
	fw.processModifiedFiles()
	fw.relinkPrograms()
}

func (fw *fileWatcher) detectAndHandleFileChanges() {
	currentCount := fw.countFiles()
	if currentCount == fw.fileCount {
		return
	}

	fw.handleNewOrDeletedFiles()
	fw.fileCount = currentCount
}

func (fw *fileWatcher) processModifiedFiles() {
	if fw.removeDeletedFiles() > 0 {
		fw.handleNewOrDeletedFiles()
		fw.fileCount = fw.countFiles()
		return
	}

	for _, f := range fw.files {
		fw.updateFileIfModified(f)
	}
}

func (fw *fileWatcher) removeDeletedFiles() int {
	validFiles := fw.files[:0]
	deletedCount := 0

	for _, f := range fw.files {
		if fw.fileWasDeleted(f) {
			fw.handleDeletedFile(f)
			deletedCount++
			continue
		}
		validFiles = append(validFiles, f)
	}

	fw.files = validFiles
	return deletedCount
}

func (fw *fileWatcher) fileWasDeleted(f *file.SourceFile) bool {
	_, err := os.Stat(f.Abs)
	return os.IsNotExist(err)
}

func (fw *fileWatcher) handleDeletedFile(f *file.SourceFile) {
	fw.logger.Info("removed " + f.Rel)
	fw.removeProgramByName(f.Name)
}

// handleNewOrDeletedFiles re-locates files and updates tracking when file
// count changes.
func (fw *fileWatcher) handleNewOrDeletedFiles() {
	fw.logger.Info("file count changed, updating...")
	oldFiles := fw.files

	files, err := locateFiles()
	if err != nil {
		fw.logger.Fatal("error locating files " + err.Error())
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
	}

	fw.updateOrAddProgram(prog)
}

// updateOrAddProgram updates an existing program or adds a new one to the linker.
func (fw *fileWatcher) updateOrAddProgram(prog *ast.Program) {
	fw.withLock(func() {
		for i := range fw.linker.Programs {
			if fw.linker.Programs[i].Name == prog.Name {
				fw.linker.Programs[i] = prog
				return
			}
		}
		fw.linker.Programs = append(fw.linker.Programs, prog)
	})
}

// relinkPrograms links all programs together and tracks any linking errors.
func (fw *fileWatcher) relinkPrograms() {
	fw.withLock(func() {
		ln := linker.New(fw.linker.Programs)
		failure := ln.LinkNodes()
		fw.linker.Programs = ln.Programs

		fw.trackLinkingError(failure)
	})
}

// trackLinkingError logs linking errors once and stores them for Template access.
func (fw *fileWatcher) trackLinkingError(failure *fail.Error) {
	if failure == nil {
		if fw.lastError != "" {
			fw.logger.Info("all templates are valid!")
			fw.lastError = ""
		}
		fw.linker.LinkError = nil
		return
	}

	errMsg := failure.Error().Error()
	if errMsg != fw.lastError {
		fw.logger.Error(errMsg)
		fw.lastError = errMsg
	}

	fw.linker.LinkError = failure
}

// cleanupDeletedPrograms removes programs from the linker that no longer exist on disk.
func (fw *fileWatcher) cleanupDeletedPrograms(oldFiles []*file.SourceFile) {
	deletedFiles := fw.findDeletedFiles(oldFiles)

	fw.withLock(func() {
		newProgs := fw.linker.Programs[:0]
		for _, prog := range fw.linker.Programs {
			if !deletedFiles[prog.Name] {
				newProgs = append(newProgs, prog)
			}
		}
		fw.linker.Programs = newProgs
	})
}

// findDeletedFiles returns a map of file names that existed in oldFiles but no longer
// exist in fw.files.
func (fw *fileWatcher) findDeletedFiles(oldFiles []*file.SourceFile) map[string]bool {
	oldSet := makeFileNameSet(oldFiles)
	deleted := make(map[string]bool)
	for name := range oldSet {
		if !fw.fileExists(name) {
			deleted[name] = true
		}
	}

	return deleted
}

// removeProgramsByName removes programs with the specified names from the linker.
func (fw *fileWatcher) removeProgramByName(name string) {
	fw.withLock(func() {
		newProgs := fw.linker.Programs[:0]
		for _, prog := range fw.linker.Programs {
			if prog.Name != name {
				newProgs = append(newProgs, prog)
			}
		}

		fw.linker.Programs = newProgs
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
	fw.linker.Lock()
	defer fw.linker.Unlock()
	fn()
}

// getFileModTime retrieves the last modification time of a file.
// Returns zero time if the file does not exist or cannot be accessed.
func (fw *fileWatcher) getFileModTime(f *file.SourceFile) time.Time {
	if _, err := os.Stat(f.Abs); os.IsNotExist(err) {
		return time.Time{}
	}

	info, err := os.Stat(f.Abs)
	if err != nil {
		fw.logger.Error("failed to stat file " + f.Abs)
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
		fw.logger.Error("error couting files " + err.Error())
	}

	return count
}
