package evaluator

import "testing"

// TestListTransformation covers the collection-style methods List gained
// alongside push/pop/first/last: map, filter, reduce, each, contains, unique,
// sort, reverse, and slice.
func TestListTransformation(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"[1, 2, 3].map(function(x) { return x * 2 }).toString()", "[2, 4, 6]"},
		{"[1, 2, 3].map(function(x, i) { return i }).toString()", "[0, 1, 2]"},
		{"[].map(function(x) { return x }).toString()", "[]"},

		{"[1, 2, 3, 4].filter(function(x) { return x % 2 == 0 }).toString()", "[2, 4]"},
		{"[1, 2, 3].filter(function(x) { return false }).toString()", "[]"},

		{"[1, 2, 3, 4].reduce(function(acc, x) { return acc + x }).toString()", "10"},
		{"[1, 2, 3].reduce(function(acc, x) { return acc + x }, 100).toString()", "106"},

		{"total = []; [1, 2, 3].each(function(x) { total.push(x * 10) }); total.toString()", "[10, 20, 30]"},

		{"[1, 2, 2, 3, 1].unique().toString()", "[1, 2, 3]"},

		{"[3, 1, 2].sort().toString()", "[1, 2, 3]"},
		{`["b", "a", "c"].sort().toString()`, "[a, b, c]"},
		{"[3, 1, 2].sort(function(a, b) { return b - a }).toString()", "[3, 2, 1]"},

		{"[1, 2, 3].reverse().toString()", "[3, 2, 1]"},

		{"[1, 2, 3, 4, 5].slice(1, 3).toString()", "[2, 3]"},
		{"[1, 2, 3].slice(1).toString()", "[2, 3]"},

		// Chaining reads the way a Collection pipeline should.
		{"[1, 2, 3, 4, 5].filter(function(x) { return x % 2 == 0 }).map(function(x) { return x * 10 }).toString()", "[20, 40]"},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isStringObject(t, result, tt.expected)
	}
}

// TestListContains covers contains(), which compares contents the same way
// `==` does - to any depth, not by identity.
func TestListContains(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"[1, 2, 3].contains(2)", true},
		{"[1, 2, 3].contains(9)", false},
		{"[[1, 2], [3, 4]].contains([1, 2])", true},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isBooleanObject(t, result, tt.expected)
	}
}

// TestListPopIsFromTheBack covers the fix to pop(): it now matches push() by
// acting on the end of the list. shift() takes over the front-removal
// behavior pop() used to have.
func TestListPopIsFromTheBack(t *testing.T) {
	result := evaluate("a = [1, 2, 3]; a.pop()")
	isNumberObject(t, result, 3)

	result = evaluate("a = [1, 2, 3]; a.pop(); a.toString()")
	isStringObject(t, result, "[1, 2]")

	result = evaluate("a = [1, 2, 3]; a.shift(); a.toString()")
	isStringObject(t, result, "[2, 3]")

	result = evaluate("a = [1, 2, 3]; a.shift()")
	isNumberObject(t, result, 1)
}

func TestListTransformationErrors(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"[1].map()", "test.ghost:1:5: argument error: `list.map()` expects 1 argument, got 0"},
		{"[].reduce(function(acc, x) { return acc })", "test.ghost:1:4: argument error: `list.reduce()` needs an initial value to reduce an empty list"},
		{"[3, [1]].sort()", "test.ghost:1:10: argument error: `list.sort()` needs a comparator to sort anything but a list of only numbers or only strings"},
		{"[1, 2].slice(5)", "test.ghost:1:8: index error: `list.slice()` start index 5 is out of range for a list of length 2"},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isErrorObject(t, result, tt.expectedMessage)
	}
}
