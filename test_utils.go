package textwire

import (
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"testing"

	"codeberg.org/textwire/textwire/v6/config"
	"codeberg.org/textwire/textwire/v6/pkg/ast"
	"codeberg.org/textwire/textwire/v6/pkg/evaluator"
	"codeberg.org/textwire/textwire/v6/pkg/fail"
	"codeberg.org/textwire/textwire/v6/pkg/linker"
	"codeberg.org/textwire/textwire/v6/pkg/value"
)

var (
	unsupportedFunc = func() {}
	unsupportedChan = make(chan int)
)

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

	actual, err := EvaluateString(inp, data)
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

// newTestFileWatcher returns a watcher that doesn't print to stdout, keeping
// test output clean when tests use intentionally broken templates.
func newTestFileWatcher(ln *linker.NodeLinker) *fileWatcher {
	fw := newFileWatcher(ln)
	fw.logger.out = io.Discard
	return fw
}

// newStartedWatcher mirrors Watch setup without starting the ticker, so
// tests can drive tick() directly.
func newStartedWatcher(t *testing.T, ln *linker.NodeLinker) *fileWatcher {
	t.Helper()

	fw := newTestFileWatcher(ln)

	files, err := locateFiles()
	if err != nil {
		t.Fatalf("failed to locate files: %v", err)
	}

	fw.files = files

	return fw
}

func useTemplateDir(t *testing.T, dir string) {
	t.Helper()

	prev := *userConf
	t.Cleanup(func() { *userConf = prev })

	userConf.Configure(&config.Config{TemplateDir: dir})
}

// evalTpl evaluates a template given as source text or as a parsed program,
// like Template.String does.
func evalTpl[T string | *ast.Program](t *testing.T, input T) string {
	t.Helper()

	var prog *ast.Program
	switch v := any(input).(type) {
	case string:
		parsed, errs := parseStr(v)
		if errs != nil {
			t.Fatalf("failed to parse %q: %v", v, errs[0].Error())
		}
		prog = parsed
	case *ast.Program:
		prog = v
	}

	global, failure := value.NewFrozenGlobal(userConf.GlobalData)
	if failure != nil {
		t.Fatalf("failed to create global: %v", failure.Error())
	}

	scope, failure := value.NewScopeFromMap(map[string]any{"name": "Kokomi"}, global)
	if failure != nil {
		t.Fatalf("failed to create scope: %v", failure.Error())
	}

	e := evaluator.New(customFunc, userConf, global)
	evaluated := e.Eval(prog, evaluator.NewContext(scope, prog.AbsPath))
	if evaluated.Is(value.ERR_VAL) {
		t.Fatalf("failed to render program: %v", evaluated.(*value.Error).Err.Error())
	}

	return evaluated.String()
}

// writeTemplate creates a file in a given location.
func writeTemplate(t *testing.T, path, content string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write template: %v", err)
	}
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
