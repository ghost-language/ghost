package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"ghostlang.org/x/ghost/ghost"
	"ghostlang.org/x/ghost/log"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/repl"
	"ghostlang.org/x/ghost/version"
)

var (
	flagHelp    bool
	flagVersion bool
	flagTime    bool
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] [<filename>]\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(0)
	}

	flag.BoolVar(&flagHelp, "h", false, "display help information")
	flag.BoolVar(&flagVersion, "v", false, "display version information")
	flag.BoolVar(&flagTime, "t", false, "display how long the program ran for")
}

func main() {
	flag.Parse()

	if flagVersion {
		fmt.Printf("%s %s\n", path.Base(os.Args[0]), version.Version)
		os.Exit(0)
	}

	if flagHelp {
		helpCommand()
		os.Exit(0)
	}

	args := flag.Args()

	if len(args) == 0 {
		fmt.Printf("Ghost (%s)\n", version.Version)
		fmt.Printf("Press Ctrl + C to exit\n\n")

		repl.Start(os.Stdin, os.Stdout)
		return
	}

	start := time.Now()
	source, err := os.ReadFile(args[0])

	if err != nil {
		log.Error("could not read %s: %s", args[0], err)

		os.Exit(1)
	}

	directory, _ := filepath.Abs(filepath.Dir(args[0]))
	fullPath, _ := filepath.Abs(args[0])
	currentFile := strings.Replace(fullPath, directory+"/", "", 1)

	instance := ghost.New()
	instance.SetSource(string(source))
	instance.SetFile(currentFile)
	instance.SetDirectory(directory)

	result := instance.Execute()

	if flagTime {
		log.Info("(executed in: %s)", time.Since(start))
	}

	// A script that failed has already had its error written out in full. What
	// is left is to say so to whatever ran Ghost, which reads the exit status
	// rather than the terminal.
	if object.IsError(result) {
		os.Exit(1)
	}
}
