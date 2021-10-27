package repl

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"ghostlang.org/x/ghost/scanner"
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
		log.Print("Error writing history file: ", err)
	} else {
		line.WriteHistory(f)
		f.Close()
	}

	for {
		source, err := line.Prompt(prompt)

		if err == liner.ErrPromptAborted {
			fmt.Println("   Exiting...")
			os.Exit(1)
		} else {
			evaluate(source)

			line.AppendHistory(source)
		}
	}
}

func evaluate(source string) {
	scanner := scanner.New(source)
	tokens := scanner.ScanTokens()

	for index, token := range tokens {
		fmt.Printf("   [%d] %s\n", index, token.String())
	}
}
