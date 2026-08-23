package object

import (
	"time"

	"ghostlang.org/x/ghost/token"
)

// Date objects hold a single instant, always in UTC. Ghost does not model
// time zones - every Date is the same instant no matter where the program
// runs, which keeps a seeded, procedurally generated, or scheduled result
// reproducible the same way a seeded random run is.
//
// A Date does not support arithmetic operators. What `date1 + date2` or
// `date * 2` would even mean is not obvious the way `date1 < date2` is, so
// that stays out - the date module's addDays, subDays, differenceInDays, and
// the rest are the one way to do date arithmetic, rather than leaving an
// operator and a function that might disagree.
type Date struct {
	Time time.Time
}

// String represents the date as an ISO 8601 timestamp, the one format that
// reads the same regardless of where it is printed.
func (date *Date) String() string {
	return date.Time.UTC().Format(time.RFC3339)
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
