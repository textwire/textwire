package textwire

import (
	"testing"
)

func TestEvalLayoutStatement(t *testing.T) {
	tests := []struct {
		inputFile string
		expected  string
	}{
		{
			"with-layout-stmt",
			"<div>This is a layout</div>",
		},
	}

	NewConfig(&Config{
		TemplateDir: "testdata",
	})

	for _, tt := range tests {
		tpl, err := ParseTemplate(tt.inputFile)

		if err != nil {
			t.Errorf("error parsing template: %s", err)
		}

		evaluated, err := tpl.Evaluate(nil)

		actual := evaluated.String()

		if err != nil {
			t.Errorf("error evaluating template: %s", err)
		}

		if actual != tt.expected {
			t.Errorf("wrong result. expected=%q, got=%q", tt.expected, actual)
		}
	}
}
