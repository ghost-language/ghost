package object

// STANDARDFUNCTION represents the object's type.
const STANDARDFUNCTION = "STANDARDFUNCTION"

type Standard struct {
	Name     string
	Function StandardFunction
}

type StandardFunction func(args ...Object) Object

// String represents the string form of the standard function object.
func (sf *StandardFunction) String() string {
	return "standard function"
}

// Type returns the standard function object type.
func (sf *StandardFunction) Type() Type {
	return STANDARDFUNCTION
}
