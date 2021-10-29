package ghost

import (
	"fmt"
	"os"

	"ghostlang.org/x/ghost/log"
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

	for index, token := range tokens {
		log.LogDebug(fmt.Sprintf("[%d] %s", index, token.String()))
	}
}
