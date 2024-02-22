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

type Author struct {
	Name string
}

type Book struct {
	Title   string
	Authors []Author
	price   float64
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	books := []Book{
		{
			Title: "Harry Potter and the Sorcerer's Stone",
			Authors: []Author{
				{Name: "J.K. Rowling"},
			},
			price: 12.99,
		},
		{
			Title: "The Lord of the Rings",
			Authors: []Author{
				{Name: "J.R.R. Tolkien"},
			},
			price: 24.99,
		},
	}

	err := tpl.Response(w, "home", map[string]interface{}{
		"title":     "Home page",
		"names":     []string{"John", "Jane", "Jack", "Jill"},
		"showNames": true,
		"books":     books,
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
