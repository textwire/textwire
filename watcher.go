package textwire

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/textwire/textwire/v3/pkg/ast"
	"github.com/textwire/textwire/v3/pkg/file"
	"github.com/textwire/textwire/v3/pkg/linker"
)

// fileWatcher is responsible for monitoring changes in template files and
// refreshing parsed AST nodes for each of this file. It's meant to be only
// for development purposes and should not be used in production due to
// potential performance implications.
type fileWatcher struct {
	oldLinker *linker.NodeLinker
	ticker    *time.Ticker
	files     []*file.SourceFile
	fileCount int
	lastError string
}

func newFileWatcher(oldLinker *linker.NodeLinker) *fileWatcher {
	return &fileWatcher{
		oldLinker: oldLinker,
		files:     nil,
		fileCount: 0,
	}
}

// Watch watches for changes in the provided directory and refreshes the linked
// programs if any changes are detected. It runs in a separate goroutine and
// checks for file modifications at regular intervals defined by the ticker.
// If a file change is detected, it parses the file again and updates the
// linked programs accordingly. If any errors occur during this process, they
// are logged. Note that this function cannot be used if the user is using
// TemplateFS, as it relies on direct file access to monitor changes.
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

	fw.fileCount = fw.countFiles()
	fw.ticker = time.NewTicker(userConf.WatcherInterval)

	go func() {
		for range fw.ticker.C {
			currentCount := fw.countFiles()
			if currentCount != fw.fileCount {
				fw.handleFileCountChange()
				fw.fileCount = currentCount
			}

			for i := range fw.files {
				fw.updateFileIfModified(fw.files[i])
			}

			fw.relinkAll()
		}
	}()
}

func (fw *fileWatcher) handleFileCountChange() {
	fw.info("File count changed, re-locating files...")
	oldFiles := fw.files

	files, err := locateFiles()
	if err != nil {
		fw.fatal("Error locating files: " + err.Error())
	}
	fw.files = files

	fw.markNewFiles(oldFiles)
	fw.removeDeletedPrograms(oldFiles)
}

func (fw *fileWatcher) updateFileIfModified(f *file.SourceFile) {
	modTime := fw.fetchModTime(f)
	if modTime.IsZero() || !modTime.After(f.ModTime) {
		return
	}

	fw.info("Refreshing file: " + f.Rel)
	f.ModTime = modTime

	prog, failure, parseErr := parseFile(f)
	if parseErr != nil {
		fw.info(parseErr.Error())
		fw.removePrograms(f.Name)
		return
	}

	if failure != nil {
		fw.info(failure.Error().Error())
	}

	fw.refreshPrograms(prog)
}

func (fw *fileWatcher) refreshPrograms(prog *ast.Program) {
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

func (fw *fileWatcher) relinkAll() {
	fw.withLock(func() {
		ln := linker.New(fw.oldLinker.Programs)
		failure := ln.LinkNodes()

		if failure != nil {
			errMsg := failure.Error().Error()
			if errMsg != fw.lastError {
				fw.info(errMsg)
				fw.lastError = errMsg
			}
			fw.oldLinker.LinkError = failure
		} else {
			if fw.lastError != "" {
				fw.info("Errors resolved")
				fw.lastError = ""
			}
			fw.oldLinker.LinkError = nil
		}

		fw.oldLinker.Programs = ln.Programs
	})
}

func (fw *fileWatcher) removeDeletedPrograms(oldFiles []*file.SourceFile) {
	oldMap := make(map[string]bool, len(oldFiles))
	for _, f := range oldFiles {
		oldMap[f.Name] = true
	}

	fw.withLock(func() {
		newProgs := fw.oldLinker.Programs[:0]
		for _, prog := range fw.oldLinker.Programs {
			if !oldMap[prog.Name] || fw.fileExists(prog.Name) {
				newProgs = append(newProgs, prog)
			}
		}
		fw.oldLinker.Programs = newProgs
	})
}

func (fw *fileWatcher) removePrograms(names ...string) {
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

func (fw *fileWatcher) fileExists(name string) bool {
	for _, f := range fw.files {
		if f.Name == name {
			return true
		}
	}
	return false
}

func (fw *fileWatcher) withLock(fn func()) {
	fw.oldLinker.Lock()
	defer fw.oldLinker.Unlock()
	fn()
}

func (fw *fileWatcher) fetchModTime(f *file.SourceFile) time.Time {
	info, err := os.Stat(f.Abs)
	if err != nil {
		fw.info("Failed to stat file: " + f.Abs)
		return time.Time{}
	}
	return info.ModTime()
}

func (fw *fileWatcher) markNewFiles(oldFiles []*file.SourceFile) {
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

func (fw *fileWatcher) countFiles() int {
	cmd := exec.Command("find", userConf.TemplateDir, "-name", "*"+userConf.TemplateExt)
	output, err := cmd.Output()
	if err != nil {
		return 0
	}
	return strings.Count(string(output), "\n")
}

func (fw *fileWatcher) info(text string) {
	fmt.Printf("[Textwire File Watch]: %s\n", text)
}

func (fw *fileWatcher) fatal(text string) {
	fw.info(text)
	os.Exit(1)
}
