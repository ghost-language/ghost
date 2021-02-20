package object

import "fmt"

// BOOLEAN represents the object's type.
const BOOLEAN = "BOOLEAN"

// Boolean objects consist of a boolean value.
type Boolean struct {
	Value bool
}

// String represents the string form of the boolean object.
func (b *Boolean) String() string {
	return fmt.Sprintf("%t", b.Value)
}

// Type returns the boolean object type.
func (b *Boolean) Type() Type {
	return NUMBER
}
