package modules

import (
	"testing"
	"time"

	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

// callDate invokes a registered date method the way the evaluator would.
func callDate(t *testing.T, name string, args ...object.Object) object.Object {
	t.Helper()

	method, ok := DateMethods[name]

	if !ok {
		t.Fatalf("date.%s is not registered", name)
	}

	return method.Function(nil, token.Token{}, args...)
}

func mustDate(t *testing.T, result object.Object) *object.Date {
	t.Helper()

	if object.IsError(result) {
		t.Fatalf("unexpected error: %s", result.String())
	}

	date, ok := result.(*object.Date)

	if !ok {
		t.Fatalf("object is not Date. got=%T (%+v)", result, result)
	}

	return date
}

func mustBoolean(t *testing.T, result object.Object) bool {
	t.Helper()

	if object.IsError(result) {
		t.Fatalf("unexpected error: %s", result.String())
	}

	boolean, ok := result.(*object.Boolean)

	if !ok {
		t.Fatalf("object is not Boolean. got=%T (%+v)", result, result)
	}

	return boolean.Value
}

func mustInt(t *testing.T, result object.Object) int64 {
	t.Helper()

	if object.IsError(result) {
		t.Fatalf("unexpected error: %s", result.String())
	}

	number, ok := result.(*object.Number)

	if !ok {
		t.Fatalf("object is not Number. got=%T (%+v)", result, result)
	}

	return number.Int64()
}

func dateOfHelper(t *testing.T, year, month, day int64, clock ...int64) *object.Date {
	t.Helper()

	args := []object.Object{object.NewInt(year), object.NewInt(month), object.NewInt(day)}

	for _, c := range clock {
		args = append(args, object.NewInt(c))
	}

	return mustDate(t, callDate(t, "of", args...))
}

func TestDateOf(t *testing.T) {
	d := dateOfHelper(t, 2024, 2, 29, 9, 30, 15)

	if d.Time.Year() != 2024 || d.Time.Month() != time.February || d.Time.Day() != 29 {
		t.Fatalf("wrong date. got=%s", d.String())
	}

	if d.Time.Hour() != 9 || d.Time.Minute() != 30 || d.Time.Second() != 15 {
		t.Fatalf("wrong time of day. got=%s", d.String())
	}

	if d.String() != "2024-02-29T09:30:15Z" {
		t.Fatalf("wrong string form. got=%s", d.String())
	}
}

func TestDateOfRejectsImpossibleDays(t *testing.T) {
	result := callDate(t, "of", object.NewInt(2023), object.NewInt(2), object.NewInt(29))

	if !object.IsError(result) {
		t.Fatalf("expected an error for February 29 in a non-leap year, got=%s", result.String())
	}
}

func TestDateOfRejectsOutOfRangeMonth(t *testing.T) {
	result := callDate(t, "of", object.NewInt(2024), object.NewInt(13), object.NewInt(1))

	if !object.IsError(result) {
		t.Fatalf("expected an error for month 13, got=%s", result.String())
	}
}

func TestDateParseISO(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"2024-06-15", "2024-06-15T00:00:00Z"},
		{"2024-06-15T13:45:30Z", "2024-06-15T13:45:30Z"},
	}

	for _, tt := range tests {
		d := mustDate(t, callDate(t, "parseISO", &object.String{Value: tt.input}))

		if d.String() != tt.expected {
			t.Errorf("parseISO(%q): got=%s, expected=%s", tt.input, d.String(), tt.expected)
		}
	}

	result := callDate(t, "parseISO", &object.String{Value: "not a date"})

	if !object.IsError(result) {
		t.Fatalf("expected an error for unparseable input, got=%s", result.String())
	}
}

func TestDateFromUnixAndToUnix(t *testing.T) {
	d := mustDate(t, callDate(t, "fromUnix", object.NewInt(1718459130)))

	if d.String() != "2024-06-15T13:45:30Z" {
		t.Fatalf("wrong date. got=%s", d.String())
	}

	back := mustInt(t, callDate(t, "toUnix", d))

	if back != 1718459130 {
		t.Fatalf("toUnix round-trip failed. got=%d", back)
	}
}

// TestDateAddMonthsClamps covers the one behavior that separates date-fns
// from naive day-rollover arithmetic: shifting into a shorter month clamps to
// that month's last day rather than spilling into the month after.
func TestDateAddMonthsClamps(t *testing.T) {
	jan31 := dateOfHelper(t, 2024, 1, 31)

	result := mustDate(t, callDate(t, "addMonths", jan31, object.NewInt(1)))

	if result.String() != "2024-02-29T00:00:00Z" {
		t.Errorf("addMonths clamp failed. got=%s, expected 2024-02-29", result.String())
	}

	result = mustDate(t, callDate(t, "addMonths", jan31, object.NewInt(2)))

	if result.String() != "2024-03-31T00:00:00Z" {
		t.Errorf("addMonths past a clamp failed. got=%s", result.String())
	}
}

