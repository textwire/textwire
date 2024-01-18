package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/textwire/textwire"
)

func main() {
	textwire.NewConfig(&textwire.Config{
		TemplateDir: "templates",
	})

	http.HandleFunc("/", homeHandler())
	http.HandleFunc("/about", aboutHandler())

	fmt.Println("Listening on http://localhost:8080")

	log.Fatalln(http.ListenAndServe(":8080", nil))
}

func homeHandler() http.HandlerFunc {
	template, err := textwire.ParseTemplate("home")

	if err != nil {
		fmt.Println(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		vars := map[string]interface{}{
			"title": "Hello, World!",
			"age":   23,
		}

		err := template.EvaluateResponse(w, vars)

		if err != nil {
			fmt.Println(err)
		}
	}
}

func aboutHandler() http.HandlerFunc {
	template, err := textwire.ParseTemplate("about")

	if err != nil {
		fmt.Println(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/about" {
			http.NotFound(w, r)
			return
		}

		vars := map[string]interface{}{
			"title": "Hello, World!",
			"age":   23,
		}

		err := template.EvaluateResponse(w, vars)

		if err != nil {
			fmt.Println(err)
		}
	}
}
