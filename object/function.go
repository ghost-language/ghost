package object

// FUNCTION represents the object's type.
const FUNCTION = "FUNCTION"

// Function objects consist of a function value.
type Function struct{}

// String represents the string form of the function object.
func (f *Function) String() string {
	return "function"
}

// Type returns the function object type.
func (f *Function) Type() Type {
	return FUNCTION
}
