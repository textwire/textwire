package config

// Config is the main configuration for Textwire
type Config struct {
	// TemplateDir is the directory where the Textwire
	// templates are located. Default is `"templates"`
	TemplateDir string

	// TemplateExt is the extension of the Textwire
	// template files. Default is `.tw.html`, recommended to use`.tw`.
	// If you use a different extension other then ".tw.html",
	// you will loose syntax highlighting in VSCode editor
	// if you use the Textwire extension
	TemplateExt string

	// ErrorPagePath is the relative path to the custom error page
	// that will be displayed when an error occurs while
	// rendering a template. Default is an internal error page.
	// It's relative to the `TemplateDi`r directory
	ErrorPagePath string

	// DebugMode is a flag to enable the debug mode. When enabled,
	// you can see error messages in the browser. Default is `false`
	DebugMode bool
}

func New(dir, ext, errPagePath string, debug bool) *Config {
	return &Config{
		TemplateDir:   dir,
		TemplateExt:   ext,
		ErrorPagePath: errPagePath,
		DebugMode:     debug,
	}
}
