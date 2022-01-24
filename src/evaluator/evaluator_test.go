package evaluator

import (
	"testing"

	"ghostlang.org/x/ghost/object"
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

// =============================================================================
// Helper functions

func evaluate(input string) object.Object {
	scanner := scanner.New(input, "test.ghost")
	parser := parser.New(scanner)
	program := parser.Parse()
	scope := &object.Scope{
		Environment: object.NewEnvironment(),
	}

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

	if number.Value.IntPart() != expected {
		t.Errorf("object has wrong value. got=%d, expected=%d", number.Value.IntPart(), expected)
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
