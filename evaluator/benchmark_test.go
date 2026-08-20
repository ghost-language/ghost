package evaluator

import (
	"testing"

	"ghostlang.org/x/ghost/library/modules"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/parser"
	"ghostlang.org/x/ghost/scanner"
)

// Benchmark programs exercise the hot paths of the tree-walking evaluator:
// function call overhead, variable lookup through the scope chain, arithmetic
// allocation, and method dispatch.
var benchmarks = []struct {
	name   string
	source string
}{
	{
		"Fib",
		`function fib(n) {
			if (n <= 1) { return n }
			return fib(n - 1) + fib(n - 2)
		}
		fib(18)`,
	},
	{
		"LoopArithmetic",
		`total = 0
		for (i = 0; i < 50000; i = i + 1) {
			total = total + i * 2 - 1
		}
		total`,
	},
	{
		"WhileLoop",
		`i = 0
		total = 0
		while (i < 50000) {
			total = total + i
			i = i + 1
		}
		total`,
	},
	{
		"FunctionCalls",
		`function add(a, b) { return a + b }
		total = 0
		for (i = 0; i < 20000; i = i + 1) {
			total = add(total, i)
		}
		total`,
	},
	{
		"NestedScopeLookup",
		`outer = 1
		function level1() {
			function level2() {
				function level3() {
					total = 0
					for (i = 0; i < 10000; i = i + 1) {
						total = total + outer
					}
					return total
				}
				return level3()
			}
			return level2()
		}
		level1()`,
	},
	{
		"ConstantExpressions",
		`total = 0
		for (i = 0; i < 20000; i = i + 1) {
			total = total + (2 * 3 + 4 * 5 - 6)
		}
		total`,
	},
	{
		"ListOperations",
		`items = []
		for (i = 0; i < 10000; i = i + 1) {
			items.push(i)
		}
		total = 0
		for (i = 0; i < 10000; i = i + 1) {
			total = total + items[i]
		}
		total`,
	},
	{
		"StringConcat",
		`s = ""
		for (i = 0; i < 5000; i = i + 1) {
			s = s + "x"
		}
		s.length()`,
	},
	{
		"ClassMethods",
		`class Counter {
			count = 0
			function increment() {
				this.count = this.count + 1
				return this.count
			}
		}
		c = Counter.new()
		for (i = 0; i < 10000; i = i + 1) {
			c.increment()
		}
		c.count`,
	},
	{
		"MapOperations",
		`m = {"a": 1, "b": 2, "c": 3}
		total = 0
		for (i = 0; i < 10000; i = i + 1) {
			total = total + m["a"] + m["b"] + m["c"]
		}
		total`,
	},
}

func BenchmarkEvaluate(b *testing.B) {
	object.RegisterEvaluator(Evaluate)
	modules.RegisterEvaluator(Evaluate)

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				scope := &object.Scope{Environment: object.NewEnvironment()}
				s := scanner.New(bm.source, "bench.ghost")
				p := parser.New(s)
				program := p.Parse()

				if len(p.Errors()) != 0 {
					b.Fatalf("parse errors: %v", p.Errors())
				}

				result := Evaluate(program, scope)

				if object.IsError(result) {
					b.Fatalf("runtime error: %s", result.(*object.Error).Message)
				}
			}
		})
	}
}
