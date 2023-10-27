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
	benchmarkFibonacci()
}

func benchmarkHelloWorld() {
	goTime := nativeHelloWorld()
	scanTime, parseTime, interpretTime, ghostTime := benchmark(`print("Hello, world!")`)

	fmt.Println("==============================")
	fmt.Println("Hello world benchmark")
	fmt.Println("==============================")
	fmt.Printf("Go:             %s\n", goTime)
	fmt.Printf("Ghost:          %s\n", ghostTime)
	fmt.Printf("Difference:     %s\n", ghostTime-goTime)
	fmt.Printf("Difference (%%): +%.2f%%\n", (float64(ghostTime-goTime)/float64(goTime))*100)
	fmt.Printf("-- Scanner:     %s\n", scanTime)
	fmt.Printf("-- Parser:      %s\n", parseTime)
	fmt.Printf("-- Interpreter: %s\n\n", interpretTime)
}

func benchmarkFibonacci() {
	goTime := nativeFibonacci()
	scanTime, parseTime, interpretTime, ghostTime := benchmark(`
		function fibonacci(n) {
			if (n <= 1) {
				return n
			}

			return fibonacci(n - 1) + fibonacci(n - 2)
		}
		
		fibonacci(20)
		`)

	fmt.Println("==============================")
	fmt.Println("Fibonacci benchmark")
	fmt.Println("==============================")
	fmt.Printf("Go:             %s\n", goTime)
	fmt.Printf("Ghost:          %s\n", ghostTime)
	fmt.Printf("Difference:     %s\n", ghostTime-goTime)
	fmt.Printf("Difference (%%): +%.2f%%\n", (float64(ghostTime-goTime)/float64(goTime))*100)
	fmt.Printf("-- Scanner:     %s\n", scanTime)
	fmt.Printf("-- Parser:      %s\n", parseTime)
	fmt.Printf("-- Interpreter: %s\n\n", interpretTime)
}

func nativeHelloWorld() time.Duration {
	start := time.Now()
	fmt.Println("Hello, world!")

	return time.Since(start)
}

func nativeFibonacci() time.Duration {
	start := time.Now()

	_ = fibonacci(20)

	return time.Since(start)
}

func benchmark(source string) (scanTime time.Duration, parseTime time.Duration, interpretTime time.Duration, ghostTime time.Duration) {
	start := time.Now()

	scope := &object.Scope{
		Environment: object.NewEnvironment(),
	}

	scanner := scanner.New(source, "benchmark.ghost")
	scanTime = time.Since(start)

	parseStart := time.Now()
	parser := parser.New(scanner)
	program := parser.Parse()
	parseTime = time.Since(parseStart)

	interpretStart := time.Now()
	evaluator.Evaluate(program, scope)
	interpretTime = time.Since(interpretStart)
	ghostTime = time.Since(start)

	return scanTime, parseTime, interpretTime, ghostTime
}

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}

	return fibonacci(n-1) + fibonacci(n-2)
}
