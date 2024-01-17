package main

import (
	"fmt"
	"net/http"

	"github.com/textwire/textwire"
)

func main() {
	textwire.NewConfig(&textwire.Config{
		TemplateDir: "templates",
	})

	http.HandleFunc("/", homeView())
	fmt.Println("Listening on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func homeView() http.HandlerFunc {
	view, err := textwire.ParseFile("home")

	if err != nil {
		fmt.Println(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		vars := map[string]interface{}{
			"title": "Hello, World!",
			"age":   23,
		}

		err := view.Evaluate(w, vars)

		if err != nil {
			fmt.Println(err)
		}
	}
}
