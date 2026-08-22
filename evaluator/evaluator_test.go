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
		{"5 + true", "test.ghost:1:3: type error: cannot use `+` between number and boolean"},
		{"5 + true; 5", "test.ghost:1:3: type error: cannot use `+` between number and boolean"},
		{"-true", "test.ghost:1:1: type error: cannot negate boolean"},
		{"true + false", "test.ghost:1:6: type error: cannot use `+` between two booleans"},
		{"5; true + false; 5", "test.ghost:1:9: type error: cannot use `+` between two booleans"},
		{"if (10 > 1) { if (10 > 1) { return true + false } return 1 }", "test.ghost:1:41: type error: cannot use `+` between two booleans"},
		{"foobar", "test.ghost:1:1: name error: `foobar` is not defined"},
		{`"Hello" - "World"`, "test.ghost:1:9: type error: cannot use `-` between two strings"},
		{`{"name": "Ghost"}[function() { 123 }]`, "test.ghost:1:18: type error: function cannot be used as a map key"},
		{`function foo() { a } foo()`, "test.ghost:1:18: name error: `a` is not defined"},
		{`class Test { function foo() { a } } test = new Test() test.foo()`, "test.ghost:1:31: name error: `a` is not defined"},
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
			bar() {
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
		{`x = 10; for (x = y; x > 0; x = x - 1) { x }`, "test.ghost:1:18: name error: `y` is not defined"},
		{`for (x = 0; x < 10; x = x + 1) { y }`, "test.ghost:1:34: name error: `y` is not defined"},
		{`bar = true; for (x = 0; x < 10; x = x + 1) { y; print(bar) }`, "test.ghost:1:46: name error: `y` is not defined"},
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
		{`list = [1, 2, 3]; for(x in lists) { x }`, "test.ghost:1:28: name error: `lists` is not defined"},
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
		constructor(area) {
			this.area = area
		}

		area() {
			return math.pi * this.area * this.area
		}
	}

	test = new Circle(10)

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
					greet() {
						return 1
					}
				}

				class Person {
					use Greet
				}

				p = new Person()
				p.greet()
			`,
			expected: 1,
		},
		{
			name: "method from second trait",
			input: `
				trait First {
					first() {
						return 1
					}
				}

				trait Second {
					second() {
						return 2
					}
				}

				class Thing {
					use First
					use Second
				}

				t = new Thing()
				t.second()
			`,
			expected: 2,
		},
		{
			name: "methods from both traits",
			input: `
				trait Add {
					add() {
						return 10
					}
				}

				trait Multiply {
					multiply() {
						return 20
					}
				}

				class Calculator {
					use Add
					use Multiply
				}

				c = new Calculator()
				c.add() + c.multiply()
			`,
			expected: 30,
		},
		{
			name: "comma-separated traits",
			input: `
				trait First {
					first() {
						return 1
					}
				}

				trait Second {
					second() {
						return 2
					}
				}

				trait Third {
					third() {
						return 3
					}
				}

				class Thing {
					use First, Second, Third
				}

				t = new Thing()
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

				d = new Dog()
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

				d = new Dog()
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

				d = new Dog()
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
					constructor(legs) {
						this.legs = legs
					}
				}

				class Dog extends Animal {
				}

				d = new Dog(4)
				d.legs
			`,
			expected: 4,
		},
		{
			name: "grandparent constructor called when child and parent have none",
			input: `
				class Animal {
					constructor(legs) {
						this.legs = legs
					}
				}

				class Mammal extends Animal {
				}

				class Dog extends Mammal {
				}

				d = new Dog(4)
				d.legs
			`,
			expected: 4,
		},
		{
			name: "child constructor overrides parent",
			input: `
				class Animal {
					constructor(legs) {
						this.legs = legs
					}
				}

				class Dog extends Animal {
					constructor(legs) {
						this.legs = legs * 2
					}
				}

				d = new Dog(4)
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
			expected: "test.ghost:1:1: name error: `this` can only be used inside a class",
		},
		{
			name:     "this in regular function",
			input:    `function foo() { return this } foo()`,
			expected: "test.ghost:1:25: name error: `this` can only be used inside a class",
		},
		{
			name:     "this.property at top level",
			input:    `this.name`,
			expected: "test.ghost:1:1: name error: `this` can only be used inside a class",
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

	if err.String() != expected {
		t.Errorf("error has wrong message. got=%s, expected=%s", err.String(), expected)
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

func TestNewExpression(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			name: "instantiation without arguments",
			input: `
				class Counter {
					value() { return 7 }
				}

				new Counter().value()
			`,
			expected: 7,
		},
		{
			name: "instantiation with arguments",
			input: `
				class Point {
					constructor(x, y) {
						this.x = x
						this.y = y
					}
				}

				p = new Point(3, 4)
				p.x + p.y
			`,
			expected: 7,
		},
		{
			name: "instantiation without parentheses",
			input: `
				class Counter {
					value() { return 7 }
				}

				c = new Counter
				c.value()
			`,
			expected: 7,
		},
		{
			name: "method call chains off the instance, not the class expression",
			input: `
				class Counter {
					value() { return 7 }
				}

				new Counter().value()
			`,
			expected: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluate(tt.input)
			isNumberObject(t, result, tt.expected)
		})
	}
}

