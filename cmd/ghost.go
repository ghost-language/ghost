package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"ghostlang.org/x/ghost/repl"
	"ghostlang.org/x/ghost/version"
)

var (
	flagHelp      bool
	flagVersion   bool
	flagBenchmark bool
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

	if len(args) > 1 {
		fmt.Println("Usage: ghost [script]")
		os.Exit(64)
	} else {
		repl.Start(os.Stdin, os.Stdout)
	}
}
