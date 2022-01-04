package main

import (
	"fmt"
	"time"

	"ghostlang.org/x/ghost/evaluator"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/parser"
	"ghostlang.org/x/ghost/scanner"
)

func benchmarkCommand() {
	benchmarkHelloWorld()
}

func benchmarkHelloWorld() {
	goTime := nativeHelloWorld()
	scanTime, parseTime, interpretTime, ghostTime := benchmark(`print("Hello, world!")`)

	fmt.Println("==============================")
	fmt.Println("Hello world benchmark")
	fmt.Println("==============================")
	fmt.Printf("Go:             %s\n", goTime)
	fmt.Printf("Ghost:          %s\n", ghostTime)
	fmt.Printf("-- Scanner:     %s\n", scanTime)
	fmt.Printf("-- Parser:      %s\n", parseTime)
	fmt.Printf("-- Interpreter: %s\n", interpretTime)
}

func nativeHelloWorld() time.Duration {
	start := time.Now()
	fmt.Println("Hello, world!")

	return time.Since(start)
}

func benchmark(source string) (scanTime time.Duration, parseTime time.Duration, interpretTime time.Duration, ghostTime time.Duration) {
	start := time.Now()

	env := object.NewEnvironment()
	scanner := scanner.New(source)
	scanTime = time.Since(start)

	parseStart := time.Now()
	parser := parser.New(scanner)
	program := parser.Parse()
	parseTime = time.Since(parseStart)

	interpretStart := time.Now()
	evaluator.Evaluate(program, env)
	interpretTime = time.Since(interpretStart)
	ghostTime = time.Since(start)

	return scanTime, parseTime, interpretTime, ghostTime
}
