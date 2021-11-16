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
	FatalError bool
	Source     string
	File       string
	Directory  string
}

func New(source string) *Engine {
	engine := &Engine{
		Source: source,
	}

	engine.resetWorkingDirectory()

	return engine
}

func (engine *Engine) resetWorkingDirectory() {
	engine.Directory, _ = os.Getwd()
}

func (engine *Engine) Execute() object.Object {
	env := environment.NewEnvironment()
	scanner := scanner.New(engine.Source)
	tokens := scanner.ScanTokens()
	parser := parser.New(tokens)
	program := parser.Parse()

	if len(parser.Errors()) != 0 {
		logParseErrors(parser.Errors())
		return nil
	}

	result, _ := interpreter.Evaluate(program, env)

	// log.Debug("Scanned tokens...")
	// for index, token := range tokens {
	// 	log.Debug(fmt.Sprintf("[%d] %s", index, token.String()))
	// }

	// log.Debug("Parsed statements...")
	// for index, statement := range statements {
	// 	log.Debug(fmt.Sprintf("[%d] %T: %q", index, statement, statement))
	// }

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
