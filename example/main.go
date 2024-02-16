package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/textwire/textwire"
)

var tpl *textwire.Template

func main() {
	var err error

	tpl, err = textwire.NewTemplate(&textwire.Config{
		TemplateDir: "templates",
	})

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/about", aboutHandler)

	fmt.Println("Listening on http://localhost:8080")

	log.Fatalln(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.Response(w, "home", map[string]interface{}{
		"title":     "Home page",
		"names":     []string{"John", "Jane", "Jack", "Jill"},
		"showNames": true,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.Response(w, "about", map[string]interface{}{
		"title": "About page",
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
