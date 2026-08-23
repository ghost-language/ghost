package modules

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

// The date module covers a single instant in time - construction, arithmetic,
// comparison, and formatting - in the spirit of the date-fns library rather
// than a language's built-in Date class: every function here takes a date and
// answers a new value, none of them mutate the one they were given, and there
// is no method-chasing through an object's own methods.
//
// A Date is always UTC; the module does not model time zones. That keeps a
// date built once and compared later the same value no matter where the
// program runs, the same way a seeded random run is reproducible regardless
// of machine.
//
// format() reads a pattern the way date-fns does - a run of the same letter
// is one token (`yyyy`, `MM`, `dd`), and anything else in the pattern is
// copied through literally - rather than Go's reference-date layout strings,
// which read like a date rather than describing one.

var DateMethods = map[string]*object.LibraryFunction{}
var DateProperties = map[string]*object.LibraryProperty{}

func init() {
	// Construction and conversion.
	RegisterMethod(DateMethods, "now", dateNow)
	RegisterMethod(DateMethods, "today", dateToday)
	RegisterMethod(DateMethods, "of", dateOf)
	RegisterMethod(DateMethods, "parseISO", dateParseISO)
	RegisterMethod(DateMethods, "fromUnix", dateFromUnix)
	RegisterMethod(DateMethods, "toUnix", dateToUnix)
	RegisterMethod(DateMethods, "toUnixNano", dateToUnixNano)
	RegisterMethod(DateMethods, "format", dateFormat)

	// Arithmetic. Each has a sub counterpart rather than accepting a negative
	// count, because "3 months before this one" reads as subMonths(d, 3), not
	// addMonths(d, -3).
	RegisterMethod(DateMethods, "addDays", dateAddDays)
	RegisterMethod(DateMethods, "subDays", dateSubDays)
	RegisterMethod(DateMethods, "addWeeks", dateAddWeeks)
	RegisterMethod(DateMethods, "subWeeks", dateSubWeeks)
	RegisterMethod(DateMethods, "addMonths", dateAddMonths)
	RegisterMethod(DateMethods, "subMonths", dateSubMonths)
	RegisterMethod(DateMethods, "addYears", dateAddYears)
	RegisterMethod(DateMethods, "subYears", dateSubYears)
	RegisterMethod(DateMethods, "addHours", dateAddHours)
	RegisterMethod(DateMethods, "addMinutes", dateAddMinutes)
	RegisterMethod(DateMethods, "addSeconds", dateAddSeconds)

	// Predicates. Ordering and exact-instant equality are `<`, `>`, and `==` -
	// see object.Date's doc comment - so only the comparisons an operator
	// cannot make are here.
	RegisterMethod(DateMethods, "isSameDay", dateIsSameDay)
	RegisterMethod(DateMethods, "isWeekend", dateIsWeekend)
	RegisterMethod(DateMethods, "isLeapYear", dateIsLeapYear)

	// Differences.
	RegisterMethod(DateMethods, "differenceInDays", dateDifferenceInDays)
	RegisterMethod(DateMethods, "differenceInHours", dateDifferenceInHours)
	RegisterMethod(DateMethods, "differenceInMinutes", dateDifferenceInMinutes)
	RegisterMethod(DateMethods, "differenceInSeconds", dateDifferenceInSeconds)

	// Start and end of a period.
	RegisterMethod(DateMethods, "startOfDay", dateStartOfDay)
	RegisterMethod(DateMethods, "endOfDay", dateEndOfDay)
	RegisterMethod(DateMethods, "startOfMonth", dateStartOfMonth)
	RegisterMethod(DateMethods, "endOfMonth", dateEndOfMonth)

	// Components.
	RegisterMethod(DateMethods, "year", dateYear)
	RegisterMethod(DateMethods, "month", dateMonth)
	RegisterMethod(DateMethods, "day", dateDay)
	RegisterMethod(DateMethods, "hour", dateHour)
	RegisterMethod(DateMethods, "minute", dateMinute)
	RegisterMethod(DateMethods, "second", dateSecond)
	RegisterMethod(DateMethods, "weekday", dateWeekday)
}

