package evaluator

import "testing"

// The tests in this file cover §13.13, §13.14 and §13.15 — one design
// question (§14 decision 9) answered three ways: assignment walks the
// enclosing chain to rebind an existing name, blocks introduce a scope of
// their own, and a loop binds its control variable once per iteration.
//
// Each behavior below was a silent wrong answer before the fix, so several of
// these tests assert what used to be a *different* value rather than what used
// to be an error.

func TestAssignmentReachesAnEnclosingScope(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			name: "a function rebinds a variable declared outside it",
			input: `
				scale = 1
				function set(n) { scale = n }
				set(4)
				scale
			`,
			expected: 4,
		},
		{
			name: "a nested function reaches past its immediate parent",
			input: `
				total = 0
				function outer() {
					function inner() { total = 9 }
					inner()
				}
				outer()
				total
			`,
			expected: 9,
		},
		{
			name: "a parameter shadows an outer name rather than rebinding it",
			input: `
				value = 1
				function take(value) { value = 99 }
				take(5)
				value
			`,
			expected: 1,
		},
		{
			name: "a name bound nowhere is declared locally, not globally",
			input: `
				function make() { fresh = 3 }
				make()
				count = 0
				count
			`,
			expected: 0,
		},
		{
			name: "compound assignment reaches an enclosing scope",
			input: `
				total = 0
				function add(n) { total += n }
				add(5)
				add(3)
				total
			`,
			expected: 8,
		},
		{
			name: "postfix increment reaches an enclosing scope",
			input: `
				count = 0
				function bump() { count++ }
				bump()
				bump()
				count
			`,
			expected: 2,
		},
		{
			name: "destructuring reaches an enclosing scope",
			input: `
				a = 0
				b = 0
				function unpack() { [a, b] = [7, 8] }
				unpack()
				a + b
			`,
			expected: 15,
		},
		{
			name: "assignment from inside a block reaches past the block",
			input: `
				found = 0
				if (true) { found = 5 }
				found
			`,
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isNumberObject(t, evaluate(tt.input), tt.expected)
		})
	}
}

func TestBlocksIntroduceAScope(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedMessage string
	}{
		{
			name: "a name first assigned in an if body does not outlive it",
			input: `
				if (true) { leaked = 1 }
				leaked
			`,
			expectedMessage: "test.gs:3:5: name error: `leaked` is not defined",
		},
		{
			name: "a name first assigned in an else body does not outlive it",
			input: `
				if (false) { x = 1 } else { leaked = 2 }
				leaked
			`,
			expectedMessage: "test.gs:3:5: name error: `leaked` is not defined",
		},
		{
			name: "a name first assigned in a while body does not outlive it",
			input: `
				run = true
				while (run) { leaked = 1; run = false }
				leaked
			`,
			expectedMessage: "test.gs:4:5: name error: `leaked` is not defined",
		},
		{
			name: "a name first assigned in a for body does not outlive it",
			input: `
				for (i = 0; i < 1; i++) { leaked = 1 }
				leaked
			`,
			expectedMessage: "test.gs:3:5: name error: `leaked` is not defined",
		},
		{
			name: "a name first assigned in a switch case does not outlive it",
			input: `
				switch (1) { case 1 { leaked = 1 } }
				leaked
			`,
			expectedMessage: "test.gs:3:5: name error: `leaked` is not defined",
		},
		{
			name: "a for loop's control variable does not outlive the loop",
			input: `
				for (i = 0; i < 3; i++) { }
				i
			`,
			expectedMessage: "test.gs:3:5: name error: `i` is not defined",
		},
		{
			name: "a for-in loop's control variables do not outlive the loop",
			input: `
				for (key, item in [1, 2]) { }
				item
			`,
			expectedMessage: "test.gs:3:5: name error: `item` is not defined",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isErrorObject(t, evaluate(tt.input), tt.expectedMessage)
		})
	}
}

