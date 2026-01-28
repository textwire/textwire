package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/textwire/textwire/v2/evaluator"
	"github.com/textwire/textwire/v2/fail"
	"github.com/textwire/textwire/v2/lexer"
	"github.com/textwire/textwire/v2/object"
	"github.com/textwire/textwire/v2/parser"
)

const PROMPT = ">>> "

func main() {
	fmt.Print("Interactive shell\n\n")

	if err := Start(os.Stdin, os.Stdout); err != nil {
		fmt.Println("ERROR: ", err)
	}
}

func Start(in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)
	env := object.NewEnv()

	for {
		fmt.Print(PROMPT)

		scanned := scanner.Scan()

		if !scanned {
			return nil
		}

		l := lexer.New(scanner.Text())
		p := parser.New(l, "")

		prog := p.Parse()

		if len(p.Errors()) != 0 {
			if err := printParserErrors(out, p.Errors()); err != nil {
				return err
			}
			continue
		}

		evaluator := evaluator.New(nil, nil)
		evaluated := evaluator.Eval(prog, env, prog.Filepath)

		if evaluated == nil {
			continue
		}

		if _, err := io.WriteString(out, evaluated.String()+"\n"); err != nil {
			return err
		}
	}
}

func printParserErrors(out io.Writer, errors []*fail.Error) error {
	if _, err := io.WriteString(out, "Textwire errors:\n"); err != nil {
		return err
	}

	for _, err := range errors {
		if _, err := io.WriteString(out, "\t"+err.String()+"\n"); err != nil {
			return err
		}
	}

	return nil
}
