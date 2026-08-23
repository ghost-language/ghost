package evaluator

import "testing"

// TestToStringIsUniversal covers toString() on the two value types that were
// missing it - boolean and null. Every other value type answers a string for
// itself; these two didn't, which was the gap "conversions read left to
// right, target last" (SPEC.md's naming conventions) implied should not exist.
func TestToStringIsUniversal(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"true.toString()", "true"},
		{"false.toString()", "false"},
		{"null.toString()", "null"},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isStringObject(t, result, tt.expected)
	}
}
