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

	fmt.Println("Listening on http://localhost:8080")
	log.Fatalln(http.ListenAndServe(":8080", nil))
}

type HomeVars struct {
	Title string
	Age   int
}

func homeHandler() http.HandlerFunc {
	template, err := textwire.ParseTemplate("home")

	if err != nil {
		fmt.Println(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		vars := map[string]interface{}{
			"title": "Hello, World!",
			"age":   23,
		}

		err := template.Evaluate(w, vars)

		if err != nil {
			fmt.Println(err)
		}
	}
}
