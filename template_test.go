package textwire

import (
	"testing"

	"github.com/textwire/textwire/fail"
)

func TestErrorHandlingEvaluatingTemplate(t *testing.T) {
	t.Skip()
	tests := []struct {
		fileName string
		err      *fail.Error
		data     map[string]interface{}
	}{
		{"1.use-inside-tpl", fail.New(1, "", "evaluator", fail.ErrUseStmtNotAllowed), nil},
	}

	for _, tt := range tests {
		tpl, _ := TemplateEngine(&Config{
			TemplateDir: "testdata/bad",
		})

		_, err := tpl.String(tt.fileName, tt.data)

		if err == nil {
			t.Errorf("expected error but got none")
			return
		}

		if err.String() != tt.err.String() {
			t.Errorf("wrong error message. EXPECTED:\n\"%s\"\n-------GOT:--------\n\"%s\"",
				tt.err, err)
		}
	}
}
