package ghost

import (
	"os"

	"ghostlang.org/x/ghost/environment"
	"ghostlang.org/x/ghost/error"
	"ghostlang.org/x/ghost/interpreter"
	"ghostlang.org/x/ghost/log"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/parser"
	"ghostlang.org/x/ghost/scanner"
)

type Engine struct {
	FatalError  bool
	Source      string
	Environment *environment.Environment
	File        string
	Directory   string
}

func New() *Engine {
	engine := &Engine{
		Environment: environment.NewEnvironment(),
	}

	engine.resetWorkingDirectory()

	return engine
}

func (engine *Engine) resetWorkingDirectory() {
	engine.Directory, _ = os.Getwd()
}

func (engine *Engine) Execute() object.Object {
	scanner := scanner.New(engine.Source)
	tokens := scanner.ScanTokens()
	parser := parser.New(tokens)
	program := parser.Parse()

	if len(parser.Errors()) != 0 {
		logParseErrors(parser.Errors())
		return nil
	}

	result, _ := interpreter.Evaluate(program, engine.Environment)

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
