package evaluator

import "testing"

// The tests in this file cover §13.25 — a library global used to win over
// every binding, including the most local ones there are. `evaluateIdentifier`
// consulted the registries before the environment, so `console` and `type`
// could not be shadowed by a parameter, a loop variable, an import, or a
// plain assignment. Nothing raised; the name simply answered with the library
// value, which for `type` — an ordinary word for an ordinary parameter — meant
// a function object turning up where a string was expected.

func TestABindingShadowsALibraryGlobal(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "a parameter named like a global function",
			input:    `function render(type) { return type }` + "\n" + `render("warning")`,
			expected: "warning",
		},
		{
			name:     "a parameter named like a global module",
			input:    `function show(console) { return console }` + "\n" + `show("text")`,
			expected: "text",
		},
		{
			name:     "a plain assignment",
			input:    `type = "warning"` + "\n" + `type`,
			expected: "warning",
		},
		{
			name:     "a field initialized from a same-named parameter",
			input:    "class Cmd { constructor(type) { this.type = type } }\nnew Cmd(\"edit\").type",
			expected: "edit",
		},
		{
			name:     "a loop variable",
			input:    "last = \"\"\nfor (type in [\"a\", \"b\"]) { last = type }\nlast",
			expected: "b",
		},
		{
			name:     "a name bound in a block",
			input:    "seen = \"\"\nif (true) { type = \"x\" seen = type }\nseen",
			expected: "x",
		},
		{
			name:     "a function's local",
			input:    "function f() { type = \"local\" return type }\nf()",
			expected: "local",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isStringObject(t, evaluate(test.input), test.expected)
		})
	}
}

// TestAnImportShadowsALibraryGlobal is the same rule reached through the other
// door: `as` names the binding, so it decides what the name means here.
func TestAnImportShadowsALibraryGlobal(t *testing.T) {
	input := `
		import "ghost:math" as console
		console.floor(2.7)
	`

	isNumberObject(t, evaluate(input), 2)
}

// TestLibraryGlobalsStillResolveUnshadowed is the other half: a global is
// still reachable with no import wherever nothing is bound over it, which is
// the whole point of §9.1's two names.
func TestLibraryGlobalsStillResolveUnshadowed(t *testing.T) {
	isStringObject(t, evaluate(`type("x")`), "string")

	t.Run("and inside a function that binds nothing of the name", func(t *testing.T) {
		isStringObject(t, evaluate("function f(v) { return type(v) }\nf(1)"), "number")
	})

	t.Run("and inside a method", func(t *testing.T) {
		isStringObject(t, evaluate("class A { probe(v) { return type(v) } }\nnew A().probe(true)"), "boolean")
	})
}

// TestAShadowedGlobalIsRestoredAfterItsScopeEnds checks that shadowing is
// scoped rather than global: rebinding inside a function must not reach the
// name everyone else reads.
func TestAShadowedGlobalIsRestoredAfterItsScopeEnds(t *testing.T) {
	input := `
		function f(type) { return type }
		f("shadowed")
		type("x")
	`

	isStringObject(t, evaluate(input), "string")
}
