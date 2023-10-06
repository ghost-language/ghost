package object

import "fmt"

const BOOLEAN = "BOOLEAN"

// Boolean objects consist of a boolean value.
type Boolean struct {
	Value bool
}

// String represents the boolean object's value as a string.
func (boolean *Boolean) String() string {
	return fmt.Sprintf("%t", boolean.Value)
}

// Type returns the boolean object type.
func (boolean *Boolean) Type() Type {
	return BOOLEAN
}

// MapKey defines a unique hash value for use as a map key.
func (boolean *Boolean) MapKey() MapKey {
	var value uint64

	if boolean.Value {
		value = 1
	} else {
		value = 0
	}

	return MapKey{Type: boolean.Type(), Value: value}
}

// Method defines the set of methods available on boolean objects.
func (boolean *Boolean) Method(method string, args []Object) (Object, bool) {
	return nil, false
}

func IsTrue(obj Object) bool {
	return isTruthy(obj)
}

func IsFalse(obj Object) bool {
	return !isTruthy(obj)
}

func isTruthy(value Object) bool {
	switch value := value.(type) {
	case *Null:
		return false
	case *Boolean:
		return value.Value
	case *String:
		return len(value.Value) > 0
	default:
		return true
	}
}
