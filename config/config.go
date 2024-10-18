package config

// Config is the main configuration for Textwire
type Config struct {
	// TemplateDir is the directory where the Textwire
	// templates are located. Default is "templates"
	TemplateDir string

	// TemplateExt is the extension of the Textwire
	// template files. Default is ".tw.html"
	// If you use a different extension other then ".tw.html",
	// you will loose syntax highlighting in VSCode editor
	// if you use the Textwire extension
	TemplateExt string
}

func New(dir, ext string) *Config {
	return &Config{
		TemplateDir: dir,
		TemplateExt: ext,
	}
}
