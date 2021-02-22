package object

// STRING represents the object's type.
const STRING = "STRING"

// String objects consist of a string value.
type String struct {
	Value string
}

// String represents the string form of the string object.
func (n *String) String() string {
	return n.Value
}

// Type returns the number object type.
func (n *String) Type() Type {
	return STRING
}
