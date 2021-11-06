package ghost

import (
	"os"

	"ghostlang.org/x/ghost/interpreter"
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

func (engine *Engine) Execute() {
	scanner := scanner.New(engine.Source)
	tokens := scanner.ScanTokens()
	parser := parser.New(tokens)
	program := parser.Parse()

	interpreter.Evaluate(program)

	// log.Debug("Scanned tokens...")
	// for index, token := range tokens {
	// 	log.Debug(fmt.Sprintf("[%d] %s", index, token.String()))
	// }

	// log.Debug("Parsed statements...")
	// for index, statement := range statements {
	// 	log.Debug(fmt.Sprintf("[%d] %T: %q", index, statement, statement))
	// }
}