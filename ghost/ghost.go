package ghost

import (
	"io"
	"os"

	"ghostlang.org/x/ghost/errors"
	"ghostlang.org/x/ghost/interpreter"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/parser"
	"ghostlang.org/x/ghost/scanner"
)

func Run(source string, writer io.Writer) {
	env := object.NewEnvironment()
	env.SetWriter(writer)

	scanner := scanner.New(source)
	tokens := scanner.ScanTokens()
	parser := parser.New(tokens)
	statements := parser.Parse()

	if errors.HadParseError {
		return
	}

	interpreter.Interpret(statements, env)

	if errors.HadParseError || errors.HadRuntimeError {
		os.Exit(1)
	}
}
