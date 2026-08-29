package repl

import (
	"io"
	"os"
	"path/filepath"

	"ghostlang.org/x/ghost/ghost"
	"ghostlang.org/x/ghost/log"
	"ghostlang.org/x/ghost/object"
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

	loadHistory(line)

	instance := ghost.New()

	directory, _ := os.Getwd()

	instance.SetDirectory(directory)
	instance.SetFile("repl.gs")

	for {
		source, err := line.Prompt(prompt)

		// Ctrl-C aborts the line being typed and Ctrl-D ends the session. A
		// session that ends this way ended the way it was asked to, so it is
		// not reported as a failure.
		if err == liner.ErrPromptAborted || err == io.EOF {
			saveHistory(line)

			return
		}

		if err != nil {
			log.Error("could not read from the terminal: %s", err)

			return
		}

		if source == "" {
			continue
		}

		line.AppendHistory(source)

		evaluate(instance, source)
	}
}

// evaluate runs one line and shows what it produced.
//
// A failure has already been written out in full by the time this sees it, and
// the session carries on: a mistyped line in a REPL is a normal part of using
// one, not a reason to stop.
func evaluate(instance *ghost.Ghost, source string) {
	instance.SetSource(source)

	result := instance.Execute()

	if result == nil || object.IsError(result) {
		return
	}

	log.Print(result.String())
}

func loadHistory(line *liner.State) {
	file, err := os.Open(history)

	if err != nil {
		return
	}

	defer file.Close()

	line.ReadHistory(file)
}

func saveHistory(line *liner.State) {
	file, err := os.Create(history)

	if err != nil {
		return
	}

	defer file.Close()

	line.WriteHistory(file)
}
