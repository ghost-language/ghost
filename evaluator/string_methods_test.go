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

		// matches() takes the pattern as its argument, the receiver is the
		// subject being searched (§13.7).
		{`"hello world".matches("wor.d")`, true},
		{`"hello world".matches("xyz")`, false},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isBooleanObject(t, result, tt.expected)
	}
}

// TestStringPatternMethods covers find()/findAll()/matches() after the
// receiver/argument flip and the findAll() bug fix (§13.7, §14 decision 3):
// the receiver is always the subject being searched, the argument is always
// the pattern - subject.find(pattern), matching JS/PHP/Python - and
// findAll() now answers every match, not just the first match's capture
// groups.
func TestStringPatternMethods(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello world".find("wor.d")`, "world"},
		{`"hello world".find("xyz")`, ""},
		{`"foo123bar456".find("[0-9]+")`, "123"},

		// findAll() answers every match's full text, in order - not, as
		// before the fix, only the first match's own capture groups.
		{`"foo123bar456".findAll("[0-9]+").toString()`, "[123, 456]"},
		{`"no digits here".findAll("[0-9]+").toString()`, "[]"},
		{`"aaa".findAll("a").toString()`, "[a, a, a]"},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isStringObject(t, result, tt.expected)
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
		{`"hello".find("(")`, "test.gs:1:9: value error: `string.find()` cannot use `(` as a pattern: missing closing )"},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isErrorObject(t, result, tt.expectedMessage)
	}
}
