package object

const STRING = "STRING"

// String objects consist of a string value.
type String struct {
	Value string
}

// String represents the string object's value as a string. So meta.
func (string *String) String() string {
	return string.Value
}

// Type returns the string object type.
func (string *String) Type() Type {
	return STRING
}
