package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"ghostlang.org/x/ghost/error"
	"ghostlang.org/x/ghost/ghost"
	"ghostlang.org/x/ghost/log"
	"ghostlang.org/x/ghost/repl"
	"ghostlang.org/x/ghost/version"
)

var (
	flagHelp      bool
	flagVersion   bool
	flagBenchmark bool
	flagTokens    bool
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] [<filename>]\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(0)
	}

	flag.BoolVar(&flagHelp, "h", false, "display help information")
	flag.BoolVar(&flagVersion, "v", false, "display version information")
	flag.BoolVar(&flagBenchmark, "b", false, "run benchmark tests against Ghost and Go")
	flag.BoolVar(&flagTokens, "t", false, "display tokens of passed source file")
}

func main() {
	flag.Parse()

	if flagVersion {
		fmt.Printf("%s %s\n", path.Base(os.Args[0]), version.Version)
		os.Exit(0)
	}

	if flagHelp {
		helpCommand()
		os.Exit(2)
	}

	if flagBenchmark {
		benchmarkCommand()
		os.Exit(2)
	}

	args := flag.Args()

	if len(args) == 0 {
		fmt.Printf("Ghost (%s)\n", version.Version)
		fmt.Printf("Press Ctrl + C to exit\n\n")

		repl.Start(os.Stdin, os.Stdout)
		return
	}

	if len(args) > 0 {
		start := time.Now()
		sourceFile, err := os.Open(args[0])

		if err != nil {
			err := error.Error{
				Reason:  error.System,
				Message: fmt.Sprintf("could not open source file %s: %s", args[0], err),
			}

			log.Error(err.Reason, err.Message)
			os.Exit(1)
		}

		defer sourceFile.Close()

		sourceBuffer := bytes.NewBuffer(nil)
		io.Copy(sourceBuffer, sourceFile)
		source := sourceBuffer.String()

		// directory, _ := filepath.Abs(filepath.Dir(args[0]))

		ghost := ghost.New()
		ghost.Source = source
		ghost.Execute()

		elapsed := time.Since(start)

		if flagTokens {
			log.Info("coming soon")
			// scanner := scanner.New(ghost.Source)
			// tokens := scanner.ScanTokens()

			// for index, token := range tokens {
			// 	log.Info("[%03d] - %s: %s", index, token.Type, token.Lexeme)
			// }
		}

		log.Info("(executed in: %s)\n", elapsed)
	}
}
