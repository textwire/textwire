package apitest

import (
	"fmt"
	"io"
	"log"
	"os"
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

// resetConfigAfterTest restores the default Textwire configuration when
// the test finishes, because NewTemplate changes the package-level config.
func resetConfigAfterTest(t *testing.T) {
	t.Helper()
	t.Cleanup(func() { textwire.Configure(config.New()) })
}
