package evaluator

import "testing"

// TestStringMethods covers the methods closing §12's string-method gap:
// contains, indexOf/lastIndexOf, repeat, padStart/padEnd, charAt, slice,
// reverse. Several mirror a List method of the same name (contains, slice,
// reverse) so the two collection types read the same way.
func TestStringMethods(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello".charAt(0)`, "h"},
		{`"hello".charAt(4)`, "o"},
		{`"hello".charAt(10)`, ""},
		{`"hello".charAt(-1)`, ""},

		{`"3".padStart(5, "0")`, "00003"},
		{`"3".padStart(5)`, "    3"},
		{`"hello".padStart(3, "0")`, "hello"},
		{`"3".padEnd(5, "0")`, "30000"},
		{`"3".padEnd(5)`, "3    "},
		{`"ab".padStart(5, "xyz")`, "xyzab"},

		{`"ab".repeat(3)`, "ababab"},
		{`"ab".repeat(0)`, ""},

		{`"hello".reverse()`, "olleh"},
		{`"".reverse()`, ""},

		{`"hello world".slice(6)`, "world"},
		{`"hello world".slice(0, 5)`, "hello"},

		// Unicode: length/charAt/slice/indexOf all count runes, not bytes.
		{`"héllo".charAt(1)`, "é"},
		{`"héllo".slice(1, 3)`, "él"},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isStringObject(t, result, tt.expected)
	}
}

func TestStringMethodNumbers(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`"hello world".indexOf("world")`, 6},
		{`"hello world".indexOf("missing")`, -1},
		{`"hello world".lastIndexOf("o")`, 7},
		{`"hello world".lastIndexOf("missing")`, -1},
		{`"héllo".indexOf("llo")`, 2},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isNumberObject(t, result, tt.expected)
	}
}

func TestStringMethodBooleans(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`"hello world".contains("world")`, true},
		{`"hello world".contains("xyz")`, false},
		{`"".isEmpty()`, true},
		{`"a".isEmpty()`, false},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isBooleanObject(t, result, tt.expected)
	}
}

func TestStringMethodErrors(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{`"ab".repeat(-1)`, "test.gs:1:6: value error: `string.repeat()` count cannot be negative, got -1"},
		{`"hello".slice(10)`, "test.gs:1:9: index error: `string.slice()` start index 10 is out of range for a string of length 5"},
		{`"hello".slice(0, 10)`, "test.gs:1:9: index error: `string.slice()` end index 10 is out of range for a string of length 5"},
		{`"a".contains()`, "test.gs:1:5: argument error: `string.contains()` expects 1 argument, got 0"},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isErrorObject(t, result, tt.expectedMessage)
	}
}
