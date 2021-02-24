package object

// ERROR represents the object's type.
const ERROR = "ERROR"

// Error objects consist of a null value.
type Error struct {
	Message string
}

// String represents the string form of the null object.
func (e *Error) String() string {
	return "Runtime error: " + e.Message
}

// Type returns the null object type.
func (e *Error) Type() Type {
	return ERROR
}
