package object

import "fmt"

const BOOLEAN = "BOOLEAN"

// Boolean objects consist of a boolean value.
type Boolean struct {
	Value bool
}

func (object *Boolean) Accept(v Visitor) {
	v.visitBoolean(object)
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
