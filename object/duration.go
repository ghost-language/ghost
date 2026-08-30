package object

import (
	"fmt"
	"strings"

	"ghostlang.org/x/ghost/token"
)

// Duration objects hold a calendar-and-clock span - years, months, days,
// hours, minutes, seconds - rather than a single reduced number: "1 month"
// has no fixed length (28-31 days), so a Duration keeps its components
// separate instead of collapsing them into one unit, the way
// date.differenceInDays/Hours/Minutes/Seconds deliberately do for the
// simpler, single-unit question. This mirrors Temporal.Duration and
// date-fns's intervalToDuration rather than a raw millisecond count.
//
// A Duration's non-zero fields all share one sign - it names a single
// direction, never a mix of "2 months forward, 3 days back" - enforced at
// construction by date.duration() (§9.5). date.durationBetween() is what
// computes one from two Dates; addDuration/subDuration are the one way to
// apply one back to a Date - see library/modules/date.go's package doc
// comment for why arithmetic is a module function rather than a method here
// too.
type Duration struct {
	Years   int64
	Months  int64
	Days    int64
	Hours   int64
	Minutes int64
	Seconds int64
}

// String renders an ISO 8601 duration - PnYnMnDTnHnMnS, only the non-zero
// components, "PT0S" for a zero duration - the same standard family
// Date.String() uses (RFC3339) for the same reason: the string alone is
// enough to reconstruct the value.
func (d *Duration) String() string {
	if d.Years == 0 && d.Months == 0 && d.Days == 0 && d.Hours == 0 && d.Minutes == 0 && d.Seconds == 0 {
		return "PT0S"
	}

	var out strings.Builder
	out.WriteByte('P')

	writeDurationUnit(&out, d.Years, 'Y')
	writeDurationUnit(&out, d.Months, 'M')
	writeDurationUnit(&out, d.Days, 'D')

	if d.Hours != 0 || d.Minutes != 0 || d.Seconds != 0 {
		out.WriteByte('T')
		writeDurationUnit(&out, d.Hours, 'H')
		writeDurationUnit(&out, d.Minutes, 'M')
		writeDurationUnit(&out, d.Seconds, 'S')
	}

	return out.String()
}

func writeDurationUnit(out *strings.Builder, value int64, suffix byte) {
	if value == 0 {
		return
	}

	fmt.Fprintf(out, "%d%c", value, suffix)
}

// Type returns the duration object type.
func (d *Duration) Type() Type {
	return DURATION
}

// Method defines the set of methods available on duration objects: reading
// back each component, and toString. Everything that builds or applies a
// Duration lives in the date module instead - see this type's doc comment.
func (d *Duration) Method(method string, tok token.Token, args []Object) (Object, bool) {
	switch method {
	case "years":
		return d.component(tok, args, "duration.years()", d.Years)
	case "months":
		return d.component(tok, args, "duration.months()", d.Months)
	case "days":
		return d.component(tok, args, "duration.days()", d.Days)
	case "hours":
		return d.component(tok, args, "duration.hours()", d.Hours)
	case "minutes":
		return d.component(tok, args, "duration.minutes()", d.Minutes)
	case "seconds":
		return d.component(tok, args, "duration.seconds()", d.Seconds)
	case "toString":
		return d.toString(tok, args)
	}

	return nil, false
}

func (d *Duration) component(tok token.Token, args []Object, name string, value int64) (Object, bool) {
	if err := Arity(name, tok, args, 0); err != nil {
		return err, true
	}

	return NewInt(value), true
}

func (d *Duration) toString(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("duration.toString()", tok, args, 0); err != nil {
		return err, true
	}

	return &String{Value: d.String()}, true
}
