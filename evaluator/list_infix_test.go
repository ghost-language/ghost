package evaluator

import (
	"testing"

	"ghostlang.org/x/ghost/object"
)

// TestListArithmetic covers the rule that makes list maths read like maths: a
// number spreads across a list, two lists pair off, and a shorter shape
// stretches across a longer one.
func TestListArithmetic(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// A number against a list, either way round.
		{"([1, 2, 3] + 10).toString()", "[11, 12, 13]"},
		{"(10 + [1, 2, 3]).toString()", "[11, 12, 13]"},
		{"([1, 2, 3] * 2).toString()", "[2, 4, 6]"},
		{"(2 * [1, 2, 3]).toString()", "[2, 4, 6]"},
		{"([1, 2, 3] - 1).toString()", "[0, 1, 2]"},
		{"([10, 20] / 4).toString()", "[2.5, 5]"},
		{"([10, 20] % 3).toString()", "[1, 2]"},

		// Two lists, paired off.
		{"([1, 2, 3] + [10, 20, 30]).toString()", "[11, 22, 33]"},
		{"([1, 2] * [3, 4]).toString()", "[3, 8]"},

		// Matrices, elementwise.
		{"([[1, 2], [3, 4]] * 2).toString()", "[[2, 4], [6, 8]]"},
		{"([[1, 2], [3, 4]] + [[10, 20], [30, 40]]).toString()", "[[11, 22], [33, 44]]"},

		// A row stretches down a matrix, and a column across it. This is the
		// case that makes a bias vector work against a batch of samples, and
		// the one a positional pairing rule gets wrong.
		{"([[1, 2], [3, 4]] + [10, 20]).toString()", "[[11, 22], [13, 24]]"},
		{"([[1, 2], [3, 4]] + [[10], [20]]).toString()", "[[11, 12], [23, 24]]"},

		// An axis of length one repeats, so a one-element list acts as a
		// number does.
		{"([1, 2] + [5]).toString()", "[6, 7]"},

		// Whole numbers stay whole.
		{"([1, 2] + 1).toString()", "[2, 3]"},

		// Compound assignment routes through the same operators.
		{"a = [1, 2, 3]; a += 10; a.toString()", "[11, 12, 13]"},
		{"a = [1, 2, 3]; a *= [2, 2, 2]; a.toString()", "[2, 4, 6]"},
		{"a = [[1, 2]]; a -= 1; a.toString()", "[[0, 1]]"},

		// Empty lists have a shape too.
		{"([] + []).toString()", "[]"},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isStringObject(t, result, tt.expected)
	}
}

// TestListEquality covers `==` comparing contents rather than identity.
func TestListEquality(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"[1, 2] == [1, 2]", true},
		{"[1, 2] == [1, 3]", false},
		{"[1, 2] == [1, 2, 3]", false},
		{"[1, 2] != [1, 3]", true},
		{"[1, 2] != [1, 2]", false},
		{"[[1], [2]] == [[1], [2]]", true},
		{"[[1], [2]] == [[1], [3]]", false},
		{`["a", "b"] == ["a", "b"]`, true},
		{`["a", "b"] == ["a", "c"]`, false},
		{`[1, "a", true] == [1, "a", true]`, true},
		{"[] == []", true},
		{"[null] == [null]", true},
		{"[1] == [true]", false},

		// A computed list equals a written one, which is the point.
		{"a = [1, 2] + [1, 2]; a == [2, 4]", true},

		// The operators and the math module's methods are the same operation
		// reached two ways, so they have to agree exactly - including on the
		// broadcasting rule.
		{"[1, 2] + [3, 4] == math.add([1, 2], [3, 4])", true},
		{"[[1, 2], [3, 4]] + [10, 20] == math.add([[1, 2], [3, 4]], [10, 20])", true},
		{"[[1, 2], [3, 4]] * 2 == math.multiply([[1, 2], [3, 4]], 2)", true},
		{"[10, 20] / 4 == math.divide([10, 20], 4)", true},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isBooleanObject(t, result, tt.expected)
	}
}

// TestListConcat covers joining, which is a method because the operators are
// arithmetic.
func TestListConcat(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"[1, 2].concat([3, 4]).toString()", "[1, 2, 3, 4]"},
		{`["a"].concat(["b"]).toString()`, "[a, b]"},
		{"[].concat([1]).toString()", "[1]"},
		{"[1].concat([]).toString()", "[1]"},

		// The operands are left alone.
		{"a = [1, 2]; a.concat([3]); a.toString()", "[1, 2]"},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isStringObject(t, result, tt.expected)
	}
}

func TestListOperatorErrors(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"[1, 2] + [1, 2, 3]", "1:8:test.ghost: runtime error: cannot evaluate LIST + LIST: shapes 2 and 3 cannot be combined"},
		{`["a"] + [1]`, "1:7:test.ghost: runtime error: cannot evaluate LIST + LIST: expected a number or a list of numbers, found STRING"},
		{"[[1, 2], [3]] + 1", "1:15:test.ghost: runtime error: cannot evaluate LIST + NUMBER: lists have to be rectangular to combine elementwise"},
		{"[1, 2] / 0", "1:8:test.ghost: runtime error: division by zero"},
		{"[1, 2] % 0", "1:8:test.ghost: runtime error: division by zero"},

		// A list against something that is neither a list nor a number stays a
		// type mismatch, which says more than a broadcasting failure would.
		{"[1, 2] + true", "1:8:test.ghost: runtime error: type mismatch: LIST + BOOLEAN"},
		{`"a" + [1]`, "1:5:test.ghost: runtime error: type mismatch: STRING + LIST"},

		// Ordering between two lists has no obvious reading, so it stays
		// unsupported rather than guessing at one.
		{"[1, 2] < [3, 4]", "1:8:test.ghost: runtime error: unknown operator: LIST < LIST"},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isErrorObject(t, result, tt.expectedMessage)
	}
}

func isStringObject(t *testing.T, obj object.Object, expected string) bool {
	t.Helper()

	str, ok := obj.(*object.String)

	if !ok {
		t.Errorf("object is not String. got=%T (%+v)", obj, obj)
		return false
	}

	if str.Value != expected {
		t.Errorf("object has wrong value. got=%s, expected=%s", str.Value, expected)
		return false
	}

	return true
}

func isBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	t.Helper()

	boolean, ok := obj.(*object.Boolean)

	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}

	if boolean.Value != expected {
		t.Errorf("object has wrong value. got=%t, expected=%t", boolean.Value, expected)
		return false
	}

	return true
}