func TestLoopVariablesDoNotDisturbTheEnclosingScope(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			name: "a for loop does not overwrite a variable of the same name",
			input: `
				i = 100
				for (i = 0; i < 3; i++) { }
				i
			`,
			expected: 100,
		},
		{
			name: "a for-in loop does not overwrite a variable of the same name",
			input: `
				item = 100
				for (key, item in [1, 2]) { }
				item
			`,
			expected: 100,
		},
		{
			name: "an accumulator declared outside the loop still accumulates",
			input: `
				total = 0
				for (n in [1, 2, 3]) { total = total + n }
				total
			`,
			expected: 6,
		},
		{
			name: "a while loop still drives a counter declared outside it",
			input: `
				i = 0
				while (i < 5) { i++ }
				i
			`,
			expected: 5,
		},
		{
			name: "the body can write the control variable and the loop sees it",
			input: `
				steps = 0
				for (i = 0; i < 10; i++) { steps++; i = i + 1 }
				steps
			`,
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isNumberObject(t, evaluate(tt.input), tt.expected)
		})
	}
}

// TestClosuresCaptureTheirIteration covers §13.14 directly, and doubles as the
// correctness test for the environment reuse that keeps block scoping from
// allocating per iteration: every case here would fail if a captured
// environment were handed to the next iteration.
func TestClosuresCaptureTheirIteration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			name: "a for-in loop's value is captured per iteration",
			input: `
				handlers = []
				for (index, name in [1, 2, 3]) {
					handlers.push(function () { return name })
				}
				handlers[0]() * 100 + handlers[1]() * 10 + handlers[2]()
			`,
			expected: 123,
		},
		{
			name: "a for-in loop's key is captured per iteration",
			input: `
				handlers = []
				for (index, name in [7, 8]) {
					handlers.push(function () { return index })
				}
				handlers[0]() * 10 + handlers[1]()
			`,
			expected: 1,
		},
		{
			name: "a for loop's control variable is captured per iteration",
			input: `
				handlers = []
				for (i = 1; i <= 3; i++) {
					handlers.push(function () { return i })
				}
				handlers[0]() * 100 + handlers[1]() * 10 + handlers[2]()
			`,
			expected: 123,
		},
		{
			name: "a local declared in a while body is captured per iteration",
			input: `
				handlers = []
				i = 1
				while (i <= 3) {
					current = i
					handlers.push(function () { return current })
					i++
				}
				handlers[0]() * 100 + handlers[1]() * 10 + handlers[2]()
			`,
			expected: 123,
		},
		{
			name: "a closure made inside a nested block captures that iteration",
			input: `
				handlers = []
				for (n in [1, 2, 3]) {
					if (true) {
						doubled = n * 2
						handlers.push(function () { return doubled })
					}
				}
				handlers[0]() * 100 + handlers[1]() * 10 + handlers[2]()
			`,
			expected: 246,
		},
		{
			name: "a closure keeps reading the variable it captured",
			input: `
				makers = []
				for (n in [1, 2]) {
					makers.push(function () { n = n + 10; return n })
				}
				makers[0]() + makers[0]() + makers[1]()
			`,
			// 11, then 21 — the same closure's own binding, mutated twice —
			// then 12 from the second iteration's separate binding.
			expected: 44,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isNumberObject(t, evaluate(tt.input), tt.expected)
		})
	}
}

// TestMethodScopingIsUnchanged pins the behavior that assignment walking
// outward must not disturb: a method's locals, its instance fields, and the
// class environment its siblings live in are three different things.
func TestMethodScopingIsUnchanged(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			name: "a method's local does not escape into the instance",
			input: `
				class Counter {
					constructor() { this.value = 1 }
					bump() { value = 99; return this.value }
				}
				new Counter().bump()
			`,
			expected: 1,
		},
		{
			name: "a method rebinds a module-level variable",
			input: `
				calls = 0
				class Tracker {
					run() { calls++ }
				}
				t = new Tracker()
				t.run()
				t.run()
				calls
			`,
			expected: 2,
		},
		{
			name: "a field assignment still reaches the instance",
			input: `
				class Box {
					constructor() { this.n = 0 }
					fill() { this.n = 42 }
				}
				b = new Box()
				b.fill()
				b.n
			`,
			expected: 42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isNumberObject(t, evaluate(tt.input), tt.expected)
		})
	}
}
