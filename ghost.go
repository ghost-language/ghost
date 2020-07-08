package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"ghostlang.org/ghost/repl"
	"ghostlang.org/ghost/version"
)

var (
	flagVersion     bool
	flagInteractive bool
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] [<filename>]\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(0)
	}

	flag.BoolVar(&flagVersion, "v", false, "display version information")
	flag.BoolVar(&flagInteractive, "i", false, "enable interactive mode")
}

func main() {
	flag.Parse()

	if flagVersion {
		fmt.Printf("%s %s\n", path.Base(os.Args[0]), version.String())
		os.Exit(0)
	}

	args := flag.Args()
	opts := &repl.Options{
		Interactive: flagInteractive,
	}

	repl := repl.New(args, opts)
	repl.Run()
}
