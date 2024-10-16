package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/textwire/textwire/v2"
	"github.com/textwire/textwire/v2/object"
	"github.com/textwire/textwire/v2/option"
)

var tpl *textwire.Template

func main() {
	var err error

	textwire.RegisterStrFunc("reverse", func(s *object.Str, args ...object.Object) object.Object {
		runes := []rune(s.Value)

		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}

		return &object.Str{Value: string(runes)}
	})

	tpl, err = textwire.NewTemplate(&option.Option{
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

type Book struct {
	ID     int     `json:"id,omitempty"`
	Isbn   string  `json:"isbn,omitempty"`
	Title  string  `json:"title,omitempty"`
	Author *Author `json:"author,omitempty"`
}

type Author struct {
	ID        int    `json:"id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	books := []Book{
		{
			ID:    1,
			Isbn:  "978-3-16-148410-0",
			Title: "The Go Programming Language",
			Author: &Author{
				ID:        1,
				FirstName: "Alan",
				LastName:  "Donovan",
			},
		},
		{
			ID:    2,
			Isbn:  "978-3-16-148410-1",
			Title: "The Rust Programming Language",
			Author: &Author{
				ID:        2,
				FirstName: "Steve",
				LastName:  "Klabnik",
			},
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
