package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/textwire/textwire/v3"
	"github.com/textwire/textwire/v3/config"
)

var (
	//go:embed templates/*
	templateFS embed.FS
)

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

var names = generateStrings(100)
var books = generateBooks(100)

func main() {
	lowerFn := func(s string, args ...any) any {
		return strings.ToLower(s)
	}

	if err := textwire.RegisterStrFunc("_lower", lowerFn); err != nil {
		log.Fatal(err)
	}

	tpl, err := textwire.NewTemplate(&config.Config{
		TemplateFS:    templateFS,
		ErrorPagePath: "error-page",
		DebugMode:     true,
		GlobalData: map[string]any{
			"env": "development",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", homeHandler(tpl))
	http.HandleFunc("/about", aboutHandler(tpl))

	fmt.Println("Listening on http://localhost:8080")

	log.Fatalln(http.ListenAndServe(":8080", nil))
}

func homeHandler(tpl *textwire.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			return
		}

		err := tpl.Response(w, "views/home", map[string]any{
			"names":     names,
			"showNames": true,
			"books":     books,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func aboutHandler(tpl *textwire.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/about" {
			return
		}

		if err := tpl.Response(w, "views/about", map[string]any{}); err != nil {
			log.Printf("Template error: %v", err)
		}
	}
}

func generateBooks(count int) []Book {
	books := make([]Book, count)

	for i := range count {
		books[i] = Book{
			ID:    i + 1,
			Isbn:  fmt.Sprintf("978-3-16-148410-%d", i),
			Title: fmt.Sprintf("Book %d", i+1),
			Author: &Author{
				ID:        i + 1,
				FirstName: fmt.Sprintf("Author%d", i+1),
				LastName:  "Smith",
			},
		}
	}
	return books
}

func generateStrings(count int) []string {
	strs := make([]string, count)

	for i := range count {
		strs[i] = fmt.Sprintf("978-3-16-148410-%d", i)
	}
	return strs
}
