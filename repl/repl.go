package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/textwire/textwire/v3/evaluator"
	"github.com/textwire/textwire/v3/fail"
	"github.com/textwire/textwire/v3/lexer"
	"github.com/textwire/textwire/v3/object"
	"github.com/textwire/textwire/v3/parser"
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
	scope := object.NewScope()

	for {
		fmt.Print(PROMPT)

		scanned := scanner.Scan()
		if !scanned {
			return nil
		}

		l := lexer.New(scanner.Text())
		p := parser.New(l, nil)
		prog := p.ParseProgram()

		if len(p.Errors()) != 0 {
			if err := printParserErrors(out, p.Errors()); err != nil {
				return err
			}
			continue
		}

		e := evaluator.New(nil, nil)
		ctx := evaluator.NewContext(scope, prog.AbsPath)
		evaluated := e.Eval(prog, ctx)
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
