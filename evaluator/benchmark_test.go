package evaluator

import (
	"testing"

	"ghostlang.org/x/ghost/library/modules"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/optimizer"
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
			increment() {
				this.count = this.count + 1
				return this.count
			}
		}
		c = new Counter()
		for (i = 0; i < 10000; i = i + 1) {
			c.increment()
		}
		c.count`,
	},
	{
		// A game loop: many globals in one flat scope, read from a tight
		// nested loop, with no function calls to push work into child scopes.
		// This shape stresses lookup in a large scope rather than call setup,
		// and the two pull environment storage in opposite directions.
		"TileRender",
		`layers = []
		for (l = 0; l < 3; l = l + 1) {
			rows = []
			for (y = 0; y < 30; y = y + 1) {
				row = []
				for (x = 0; x < 40; x = x + 1) { row[x] = (x + y + l) % 8 }
				rows[y] = row
			}
			layers[l] = rows
		}
		total = 0
		for (frame = 0; frame < 10; frame = frame + 1) {
			for (l = 0; l < 3; l = l + 1) {
				rows = layers[l]
				for (y = 0; y < 30; y = y + 1) {
					row = rows[y]
					for (x = 0; x < 40; x = x + 1) {
						tile = row[x]
						if (tile != 0) { total = total + tile * 16 + x * 16 + y * 16 }
					}
				}
			}
		}
		total`,
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
				s := scanner.New(bm.source, "bench.gs")
				p := parser.New(s)
				program := p.Parse()

				if len(p.Errors()) != 0 {
					b.Fatalf("parse errors: %v", p.Errors())
				}

				// Mirror the real pipeline in ghost.Execute, which optimizes
				// the program before evaluating it.
				program = optimizer.Optimize(program)

				result := Evaluate(program, scope)

				if object.IsError(result) {
					b.Fatalf("runtime error: %s", result.(*object.Error).String())
				}
			}
		})
	}
}
