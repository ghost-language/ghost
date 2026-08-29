package evaluator

import "testing"

// TestNumberMethods covers the methods closing §12's number-method gap:
// ceil() (round/floor already existed) and, for parity with the math
// module, abs()/pow()/clamp() and the isX predicates as instance methods.
func TestNumberMethods(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`2.4.ceil()`, 3},
		{`2.6.ceil()`, 3},
		{`(-2.4).ceil()`, -2},
		{`3.ceil()`, 3},

		{`5.abs()`, 5},
		{`(-5).abs()`, 5},
		{`0.abs()`, 0},

		{`2.pow(10)`, 1024},
		{`3.pow(0)`, 1},
		{`(-2).pow(3)`, -8},

		{`5.clamp(1, 10)`, 5},
		{`(-5).clamp(1, 10)`, 1},
		{`50.clamp(1, 10)`, 10},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isNumberObject(t, result, tt.expected)
	}
}

// TestNumberMethodFloats checks the methods above where the result carries a
// fractional part, comparing against the rendered string since isNumberObject
// truncates through Int64().
func TestNumberMethodFloats(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`2.5.pow(2)`, "6.25"},
		{`2.5.abs()`, "2.5"},
		{`(-2.5).abs()`, "2.5"},
		{`2.5.clamp(0, 2)`, "2"},
		{`0.5.clamp(1.5, 3.5)`, "1.5"},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		if result.String() != tt.expected {
			t.Errorf("wrong result for %q. got=%s, expected=%s", tt.input, result.String(), tt.expected)
		}
	}
}

func TestNumberMethodBooleans(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`import "ghost:math" math.nan.isNaN()`, true},
		{`1.isNaN()`, false},

		{`import "ghost:math" math.infinity.isInfinite()`, true},
		{`1.isInfinite()`, false},

		{`1.isFinite()`, true},
		{`import "ghost:math" math.infinity.isFinite()`, false},
		{`import "ghost:math" math.nan.isFinite()`, false},

		{`1.isInteger()`, true},
		{`1.5.isInteger()`, false},
		{`2.0.isInteger()`, true},

		{`2.isEven()`, true},
		{`3.isEven()`, false},
		{`3.isOdd()`, true},
		{`2.isOdd()`, false},

		{`(-1).isNegative()`, true},
		{`1.isNegative()`, false},
		{`1.isPositive()`, true},
		{`(-1).isPositive()`, false},
		{`0.isZero()`, true},
		{`1.isZero()`, false},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isBooleanObject(t, result, tt.expected)
	}
}

func TestNumberMethodErrors(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{`1.pow()`, "test.gs:1:3: argument error: `number.pow()` expects 1 argument, got 0"},
		{`1.clamp(1)`, "test.gs:1:3: argument error: `number.clamp()` expects 2 arguments, got 1"},
		{`5.clamp(10, 1)`, "test.gs:1:3: value error: `number.clamp()` expects the lower bound to be no greater than the upper bound"},
		{`1.ceil(1)`, "test.gs:1:3: argument error: `number.ceil()` expects 0 arguments, got 1"},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isErrorObject(t, result, tt.expectedMessage)
	}
}
