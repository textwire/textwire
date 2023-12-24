package main

import (
	"fmt"
	"os"

	"github.com/textwire/textwire/repl"
)

func main() {
	fmt.Printf("Interactive shell\n\n")

	repl.Start(os.Stdin, os.Stdout)
}