// newDate wraps a time.Time as a Date, normalizing it to UTC so that every
// Date the module hands back agrees with every other one on what "the same
// instant" means.
func newDate(t time.Time) *object.Date {
	return &object.Date{Time: t.UTC()}
}

// =============================================================================
// Construction and conversion

func dateNow(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("date.now", tok, args, 0); err != nil {
		return err
	}

	return newDate(time.Now())
}

func dateToday(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("date.today", tok, args, 0); err != nil {
		return err
	}

	now := time.Now().UTC()

	return newDate(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC))
}

// dateOf builds a date from its calendar components. Month is 1-12, matching
// the numbers on a calendar rather than JavaScript's 0-11.
func dateOf(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arityRange("date.of", tok, args, 3, 6); err != nil {
		return err
	}

	year, err := integerAt("date.of", tok, args, 0)

	if err != nil {
		return err
	}

	month, err := integerAt("date.of", tok, args, 1)

	if err != nil {
		return err
	}

	if month < 1 || month > 12 {
		return object.NewError(fault.Value, tok, "`date.of()` expects a month between 1 and 12, got %d", month)
	}

	day, err := integerAt("date.of", tok, args, 2)

	if err != nil {
		return err
	}

	hour, minute, second, timeErr := dateOfTimeComponents(tok, args)

	if timeErr != nil {
		return timeErr
	}

	built := time.Date(int(year), time.Month(month), int(day), int(hour), int(minute), int(second), 0, time.UTC)

	// time.Date normalizes an out-of-range day or component by rolling it into
	// the following period rather than reporting it, which would turn a typo
	// into a silently different date. Catching the roll here reports it
	// instead.
	if built.Day() != int(day) || built.Month() != time.Month(month) {
		return object.NewError(fault.Value, tok, "`date.of()` has no day %d in month %d", day, month)
	}

	return newDate(built)
}

func dateOfTimeComponents(tok token.Token, args []object.Object) (int64, int64, int64, *object.Error) {
	var hour, minute, second int64

	if len(args) >= 4 {
		value, err := integerAt("date.of", tok, args, 3)

		if err != nil {
			return 0, 0, 0, err
		}

		if value < 0 || value > 23 {
			return 0, 0, 0, object.NewError(fault.Value, tok, "`date.of()` expects an hour between 0 and 23, got %d", value)
		}

		hour = value
	}

	if len(args) >= 5 {
		value, err := integerAt("date.of", tok, args, 4)

		if err != nil {
			return 0, 0, 0, err
		}

		if value < 0 || value > 59 {
			return 0, 0, 0, object.NewError(fault.Value, tok, "`date.of()` expects a minute between 0 and 59, got %d", value)
		}

		minute = value
	}

	if len(args) == 6 {
		value, err := integerAt("date.of", tok, args, 5)

		if err != nil {
			return 0, 0, 0, err
		}

		if value < 0 || value > 59 {
			return 0, 0, 0, object.NewError(fault.Value, tok, "`date.of()` expects a second between 0 and 59, got %d", value)
		}

		second = value
	}

	return hour, minute, second, nil
}

// dateParseISO reads an ISO 8601 timestamp, with or without a time-of-day
// component, the way date-fns's parseISO does.
func dateParseISO(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("date.parseISO", tok, args, 1); err != nil {
		return err
	}

	text, err := stringAt("date.parseISO", tok, args, 0)

	if err != nil {
		return err
	}

	if parsed, parseErr := time.Parse(time.RFC3339, text); parseErr == nil {
		return newDate(parsed)
	}

	if parsed, parseErr := time.Parse("2006-01-02", text); parseErr == nil {
		return newDate(parsed)
	}

	return object.NewError(fault.Value, tok, "`date.parseISO()` cannot read `%s` as an ISO 8601 date", text).
		WithHelp("an ISO date looks like 2024-01-15 or 2024-01-15T09:30:00Z")
}

func dateFromUnix(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("date.fromUnix", tok, args, 1); err != nil {
		return err
	}

	seconds, err := integerAt("date.fromUnix", tok, args, 0)

	if err != nil {
		return err
	}

	return newDate(time.Unix(seconds, 0))
}

func dateToUnix(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("date.toUnix", tok, args, 1); err != nil {
		return err
	}

	date, err := dateAt("date.toUnix", tok, args, 0)

	if err != nil {
		return err
	}

	return object.NewInt(date.Time.Unix())
}

