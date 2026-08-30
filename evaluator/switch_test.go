package evaluator

import (
	"testing"
)

func TestSwitchMatchesCases(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`switch (1) { case 1 { "one" } case 2 { "two" } }`, "one"},
		{`switch (2) { case 1 { "one" } case 2 { "two" } }`, "two"},
		{`switch (9) { case 1 { "one" } case 2 { "two" } default { "other" } }`, "other"},
		// A case can list several values, matching any of them.
		{`switch (3) { case 1, 2, 3 { "small" } default { "other" } }`, "small"},
		{`switch (9) { case 1, 2, 3 { "small" } default { "other" } }`, "other"},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isStringObject(t, result, tt.expected)
	}
}

// TestSwitchComparesByValueNotString covers §13.3's second defect: case
// matching used to compare Type() and String() rather than going through
// object.ValuesEqual, and every Function stringifies to the literal
// "function" (every Class to "class") regardless of which one it actually
// is - so a switch on one function used to match a case naming a
// completely different one. It also confirms the fix gives switch the same
// content-equality a map or list already gets from `==` (§13.2).
func TestSwitchComparesByValueNotString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"f = function() { return 1 }\ng = function() { return 2 }\nswitch (f) { case g { \"wrong\" } default { \"correct\" } }",
			"correct",
		},
		{
			"f = function() { return 1 }\nswitch (f) { case f { \"correct\" } default { \"wrong\" } }",
			"correct",
		},
		{
			`switch ({x: 1}) { case {x: 1} { "matched by contents" } default { "wrong" } }`,
			"matched by contents",
		},
		{
			`switch ({x: 1}) { case {x: 2} { "wrong" } default { "correctly unmatched" } }`,
			"correctly unmatched",
		},
		{
			`switch ([1, 2]) { case [1, 2] { "matched by contents" } default { "wrong" } }`,
			"matched by contents",
		},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isStringObject(t, result, tt.expected)
	}
}

// TestSwitchPropagatesErrors covers §13.3's first defect: an error in the
// subject or a case expression used to be compared like any other value
// (falling through to default, or simply never matching) instead of
// propagating the way every other construct in the language does.
func TestSwitchPropagatesErrors(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			`switch (undefinedVar) { case 1 { "one" } default { "default" } }`,
			"test.gs:1:9: name error: `undefinedVar` is not defined",
		},
		{
			`switch (1) { case undefinedVar { "one" } default { "default" } }`,
			"test.gs:1:19: name error: `undefinedVar` is not defined",
		},
		{
			// The error surfaces even when an earlier, non-erroring case
			// would never have matched anyway - every case value is still
			// evaluated in order until a match or an error is found.
			`switch (99) { case 1 { "one" } case undefinedVar { "two" } default { "default" } }`,
			"test.gs:1:37: name error: `undefinedVar` is not defined",
		},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isErrorObject(t, result, tt.expectedMessage)
	}
}

// TestSwitchWithNoMatchAndNoDefaultIsEmpty confirms a switch that matches
// nothing, and has no default to fall back to, produces no value - unrelated
// to §13.3, just the existing, unchanged behavior for that shape of switch.
func TestSwitchWithNoMatchAndNoDefaultIsEmpty(t *testing.T) {
	result := evaluate(`switch (99) { case 1 { "one" } }`)

	if result != nil {
		t.Errorf("expected no value, got=%#v", result)
	}
}
