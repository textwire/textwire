package textwire

import (
	"testing"

	"github.com/textwire/textwire/fail"
)

func TestErrorHandlingEvaluatingTemplate(t *testing.T) {
	tpl, _ := NewTemplate(&Config{
		TemplateDir: "testdata/bad",
	})

	path, err := getFullPath("", false)
	path += "/"

	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	tests := []struct {
		fileName string
		err      *fail.Error
		data     map[string]interface{}
	}{
		{
			"1.use-inside-tpl",
			fail.New(1, path+"1.use-inside-tpl.tw.html", "evaluator", fail.ErrUseStmtNotAllowed),
			nil,
		},
	}

	for _, tt := range tests {
		_, err := tpl.String(tt.fileName, tt.data)

		if err == nil {
			t.Errorf("expected error but got none")
			return
		}

		if err.String() != tt.err.String() {
			t.Errorf("wrong error message. EXPECTED:\n\"%s\"\nGOT:\n\"%s\"", tt.err, err)
		}
	}
}
