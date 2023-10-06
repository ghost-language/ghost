package object

const RETURN = "RETURN"

// Return objects consist of a value.
type Return struct {
	Value Object
}

// String represents the return object's value as a string.
func (obj *Return) String() string {
	return "return"
}

// Type returns the return object type.
func (obj *Return) Type() Type {
	return RETURN
}

// Method defines the set of methods available on return objects.
func (obj *Return) Method(method string, args []Object) (Object, bool) {
	return nil, false
}
