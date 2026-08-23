package evaluator

import "testing"

// TestDateComparisons covers the operators Date supports directly - ordering
// and exact-instant equality - which is what lets `d1 < d2` read the way
// date-fns's isBefore(d1, d2) does, without needing a function for it.
func TestDateComparisons(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"date.of(2024, 1, 1) < date.of(2024, 1, 2)", true},
		{"date.of(2024, 1, 2) < date.of(2024, 1, 1)", false},
		{"date.of(2024, 1, 2) > date.of(2024, 1, 1)", true},
		{"date.of(2024, 1, 1) <= date.of(2024, 1, 1)", true},
		{"date.of(2024, 1, 1) >= date.of(2024, 1, 1)", true},
		{"date.of(2024, 1, 1) == date.of(2024, 1, 1)", true},
		{"date.of(2024, 1, 1, 9, 0, 0) == date.of(2024, 1, 1, 9, 0, 1)", false},
		{"date.of(2024, 1, 1) != date.of(2024, 1, 2)", true},
		{"date.now() > date.fromUnix(0)", true},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isBooleanObject(t, result, tt.expected)
	}
}

func TestDateArithmeticOperatorsAreRejected(t *testing.T) {
	result := evaluate("date.now() + date.now()")

	isErrorObject(t, result, "test.ghost:1:12: type error: cannot use `+` between two dates")
}

func TestDateToString(t *testing.T) {
	result := evaluate("date.of(2024, 6, 15, 9, 30, 0).toString()")

	isStringObject(t, result, "2024-06-15T09:30:00Z")
}