// TestDateAddYearsClamps covers the same clamping for February 29 crossing
// into a non-leap year.
func TestDateAddYearsClamps(t *testing.T) {
	leapDay := dateOfHelper(t, 2024, 2, 29)

	result := mustDate(t, callDate(t, "addYears", leapDay, object.NewInt(1)))

	if result.String() != "2025-02-28T00:00:00Z" {
		t.Errorf("addYears clamp failed. got=%s, expected 2025-02-28", result.String())
	}

	backToLeap := mustDate(t, callDate(t, "addYears", leapDay, object.NewInt(4)))

	if backToLeap.String() != "2028-02-29T00:00:00Z" {
		t.Errorf("addYears into another leap year failed. got=%s", backToLeap.String())
	}
}

func TestDateArithmetic(t *testing.T) {
	base := dateOfHelper(t, 2024, 1, 15, 10, 0, 0)

	tests := []struct {
		method   string
		amount   int64
		expected string
	}{
		{"addDays", 20, "2024-02-04T10:00:00Z"},
		{"subDays", 20, "2023-12-26T10:00:00Z"},
		{"addWeeks", 2, "2024-01-29T10:00:00Z"},
		{"subWeeks", 2, "2024-01-01T10:00:00Z"},
		{"subMonths", 1, "2023-12-15T10:00:00Z"},
		{"subYears", 1, "2023-01-15T10:00:00Z"},
		{"addHours", 5, "2024-01-15T15:00:00Z"},
		{"addMinutes", 90, "2024-01-15T11:30:00Z"},
		{"addSeconds", 3661, "2024-01-15T11:01:01Z"},
	}

	for _, tt := range tests {
		result := mustDate(t, callDate(t, tt.method, base, object.NewInt(tt.amount)))

		if result.String() != tt.expected {
			t.Errorf("%s(base, %d): got=%s, expected=%s", tt.method, tt.amount, result.String(), tt.expected)
		}
	}
}

func TestDatePredicates(t *testing.T) {
	if !mustBoolean(t, callDate(t, "isLeapYear", dateOfHelper(t, 2024, 1, 1))) {
		t.Error("2024 should be a leap year")
	}

	if mustBoolean(t, callDate(t, "isLeapYear", dateOfHelper(t, 2023, 1, 1))) {
		t.Error("2023 should not be a leap year")
	}

	if mustBoolean(t, callDate(t, "isLeapYear", dateOfHelper(t, 1900, 1, 1))) {
		t.Error("1900 should not be a leap year (divisible by 100, not 400)")
	}

	if !mustBoolean(t, callDate(t, "isLeapYear", dateOfHelper(t, 2000, 1, 1))) {
		t.Error("2000 should be a leap year (divisible by 400)")
	}

	saturday := dateOfHelper(t, 2024, 1, 6)
	monday := dateOfHelper(t, 2024, 1, 8)

	if !mustBoolean(t, callDate(t, "isWeekend", saturday)) {
		t.Error("January 6 2024 was a Saturday")
	}

	if mustBoolean(t, callDate(t, "isWeekend", monday)) {
		t.Error("January 8 2024 was a Monday")
	}

	morning := dateOfHelper(t, 2024, 1, 1, 6, 0, 0)
	night := dateOfHelper(t, 2024, 1, 1, 23, 0, 0)
	otherDay := dateOfHelper(t, 2024, 1, 2, 6, 0, 0)

	if !mustBoolean(t, callDate(t, "isSameDay", morning, night)) {
		t.Error("same calendar day at different times should be isSameDay")
	}

	if mustBoolean(t, callDate(t, "isSameDay", morning, otherDay)) {
		t.Error("different calendar days should not be isSameDay")
	}
}

func TestDateDifferences(t *testing.T) {
	earlier := dateOfHelper(t, 2024, 1, 1, 0, 0, 0)
	later := dateOfHelper(t, 2024, 1, 3, 12, 30, 0)

	tests := []struct {
		method   string
		expected int64
	}{
		{"differenceInDays", 2},
		{"differenceInHours", 60},
		{"differenceInMinutes", 3630},
		{"differenceInSeconds", 217800},
	}

	for _, tt := range tests {
		got := mustInt(t, callDate(t, tt.method, later, earlier))

		if got != tt.expected {
			t.Errorf("%s: got=%d, expected=%d", tt.method, got, tt.expected)
		}
	}

	// The reverse difference is the exact negative - truncation never makes
	// the two readings disagree by one.
	forward := mustInt(t, callDate(t, "differenceInDays", later, earlier))
	backward := mustInt(t, callDate(t, "differenceInDays", earlier, later))

	if forward != -backward {
		t.Errorf("differenceInDays is not antisymmetric: forward=%d, backward=%d", forward, backward)
	}
}

