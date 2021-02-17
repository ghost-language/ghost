package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

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
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">>> ")
		data, _ := reader.ReadBytes('\n')

		run(string(data))
	}
}

func runFile(file string) {
	//
}

func run(source string) {
	fmt.Println(source)
}
