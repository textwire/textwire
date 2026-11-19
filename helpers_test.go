package textwire

import (
	"os"
	"testing"

	"github.com/textwire/textwire/v6/config"
	"github.com/textwire/textwire/v6/pkg/ast"
	"github.com/textwire/textwire/v6/pkg/evaluator"
	"github.com/textwire/textwire/v6/pkg/value"
)

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