func TestSuperMethodCalls(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			name: "super calls the overridden parent method",
			input: `
				class Animal {
					legs() { return 4 }
				}

				class Dog extends Animal {
					legs() { return super.legs() + 1 }
				}

				new Dog().legs()
			`,
			expected: 5,
		},
		{
			name: "super resolves from the declaring class, not the receiver",
			input: `
				class A {
					value() { return 1 }
				}

				class B extends A {
					value() { return super.value() + 10 }
				}

				class C extends B {
					value() { return super.value() + 100 }
				}

				new C().value()
			`,
			expected: 111,
		},
		{
			name: "super skips a level that does not override",
			input: `
				class A {
					value() { return 1 }
				}

				class B extends A {
				}

				class C extends B {
					value() { return super.value() + 10 }
				}

				new C().value()
			`,
			expected: 11,
		},
		{
			name: "super.constructor chains initialization",
			input: `
				class Animal {
					constructor(legs) {
						this.legs = legs
					}
				}

				class Dog extends Animal {
					constructor() {
						super.constructor(4)
						this.tails = 1
					}
				}

				d = new Dog()
				d.legs + d.tails
			`,
			expected: 5,
		},
		{
			name: "an inherited method still dispatches to the override",
			input: `
				class Animal {
					legs() { return 4 }
					total() { return this.legs() }
				}

				class Dog extends Animal {
					legs() { return super.legs() * 2 }
				}

				new Dog().total()
			`,
			expected: 8,
		},
		{
			name: "super reads a parent property",
			input: `
				class Animal {
					sound() { return 1 }
				}

				class Dog extends Animal {
					sound() { return 2 }
					parentSound() { return super.sound() }
				}

				new Dog().parentSound()
			`,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluate(tt.input)
			isNumberObject(t, result, tt.expected)
		})
	}
}

func TestPerInstanceFields(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			name: "instances do not share a mutable field",
			input: `
				class Bag {
					items = []

					add(item) { this.items.push(item) }
				}

				a = new Bag()
				b = new Bag()

				a.add(1)
				a.add(2)

				a.items.length() * 10 + b.items.length()
			`,
			expected: 20,
		},
		{
			name: "instances do not share a scalar field",
			input: `
				class Counter {
					count = 0

					bump() { this.count = this.count + 1 }
				}

				a = new Counter()
				b = new Counter()

				a.bump()
				a.bump()

				a.count * 10 + b.count
			`,
			expected: 20,
		},
		{
			name: "field initializers run parent-first",
			input: `
				class Animal {
					legs = 4
					total = 4
				}

				class Dog extends Animal {
					total = 5
				}

				d = new Dog()
				d.legs * 10 + d.total
			`,
			expected: 45,
		},
		{
			name: "trait fields are per-instance too",
			input: `
				trait Countable {
					count = 0

					bump() { this.count = this.count + 1 }
				}

				class Widget {
					use Countable
				}

				a = new Widget()
				b = new Widget()

				a.bump()

				a.count * 10 + b.count
			`,
			expected: 10,
		},
		{
			name: "a subclass field overrides the parent declaration",
			input: `
				class Animal {
					sound = 1
				}

				class Dog extends Animal {
					sound = 2
				}

				new Dog().sound
			`,
			expected: 2,
		},
		{
			name: "a field initializer can reference the enclosing scope",
			input: `
				base = 40

				class Thing {
					value = base + 2
				}

				new Thing().value
			`,
			expected: 42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluate(tt.input)
			isNumberObject(t, result, tt.expected)
		})
	}
}

