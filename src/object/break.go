package object

const BREAK = "BREAK"

// Break objects consist of a nil value.
type Break struct{}

// String represents the break object's value as a string.
func (obj *Break) String() string {
	return "break"
}

// Type returns the break object type.
func (obj *Break) Type() Type {
	return BREAK
}

// Method defines the set of methods available on break objects.
func (obj *Break) Method(method string, args []Object) (Object, bool) {
	return nil, false
}
