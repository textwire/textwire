package textwire

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/textwire/textwire/v3/pkg/ast"
	"github.com/textwire/textwire/v3/pkg/fail"
	"github.com/textwire/textwire/v3/pkg/file"
	"github.com/textwire/textwire/v3/pkg/linker"
)

// fileWatcher is responsible for monitoring changes in template files and
// refreshing parsed AST nodes for each of this file. It's meant to be only
// for development purposes and should not be used in production due to
// potential performance implications.
type fileWatcher struct {
	ticker    *time.Ticker
	oldLinker *linker.NodeLinker
}

// Watch watches for changes in the provided files and refreshes the linked
// programs if any changes are detected. It runs in a separate goroutine and
// checks for file modifications at regular intervals defined by the ticker.
// If a file change is detected, it parses the file again and updates the
// linked programs accordingly. If any errors occur during this process, they
// are logged and the application is terminated. Note that this function
// cannot be used if the user is using TemplateFS, as it relies on direct file
// access to monitor changes.
func (fw *fileWatcher) Watch(files []*file.SourceFile) error {
	if userConf.UsesFS() {
		return errors.New("cannot use config.RefreshFiles when using config.TemplateFS")
	}

	fw.ticker = time.NewTicker(userConf.RefreshInterval)

	go func() {
		for range fw.ticker.C {
			for i := range files {
				if failure := fw.updateFileIfModified(files[i]); failure != nil {
					log.Fatalln(failure)
				}
			}
		}
	}()

	return nil
}

func (fw *fileWatcher) updateFileIfModified(f *file.SourceFile) *fail.Error {
	modTime, err := fw.fetchModTime(f)
	if err != nil {
		log.Fatalln(err)
	}

	if f.ModTime.Equal(modTime) {
		return nil
	}

	log.Printf("Refreshing %s", f.Rel)

	f.ModTime = modTime

	prog, failure, parseErr := parseFile(f)
	if parseErr != nil {
		return fail.FromError(parseErr, 0, f.Abs, "template")
	}

	if failure != nil {
		return failure
	}

	if failure := fw.refreshPrograms(prog); failure != nil {
		return failure
	}

	return nil
}

func (fw *fileWatcher) refreshPrograms(prog *ast.Program) *fail.Error {
	fw.oldLinker.Lock()
	defer fw.oldLinker.Unlock()

	for i := range fw.oldLinker.Programs {
		if fw.oldLinker.Programs[i].Name == prog.Name {
			fw.oldLinker.Programs[i] = prog
		}
	}

	ln := linker.New(fw.oldLinker.Programs)
	if failure := ln.LinkNodes(); failure != nil {
		return failure
	}

	fw.oldLinker.Programs = ln.Programs

	return nil
}

// fetchModTime fetches the file's info and retrieves last modified date.
func (fw *fileWatcher) fetchModTime(f *file.SourceFile) (time.Time, error) {
	var fileInfo os.FileInfo
	fileInfo, err := os.Stat(f.Abs)
	if err != nil {
		return time.Now(), err
	}

	return fileInfo.ModTime(), nil
}
