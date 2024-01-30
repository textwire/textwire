package textwire

var config = &Config{
	TemplateDir: "templates",
	TemplateExt: ".textwire.html",
}

// Config is the main configuration for Textwire
type Config struct {
	// TemplateDir is the directory where the Textwire
	// templates are located
	TemplateDir string

	// TemplateExt is the extension of the Textwire
	// template files
	TemplateExt string
}

func New(c *Config) (*Template, error) {
	applyConfig(c)

	paths, err := findTextwireFiles()

	if err != nil {
		return nil, err
	}

	programs, err := parsePrograms(paths)

	return &Template{programs: programs}, err
}
