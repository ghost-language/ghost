package evaluator

import (
	"testing"

	"ghostlang.org/x/ghost/library/modules"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/optimizer"
	"ghostlang.org/x/ghost/parser"
	"ghostlang.org/x/ghost/scanner"
)

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"5 + true", "1:3:test.ghost: runtime error: type mismatch: NUMBER + BOOLEAN"},
		{"5 + true; 5", "1:3:test.ghost: runtime error: type mismatch: NUMBER + BOOLEAN"},
		{"-true", "1:1:test.ghost: runtime error: unknown operator: -BOOLEAN"},
		{"true + false", "1:6:test.ghost: runtime error: unknown operator: BOOLEAN + BOOLEAN"},
		{"5; true + false; 5", "1:9:test.ghost: runtime error: unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { if (10 > 1) { return true + false } return 1 }", "1:41:test.ghost: runtime error: unknown operator: BOOLEAN + BOOLEAN"},
		{"foobar", "1:1:test.ghost: runtime error: unknown identifier: foobar"},
		{`"Hello" - "World"`, "1:9:test.ghost: runtime error: unknown operator: STRING - STRING"},
		{`{"name": "Ghost"}[function() { 123 }]`, "1:18:test.ghost: runtime error: unusable as map key: FUNCTION"},
		{`function foo() { a } foo()`, "1:18:test.ghost: runtime error: unknown identifier: a"},
		{`class Test { function foo() { a } } test = Test.new() test.foo()`, "1:31:test.ghost: runtime error: unknown identifier: a"},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isErrorObject(t, result, tt.expectedMessage)
	}
}

func TestAssign(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"a = 5; a", 5},
		{"a = 5 * 5; a", 25},
		{"a = 5; b = a; b", 5},
		{"a = 5; b = a; c = a + b + 5; c", 15},
		{"a = 5; a = 10; a", 10},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isNumberObject(t, result, tt.expected)
	}
}

func TestNumbers(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"x = 5; x += 1; x", 6},
		{"x = 5; x -= 1; x", 4},
		{"x = 5; x *= 2; x", 10},
		{"x = 10; x /= 2; x", 5},
		{"x = 0; x++; x", 1},
		{"x = 6; x--; x", 5},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isNumberObject(t, result, tt.expected)
	}
}

func TestClassStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`class Foo {}`, "Foo"},
		{`class Foo {
			function bar() {
				true
			}
		}`, "Foo"},
	}

	for _, tt := range tests {
		evaluated := evaluate(tt.input)

		isClassObject(t, evaluated, tt.expected)
	}
}

func TestForExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`x = 10; for (x = y; x > 0; x = x - 1) { x }`, "1:18:test.ghost: runtime error: unknown identifier: y"},
		{`for (x = 0; x < 10; x = x + 1) { y }`, "1:34:test.ghost: runtime error: unknown identifier: y"},
		{`bar = true; for (x = 0; x < 10; x = x + 1) { y; print(bar) }`, "1:46:test.ghost: runtime error: unknown identifier: y"},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)
		number, ok := tt.expected.(int64)

		if ok {
			isNumberObject(t, result, number)
		} else {
			isErrorObject(t, result, tt.expected.(string))
		}
	}
}

func TestForInExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`list = [1, 2, 3]; for(x in lists) { x }`, "1:28:test.ghost: runtime error: unknown identifier: lists"},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)
		number, ok := tt.expected.(int64)

		if ok {
			isNumberObject(t, result, number)
		} else {
			isErrorObject(t, result, tt.expected.(string))
		}
	}
}

func TestRangeExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`1 .. 0`, []int{}},
		{`-1 .. 0`, []int{-1, 0}},
		{`1 .. 1`, []int{1}},
		{`1 .. 5`, []int{1, 2, 3, 4, 5}},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		list, ok := result.(*object.List)

		if !ok {
			t.Errorf("object not List. got=%T (+%v)", result, result)
		}

		if len(list.Elements) != len(tt.expected.([]int)) {
			t.Errorf("wrong number of elements. wanted=%d, got=%d", len(tt.expected.([]int)), len(list.Elements))
		}
	}
}

func TestWhileExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`while (false) { }`, nil},
		{`n = 0; while (n < 10) { n = n + 1 }; n`, 10},
		{"n = 10; while (n > 0) { n = n - 1 }; n", 0},
		{"n = 0; while (n < 10) { n = n + 1 }", nil},
		{"n = 10; while (n > 0) { n = n - 1 }", nil},
		{"while (true) { break }", nil},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)
		number, ok := tt.expected.(int)

		if ok {
			isNumberObject(t, result, int64(number))
		} else {
			isNil(t, result)
		}
	}
}

func TestClassProperties(t *testing.T) {
	input := `
	class Circle {
		function constructor(area) {
			this.area = area
		}

		function area() {
			return math.pi * this.area * this.area
		}
	}

	test = Circle.new(10)

	return test.area()
	`

	result := evaluate(input)

	isNumberObject(t, result, 314)
}

