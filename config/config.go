package config

import "github.com/textwire/textwire/v2/object"

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
			Str:   make(map[string]object.BuiltinFunction),
			Arr:   make(map[string]object.BuiltinFunction),
			Int:   make(map[string]object.BuiltinFunction),
			Float: make(map[string]object.BuiltinFunction),
			Bool:  make(map[string]object.BuiltinFunction),
		},
	}
}

type Funcs struct {
	Str   map[string]object.BuiltinFunction
	Arr   map[string]object.BuiltinFunction
	Int   map[string]object.BuiltinFunction
	Float map[string]object.BuiltinFunction
	Bool  map[string]object.BuiltinFunction
}
