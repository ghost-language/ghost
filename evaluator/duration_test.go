package evaluator

import "testing"

// TestDurationToString covers Duration's ISO 8601 string form, reached
// through Ghost's toString() the way every other value's is.
func TestDurationToString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"date.duration(1, 2, 3, 4, 5, 6).toString()", "P1Y2M3DT4H5M6S"},
		{"date.duration(0, 1, 0).toString()", "P1M"},
		{"date.duration(0, 0, 0).toString()", "PT0S"},
		{"date.duration(-1, -2, -3).toString()", "P-1Y-2M-3D"},
	}

	for _, tt := range tests {
		result := evaluate(dateImport + tt.input)

		isStringObject(t, result, tt.expected)
	}
}

func TestDurationType(t *testing.T) {
	result := evaluate(dateImport + "type(date.duration(1, 2, 3))")

	isStringObject(t, result, "duration")
}

// TestDurationComponentAccessors covers reading each field back out through
// Ghost method syntax.
func TestDurationComponentAccessors(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"date.duration(1, 2, 3, 4, 5, 6).years()", 1},
		{"date.duration(1, 2, 3, 4, 5, 6).months()", 2},
		{"date.duration(1, 2, 3, 4, 5, 6).days()", 3},
		{"date.duration(1, 2, 3, 4, 5, 6).hours()", 4},
		{"date.duration(1, 2, 3, 4, 5, 6).minutes()", 5},
		{"date.duration(1, 2, 3, 4, 5, 6).seconds()", 6},
	}

	for _, tt := range tests {
		result := evaluate(dateImport + tt.input)

		isNumberObject(t, result, tt.expected)
	}
}

// TestDurationBetweenAndApply exercises the full round trip through Ghost
// syntax: computing a Duration between two dates and applying it back.
func TestDurationBetweenAndApply(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`date.addDuration(date.of(2024, 1, 15), date.durationBetween(date.of(2024, 8, 3, 22, 15, 45), date.of(2024, 1, 15))) == date.of(2024, 8, 3, 22, 15, 45)`, true},
		{`date.subDuration(date.of(2024, 8, 3, 22, 15, 45), date.durationBetween(date.of(2024, 8, 3, 22, 15, 45), date.of(2024, 1, 15))) == date.of(2024, 1, 15)`, true},
	}

	for _, tt := range tests {
		result := evaluate(dateImport + tt.input)

		isBooleanObject(t, result, tt.expected)
	}
}

func TestDurationArgumentErrors(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"date.duration(1, -1, 0)", "test.gs:2:6: value error: `date.duration()` expects every component to point the same direction - all positive, all negative, or zero"},
		{"date.duration(1, 2)", "test.gs:2:6: argument error: `date.duration()` expects between 3 and 6 arguments, got 2"},
		{"date.durationBetween(date.now())", "test.gs:2:6: argument error: `date.durationBetween()` expects 2 arguments, got 1"},
		{`date.addDuration(date.now(), "not a duration")`, "test.gs:2:6: argument error: `date.addDuration()` expects argument 2 to be a duration, got string"},
	}

	for _, tt := range tests {
		result := evaluate(dateImport + tt.input)

		isErrorObject(t, result, tt.expectedMessage)
	}
}
