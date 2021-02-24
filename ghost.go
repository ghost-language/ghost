package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"ghostlang.org/x/ghost/ghost"
	"ghostlang.org/x/ghost/interpreter"
	"ghostlang.org/x/ghost/parser"
	"ghostlang.org/x/ghost/scanner"
	"ghostlang.org/x/ghost/version"
	"github.com/peterh/liner"
)

var (
	flagVersion     bool
	flagInteractive bool
	flagHelp        bool
	flagTokens      bool
	history         = filepath.Join(os.TempDir(), ".ghost_history")
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
	flag.BoolVar(&flagTokens, "t", false, "output scanned token information")
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

	if len(args) > 1 {
		fmt.Println("Usage: ghost [script]")
		os.Exit(64)
	} else if len(args) == 1 {
		runFile(args[0])
	} else {
		runPrompt()
	}
}

func runPrompt() {
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
		source, err := line.Prompt(">> ")

		if err == liner.ErrPromptAborted {
			fmt.Println("   Exiting...")
			os.Exit(1)
		} else {
			run(source)
			ghost.HadParseError = false
			line.AppendHistory(source)
		}
	}
}

func runFile(file string) {
	//
}

func run(source string) {
	scanner := scanner.New(source)
	tokens := scanner.ScanTokens()
	parser := parser.New(tokens)
	statements := parser.Parse()
	interpreter.Interpret(statements)

	if flagTokens {
		fmt.Printf("   =====\n")

		for index, token := range tokens {
			fmt.Printf("   [%d] %s\n", index, token.String())
		}
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
