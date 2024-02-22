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
	Title  string
	Author Author
	price  float64
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	books := []Book{
		{Title: "The Catcher in the Rye", Author: Author{Name: "J.D. Salinger"}, price: 10.99},
		{Title: "To Kill a Mockingbird", Author: Author{Name: "Harper Lee"}, price: 12.99},
		{Title: "1984", Author: Author{Name: "George Orwell"}, price: 9.99},
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
