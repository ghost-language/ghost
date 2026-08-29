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
		{"[1].map()", "test.gs:1:5: argument error: `list.map()` expects 1 argument, got 0"},
		{"[].reduce(function(acc, x) { return acc })", "test.gs:1:4: argument error: `list.reduce()` needs an initial value to reduce an empty list"},
		{"[3, [1]].sort()", "test.gs:1:10: argument error: `list.sort()` needs a comparator to sort anything but a list of only numbers or only strings"},
		{"[1, 2].slice(5)", "test.gs:1:8: index error: `list.slice()` start index 5 is out of range for a list of length 2"},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isErrorObject(t, result, tt.expectedMessage)
	}
}

// TestListSearch covers the value-based (indexOf) and predicate-based
// (find/findIndex/some/every) search methods, plus isEmpty.
func TestListSearch(t *testing.T) {
	numberTests := []struct {
		input    string
		expected int64
	}{
		{"[1, 2, 3].indexOf(2)", 1},
		{"[1, 2, 3].indexOf(9)", -1},
		{"[[1, 2], [3, 4]].indexOf([3, 4])", 1},

		{"[1, 2, 3, 4].findIndex(function(x) { return x % 2 == 0 })", 1},
		{"[1, 3, 5].findIndex(function(x) { return x % 2 == 0 })", -1},
	}

	for _, tt := range numberTests {
		result := evaluate(tt.input)

		isNumberObject(t, result, tt.expected)
	}

	stringTests := []struct {
		input    string
		expected string
	}{
		{"[1, 2, 3, 4].find(function(x) { return x % 2 == 0 }).toString()", "2"},
		{"[1, 3, 5].find(function(x) { return x % 2 == 0 }).toString()", "null"},
	}

	for _, tt := range stringTests {
		result := evaluate(tt.input)

		isStringObject(t, result, tt.expected)
	}

	boolTests := []struct {
		input    string
		expected bool
	}{
		{"[1, 2, 3].some(function(x) { return x > 2 })", true},
		{"[1, 2, 3].some(function(x) { return x > 5 })", false},
		{"[1, 2, 3].every(function(x) { return x > 0 })", true},
		{"[1, 2, 3].every(function(x) { return x > 1 })", false},
		{"[].isEmpty()", true},
		{"[1].isEmpty()", false},
	}

	for _, tt := range boolTests {
		result := evaluate(tt.input)

		isBooleanObject(t, result, tt.expected)
	}
}

// TestListFlattenAndChunk covers flatten, flatMap, and chunk - the methods
// that reshape a list's structure rather than filtering or reducing it.
func TestListFlattenAndChunk(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"[1, [2, 3], [4, [5, 6]]].flatten().toString()", "[1, 2, 3, 4, 5, 6]"},
		{"[[1], [2], [3]].flatten().toString()", "[1, 2, 3]"},
		{"[].flatten().toString()", "[]"},

		{"[1, 2, 3].flatMap(function(x) { return [x, x * 10] }).toString()", "[1, 10, 2, 20, 3, 30]"},
		{"[1, 2, 3].flatMap(function(x) { return x * 2 }).toString()", "[2, 4, 6]"},

		{"[1, 2, 3, 4, 5].chunk(2).toString()", "[[1, 2], [3, 4], [5]]"},
		{"[1, 2, 3].chunk(1).toString()", "[[1], [2], [3]]"},
		{"[].chunk(3).toString()", "[]"},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isStringObject(t, result, tt.expected)
	}
}

// TestListFill covers fill(), which - unlike push/pop/shift/unshift/
// insertAt/removeAt - does not mutate, matching slice()/sort()/reverse().
func TestListFill(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"[1, 2, 3, 4].fill(0).toString()", "[0, 0, 0, 0]"},
		{"[1, 2, 3, 4].fill(0, 1).toString()", "[1, 0, 0, 0]"},
		{"[1, 2, 3, 4].fill(0, 1, 3).toString()", "[1, 0, 0, 4]"},
		{"a = [1, 2, 3]; a.fill(9); a.toString()", "[1, 2, 3]"},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isStringObject(t, result, tt.expected)
	}
}

// TestListFrontAndPositionMutation covers unshift, insertAt, and removeAt -
// the mutating position-based methods that round out push/pop/shift.
func TestListFrontAndPositionMutation(t *testing.T) {
	result := evaluate("a = [2, 3]; a.unshift(1); a.toString()")
	isStringObject(t, result, "[1, 2, 3]")

	result = evaluate("a = [2, 3]; a.unshift(1)")
	isNumberObject(t, result, 3)

	result = evaluate(`a = [1, 2, 4]; a.insertAt(2, 3); a.toString()`)
	isStringObject(t, result, "[1, 2, 3, 4]")

	result = evaluate(`a = [1, 2, 4]; a.insertAt(2, 3)`)
	isNumberObject(t, result, 4)

	// An out-of-range index clamps rather than erroring.
	result = evaluate(`a = [1, 2]; a.insertAt(99, 3); a.toString()`)
	isStringObject(t, result, "[1, 2, 3]")

	result = evaluate(`a = [1, 2]; a.insertAt(-5, 0); a.toString()`)
	isStringObject(t, result, "[0, 1, 2]")

	result = evaluate(`a = [1, 2, 3]; a.removeAt(1)`)
	isNumberObject(t, result, 2)

	result = evaluate(`a = [1, 2, 3]; a.removeAt(1); a.toString()`)
	isStringObject(t, result, "[1, 3]")

	// An out-of-range index is lenient, like pop()/shift() on an empty list.
	result = evaluate(`a = [1, 2, 3]; a.removeAt(9).toString()`)
	isStringObject(t, result, "null")
}

func TestListSearchAndMutationErrors(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"[1].find()", "test.gs:1:5: argument error: `list.find()` expects 1 argument, got 0"},
		{"[1].indexOf()", "test.gs:1:5: argument error: `list.indexOf()` expects 1 argument, got 0"},
		{"[1, 2].chunk(0)", "test.gs:1:8: value error: `list.chunk()` size has to be positive, got 0"},
		{"[1, 2].chunk(-1)", "test.gs:1:8: value error: `list.chunk()` size has to be positive, got -1"},
		{"[1, 2].fill(0, 5)", "test.gs:1:8: index error: `list.fill()` start index 5 is out of range for a list of length 2"},
		{"[1, 2].fill(0, 0, 5)", "test.gs:1:8: index error: `list.fill()` end index 5 is out of range for a list of length 2"},
		{"[1].insertAt(0)", "test.gs:1:5: argument error: `list.insertAt()` expects 2 arguments, got 1"},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isErrorObject(t, result, tt.expectedMessage)
	}
}
