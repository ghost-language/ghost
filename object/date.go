package object

import (
	"time"

	"ghostlang.org/x/ghost/token"
)

// Date objects hold a single instant together with the time zone it should
// be read in, defaulting to UTC. The instant is what `<`, `>`, and `==`
// compare (evaluateDateInfix) and stays independent of the attached zone, so
// a comparison is reproducible no matter where the program runs, the same
// guarantee a seeded random run has. The zone only governs what reading a
// *calendar* position out of that instant answers - String() included - the
// date module's inTimeZone/ofInZone are the one way to change it (§9.5).
//
// A Date does not support arithmetic operators. What `date1 + date2` or
// `date * 2` would even mean is not obvious the way `date1 < date2` is, so
// that stays out - the date module's addDays, subDays, differenceInDays, and
// the rest are the one way to do date arithmetic, rather than leaving an
// operator and a function that might disagree.
type Date struct {
	Time time.Time
}

// String represents the date as an ISO 8601 timestamp in whatever zone the
// Date is attached to - "Z" for the UTC default, an explicit offset
// otherwise (RFC3339 always writes one or the other, never a bare local
// time), so the string alone is always enough to reconstruct the instant.
func (date *Date) String() string {
	return date.Time.Format(time.RFC3339)
}

// Type returns the date object type.
func (date *Date) Type() Type {
	return DATE
}

// MapKey defines a unique hash value for use as a map key.
func (date *Date) MapKey() MapKey {
	return MapKey{Type: date.Type(), Value: uint64(date.Time.UnixNano())}
}

// Method defines the set of methods available on date objects. Everything
// that reads or transforms a date beyond this lives in the date module - see
// object.Date's doc comment for why arithmetic is not here either.
func (date *Date) Method(method string, tok token.Token, args []Object) (Object, bool) {
	switch method {
	case "toString":
		return date.toString(tok, args)
	}

	return nil, false
}

func (date *Date) toString(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("date.toString()", tok, args, 0); err != nil {
		return err, true
	}

	return &String{Value: date.String()}, true
}
