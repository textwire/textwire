package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/textwire/textwire/evaluator"
	"github.com/textwire/textwire/fail"
	"github.com/textwire/textwire/lexer"
	"github.com/textwire/textwire/object"
	"github.com/textwire/textwire/parser"
)

const PROMPT = ">>> "

func main() {
	fmt.Print("Interactive shell\n\n")

	Start(os.Stdin, os.Stdout)
}

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnv()

	for {
		fmt.Print(PROMPT)

		scanned := scanner.Scan()

		if !scanned {
			return
		}

		l := lexer.New(scanner.Text())
		p := parser.New(l, "")

		prog := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluator := evaluator.New(nil)
		evaluated := evaluator.Eval(prog, env)

		if evaluated != nil {
			io.WriteString(out, evaluated.String())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []*fail.Error) {
	io.WriteString(out, "Textwire errors:\n")

	for _, err := range errors {
		io.WriteString(out, "\t"+err.String()+"\n")
	}
}
