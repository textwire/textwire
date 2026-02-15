package config

import (
	"io/fs"
	"os"
	"strings"
	"time"
)

// Config holds the configuration settings for Textwire template engine.
type Config struct {
	// TemplateDir specifies the directory containing Textwire template files.
	// Note: If TemplatesFS is provided, TemplateDir is ignored because there
	// are no absolute paths for embeded files.
	// Default: "templates"
	TemplateDir string

	// TemplateFS provides an optional fs.FS filesystem for template access.
	// Use this field to embed templates into your binary using Go's embed
	// package. When provided, TemplateDir is not used for file access.
	// Default: os.DirFS(TemplateDir)
	TemplateFS fs.FS

	// TemplateExt defines the file extension for Textwire template files.
	// Note: Using a different extension may disable syntax highlighting
	// in editors like VSCode when using the Textwire extension.
	// Default: ".tw"
	TemplateExt string

	// ErrorPagePath sets the relative path to a custom error page template.
	// The path is relative to the template directory
	// (TemplateDir or TemplatesFS root).
	// Default: embed/default-error-page.tw (default error page)
	ErrorPagePath string

	// DebugMode enables detailed error reporting in the browser and server
	// logs. When true, error messages with file paths and line numbers are
	// displayed during development.
	// Default: false
	DebugMode bool

	// GlobalData stores shared data accessible across all templates.
	// Access these values in templates using the `global` object
	// (e.g., `global.authUser`). Useful for storing environment variables,
	// configuration, or common data.
	GlobalData map[string]any

	// FileReload watches all of your template files for changes and
	// automatically reparses them when they are modified. This is intended
	// for development use only and should not be enabled in production due to
	// performance implications. It doesn't work with TemplateFS configuration
	// enabled! Disable it to use file reload functionality.
	// Default: false
	FileReload bool

	// FileReloadInterval specifies how often Textwire checks for changes in
	// template files when FileReload is enabled. The higher the interval,
	// the less frequently Textwire checks for file changes, which can reduce
	// CPU usage but may delay updates. Values less than 1 second will be
	// treated as the default (1 second). Adjust this value based on your
	// development needs.
	// Default: time.Second (1 second)
	FileReloadInterval time.Duration

	// usesFS is a flag to determine if user uses TemplateFS or not.
	usesFS bool
}

func New(dir, ext, errPagePath string, debug bool) *Config {
	return &Config{
		TemplateDir:        dir,
		TemplateExt:        ext,
		ErrorPagePath:      errPagePath,
		DebugMode:          debug,
		GlobalData:         map[string]any{},
		FileReloadInterval: time.Second,
	}
}

// UsesFS returns value of usesFS field since that field is private.
func (c *Config) UsesFS() bool {
	return c.usesFS
}

func (c *Config) Configure(opt *Config) {
	if opt == nil {
		return
	}

	if opt.TemplateDir != "" {
		c.TemplateDir = strings.Trim(opt.TemplateDir, "/")
	}

	if opt.TemplateExt != "" {
		c.TemplateExt = opt.TemplateExt
	}

	if opt.TemplateFS == nil {
		c.TemplateFS = os.DirFS(c.TemplateDir)
	} else {
		c.TemplateFS = opt.TemplateFS
	}

	if opt.ErrorPagePath != "" {
		c.ErrorPagePath = opt.ErrorPagePath
	}

	if opt.FileReloadInterval >= time.Second {
		c.FileReloadInterval = opt.FileReloadInterval
	}

	if opt.GlobalData != nil {
		c.GlobalData = opt.GlobalData
	}

	c.FileReload = opt.FileReload
	c.DebugMode = opt.DebugMode
	c.usesFS = opt.TemplateFS != nil
}
