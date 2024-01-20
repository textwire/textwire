package textwire

import (
	"io"
	"os"
	"testing"
)

func TestEvalLayoutStatement(t *testing.T) {
	tests := []struct {
		fileName string
		vars     map[string]interface{}
	}{
		{"1.no-stmts", nil},
		{"2.with-inserts", nil},
		{"3.without-layout", map[string]interface{}{
			"pageTitle": "Test Page",
			"NAME_1":    "Anna Korotchaeva",
			"name_2":    "Serhii Cho",
		}},
	}

	NewConfig(&Config{
		TemplateDir: "testdata/before",
	})

	for _, tt := range tests {
		tpl, err := ParseTemplate(tt.fileName)

		if err != nil {
			t.Errorf("error parsing template: %s", err)
		}

		evaluated, err := tpl.Evaluate(tt.vars)

		actual := evaluated.String()

		if err != nil {
			t.Errorf("error evaluating template: %s", err)
		}

		expected, err := readFile("testdata/expected/" + tt.fileName + ".html")

		if err != nil {
			t.Errorf("error reading expected file: %s", err)
			return
		}

		if actual != expected {
			t.Errorf("wrong result. EXPECTED:\n\"%s\"\n-------GOT:--------\n\"%s\"", expected, actual)
		}
	}
}

func readFile(fileName string) (string, error) {
	file, err := os.Open(fileName)

	if err != nil {
		return "", err
	}

	defer file.Close()

	bytes, err := io.ReadAll(file)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
