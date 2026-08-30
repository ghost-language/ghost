package evaluator

import "testing"

// TestMapMethods covers the methods Map gained - it had none before. keys()
// and values() are asserted through sort() here for a reason unrelated to
// map ordering: these particular cases only care that the right elements
// came back, not what order they're in. TestMapPreservesInsertionOrder below
// covers order itself (§13.5, §14 decision 2) - Map does guarantee one now.
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

		// entries() answers [key, value] pairs; sort by the joined string,
		// since element order across pairs is not something a test should
		// depend on either.
		{`{"a": 1, "b": 2}.entries().map(function(pair) { return pair.join(":") }).sort().toString()`, "[a:1, b:2]"},
		{`{}.entries().toString()`, "[]"},

		// remove() mutates in place and answers the value that was there,
		// or null when the key was never present.
		{`m = {"a": 1, "b": 2}; m.remove("a").toString()`, "1"},
		{`m = {"a": 1, "b": 2}; m.remove("a"); m.keys().sort().toString()`, "[b]"},
		{`{"a": 1}.remove("missing").toString()`, "null"},
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
		{`{}.get()`, "test.gs:1:4: argument error: `map.get()` expects between 1 and 2 arguments, got 0"},
		{`{}.set("a")`, "test.gs:1:4: argument error: `map.set()` expects 2 arguments, got 1"},
		{`{}.merge(1)`, "test.gs:1:4: argument error: `map.merge()` expects argument 1 to be a map, got number"},
		{`{}.remove()`, "test.gs:1:4: argument error: `map.remove()` expects 1 argument, got 0"},
		{`{}.remove([1])`, "test.gs:1:4: type error: `map.remove()` cannot use list as a map key"},
		{`{}.entries(1)`, "test.gs:1:4: argument error: `map.entries()` expects 0 arguments, got 1"},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isErrorObject(t, result, tt.expectedMessage)
	}
}

// TestMapPreservesInsertionOrder covers §13.5 (§14 decision 2): keys(),
// values(), entries(), String(), and `for ... in` all agree on the order
// keys were first inserted in, the same guarantee a JS object or a PHP
// associative array already gives their users - rather than a bare Go map's
// randomized order, which could (and did) differ from call to call on the
// very same map, not just from run to run.
func TestMapPreservesInsertionOrder(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`{z: 1, a: 2, m: 3, b: 4}.keys().toString()`, "[z, a, m, b]"},
		{`{z: 1, a: 2, m: 3, b: 4}.values().toString()`, "[1, 2, 3, 4]"},
		{"`${ {z: 1, a: 2, m: 3, b: 4} }`", "{z: 1, a: 2, m: 3, b: 4}"},

		// A repeated key keeps the position of its first appearance; only
		// its value changes - the same rule set() and assignment follow.
		{`{x: 1, x: 2}.keys().toString()`, "[x]"},
		{`{x: 1, x: 2}.values().toString()`, "[2]"},

		// set() on a new key appends; set() on an existing key updates the
		// value in place without moving it.
		{`m = {}; m.set("z", 1); m.set("a", 2); m.set("z", 99); m.keys().toString()`, "[z, a]"},
		{`m = {}; m.set("z", 1); m.set("a", 2); m.set("z", 99); m.values().toString()`, "[99, 2]"},

		// Index and property assignment follow the same rule as set().
		{`m = {}; m["x"] = 1; m.y = 2; m["x"] = 99; m.keys().toString()`, "[x, y]"},
		{`m = {}; m["x"] = 1; m.y = 2; m["x"] = 99; m.values().toString()`, "[99, 2]"},

		// remove() drops a key from the order entirely, rather than leaving
		// a gap; a key added afterward is appended at the end as usual.
		{`m = {a: 1, b: 2, c: 3}; m.remove("b"); m.set("d", 4); m.keys().toString()`, "[a, c, d]"},

		// merge() answers this map's keys in their own order first, then
		// other's remaining keys in its order - a shared key keeps this
		// map's position for it and takes other's value, the same result a
		// plain object spread (`{...left, ...right}`) gives in JS.
		{`{a: 1, b: 2}.merge({b: 99, c: 3}).keys().toString()`, "[a, b, c]"},
		{`{a: 1, b: 2}.merge({b: 99, c: 3}).values().toString()`, "[1, 99, 3]"},

		// entries() zips keys and values together in the same order.
		{`{a: 1, b: 2}.entries().toString()`, "[[a, 1], [b, 2]]"},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isStringObject(t, result, tt.expected)
	}

	// The order is stable across repeated reads of the same map, not just
	// consistent within one - a bare Go map's iteration order can (and did)
	// differ between two calls in the same run, not only between runs.
	m := evaluate(`{z: 1, a: 2, m: 3, b: 4}`)

	first := m.String()

	for i := 0; i < 20; i++ {
		if again := m.String(); again != first {
			t.Fatalf("map order changed between reads: got=%s, expected=%s", again, first)
		}
	}
}

// TestMapForInPreservesInsertionOrder covers the same guarantee for
// `for ... in`, which reads a map's pairs the same way keys()/values() do.
func TestMapForInPreservesInsertionOrder(t *testing.T) {
	result := evaluate(`
	m = {z: 1, a: 2, m: 3, b: 4}
	order = []
	for (k, v in m) {
		order.push(k)
	}
	order.toString()
	`)

	isStringObject(t, result, "[z, a, m, b]")
}
