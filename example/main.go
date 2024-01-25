package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/textwire/textwire"
)

var tpl *textwire.Template

func main() {
	template, err := textwire.New(&textwire.Config{
		TemplateDir: "templates",
	})

	if err != nil {
		log.Fatalln(err)
	}

	tpl = template

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/about", aboutHandler)

	fmt.Println("Listening on http://localhost:8080")

	log.Fatalln(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	vars := map[string]interface{}{
		"title":     "Home page",
		"names":     []string{"John", "Jane", "Jack", "Jill"},
		"showNames": true,
	}

	err := tpl.EvaluateResponse(w, "home", vars)

	if err != nil {
		fmt.Println(err)
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/about" {
		http.NotFound(w, r)
		return
	}

	vars := map[string]interface{}{
		"title": "About page",
	}

	err := tpl.EvaluateResponse(w, "about", vars)

	if err != nil {
		fmt.Println(err)
	}
}
