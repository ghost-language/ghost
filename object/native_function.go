package object

// NATIVE_FUNCTION represents the object's type.
const NATIVE_FUNCTION = "NATIVE_FUNCTION"

type NativeFunction struct {
	Name string
	Function GhostFunction
}

// String represents the string form of the native function object.
func (nf *NativeFunction) String() string {
	return nf.Name
}

// Type returns the native function object type.
func (nf *NativeFunction) Type() Type {
	return NATIVE_FUNCTION
}