func dateToUnixNano(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("date.toUnixNano", tok, args, 1); err != nil {
		return err
	}

	date, err := dateAt("date.toUnixNano", tok, args, 0)

	if err != nil {
		return err
	}

	return object.NewInt(date.Time.UnixNano())
}

func dateFormat(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("date.format", tok, args, 2); err != nil {
		return err
	}

	date, err := dateAt("date.format", tok, args, 0)

	if err != nil {
		return err
	}

	pattern, patternErr := stringAt("date.format", tok, args, 1)

	if patternErr != nil {
		return patternErr
	}

	return &object.String{Value: formatDate(date.Time, pattern)}
}

// =============================================================================
// Arithmetic

func dateAddDays(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	return dateShift("date.addDays", tok, args, func(t time.Time, n int64) time.Time { return t.AddDate(0, 0, int(n)) })
}

func dateSubDays(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	return dateShift("date.subDays", tok, args, func(t time.Time, n int64) time.Time { return t.AddDate(0, 0, -int(n)) })
}

func dateAddWeeks(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	return dateShift("date.addWeeks", tok, args, func(t time.Time, n int64) time.Time { return t.AddDate(0, 0, 7*int(n)) })
}

func dateSubWeeks(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	return dateShift("date.subWeeks", tok, args, func(t time.Time, n int64) time.Time { return t.AddDate(0, 0, -7*int(n)) })
}

func dateAddMonths(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	return dateShift("date.addMonths", tok, args, func(t time.Time, n int64) time.Time { return addCalendarMonths(t, int(n)) })
}

func dateSubMonths(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	return dateShift("date.subMonths", tok, args, func(t time.Time, n int64) time.Time { return addCalendarMonths(t, -int(n)) })
}

func dateAddYears(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	return dateShift("date.addYears", tok, args, func(t time.Time, n int64) time.Time { return addCalendarMonths(t, 12*int(n)) })
}

func dateSubYears(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	return dateShift("date.subYears", tok, args, func(t time.Time, n int64) time.Time { return addCalendarMonths(t, -12*int(n)) })
}

// addCalendarMonths shifts a date by a number of months the way date-fns
// does: when the day of month does not exist in the target month - January
// 31 plus one month has no February 31 - it clamps to the target month's
// last day rather than rolling over into the month after, the way naive
// day-arithmetic would turn January 31 into March 2 or 3.
func addCalendarMonths(t time.Time, months int) time.Time {
	year, month, day := t.Date()
	hour, minute, second := t.Clock()

	totalMonths := int(month) - 1 + months
	targetYear := year + totalMonths/12
	targetMonth := totalMonths % 12

	if targetMonth < 0 {
		targetMonth += 12
		targetYear--
	}

	lastDayOfTargetMonth := time.Date(targetYear, time.Month(targetMonth+2), 0, 0, 0, 0, 0, time.UTC).Day()

	if day > lastDayOfTargetMonth {
		day = lastDayOfTargetMonth
	}

	return time.Date(targetYear, time.Month(targetMonth+1), day, hour, minute, second, t.Nanosecond(), time.UTC)
}

func dateAddHours(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	return dateShift("date.addHours", tok, args, func(t time.Time, n int64) time.Time { return t.Add(time.Duration(n) * time.Hour) })
}

func dateAddMinutes(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	return dateShift("date.addMinutes", tok, args, func(t time.Time, n int64) time.Time { return t.Add(time.Duration(n) * time.Minute) })
}

func dateAddSeconds(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	return dateShift("date.addSeconds", tok, args, func(t time.Time, n int64) time.Time { return t.Add(time.Duration(n) * time.Second) })
}

// dateShift is what every add/sub function above is: read a date and a
// count, shift the date by the count, answer a new one. Writing the shared
// shape once keeps the eleven of them one line each.
func dateShift(name string, tok token.Token, args []object.Object, shift func(time.Time, int64) time.Time) object.Object {
	if err := arity(name, tok, args, 2); err != nil {
		return err
	}

	date, err := dateAt(name, tok, args, 0)

	if err != nil {
		return err
	}

	count, countErr := integerAt(name, tok, args, 1)

	if countErr != nil {
		return countErr
	}

	return newDate(shift(date.Time, count))
}

