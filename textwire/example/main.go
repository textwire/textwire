package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/textwire/textwire/v2"
	"github.com/textwire/textwire/v2/config"
)

var tpl *textwire.Template

func main() {
	var err error

	textwire.RegisterStrFunc("reverse", func(s string, args ...interface{}) string {
		runes := []rune(s)

		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}

		return string(runes)
	})

	tpl, err = textwire.NewTemplate(&config.Config{
		TemplateExt:   ".tw",
		ErrorPagePath: "error-page",
		DebugMode:     true,
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
	if r.URL.Path != "/" {
		return
	}

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
		"names":     []string{"John", "Jane", "Jack", "Jill"},
		"showNames": true,
		"books":     books,
	})
	if err != nil {
		log.Println(err.Error())
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/about" {
		return
	}

	err := tpl.Response(w, "about", map[string]interface{}{})
	if err != nil {
		log.Println(err.Error())
	}
}
