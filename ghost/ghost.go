package ghost

import (
	"io"
	"os"

	"ghostlang.org/x/ghost/environment"
	"ghostlang.org/x/ghost/glitch"
	"ghostlang.org/x/ghost/interpreter"
	"ghostlang.org/x/ghost/parser"
	"ghostlang.org/x/ghost/scanner"
)

func Run(source string, writer io.Writer) {
	env := environment.New()
	env.SetWriter(writer)

	scanner := scanner.New(source)
	tokens := scanner.ScanTokens()
	parser := parser.New(tokens)
	statements := parser.Parse()

	if glitch.HadParseError {
		return
	}

	interpreter.Interpret(statements, env)

	if glitch.HadParseError || glitch.HadRuntimeError {
		os.Exit(1)
	}
}