// =============================================================================
// Predicates

func dateIsSameDay(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("date.isSameDay", tok, args, 2); err != nil {
		return err
	}

	left, err := dateAt("date.isSameDay", tok, args, 0)

	if err != nil {
		return err
	}

	right, rightErr := dateAt("date.isSameDay", tok, args, 1)

	if rightErr != nil {
		return rightErr
	}

	sameYear := left.Time.Year() == right.Time.Year()
	sameMonth := left.Time.Month() == right.Time.Month()
	sameDay := left.Time.Day() == right.Time.Day()

	return &object.Boolean{Value: sameYear && sameMonth && sameDay}
}

func dateIsWeekend(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("date.isWeekend", tok, args, 1); err != nil {
		return err
	}

	date, err := dateAt("date.isWeekend", tok, args, 0)

	if err != nil {
		return err
	}

	weekday := date.Time.Weekday()

	return &object.Boolean{Value: weekday == time.Saturday || weekday == time.Sunday}
}

func dateIsLeapYear(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("date.isLeapYear", tok, args, 1); err != nil {
		return err
	}

	date, err := dateAt("date.isLeapYear", tok, args, 0)

	if err != nil {
		return err
	}

	year := date.Time.Year()
	leap := year%4 == 0 && (year%100 != 0 || year%400 == 0)

	return &object.Boolean{Value: leap}
}

// =============================================================================
// Differences

func dateDifferenceInDays(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	return dateDifference("date.differenceInDays", tok, args, 24*time.Hour)
}

func dateDifferenceInHours(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	return dateDifference("date.differenceInHours", tok, args, time.Hour)
}

func dateDifferenceInMinutes(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	return dateDifference("date.differenceInMinutes", tok, args, time.Minute)
}

func dateDifferenceInSeconds(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	return dateDifference("date.differenceInSeconds", tok, args, time.Second)
}

// dateDifference answers how many whole units of the given size separate two
// dates, truncated toward zero the way date-fns's differenceInX functions
// are - so differenceInDays(a, b) and differenceInDays(b, a) are negatives of
// each other, never off by one because of which way the truncation went.
func dateDifference(name string, tok token.Token, args []object.Object, unit time.Duration) object.Object {
	if err := arity(name, tok, args, 2); err != nil {
		return err
	}

	left, err := dateAt(name, tok, args, 0)

	if err != nil {
		return err
	}

	right, rightErr := dateAt(name, tok, args, 1)

	if rightErr != nil {
		return rightErr
	}

	elapsed := left.Time.Sub(right.Time)

	return object.NewInt(int64(elapsed / unit))
}

// =============================================================================
// Start and end of a period

func dateStartOfDay(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	date, err := dateForBoundary("date.startOfDay", tok, args)

	if err != nil {
		return err
	}

	return newDate(time.Date(date.Time.Year(), date.Time.Month(), date.Time.Day(), 0, 0, 0, 0, time.UTC))
}

func dateEndOfDay(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	date, err := dateForBoundary("date.endOfDay", tok, args)

	if err != nil {
		return err
	}

	return newDate(time.Date(date.Time.Year(), date.Time.Month(), date.Time.Day(), 23, 59, 59, 999999999, time.UTC))
}

func dateStartOfMonth(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	date, err := dateForBoundary("date.startOfMonth", tok, args)

	if err != nil {
		return err
	}

	return newDate(time.Date(date.Time.Year(), date.Time.Month(), 1, 0, 0, 0, 0, time.UTC))
}

func dateEndOfMonth(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	date, err := dateForBoundary("date.endOfMonth", tok, args)

	if err != nil {
		return err
	}

	startOfNextMonth := time.Date(date.Time.Year(), date.Time.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, 0)

	return newDate(startOfNextMonth.Add(-time.Nanosecond))
}

func dateForBoundary(name string, tok token.Token, args []object.Object) (*object.Date, *object.Error) {
	if err := arity(name, tok, args, 1); err != nil {
		return nil, err
	}

	return dateAt(name, tok, args, 0)
}

// =============================================================================
// Components

