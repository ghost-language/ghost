package repl

import (
	"io"
	"os"
	"path/filepath"

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
	line.SetTabCompletionStyle(liner.TabPrints)

	_, err := os.Open(history)
	if err != nil {
		f, err := os.Create(history)
		if err != nil {
			log.Error("system error: unable to write to history file: %s", err)
		} else {
			line.WriteHistory(f)
			f.Close()
		}
	}

	ghost := ghost.New()

	for {
		source, err := line.Prompt(prompt)

		if err == liner.ErrPromptAborted {
			log.Info("Exiting...")
			os.Exit(1)
		} else {
			evaluate(ghost, source)

			line.AppendHistory(source)
		}
	}
}

func evaluate(ghost *ghost.Ghost, source string) {
	directory, _ := os.Getwd()

	ghost.SetSource(source)
	ghost.SetFile("repl.ghost")
	ghost.SetDirectory(directory)

	result := ghost.Execute()

	if result != nil {
		log.Info(result.String())
	}
}
