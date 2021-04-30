package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"ghostlang.org/x/ghost/repl"
	"ghostlang.org/x/ghost/server"
	"ghostlang.org/x/ghost/version"
)

var (
	flagVersion     bool
	flagInteractive bool
	flagHelp        bool
	flagServer      bool
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] [<filename>]\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(0)
	}

	flag.BoolVar(&flagHelp, "h", false, "display help information")
	flag.BoolVar(&flagVersion, "v", false, "display version information")
	flag.BoolVar(&flagInteractive, "i", false, "enable interactive mode")
	flag.BoolVar(&flagServer, "s", false, "run as a server")
}

func main() {
	flag.Parse()

	if flagVersion {
		fmt.Printf("%s %s\n", path.Base(os.Args[0]), version.String())
		os.Exit(0)
	}

	if flagHelp {
		showHelp()
		os.Exit(2)
	}

	args := flag.Args()

	if flagServer {
		server := server.New(args)
		server.Run()
	} else {
		repl := repl.New(args, &repl.Options{
			Interactive: flagInteractive,
		})

		repl.Run()
	}
}

func showHelp() {
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Println("    ghost [flags] {file}")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println()
	fmt.Println("    -h  show help")
	fmt.Println("    -i  enter interactive mode after executing file")
	fmt.Println("    -v  show version")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println()
	fmt.Println("    ghost")
	fmt.Println()
	fmt.Println("            Start Ghost REPL")
	fmt.Println()
	fmt.Println("    ghost example.ghost")
	fmt.Println()
	fmt.Println("            Execute source file (example.ghost)")
	fmt.Println()
	fmt.Println("    ghost -i example.ghost")
	fmt.Println()
	fmt.Println("            Execute source file (example.ghost)")
	fmt.Println("            and enter interactive mode (REPL)")
	fmt.Println("            with the scripts environment intact")
	fmt.Println()
	fmt.Println()
}
