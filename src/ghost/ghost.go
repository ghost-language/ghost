package ghost

import (
	"os"

	"ghostlang.org/x/ghost/error"
	"ghostlang.org/x/ghost/interpreter"
	"ghostlang.org/x/ghost/log"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/parser"
	"ghostlang.org/x/ghost/scanner"
)

type Ghost struct {
	FatalError  bool
	Source      string
	Environment *object.Environment
	File        string
	Directory   string
}

func New() *Ghost {
	ghost := &Ghost{
		Environment: object.NewEnvironment(),
	}

	ghost.resetWorkingDirectory()

	return ghost
}

func (ghost *Ghost) resetWorkingDirectory() {
	ghost.Directory, _ = os.Getwd()
}

func (ghost *Ghost) Execute() object.Object {
	scanner := scanner.New(ghost.Source)
	parser := parser.New(scanner)
	program := parser.Parse()

	if len(parser.Errors()) != 0 {
		logParseErrors(parser.Errors())
		return nil
	}

	result, _ := interpreter.Evaluate(program, ghost.Environment)

	return result
}

func logParseErrors(errors []string) {
	for _, message := range errors {
		err := error.Error{
			Reason:  error.Syntax,
			Message: message,
		}

		log.Error(err.Reason, err.Message)
	}
}
