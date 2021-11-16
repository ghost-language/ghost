package repl

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"ghostlang.org/x/ghost/error"
	"ghostlang.org/x/ghost/ghost"
	"ghostlang.org/x/ghost/log"
	"github.com/peterh/liner"
)

var (
	prompt  = ">> "
	history = filepath.Join(os.TempDir(), ".ghost_history")
)

func Start(in io.Reader, out io.Writer) {
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

	if f, err := os.Open(history); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	if f, err := os.Create(history); err != nil {
		err := error.Error{
			Reason:  error.System,
			Message: fmt.Sprintf("unable to write to history file: %s", err),
		}

		log.Error(err.Reason, err.Message)
	} else {
		line.WriteHistory(f)
		f.Close()
	}

	engine := ghost.New()

	for {
		source, err := line.Prompt(prompt)

		if err == liner.ErrPromptAborted {
			log.Info("Exiting...")
			os.Exit(1)
		} else {
			evaluate(engine, source)

			line.AppendHistory(source)
		}
	}
}

func evaluate(engine *ghost.Engine, source string) {
	engine.Source = source
	result := engine.Execute()

	if result != nil {
		log.Info(result.String())
	}
}