func dateYear(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	return dateComponent("date.year", tok, args, func(t time.Time) int64 { return int64(t.Year()) })
}

func dateMonth(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	return dateComponent("date.month", tok, args, func(t time.Time) int64 { return int64(t.Month()) })
}

func dateDay(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	return dateComponent("date.day", tok, args, func(t time.Time) int64 { return int64(t.Day()) })
}

func dateHour(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	return dateComponent("date.hour", tok, args, func(t time.Time) int64 { return int64(t.Hour()) })
}

func dateMinute(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	return dateComponent("date.minute", tok, args, func(t time.Time) int64 { return int64(t.Minute()) })
}

func dateSecond(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	return dateComponent("date.second", tok, args, func(t time.Time) int64 { return int64(t.Second()) })
}

// dateWeekday answers the day of the week as a number from 0 (Sunday) to 6
// (Saturday), the same numbering date-fns's getDay uses.
func dateWeekday(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	return dateComponent("date.weekday", tok, args, func(t time.Time) int64 { return int64(t.Weekday()) })
}

func dateComponent(name string, tok token.Token, args []object.Object, read func(time.Time) int64) object.Object {
	if err := arity(name, tok, args, 1); err != nil {
		return err
	}

	date, err := dateAt(name, tok, args, 0)

	if err != nil {
		return err
	}

	return object.NewInt(read(date.Time))
}

// =============================================================================
// Formatting

var weekdayNames = [...]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
var monthNames = [...]string{
	"January", "February", "March", "April", "May", "June",
	"July", "August", "September", "October", "November", "December",
}

// formatDate renders a date against a date-fns-style pattern: a run of the
// same letter is one token (`yyyy`, `MM`, `EEEE`), and any character that is
// not part of a recognized run - `-`, `:`, a space, `T` - is copied through
// literally. This is deliberately not Go's reference-date layout string,
// which reads like a specific date rather than describing the shape of one.
func formatDate(t time.Time, pattern string) string {
	var out strings.Builder

	runes := []rune(pattern)
	index := 0

	for index < len(runes) {
		letter := runes[index]
		end := index

		for end < len(runes) && runes[end] == letter {
			end++
		}

		run := end - index

		switch letter {
		case 'y':
			out.WriteString(formatYear(t.Year(), run))
		case 'M':
			out.WriteString(formatMonth(int(t.Month()), run))
		case 'd':
			out.WriteString(formatPadded(t.Day(), run))
		case 'E':
			out.WriteString(formatWeekday(t.Weekday(), run))
		case 'H':
			out.WriteString(formatPadded(t.Hour(), run))
		case 'h':
			out.WriteString(formatPadded(hour12(t.Hour()), run))
		case 'm':
			out.WriteString(formatPadded(t.Minute(), run))
		case 's':
			out.WriteString(formatPadded(t.Second(), run))
		case 'a':
			out.WriteString(amPm(t.Hour()))
		default:
			out.WriteString(string(runes[index:end]))
		}

		index = end
	}

	return out.String()
}

func formatYear(year int, run int) string {
	switch {
	case run >= 4:
		return fmt.Sprintf("%04d", year)
	case run == 2:
		return fmt.Sprintf("%02d", year%100)
	default:
		return strconv.Itoa(year)
	}
}

func formatMonth(month int, run int) string {
	switch {
	case run >= 4:
		return monthNames[month-1]
	case run == 3:
		return monthNames[month-1][:3]
	case run == 2:
		return fmt.Sprintf("%02d", month)
	default:
		return strconv.Itoa(month)
	}
}

func formatWeekday(day time.Weekday, run int) string {
	if run >= 4 {
		return weekdayNames[day]
	}

	return weekdayNames[day][:3]
}

// formatPadded renders a two-digit clock or calendar field: zero-padded when
// the pattern asked for two letters or more, bare otherwise.
func formatPadded(value int, run int) string {
	if run >= 2 {
		return fmt.Sprintf("%02d", value)
	}

	return strconv.Itoa(value)
}

func hour12(hour int) int {
	twelve := hour % 12

	if twelve == 0 {
		twelve = 12
	}

	return twelve
}

func amPm(hour int) string {
	if hour < 12 {
		return "AM"
	}

	return "PM"
}
