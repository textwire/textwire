package textwire

import (
	"fmt"
	"os"
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
}

func newFileWatcher(oldLinker *linker.NodeLinker) *fileWatcher {
	return &fileWatcher{oldLinker: oldLinker}
}

// Watch watches for changes in the provided files and refreshes the linked
// programs if any changes are detected. It runs in a separate goroutine and
// checks for file modifications at regular intervals defined by the ticker.
// If a file change is detected, it parses the file again and updates the
// linked programs accordingly. If any errors occur during this process, they
// are logged and the application is terminated. Note that this function
// cannot be used if the user is using TemplateFS, as it relies on direct file
// access to monitor changes.
func (fw *fileWatcher) Watch(files []*file.SourceFile) {
	if userConf.UsesFS() {
		fw.fatal("cannot use config.RefreshFiles when using config.TemplateFS")
	}

	fw.info("Watching files for changes...")

	fw.ticker = time.NewTicker(userConf.WatcherInterval)

	go func() {
		for range fw.ticker.C {
			for i := range files {
				fw.updateFileIfModified(files[i])
			}
		}
	}()
}

func (fw *fileWatcher) updateFileIfModified(f *file.SourceFile) {
	modTime := fw.fetchModTime(f)
	if !modTime.After(f.ModTime) {
		return
	}

	fw.info("Refreshing file: " + f.Rel)

	f.ModTime = modTime

	prog, failure, parseErr := parseFile(f)
	if parseErr != nil {
		fw.info(parseErr.Error())
		return
	}

	if failure != nil {
		fw.fatal(failure.String())
	}

	fw.refreshPrograms(prog)
}

func (fw *fileWatcher) refreshPrograms(prog *ast.Program) {
	fw.oldLinker.Lock()
	defer fw.oldLinker.Unlock()

	for i := range fw.oldLinker.Programs {
		if fw.oldLinker.Programs[i].Name == prog.Name {
			fw.oldLinker.Programs[i] = prog
		}
	}

	ln := linker.New(fw.oldLinker.Programs)
	if failure := ln.LinkNodes(); failure != nil {
		fw.fatal(failure.String())
	}

	fw.oldLinker.Programs = ln.Programs
}

// fetchModTime fetches the file's info and retrieves last modified date.
func (fw *fileWatcher) fetchModTime(f *file.SourceFile) time.Time {
	fileInfo, err := os.Stat(f.Abs)
	if err != nil {
		fw.fatal(err.Error())
	}

	return fileInfo.ModTime()
}

func (fw *fileWatcher) info(text string) {
	fmt.Printf("[Textwire File Watch]: %s\n", text)
}

func (fw *fileWatcher) fatal(text string) {
	fw.info(text)
	os.Exit(1)
}