func TestTraitMethodLookup(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			name: "method from single trait",
			input: `
				trait Greet {
					function greet() {
						return 1
					}
				}

				class Person {
					use Greet
				}

				p = Person.new()
				p.greet()
			`,
			expected: 1,
		},
		{
			name: "method from second trait",
			input: `
				trait First {
					function first() {
						return 1
					}
				}

				trait Second {
					function second() {
						return 2
					}
				}

				class Thing {
					use First
					use Second
				}

				t = Thing.new()
				t.second()
			`,
			expected: 2,
		},
		{
			name: "methods from both traits",
			input: `
				trait Add {
					function add() {
						return 10
					}
				}

				trait Multiply {
					function multiply() {
						return 20
					}
				}

				class Calculator {
					use Add
					use Multiply
				}

				c = Calculator.new()
				c.add() + c.multiply()
			`,
			expected: 30,
		},
		{
			name: "comma-separated traits",
			input: `
				trait First {
					function first() {
						return 1
					}
				}

				trait Second {
					function second() {
						return 2
					}
				}

				trait Third {
					function third() {
						return 3
					}
				}

				class Thing {
					use First, Second, Third
				}

				t = Thing.new()
				t.first() + t.second() + t.third()
			`,
			expected: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluate(tt.input)
			isNumberObject(t, result, tt.expected)
		})
	}
}

func TestInheritedProperties(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			name: "access parent class property",
			input: `
				class Animal {
					species = 42
				}

				class Dog extends Animal {
				}

				d = Dog.new()
				d.species
			`,
			expected: 42,
		},
		{
			name: "access grandparent class property",
			input: `
				class Animal {
					legs = 4
				}

				class Mammal extends Animal {
				}

				class Dog extends Mammal {
				}

				d = Dog.new()
				d.legs
			`,
			expected: 4,
		},
		{
			name: "child property overrides parent",
			input: `
				class Animal {
					sound = 1
				}

				class Dog extends Animal {
					sound = 2
				}

				d = Dog.new()
				d.sound
			`,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluate(tt.input)
			isNumberObject(t, result, tt.expected)
		})
	}
}

func TestInheritedConstructor(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			name: "parent constructor called when child has none",
			input: `
				class Animal {
					function constructor(legs) {
						this.legs = legs
					}
				}

				class Dog extends Animal {
				}

				d = Dog.new(4)
				d.legs
			`,
			expected: 4,
		},
		{
			name: "grandparent constructor called when child and parent have none",
			input: `
				class Animal {
					function constructor(legs) {
						this.legs = legs
					}
				}

				class Mammal extends Animal {
				}

				class Dog extends Mammal {
				}

				d = Dog.new(4)
				d.legs
			`,
			expected: 4,
		},
		{
			name: "child constructor overrides parent",
			input: `
				class Animal {
					function constructor(legs) {
						this.legs = legs
					}
				}

				class Dog extends Animal {
					function constructor(legs) {
						this.legs = legs * 2
					}
				}

				d = Dog.new(4)
				d.legs
			`,
			expected: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluate(tt.input)
			isNumberObject(t, result, tt.expected)
		})
	}
}

func TestThisOutsideClass(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "this at top level",
			input:    `this`,
			expected: "1:1:test.ghost: runtime error: 'this' used outside of class context",
		},
		{
			name:     "this in regular function",
			input:    `function foo() { return this } foo()`,
			expected: "1:25:test.ghost: runtime error: 'this' used outside of class context",
		},
		{
			name:     "this.property at top level",
			input:    `this.name`,
			expected: "1:1:test.ghost: runtime error: 'this' used outside of class context",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluate(tt.input)
			isErrorObject(t, result, tt.expected)
		})
	}
}

// =============================================================================
// Helper functions

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!null", true},
		{"!!true", true},
		{"!!false", false},

		// Comparisons and library functions produce boolean objects that are
		// not the shared TRUE/FALSE singletons. The bang operator must compare
		// them by value, not by pointer identity.
		{`!("a" != "a")`, true},
		{`!("a" == "a")`, false},
		{`!("a" == "b")`, true},
		{"!(1 != 1)", true},
		{"!(1 == 1)", false},
		{`!("abc".startsWith("x"))`, true},
		{`!("abc".startsWith("a"))`, false},

		// Non-boolean, non-null operands remain falsy under bang.
		{"!5", false},
		{`!"abc"`, false},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		boolean, ok := result.(*object.Boolean)

		if !ok {
			t.Fatalf("evaluate(%q) is not object.Boolean. got=%T (%+v)", tt.input, result, result)
		}

		if boolean.Value != tt.expected {
			t.Errorf("evaluate(%q) is not %t. got=%t", tt.input, tt.expected, boolean.Value)
		}
	}
}