func TestDateStartAndEndOf(t *testing.T) {
	mid := dateOfHelper(t, 2024, 2, 15, 13, 45, 30)

	startOfDay := mustDate(t, callDate(t, "startOfDay", mid))
	if startOfDay.String() != "2024-02-15T00:00:00Z" {
		t.Errorf("startOfDay: got=%s", startOfDay.String())
	}

	endOfDay := mustDate(t, callDate(t, "endOfDay", mid))
	if endOfDay.Time.Hour() != 23 || endOfDay.Time.Minute() != 59 || endOfDay.Time.Second() != 59 {
		t.Errorf("endOfDay: got=%s", endOfDay.String())
	}

	startOfMonth := mustDate(t, callDate(t, "startOfMonth", mid))
	if startOfMonth.String() != "2024-02-01T00:00:00Z" {
		t.Errorf("startOfMonth: got=%s", startOfMonth.String())
	}

	endOfMonth := mustDate(t, callDate(t, "endOfMonth", mid))
	if endOfMonth.Time.Day() != 29 || endOfMonth.Time.Hour() != 23 {
		t.Errorf("endOfMonth in a leap February: got=%s", endOfMonth.String())
	}

	nonLeapFeb := dateOfHelper(t, 2023, 2, 10)
	endOfNonLeapMonth := mustDate(t, callDate(t, "endOfMonth", nonLeapFeb))

	if endOfNonLeapMonth.Time.Day() != 28 {
		t.Errorf("endOfMonth in a non-leap February: got=%s", endOfNonLeapMonth.String())
	}
}

func TestDateComponents(t *testing.T) {
	d := dateOfHelper(t, 2024, 3, 14, 15, 9, 26)

	tests := []struct {
		method   string
		expected int64
	}{
		{"year", 2024},
		{"month", 3},
		{"day", 14},
		{"hour", 15},
		{"minute", 9},
		{"second", 26},
		{"weekday", 4}, // March 14 2024 was a Thursday
	}

	for _, tt := range tests {
		got := mustInt(t, callDate(t, tt.method, d))

		if got != tt.expected {
			t.Errorf("%s: got=%d, expected=%d", tt.method, got, tt.expected)
		}
	}
}

func TestDateFormat(t *testing.T) {
	d := dateOfHelper(t, 2024, 3, 5, 9, 5, 3)

	tests := []struct {
		pattern  string
		expected string
	}{
		{"yyyy-MM-dd", "2024-03-05"},
		{"yy-M-d", "24-3-5"},
		{"MMMM d, yyyy", "March 5, 2024"},
		{"MMM", "Mar"},
		{"EEEE", "Tuesday"},
		{"EEE", "Tue"},
		{"HH:mm:ss", "09:05:03"},
		{"H:m:s", "9:5:3"},
		{"h:mm a", "9:05 AM"},

		// T, Z, :, -, and space are not tokens, so an ISO-shaped pattern needs no
		// escaping to come out the way it reads.
		{"yyyy-MM-ddTHH:mm:ss", "2024-03-05T09:05:03"},
	}

	for _, tt := range tests {
		result := callDate(t, "format", d, &object.String{Value: tt.pattern})
		str, ok := result.(*object.String)

		if !ok {
			t.Fatalf("format(%q): expected a string, got=%T", tt.pattern, result)
		}

		if str.Value != tt.expected {
			t.Errorf("format(%q): got=%s, expected=%s", tt.pattern, str.Value, tt.expected)
		}
	}
}

func TestDateArgumentErrors(t *testing.T) {
	tests := []struct {
		name string
		args []object.Object
	}{
		{"addDays", []object.Object{&object.String{Value: "not a date"}, object.NewInt(1)}},
		{"toUnix", []object.Object{object.NewInt(1)}},
		{"format", []object.Object{dateOfHelper(t, 2024, 1, 1), object.NewInt(1)}},
	}

	for _, tt := range tests {
		result := callDate(t, tt.name, tt.args...)

		if !object.IsError(result) {
			t.Errorf("date.%s: expected an error, got=%v", tt.name, result)
		}
	}
}
