package object

import "ghostlang.org/x/ghost/token"

// Continue objects consist of a nil value.
type Continue struct{}

// String represents the continue object's value as a string.
func (obj *Continue) String() string {
	return "continue"
}

// Type returns the continue object type.
func (obj *Continue) Type() Type {
	return CONTINUE
}

// Method defines the set of methods available on continue objects.
func (obj *Continue) Method(method string, tok token.Token, args []Object) (Object, bool) {
	return nil, false
}
