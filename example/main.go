package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/textwire/textwire"
)

var tpl *textwire.Template

func main() {
	tpl = textwire.New(&textwire.Config{
		TemplateDir: "templates",
	})

	if tpl.HasErrors() {
		log.Fatal(tpl.FirstError())
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/about", aboutHandler)

	fmt.Println("Listening on http://localhost:8080")

	log.Fatalln(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"title":     "Home page",
		"names":     []string{"John", "Jane", "Jack", "Jill"},
		"showNames": true,
	}

	err := tpl.View(w, "home", data)

	if err != nil {
		fmt.Println(err)
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"title": "About page",
	}

	err := tpl.View(w, "about", data)

	if err != nil {
		fmt.Println(err)
	}
}
