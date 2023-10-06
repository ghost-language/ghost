package object

import (
	"fmt"
)

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

// IsError determines if the referenced object is an error.
func IsError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR
	}

	return false
}

func NewError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}
