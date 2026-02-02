package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/textwire/textwire/v3"
	"github.com/textwire/textwire/v3/config"
)

var (
	tpl *textwire.Template

	//go:embed templates/*
	templateFS embed.FS
)

func main() {
	var err error

	err = textwire.RegisterStrFunc("_isCool", func(s string, args ...any) any {
		return s == "John Wick"
	})
	if err != nil {
		log.Fatal(err)
	}

	tpl, err = textwire.NewTemplate(&config.Config{
		TemplateFS:    templateFS,
		TemplateDir:   "templates",
		ErrorPagePath: "error-page",
		DebugMode:     true,
		GlobalData: map[string]any{
			"env": "development",
		},
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
				FirstName: "Amy",
				LastName:  "Adams",
			},
		},
	}

	err := tpl.Response(w, "views/home", map[string]any{
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

	err := tpl.Response(w, "views/about", map[string]any{})
	if err != nil {
		log.Println(err.Error())
	}
}
