package object

// NULL represents the object's type.
const NULL = "NULL"

// Null objects consist of a null value.
type Null struct{}

// String represents the string form of the null object.
func (n *Null) String() string {
	return "null"
}

// Type returns the null object type.
func (n *Null) Type() Type {
	return NULL
}
