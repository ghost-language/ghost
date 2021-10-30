package object

const NULL = "NULL"

// Null objects consist of a nil value.
type Null struct{}

// String represents the null object's value as a string.
func (null *Null) String() string {
	return "null"
}

// Type returns the null object type.
func (null *Null) Type() Type {
	return NULL
}
