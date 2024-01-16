package main

import (
	"fmt"
	"net/http"

	"github.com/textwire/textwire"
)

func main() {
	textwire.SetConfig(&textwire.Config{
		TemplateDir: "templates",
	})

	http.HandleFunc("/", homeView)
	fmt.Println("Listening on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func homeView(w http.ResponseWriter, r *http.Request) {
	vars := map[string]interface{}{
		"title": "Hello, World!",
		"age":   23,
	}

	err := textwire.View(w, "home", vars)

	if err != nil {
		fmt.Println(err)
	}
}
