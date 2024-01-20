package textwire

import (
	"io"
	"os"
	"testing"
)

func TestEvalLayoutStatement(t *testing.T) {
	tests := []string{
		"1.no-stmts",
		"2.with-inserts",
	}

	NewConfig(&Config{
		TemplateDir: "testdata/before",
	})

	for _, fileName := range tests {
		tpl, err := ParseTemplate(fileName)

		if err != nil {
			t.Errorf("error parsing template: %s", err)
		}

		evaluated, err := tpl.Evaluate(nil)

		actual := evaluated.String()

		if err != nil {
			t.Errorf("error evaluating template: %s", err)
		}

		expected, err := readFile("testdata/expected/" + fileName + ".html")

		if err != nil {
			t.Errorf("error reading expected file: %s", err)
			return
		}

		if actual != expected {
			t.Errorf("wrong result. EXPECTED:\n%q\n-------GOT:--------\n%q", expected, actual)
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
