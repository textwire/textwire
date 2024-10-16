package config

type StrFunc func(string) string
type ArrFunc func([]interface{}) []string
type IntFunc func(int) int
type FloatFunc func(float64) float64
type BoolFunc func(bool) bool

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

	// Funcs is the custom functions that can be used
	Funcs Funcs
}

func New(dir, ext string) *Config {
	return &Config{
		TemplateDir: dir,
		TemplateExt: ext,
		Funcs: Funcs{
			Str:   make(map[string]StrFunc),
			Arr:   make(map[string]ArrFunc),
			Int:   make(map[string]IntFunc),
			Float: make(map[string]FloatFunc),
			Bool:  make(map[string]BoolFunc),
		},
	}
}

type Funcs struct {
	Str   map[string]StrFunc
	Arr   map[string]ArrFunc
	Int   map[string]IntFunc
	Float map[string]FloatFunc
	Bool  map[string]BoolFunc
}
