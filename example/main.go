package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/textwire/textwire"
	"github.com/textwire/textwire/fail"
)

var tpl *textwire.Template

func main() {
	var err *fail.Error

	tpl, err = textwire.TemplateEngine(&textwire.Config{
		TemplateDir: "templates",
	})

	err.IfErrorFatal()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/about", aboutHandler)

	fmt.Println("Listening on http://localhost:8080")

	log.Fatalln(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.View(w, "home", map[string]interface{}{
		"title":     "Home page",
		"names":     []string{"John", "Jane", "Jack", "Jill"},
		"showNames": true,
	})

	err.IfErrorFatal()
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.View(w, "about", map[string]interface{}{
		"title": "About page",
	})

	err.IfErrorFatal()
}
