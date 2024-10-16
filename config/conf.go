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

	// Configurations for custom functions
	StrFuncs   map[string]StrFunc
	ArrFuncs   map[string]ArrFunc
	IntFuncs   map[string]IntFunc
	FloatFuncs map[string]FloatFunc
	BoolFuncs  map[string]BoolFunc
}
