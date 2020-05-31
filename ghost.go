package main

import (
	"fmt"
	"os"

	"ghostlang.org/ghost/repl"
)

func main() {
	args := os.Args[1:]

	if len(args) > 1 {
		fmt.Println("Usage: ghost [script]")
		os.Exit(64)
	} else if len(args) == 1 {
		// runFile(args[0])
	} else {
		repl.Start(os.Stdin, os.Stdout)
	}
}
