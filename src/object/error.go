package object

const ERROR = "ERROR"

// Error objects consist of a nil value.
type Error struct {
	Message string
}

// String represents the error object's value as a string.
func (err *Error) String() string {
	return "error"
}

// Type returns the error object type.
func (err *Error) Type() Type {
	return ERROR
}

// Method defines the set of methods available on error objects.
func (err *Error) Method(method string, args []Object) (Object, bool) {
	return nil, false
}