func TestClassMemberResolution(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			name: "a method calls a sibling method by bare name",
			input: `
				class Math {
					double(n) { return n * 2 }
					quadruple(n) { return double(double(n)) }
				}

				new Math().quadruple(3)
			`,
			expected: 12,
		},
		{
			name: "a closure inside a method keeps the receiver",
			input: `
				class Counter {
					count = 21

					doubled() {
						f = function() { return this.count * 2 }

						return f()
					}
				}

				new Counter().doubled()
			`,
			expected: 42,
		},
		{
			name: "a class method beats a trait method of the same name",
			input: `
				trait T {
					value() { return 1 }
				}

				class A {
					use T

					value() { return 2 }
				}

				new A().value()
			`,
			expected: 2,
		},
		{
			name: "a trait on a subclass beats the superclass method",
			input: `
				trait T {
					value() { return 2 }
				}

				class A {
					value() { return 1 }
				}

				class B extends A {
					use T
				}

				new B().value()
			`,
			expected: 2,
		},
		{
			name: "an instance field does not leak in from an enclosing scope",
			input: `
				name = 99

				class Thing {
					value() {
						if (this.name) {
							return 0
						}

						return 1
					}
				}

				new Thing().value()
			`,
			expected: 1,
		},
		{
			name: "methods declared with the function keyword still work",
			input: `
				class Legacy {
					function value() { return 5 }
				}

				new Legacy().value()
			`,
			expected: 5,
		},
		{
			name: "default parameters work on shorthand methods",
			input: `
				class Greeter {
					value(n = 6) { return n }
				}

				new Greeter().value()
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

// TestClassRuntimeErrors covers the cases that previously crashed the
// interpreter by returning a bare Go nil into an expression.
func TestClassRuntimeErrors(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "property access on a class",
			input:    "class A { value() { return 1 } }\nA.value",
			expected: "test.ghost:2:2: property error: class `A` has no property `value` to read on the class itself",
		},
		{
			name:     "method call on a class",
			input:    "class A { value() { return 1 } }\nA.value()",
			expected: "test.ghost:2:3: property error: class `A` has no method `value` to call on the class itself",
		},
		{
			name:     "unknown method on a primitive",
			input:    "5.nope()",
			expected: "test.ghost:1:3: property error: number has no method `nope`",
		},
		{
			name:     "extending an undefined identifier",
			input:    "class A extends Nope {}",
			expected: "test.ghost:1:17: name error: `Nope` is not defined",
		},
		{
			name:     "extending a non-class",
			input:    "x = 5\nclass A extends x {}",
			expected: "test.ghost:2:17: type error: cannot extend `x`, which is a number, not a class",
		},
		{
			name:     "using a non-trait",
			input:    "x = 5\nclass A { use x }",
			expected: "test.ghost:2:15: type error: cannot use `x`, which is a number, not a trait",
		},
		{
			name:     "declaring constructor as a field",
			input:    "class A { constructor = 5 }",
			expected: "test.ghost:1:23: syntax error: `constructor` has to be declared as a method, not a field",
		},
		{
			name:     "instantiating a non-class",
			input:    "x = 5\nnew x()",
			expected: "test.ghost:2:1: type error: cannot instantiate number, which is not a class",
		},
		{
			name:     "super outside a class",
			input:    "super",
			expected: "test.ghost:1:1: name error: `super` can only be used inside a class",
		},
		{
			name:     "super in a class with no superclass",
			input:    "class A { value() { return super.value() } }\nnew A().value()",
			expected: "test.ghost:1:28: name error: class `A` has no superclass",
		},
		{
			name:     "calling an undefined method",
			input:    "class A {}\nnew A().nope()",
			expected: "test.ghost:2:9: property error: class `A` has no method `nope`",
		},
		{
			name:     "calling a field that is not a function",
			input:    "class A { value = 5 }\nnew A().value()",
			expected: "test.ghost:2:9: property error: class `A` has no method `value`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluate(tt.input)
			isErrorObject(t, result, tt.expected)
		})
	}
}

func TestEqualityComparisons(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"null equals null", `null == null`, true},
		{"null does not differ from null", `null != null`, false},
		{"an unset variable equals null", `x = null  x == null`, true},
		{"a value does not equal null", `5 == null`, false},
		{"a value differs from null", `5 != null`, true},
		{"null compares from the left", `null != 5`, true},
		{
			name:     "an instance equals itself",
			input:    "class A {}\na = new A()\na == a",
			expected: true,
		},
		{
			name:     "an instance equals another name for itself",
			input:    "class A {}\na = new A()\nb = a\na == b",
			expected: true,
		},
		{
			name:     "instances compare by identity, not by field values",
			input:    "class A { constructor(n) { this.n = n } }\nnew A(1) == new A(1)",
			expected: false,
		},
		{
			name:     "distinct instances differ",
			input:    "class A {}\nnew A() != new A()",
			expected: true,
		},
		{
			name:     "an instance does not equal null",
			input:    "class A {}\nnew A() == null",
			expected: false,
		},
		{
			name:     "an unset property reads as null",
			input:    "class A {}\nnew A().missing == null",
			expected: true,
		},
		{
			name:     "a set property does not read as null",
			input:    "class A { value = 5 }\nnew A().value == null",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluate(tt.input)

			boolean, ok := result.(*object.Boolean)

			if !ok {
				t.Fatalf("object is not Boolean. got=%T (%+v)", result, result)
			}

			if boolean.Value != tt.expected {
				t.Fatalf("object has wrong value. got=%t, expected=%t", boolean.Value, tt.expected)
			}
		})
	}
}

// TestEqualityTypeMismatch confirms that comparing two unrelated non-null types
// is still an error rather than silently false.
func TestEqualityTypeMismatch(t *testing.T) {
	result := evaluate(`1 == "a"`)

	isErrorObject(t, result, "test.ghost:1:3: type error: cannot use `==` between number and string")
}