// TestOptimizerPreservesSemantics evaluates each program twice, once on the raw
// AST and once on the optimized AST, and requires the two to agree. Constant
// folding is only correct if it is invisible, including for the expressions it
// deliberately refuses to fold because they raise runtime errors.
func TestOptimizerPreservesSemantics(t *testing.T) {
	programs := []string{
		// Arithmetic, including the integer/float promotion rules.
		"1 + 2", "10 - 4", "6 * 7", "7 % 3", "-5", "-(3 * 4)",
		"2 * 3 + 4 * 5 - 6", "1.5 + 2.5", "1 + 2.5", "6 / 3", "7 / 2",
		"2147483647 * 2147483647", "-9223372036854775807 - 1",

		// Comparisons and booleans.
		"1 < 2", "2 <= 2", "3 > 4", "1 == 1", "1 != 1", "1.0 == 1",
		`"a" == "a"`, `"a" < "b"`, `"a" != "b"`,
		"true and false", "true or false", "true == true", "false != true",
		"!true", "!false", "!null", "!(1 == 2)", `!("a" != "a")`,

		// Strings.
		`"hello" + " " + "world"`, `"" + "a"`,

		// Expressions the optimizer must leave alone: each raises a runtime
		// error whose message and position must survive.
		"1 / 0", "1 % 0", "1.0 / 0.0", "1 + true", `1 + "a"`, "-true",
		"true + false", `"a" - "b"`,

		// Ranges build list objects and must not be folded.
		"1 .. 5", "(1 + 1) .. (2 + 3)",

		// Folding inside larger constructs.
		`total = 0
		 for (i = 0; i < 2 * 5; i = i + 1) { total = total + (3 * 4) }
		 total`,
		`function f(a = 2 * 3) { return a + (1 + 1) } f()`,
		`if (1 + 1 == 2) { "yes" } else { "no" }`,
		`x = 5 while (x > 5 - 5) { x = x - 1 } x`,
		`[1 + 1, 2 * 2, "a" + "b"]`,
		`m = {"k": 3 * 3} m["k"]`,
		`(1 == 1) ? "t" : "f"`,

		// Library globals are classified by the optimizer, so confirm modules,
		// functions, and their precedence over scope bindings all survive it.
		`math.pi > 3`,
		`math.floor(3.7)`,
		`type(5)`,
		`type("a")`,
		`type([1])`,
		`math.abs(-4)`,
		`x = math.pi x > 3`,
		// A library name still wins over a same-named variable, as before.
		`type = 5 type("a")`,
	}

	for _, source := range programs {
		plain := evaluate(source)
		optimized := evaluateOptimized(source)

		if (plain == nil) != (optimized == nil) {
			t.Errorf("%q: nil mismatch: plain=%v optimized=%v", source, plain, optimized)
			continue
		}

		if plain == nil {
			continue
		}

		if plain.Type() != optimized.Type() {
			t.Errorf("%q: type mismatch: plain=%s optimized=%s", source, plain.Type(), optimized.Type())
			continue
		}

		if plain.String() != optimized.String() {
			t.Errorf("%q: value mismatch: plain=%q optimized=%q", source, plain.String(), optimized.String())
		}
	}
}

func evaluateOptimized(input string) object.Object {
	scope := &object.Scope{
		Environment: object.NewEnvironment(),
	}

	object.RegisterEvaluator(Evaluate)
	modules.RegisterEvaluator(Evaluate)

	p := parser.New(scanner.New(input, "test.ghost"))

	return Evaluate(optimizer.Optimize(p.Parse()), scope)
}

func evaluate(input string) object.Object {
	scope := &object.Scope{
		Environment: object.NewEnvironment(),
	}

	evaluatorInstance := Evaluate

	object.RegisterEvaluator(evaluatorInstance)
	modules.RegisterEvaluator(evaluatorInstance)

	scanner := scanner.New(input, "test.ghost")
	parser := parser.New(scanner)
	program := parser.Parse()

	result := Evaluate(program, scope)

	return result
}

func isErrorObject(t *testing.T, obj object.Object, expected string) bool {
	err, ok := obj.(*object.Error)

	if !ok {
		t.Errorf("object is not Error. got=%T (%+v", obj, obj)
		return false
	}

	if err.Message != expected {
		t.Errorf("error has wrong message. got=%s, expected=%s", err.Message, expected)
		return false
	}

	return true
}

func isNumberObject(t *testing.T, obj object.Object, expected int64) bool {
	number, ok := obj.(*object.Number)

	if !ok {
		t.Errorf("object is not Number. got=%T (%+v", obj, obj)
		return false
	}

	if number.Int64() != expected {
		t.Errorf("object has wrong value. got=%d, expected=%d", number.Int64(), expected)
		return false
	}

	return true
}

func isNil(t *testing.T, obj object.Object) bool {
	if obj != nil {
		t.Errorf("object is not nil. got=%T (%+v", obj, obj)
		return false
	}

	return true
}

func isClassObject(t *testing.T, obj object.Object, expected string) bool {
	class, ok := obj.(*object.Class)

	if !ok {
		t.Errorf("object is not Class. got=%T (%+v", obj, obj)
		return false
	}

	if class.Name.Value != expected {
		t.Errorf("class has wrong name. got=%s, expected=%s", class.Name.Value, expected)
		return false
	}

	return true
}
