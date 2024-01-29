package textwire

import (
	"io"
	"os"
	"testing"
)

func TestEvalUseStatement(t *testing.T) {
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

	tpl, err := New(&Config{
		TemplateDir: "testdata/before",
	})

	if err != nil {
		t.Errorf("error creating template: %s", err)
	}

	for _, tt := range tests {
		evaluated, err := tpl.Evaluate(tt.fileName, tt.vars)

		if err != nil {
			t.Errorf("error evaluating template: %s", err)
		}

		actual := evaluated.String()

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
