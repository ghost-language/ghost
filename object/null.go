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
	switch method {
	case "toString":
		return null.toString(tok, args)
	}

	return nil, false
}

func (null *Null) toString(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("null.toString()", tok, args, 0); err != nil {
		return err, true
	}

	return &String{Value: null.String()}, true
}
