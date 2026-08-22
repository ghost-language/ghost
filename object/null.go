package object

import "ghostlang.org/x/ghost/token"

// Null objects consist of a nil value.
type Null struct{}

// String represents the null object's value as a string.
func (null *Null) String() string {
	return "null"
}

// Type returns the null object type.
func (null *Null) Type() Type {
	return NULL
}

// Method defines the set of methods available on null objects.
func (null *Null) Method(method string, tok token.Token, args []Object) (Object, bool) {
	return nil, false
}
