package evaluator

import "testing"

// dateImport prefixes a test's source with the import date's module methods
// now need, since the standard library (console/type excepted) is no longer
// ambiently available.
const dateImport = "import \"ghost:date\"\n"

// TestDateComparisons covers the operators Date supports directly - ordering
// and exact-instant equality - which is what lets `d1 < d2` read the way
// date-fns's isBefore(d1, d2) does, without needing a function for it.
func TestDateComparisons(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"date.of(2024, 1, 1) < date.of(2024, 1, 2)", true},
		{"date.of(2024, 1, 2) < date.of(2024, 1, 1)", false},
		{"date.of(2024, 1, 2) > date.of(2024, 1, 1)", true},
		{"date.of(2024, 1, 1) <= date.of(2024, 1, 1)", true},
		{"date.of(2024, 1, 1) >= date.of(2024, 1, 1)", true},
		{"date.of(2024, 1, 1) == date.of(2024, 1, 1)", true},
		{"date.of(2024, 1, 1, 9, 0, 0) == date.of(2024, 1, 1, 9, 0, 1)", false},
		{"date.of(2024, 1, 1) != date.of(2024, 1, 2)", true},
		{"date.now() > date.fromUnix(0)", true},
	}

	for _, tt := range tests {
		result := evaluate(dateImport + tt.input)

		isBooleanObject(t, result, tt.expected)
	}
}

func TestDateArithmeticOperatorsAreRejected(t *testing.T) {
	result := evaluate(dateImport + "date.now() + date.now()")

	isErrorObject(t, result, "test.gs:2:12: type error: cannot use `+` between two dates")
}

func TestDateToString(t *testing.T) {
	result := evaluate(dateImport + "date.of(2024, 6, 15, 9, 30, 0).toString()")

	isStringObject(t, result, "2024-06-15T09:30:00Z")
}

// TestDateZonedToString confirms toString() reads back through the Date's
// attached zone, not always UTC (§9.5, object.Date's doc comment).
func TestDateZonedToString(t *testing.T) {
	result := evaluate(dateImport + `date.inTimeZone(date.of(2024, 1, 15, 9, 30, 0), "America/New_York").toString()`)

	isStringObject(t, result, "2024-01-15T04:30:00-05:00")
}

// TestDateComparisonsAreZoneIndependent confirms moving a Date to a named
// zone never changes what it compares equal, before, or after to - the same
// instant read through two zones is still the same instant.
func TestDateComparisonsAreZoneIndependent(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`date.of(2024, 1, 15, 9, 30, 0) == date.inTimeZone(date.of(2024, 1, 15, 9, 30, 0), "America/New_York")`, true},
		{`date.inTimeZone(date.of(2024, 1, 15, 9, 30, 0), "America/New_York") < date.of(2024, 1, 15, 9, 30, 1)`, true},
		{`date.ofInZone(2024, 7, 15, 9, 30, 0, "America/New_York") == date.of(2024, 7, 15, 9, 30, 0)`, false},
	}

	for _, tt := range tests {
		result := evaluate(dateImport + tt.input)

		isBooleanObject(t, result, tt.expected)
	}
}
