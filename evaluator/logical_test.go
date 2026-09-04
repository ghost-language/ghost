package evaluator

import (
	"testing"

	"ghostlang.org/x/ghost/object"
)

// The tests in this file cover §13.21 — `and` and `or` short-circuit, per
// §14 decision 11, which reversed §8.4's original stance that both operands
// are always evaluated.
//
// The behavior that matters is not the truth table, which is unchanged, but
// which operands get evaluated: a guard whose left side already settles the
// answer must never reach its right side. That is what makes
// `x == null or x.field` safe to write, and it is the only reason this
// change exists.

func TestLogicalOperatorsAnswerTheSameTruthTable(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{name: "true and true", input: "true and true", expected: true},
		{name: "true and false", input: "true and false", expected: false},
		{name: "false and true", input: "false and true", expected: false},
		{name: "false and false", input: "false and false", expected: false},
		{name: "true or true", input: "true or true", expected: true},
		{name: "true or false", input: "true or false", expected: true},
		{name: "false or true", input: "false or true", expected: true},
		{name: "false or false", input: "false or false", expected: false},
		{
			name:     "a computed left operand still decides",
			input:    "(1 == 2) or (3 < 4)",
			expected: true,
		},
		{
			name:     "chained operators associate left to right",
			input:    "false and true or true",
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isBooleanObject(t, evaluate(test.input), test.expected)
		})
	}
}

// TestLogicalOperatorsShortCircuit is the point of §13.21: the right operand
// is not evaluated when the left one has already settled the answer. Each
// case proves it by making evaluation of the right operand observable — it
// either raises, or it appends to a list the test then measures.
func TestLogicalOperatorsShortCircuit(t *testing.T) {
	t.Run("a null guard does not reach the dereference it guards", func(t *testing.T) {
		// The exact shape that crashed Chisel's Ui.paintTooltip() twice.
		input := `
			target = null
			target == null or target.hint == ""
		`

		isBooleanObject(t, evaluate(input), true)
	})

	t.Run("an and-guard does not reach its right side when the left is false", func(t *testing.T) {
		input := `
			target = null
			target != null and target.hint == ""
		`

		isBooleanObject(t, evaluate(input), false)
	})

	t.Run("a raising right operand is never reached", func(t *testing.T) {
		isBooleanObject(t, evaluate("false and (1 / 0) == 0"), false)
		isBooleanObject(t, evaluate("true or (1 / 0) == 0"), true)
	})

	t.Run("a side effect in the right operand does not happen", func(t *testing.T) {
		input := `
			calls = []
			function touched() { calls.push(1) return true }
			false and touched()
			true or touched()
			calls.length()
		`

		isNumberObject(t, evaluate(input), 0)
	})

	t.Run("the right operand is still evaluated when the left leaves it open", func(t *testing.T) {
		input := `
			calls = []
			function touched() { calls.push(1) return true }
			true and touched()
			false or touched()
			calls.length()
		`

		isNumberObject(t, evaluate(input), 2)
	})

	t.Run("an unreached operand is not type-checked", func(t *testing.T) {
		// The one behavior §14 decision 11 loosens: these were type errors
		// when both sides were always evaluated.
		isBooleanObject(t, evaluate("false and 1"), false)
		isBooleanObject(t, evaluate("true or 1"), true)
	})
}

// TestLogicalOperandErrors pins the wording of a non-boolean operand that IS
// reached. `and`/`or` stay boolean-only (§14 decision 11), and the fault names
// the side at fault rather than both types, since short-circuiting means the
// other side may deliberately never have been evaluated.
func TestLogicalOperandErrors(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "a null left operand, the truthy-guard mistake",
			input:    "x = null x and true",
			expected: "test.gs:1:12: type error: cannot use `and` with null on the left",
		},
		{
			name:     "a number left operand",
			input:    "1 and 2",
			expected: "test.gs:1:3: type error: cannot use `and` with number on the left",
		},
		{
			name:     "a string left operand of or",
			input:    `"a" or false`,
			expected: "test.gs:1:5: type error: cannot use `or` with string on the left",
		},
		{
			name:     "a number right operand, reached because the left leaves it open",
			input:    "true and 1",
			expected: "test.gs:1:6: type error: cannot use `and` with number on the right",
		},
		{
			name:     "a null right operand of or, reached because the left is false",
			input:    "y = null false or y",
			expected: "test.gs:1:16: type error: cannot use `or` with null on the right",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isErrorObject(t, evaluate(test.input), test.expected)
		})
	}
}

// TestLogicalOperandErrorHelp checks the help line the null case carries,
// since a null operand is nearly always a guard written in the idiom of a
// language where `and` is truthy rather than boolean.
func TestLogicalOperandErrorHelp(t *testing.T) {
	result := evaluate("x = null x and true")

	err, ok := result.(*object.Error)

	if !ok {
		t.Fatalf("object is not Error. got=%T (%+v)", result, result)
	}

	if err.Fault.Help != "compare it first, as in `x != null`" {
		t.Errorf("help has wrong text. got=%q", err.Fault.Help)
	}
}
