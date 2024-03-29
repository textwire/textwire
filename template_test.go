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

func TestEvalUseStmt(t *testing.T) {
	tests := []struct {
		fileName string
		data     map[string]interface{}
	}{
		{"1.no-stmts", nil},
		{"2.with-inserts", nil},
		{"3.without-layout", map[string]interface{}{
			"pageTitle": "Test Page",
			"NAME_1":    "Anna Korotchaeva",
			"name_2":    "Serhii Cho",
		}},
		{"4.loops", map[string]interface{}{
			"names": []string{"Anna", "Serhii", "Vladimir"},
		}},
		{"5.with-component", map[string]interface{}{
			"names": []string{"Anna", "Serhii", "Vladimir"},
		}},
		{"6.use-inside-if", nil},
		{"7.insert-without-use", nil},
		{"8.with-component", nil},
	}

	tpl, err := NewTemplate(&Config{
		TemplateDir: "testdata/good/before",
	})

	if err != nil {
		t.Errorf("error creating template: %s", err)
		return
	}

	for _, tt := range tests {
		actual, evalErr := tpl.String(tt.fileName, tt.data)

		if evalErr != nil {
			t.Errorf("error evaluating template: %s", evalErr)
			return
		}

		expected, err := readFile("testdata/good/expected/" + tt.fileName + ".html")

		if err != nil {
			t.Errorf("error reading expected file: %s", err)
			return
		}

		if actual != expected {
			t.Errorf("wrong result. EXPECTED:\n\"%s\"\nGOT:\n\"%s\"",
				expected, actual)
		}
	}
}
