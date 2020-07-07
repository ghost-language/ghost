package repl

import (
	"bufio"
	"fmt"
	"io"

	"ghostlang.org/ghost/evaluator"
	"ghostlang.org/ghost/lexer"
	"ghostlang.org/ghost/object"
	"ghostlang.org/ghost/parser"
)

// PROMPT designates the REPL prompt characters to accept
// user input.
const PROMPT = ">> "

// OUTPUT designates the REPL output characters to display
// program results.
const OUTPUT = "   "

// Start will initiate a new REPL session.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)

		if evaluated != nil {
			io.WriteString(out, OUTPUT+evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "\tPARSE ERROR:\n")

	for _, message := range errors {
		io.WriteString(out, "\t"+message+"\n")
	}
}
