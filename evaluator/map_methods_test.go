package evaluator

import "testing"

// TestMapMethods covers the methods Map gained - it had none before. keys()
// and values() are asserted through sort(), since a Go map's iteration order
// is not something a test should depend on.
func TestMapMethods(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`{"a": 1, "b": 2}.keys().sort().toString()`, "[a, b]"},
		{`{"a": 1, "b": 2}.values().sort().toString()`, "[1, 2]"},
		{`{"a": 1}.get("a").toString()`, "1"},
		{`{"a": 1}.get("missing", 42).toString()`, "42"},
		{`{"a": 1}.get("missing").toString()`, "null"},
		{`{}.length().toString()`, "0"},
		{`{"a": 1, "b": 2}.length().toString()`, "2"},

		// set() answers the map itself, so it chains.
		{`m = {}; m.set("a", 1); m.get("a").toString()`, "1"},
		{`m = {"a": 1}; m.set("a", 2).set("b", 3); m.keys().sort().toString()`, "[a, b]"},

		// merge() leaves both operands alone and lets the argument win on a
		// shared key.
		{`{"a": 1}.merge({"b": 2}).keys().sort().toString()`, "[a, b]"},
		{`{"a": 1}.merge({"a": 2}).get("a").toString()`, "2"},
		{`x = {"a": 1}; x.merge({"b": 2}); x.keys().sort().toString()`, "[a]"},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isStringObject(t, result, tt.expected)
	}
}

func TestMapMethodBooleans(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`{"a": 1}.has("a")`, true},
		{`{"a": 1}.has("b")`, false},

		// A key mapped to null is still present.
		{`{"a": null}.has("a")`, true},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isBooleanObject(t, result, tt.expected)
	}
}

func TestMapMethodErrors(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{`{}.get()`, "test.ghost:1:4: argument error: `map.get()` expects between 1 and 2 arguments, got 0"},
		{`{}.set("a")`, "test.ghost:1:4: argument error: `map.set()` expects 2 arguments, got 1"},
		{`{}.merge(1)`, "test.ghost:1:4: argument error: `map.merge()` expects argument 1 to be a map, got number"},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isErrorObject(t, result, tt.expectedMessage)
	}
}
