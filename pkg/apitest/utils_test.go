package apitest

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"codeberg.org/textwire/textwire/v6"
	"codeberg.org/textwire/textwire/v6/config"
	"codeberg.org/textwire/textwire/v6/pkg/fail"
)

var (
	unsupportedFunc = func() {}
	unsupportedChan = make(chan int)
)

type tplFile struct {
	relPath string
	content string
}

func assertNoFailure(t *testing.T, failure *fail.Error) {
	t.Helper()

	if failure != nil {
		t.Fatalf("error registering function: %s", failure)
	}
}

func assertUnsupportedResult(t *testing.T, failure *fail.Error, value any) {
	t.Helper()

	if failure == nil {
		t.Fatal("expected error for unsupported return type, got none")
	}

	expect := fmt.Sprintf(fail.ErrUnsupportedType, fmt.Sprintf("%T", value))
	if failure.Message() != expect {
		t.Fatalf("wrong error message. Expect %q, got %q", expect, failure.Message())
	}
}

func assertEvalOutput(t *testing.T, inp string, data map[string]any, expect string) {
	t.Helper()

	actual, err := textwire.EvaluateString(inp, data)
	if err != nil {
		t.Fatalf("error evaluating template: %s", err)
	}

	if actual != expect {
		t.Fatalf("wrong result. Expect %q, got %q", expect, actual)
	}
}

func compareFailures(got, expect *fail.Error) error {
	if got.String() != expect.String() {
		return fmt.Errorf("wrong error message! Expect:\n%q\ngot:\n%q", expect, got)
	}

	if got.Origin() != expect.Origin() {
		return fmt.Errorf(
			"wrong origin on error message, expect %s, got: %s in error message:\n%q",
			expect.Origin(),
			got.Origin(),
			got,
		)
	}

	if !reflect.DeepEqual(got.Pos(), expect.Pos()) {
		return fmt.Errorf(
			"wrong position on error message, expect %v, got: %v in error message:\n%q",
			expect.Pos(),
			got.Pos(),
			got,
		)
	}

	return nil
}

func readFile(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Function readFile failed to close %s: %v", fileName, err)
		}
	}()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func assertReadOnlyValueError(t *testing.T, failure *fail.Error) {
	t.Helper()

	if failure == nil {
		t.Fatal("expected read-only error, got none")
	}

	if failure.Message() != fail.ErrReadOnlyValue {
		t.Fatalf("wrong error message. Expect %q, got %q", fail.ErrReadOnlyValue, failure)
	}
}

// resetConfigAfterTest restores the default Textwire configuration when
// the test finishes, because NewTemplate changes the package-level config.
func resetConfigAfterTest(t *testing.T) {
	t.Helper()
	t.Cleanup(func() { textwire.Configure(config.New()) })
}

// setupTemplate writes files into a temporary template directory and
// creates a template from it. A nil conf uses the default config.
// The global config is reset after the test.
func setupTemplate(t *testing.T, files []tplFile, conf *config.Config) *textwire.Template {
	t.Helper()

	resetConfigAfterTest(t)

	if conf == nil {
		conf = &config.Config{}
	}

	dir := t.TempDir()

	for _, f := range files {
		path := filepath.Join(dir, filepath.FromSlash(f.relPath))

		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("failed to create directory for %s: %v", f.relPath, err)
		}

		if err := os.WriteFile(path, []byte(f.content), 0o644); err != nil {
			t.Fatalf("failed to write template %s: %v", f.relPath, err)
		}
	}

	conf.TemplateDir = dir

	tpl, failure := textwire.NewTemplate(conf)
	if failure != nil {
		t.Fatalf("failed to create template: %s", failure)
	}

	return tpl
}

func assertOutput(t *testing.T, actual string, failure *fail.Error, expect string) {
	t.Helper()

	if failure != nil {
		t.Fatalf("unexpected error: %s", failure)
	}

	if actual != expect {
		t.Fatalf("wrong result. Expect %q, got %q", expect, actual)
	}
}

func assertDepthError(t *testing.T, failure *fail.Error, compName string, maxDepth int) {
	t.Helper()

	if failure == nil {
		t.Fatal("expected depth error, got none")
	}

	expect := fmt.Sprintf(fail.ErrMaxComponentDepth, compName, maxDepth)
	if failure.Message() != expect {
		t.Fatalf("wrong error message. Expect %q, got %q", expect, failure.Message())
	}
}
