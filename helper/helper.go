package helper

import (
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

// NativeBooleanToObject converts a native Go boolean value to a Ghost boolean
// value.
func NativeBooleanToObject(boolean bool) *object.Boolean {
	if boolean {
		return value.TRUE
	}

	return value.FALSE
}

// IsEqual determines if the based objects are of equal value.
func IsEqual(left object.Object, right object.Object) bool {
	switch {
	case left.Type() == object.NUMBER && right.Type() == object.NUMBER:
		return left.(*object.Number).Value.Equal(right.(*object.Number).Value)
	case left.Type() == object.STRING && right.Type() == object.STRING:
		return left.(*object.String).Value == right.(*object.String).Value
	case left.Type() == object.NULL && right.Type() == object.NULL:
		return true
	}

	return false
}

// IsTruthy returns the truthy value of the passed object.
func IsTruthy(obj object.Object) bool {
	switch obj {
	case value.NULL:
		return false
	case value.TRUE:
		return true
	case value.FALSE:
		return false
	default:
		return true
	}
}
