package config

import (
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"testing/fstest"
	"time"
)

func TestConfigure_MaxComponentDepth(t *testing.T) {
	cases := []struct {
		name   string
		value  int
		expect int
	}{
		{name: "zero keeps the default", value: 0, expect: DefaultMaxComponentDepth},
		{name: "negative keeps the default", value: -1, expect: DefaultMaxComponentDepth},
		{name: "positive overrides the default", value: 50, expect: 50},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conf := New()
			conf.Configure(&Config{TemplateDir: t.TempDir(), MaxComponentDepth: tc.value})

			if conf.MaxComponentDepth != tc.expect {
				t.Fatalf("expected %d, got %d", tc.expect, conf.MaxComponentDepth)
			}
		})
	}
}

func TestConfigure_NilKeepsTemplateFS(t *testing.T) {
	conf := New()
	conf.Configure(&Config{
		TemplateDir: "tpl",
		TemplateFS: fstest.MapFS{
			"tpl/index.tw": &fstest.MapFile{Data: []byte("embedded")},
		},
	})

	conf.Configure(nil)

	content, err := fs.ReadFile(conf.TemplateFS, "tpl/index.tw")
	if err != nil {
		t.Fatalf("expected the embedded template source to be kept: %v", err)
	}

	if string(content) != "embedded" {
		t.Fatalf("expected %q, got %q", "embedded", content)
	}

	if !conf.UsesFS() {
		t.Fatal("expected the embedded source to stay in use")
	}
}

func TestConfigure_NilDerivesTemplateFSWhenMissing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "index.tw")

	if err := os.WriteFile(path, []byte("disk"), 0o644); err != nil {
		t.Fatalf("failed to write template: %v", err)
	}

	conf := New()
	conf.TemplateDir = dir
	conf.Configure(nil)

	if conf.TemplateFS == nil {
		t.Fatal("expected a template source to be derived from TemplateDir")
	}

	if conf.UsesFS() {
		t.Fatal("expected a disk source, not a custom file system")
	}

	content, err := fs.ReadFile(conf.TemplateFS, "index.tw")
	if err != nil {
		t.Fatalf("expected to read from the derived source: %v", err)
	}

	if string(content) != "disk" {
		t.Fatalf("expected %q, got %q", "disk", content)
	}
}

func TestConfigure_DoesNotAliasGlobalData(t *testing.T) {
	conf := New()
	data := map[string]any{"env": "dev"}

	conf.Configure(&Config{GlobalData: data})
	data["env"] = "prod"

	if conf.GlobalData["env"] != "dev" {
		t.Fatalf("expected %q, got %v", "dev", conf.GlobalData["env"])
	}
}

// TestConfigure_NilKeepsOtherSettings makes sure Configure(nil) leaves every
// setting except the template source alone. TemplateFS itself is covered by
// TestConfigure_NilKeepsTemplateFS and TestConfigure_NilDerivesTemplateFSWhenMissing.
func TestConfigure_NilKeepsOtherSettings(t *testing.T) {
	conf := New()
	conf.TemplateDir = t.TempDir()
	conf.TemplateExt = ".twig"
	conf.ErrorPagePath = "error"
	conf.DebugMode = true
	conf.GlobalData = map[string]any{"env": "prod"}
	conf.FileWatcher = true
	conf.WatcherInterval = 2 * time.Second
	conf.MaxComponentDepth = 7

	newConf := conf.Clone()
	conf.Configure(nil)

	// Ignore the template source derived from TemplateDir. Every other field
	// has to stay exactly as it was.
	conf.TemplateFS = nil
	conf.usesFS = false

	if !reflect.DeepEqual(conf, newConf) {
		t.Fatalf("Configure(nil) changed settings.\nbefore: %+v\nafter:  %+v", newConf, conf)
	}
}

func TestClone_CopiesGlobalData(t *testing.T) {
	orig := New()
	orig.GlobalData = map[string]any{"env": "dev"}

	clone := orig.Clone()
	clone.GlobalData["env"] = "prod"

	if orig.GlobalData["env"] != "dev" {
		t.Fatalf("expected the original to keep %q, got %v", "dev", orig.GlobalData["env"])
	}
}

// TestClone_CopiesMutableContainers guards Clone against new map or slice
// fields that are not copied.
func TestClone_CopiesMutableContainers(t *testing.T) {
	orig := New()
	fillMutableContainers(t, orig)

	clone := orig.Clone()

	origVal := reflect.ValueOf(orig).Elem()
	cloneVal := reflect.ValueOf(clone).Elem()

	for i := range origVal.NumField() {
		field := origVal.Type().Field(i)
		origField := origVal.Field(i)

		if !isMutableContainer(origField) {
			continue
		}

		if origField.IsNil() {
			t.Fatalf("test setup did not fill the %q field", field.Name)
		}

		if origField.Pointer() == cloneVal.Field(i).Pointer() {
			t.Fatalf("Clone shares the %q field with the original", field.Name)
		}
	}
}

// fillMutableContainers sets every map and slice field to a non-empty value.
func fillMutableContainers(t *testing.T, conf *Config) {
	t.Helper()

	v := reflect.ValueOf(conf).Elem()

	for i := range v.NumField() {
		field := v.Field(i)

		switch field.Kind() {
		case reflect.Map:
			m := reflect.MakeMap(field.Type())
			m.SetMapIndex(reflect.Zero(field.Type().Key()), reflect.Zero(field.Type().Elem()))
			field.Set(m)
		case reflect.Slice:
			field.Set(reflect.MakeSlice(field.Type(), 1, 1))
		}
	}
}

func isMutableContainer(v reflect.Value) bool {
	return v.Kind() == reflect.Map || v.Kind() == reflect.Slice
}
