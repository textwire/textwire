package config

import "io/fs"

// Config holds the configuration settings for Textwire template engine.
type Config struct {
	// TemplateDir specifies the directory containing Textwire template files.
	// Default: "templates"
	// Note: If TemplatesFS is provided, it will be used for templates path
	// instead of TemplateDir.
	TemplateDir string

	// TemplatesFS provides an optional fs.FS filesystem for template access.
	// Default: os.DirFS(TemplateDir)
	// Use this field to embed templates into your binary using Go's embed package.
	// When provided, TemplateDir is not used for file access.
	TemplatesFS fs.FS

	// TemplateExt defines the file extension for Textwire template files.
	// Default: ".tw"
	// Note: Using a different extension may disable syntax highlighting
	// in editors like VSCode when using the Textwire extension.
	TemplateExt string

	// ErrorPagePath sets the relative path to a custom error page template.
	// Default: internal error page
	// The path is relative to the template directory (TemplateDir or TemplatesFS root).
	ErrorPagePath string

	// DebugMode enables detailed error reporting in the browser and server logs.
	// Default: false (keep false in production)
	// When true, error messages with file paths and line numbers are displayed
	// during development.
	DebugMode bool

	// GlobalData stores shared data accessible across all templates.
	// Access these values in templates using the `global` object (e.g., `global.authUser`).
	// Useful for storing environment variables, configuration, or common data.
	GlobalData map[string]any
}

func New(dir, ext, errPagePath string, debug bool) *Config {
	return &Config{
		TemplateDir:   dir,
		TemplateExt:   ext,
		ErrorPagePath: errPagePath,
		DebugMode:     debug,
		GlobalData:    map[string]any{},
	}
}
