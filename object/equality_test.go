package object

import (
	"testing"
	"time"

	"ghostlang.org/x/ghost/token"
)

// stubObject is a minimal Object with no case of its own in ValuesEqual, so
// two of them exercise the identity fallback the way a Function, Class, or
// Instance would - without needing to build a full closure or class chain
// just to ask "are these the same one". The unused field keeps two instances
// from sharing an address - Go may otherwise give every zero-size value the
// same one, which would make two separately built stubs compare equal for a
// reason that has nothing to do with what this test is checking.
type stubObject struct{ _ int }

func (s *stubObject) Type() Type                                          { return Type(9999) }
func (s *stubObject) String() string                                      { return "stub" }
func (s *stubObject) Method(string, token.Token, []Object) (Object, bool) { return nil, false }

func TestValuesEqualScalarsAndNull(t *testing.T) {
	tests := []struct {
		name     string
		left     Object
		right    Object
		expected bool
	}{
		{"equal numbers", NewInt(5), NewInt(5), true},
		{"int equals equivalent float", NewInt(5), NewFloat(5), true},
		{"different numbers", NewInt(5), NewInt(6), false},
		{"equal strings", &String{Value: "a"}, &String{Value: "a"}, true},
		{"different strings", &String{Value: "a"}, &String{Value: "b"}, false},
		{"equal booleans", &Boolean{Value: true}, &Boolean{Value: true}, true},
		{"different booleans", &Boolean{Value: true}, &Boolean{Value: false}, false},
		{"null equals null", &Null{}, &Null{}, true},
		{"different types are never equal", NewInt(5), &String{Value: "5"}, false},
		{"nil equals nil", nil, nil, true},
		{"nil does not equal a value", nil, NewInt(5), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValuesEqual(tt.left, tt.right); got != tt.expected {
				t.Errorf("ValuesEqual(%v, %v) = %v, expected %v", tt.left, tt.right, got, tt.expected)
			}
		})
	}
}

func TestValuesEqualLists(t *testing.T) {
	a := &List{Elements: []Object{NewInt(1), NewInt(2)}}
	b := &List{Elements: []Object{NewInt(1), NewInt(2)}}
	c := &List{Elements: []Object{NewInt(1), NewInt(3)}}
	shorter := &List{Elements: []Object{NewInt(1)}}
	nested := &List{Elements: []Object{&List{Elements: []Object{NewInt(1)}}}}
	nestedSame := &List{Elements: []Object{&List{Elements: []Object{NewInt(1)}}}}

	if !ValuesEqual(a, b) {
		t.Error("lists with the same contents should be equal")
	}

	if ValuesEqual(a, c) {
		t.Error("lists with different contents should not be equal")
	}

	if ValuesEqual(a, shorter) {
		t.Error("lists of different lengths should not be equal")
	}

	if !ValuesEqual(nested, nestedSame) {
		t.Error("nested lists should compare to any depth")
	}
}

// TestValuesEqualMaps covers §13.2: maps used to have no case in ValuesEqual
// at all, and fell to the identity fallback below - two maps built separately
// with identical contents compared unequal. They now compare like List: by
// contents, to any depth, key order irrelevant.
func TestValuesEqualMaps(t *testing.T) {
	key := func(s string) MapKey { return (&String{Value: s}).MapKey() }

	a := &Map{Pairs: map[MapKey]MapPair{
		key("x"): {Key: &String{Value: "x"}, Value: NewInt(1)},
		key("y"): {Key: &String{Value: "y"}, Value: NewInt(2)},
	}}
	b := &Map{Pairs: map[MapKey]MapPair{
		key("y"): {Key: &String{Value: "y"}, Value: NewInt(2)},
		key("x"): {Key: &String{Value: "x"}, Value: NewInt(1)},
	}}
	differentValue := &Map{Pairs: map[MapKey]MapPair{
		key("x"): {Key: &String{Value: "x"}, Value: NewInt(1)},
		key("y"): {Key: &String{Value: "y"}, Value: NewInt(99)},
	}}
	fewerKeys := &Map{Pairs: map[MapKey]MapPair{
		key("x"): {Key: &String{Value: "x"}, Value: NewInt(1)},
	}}
	nested := &Map{Pairs: map[MapKey]MapPair{
		key("inner"): {Key: &String{Value: "inner"}, Value: &Map{Pairs: map[MapKey]MapPair{
			key("z"): {Key: &String{Value: "z"}, Value: NewInt(3)},
		}}},
	}}
	nestedSame := &Map{Pairs: map[MapKey]MapPair{
		key("inner"): {Key: &String{Value: "inner"}, Value: &Map{Pairs: map[MapKey]MapPair{
			key("z"): {Key: &String{Value: "z"}, Value: NewInt(3)},
		}}},
	}}

	if !ValuesEqual(a, b) {
		t.Error("maps with the same contents should be equal regardless of insertion order")
	}

	if ValuesEqual(a, differentValue) {
		t.Error("maps with a differing value should not be equal")
	}

	if ValuesEqual(a, fewerKeys) {
		t.Error("maps with different key counts should not be equal")
	}

	if !ValuesEqual(nested, nestedSame) {
		t.Error("nested maps should compare to any depth")
	}
}

func TestValuesEqualDurations(t *testing.T) {
	a := &Duration{Years: 1, Months: 2, Days: 3}
	b := &Duration{Years: 1, Months: 2, Days: 3}
	c := &Duration{Years: 1, Months: 2, Days: 4}

	if !ValuesEqual(a, b) {
		t.Error("durations with the same components should be equal")
	}

	if ValuesEqual(a, c) {
		t.Error("durations with different components should not be equal")
	}
}

// TestValuesEqualDates confirms ValuesEqual reaches the same answer
// evaluateDateInfix does for a direct `date1 == date2`: the instant, not the
// attached zone - a Date reads as equal from list.contains()/unique() under
// exactly the same rule `==` already uses.
func TestValuesEqualDates(t *testing.T) {
	instant := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	utc := &Date{Time: instant}
	sameInstantInEST := &Date{Time: instant.In(time.FixedZone("EST", -5*60*60))}
	different := &Date{Time: instant.Add(time.Hour)}

	if !ValuesEqual(utc, sameInstantInEST) {
		t.Error("dates naming the same instant should be equal regardless of attached zone")
	}

	if ValuesEqual(utc, different) {
		t.Error("dates naming different instants should not be equal")
	}
}

// TestValuesEqualIdentityFallback covers everything ValuesEqual has no
// content-comparison case for - instances, functions, classes, and the rest:
// the same object is equal to itself, and two separately built ones are not,
// mirroring how `==` already treated Instance before §13.2 gave every other
// such type the same treatment instead of a type error.
func TestValuesEqualIdentityFallback(t *testing.T) {
	a := &stubObject{}
	b := &stubObject{}

	if !ValuesEqual(a, a) {
		t.Error("a value should equal itself")
	}

	if ValuesEqual(a, b) {
		t.Error("two separately built values should not be equal by default")
	}
}
