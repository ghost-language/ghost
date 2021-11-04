package main

import "fmt"

func helpCommand() {
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
