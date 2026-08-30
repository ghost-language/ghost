package evaluator

import (
	"testing"

	"ghostlang.org/x/ghost/object"
)

// TestIndexingIsLenientOnPosition covers §13.6: `list[i]`/`map[k]`/`string[i]`
// all name a position, not a range, so an out-of-range read answers null
// rather than raising an Index error - the same leniency charAt()/get() give
// deliberately, contrasted with slice()'s range validation.
func TestIndexingIsLenientOnPosition(t *testing.T) {
	tests := []struct {
		input string
	}{
		{`[1, 2, 3][3]`},
		{`[1, 2, 3][-1]`},
		{`[][0]`},
		{`{"a": 1}["b"]`},
		{`{}["missing"]`},
		{`"abc"[3]`},
		{`"abc"[-1]`},
		{`""[0]`},

		// A receiver with a multi-byte rune has more bytes than runes -
		// reproduces the bug where the bounds check compared idx against
		// len(str.Value) (a byte count) instead of the rune count, so an
		// index between the two lengths passed the check and then panicked
		// indexing past the end of the []rune conversion.
		{`"héllo"[5]`},
		{`"héllo"[99]`},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isNullObject(t, result)
	}
}

// TestStringIndexUsesRunePositions confirms `string[i]` agrees with
// length()/charAt() about what a "character" is - one rune, not one byte -
// for a receiver containing multi-byte runes.
func TestStringIndexUsesRunePositions(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"héllo"[0]`, "h"},
		{`"héllo"[1]`, "é"},
		{`"héllo"[4]`, "o"},
		{`"日本語"[1]`, "本"},
	}

	for _, tt := range tests {
		result := evaluate(tt.input)

		isStringObject(t, result, tt.expected)
	}
}

func isNullObject(t *testing.T, obj object.Object) bool {
	t.Helper()

	if _, ok := obj.(*object.Null); !ok {
		t.Errorf("object is not Null. got=%T (%+v)", obj, obj)
		return false
	}

	return true
}
