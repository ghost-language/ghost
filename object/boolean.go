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
